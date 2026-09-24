package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nizartuanku/posture-report/core"
	"github.com/nizartuanku/posture-report/posture"
)

func aiReport() posture.Report {
	f := core.Finding{Module: "certwatch", Fingerprint: "fp1", Target: "mail.example.com:443",
		Check: "tls.expired", Title: "Certificate expired", Severity: core.SeverityCritical,
		Remediation: "Replace the certificate.", Evidence: map[string]any{"days": 3, "session_key": "SECRET"}}
	return posture.Build("Acme", []posture.Item{{Product: "CertLight", Module: "certwatch", Findings: []core.Finding{f}}}, time.Now())
}

func fakeAISidecar(t *testing.T, calls *atomic.Int32, lastBody *atomic.Value) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var body strings.Builder
		b := make([]byte, 1<<16)
		n, _ := r.Body.Read(b)
		body.Write(b[:n])
		lastBody.Store(body.String())
		content, _ := json.Marshal(map[string]any{"explanation": "It expired.", "what_to_verify": []string{"Check renewal."}, "disclaimer": "x"})
		json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]any{"role": "assistant", "content": string(content)}}}})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func post(t *testing.T, url, body string) (int, explainResponse) {
	t.Helper()
	resp, err := http.Post(url+"/api/findings/explain", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out explainResponse
	json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func TestPostureAI_ExplainPriorityAndGating(t *testing.T) {
	var calls atomic.Int32
	var last atomic.Value
	side := fakeAISidecar(t, &calls, &last)

	// Off by default.
	off := httptest.NewServer((&Server{Load: aiReport, Tier: "free"}).Handler())
	defer off.Close()
	if code, out := post(t, off.URL, `{"fingerprint":"fp1"}`); code != 200 || out.Available {
		t.Fatalf("AI off: got %d %+v", code, out)
	}

	// Free + unkeyed sidecar works; evidence is sanitised; priorities listed.
	ai, err := NewAIAssist(AIConfig{URL: side.URL})
	if err != nil {
		t.Fatal(err)
	}
	free := httptest.NewServer((&Server{Load: aiReport, Tier: "free", AI: ai}).Handler())
	defer free.Close()
	code, out := post(t, free.URL, `{"fingerprint":"fp1"}`)
	if code != 200 || !out.Available || out.Disclaimer != "AI-generated summary — verify against raw findings" {
		t.Fatalf("free+sidecar: got %d %+v", code, out)
	}
	body, _ := last.Load().(string)
	if strings.Contains(body, "SECRET") || !strings.Contains(body, "certwatch") {
		t.Errorf("unexpected packet sent: %s", body)
	}
	if code, _ := post(t, free.URL, `{"fingerprint":"nope"}`); code != http.StatusNotFound {
		t.Errorf("unknown fingerprint: got %d, want 404", code)
	}
	resp, _ := http.Get(free.URL + "/api/priorities")
	var pr []priorityJSON
	json.NewDecoder(resp.Body).Decode(&pr)
	resp.Body.Close()
	if len(pr) != 1 || pr[0].Fingerprint != "fp1" || pr[0].Module != "certwatch" {
		t.Errorf("priorities = %+v", pr)
	}

	// Keyed endpoint: refused on free, allowed on team.
	kf := filepath.Join(t.TempDir(), "k")
	os.WriteFile(kf, []byte("0123456789abcdef0123456789abcdef"), 0o600)
	keyed, _ := NewAIAssist(AIConfig{URL: side.URL, KeyFile: kf})
	before := calls.Load()
	fk := httptest.NewServer((&Server{Load: aiReport, Tier: "free", AI: keyed}).Handler())
	defer fk.Close()
	if _, out := post(t, fk.URL, `{"fingerprint":"fp1"}`); out.Available || !strings.Contains(out.Reason, "Pro or Team") {
		t.Errorf("free+keyed must be refused, got %+v", out)
	}
	if calls.Load() != before {
		t.Error("free tier must not contact a keyed endpoint")
	}
	tk := httptest.NewServer((&Server{Load: aiReport, Tier: "team", AI: keyed}).Handler())
	defer tk.Close()
	if _, out := post(t, tk.URL, `{"fingerprint":"fp1"}`); !out.Available {
		t.Errorf("team+keyed must work, got %+v", out)
	}
}
