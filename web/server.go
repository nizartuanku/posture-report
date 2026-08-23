// Package web serves Posture Report's small dashboard: a landing page with the
// current score and per-tool breakdown, and the full two-model report at
// /report (which prints straight to PDF). Everything is re-read live from the
// source databases on each request, so the report is always current.
package web

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/nizartuanku/posture-report/posture"
	"github.com/nizartuanku/posture-report/report"
)

//go:embed static
var staticFS embed.FS

// Server renders the dashboard and report from a live report provider.
type Server struct {
	// Load re-reads the sources and builds a fresh report. Called per request.
	Load   func() posture.Report
	Tier   string
	Notice string
}

// Handler builds the HTTP handler.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/summary", s.handleSummary)
	mux.HandleFunc("GET /report", s.handleReport)
	mux.HandleFunc("GET /report.html", s.handleDownload)
	sub, _ := fs.Sub(staticFS, "static")
	mux.Handle("GET /", http.FileServer(http.FS(sub)))
	return mux
}

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	rep := s.Load()
	type src struct {
		Product string `json:"product"`
		Open    int    `json:"open"`
	}
	out := struct {
		Title     string         `json:"title"`
		Tier      string         `json:"tier"`
		Notice    string         `json:"notice,omitempty"`
		Score     int            `json:"score"`
		Rating    string         `json:"rating"`
		OpenTotal int            `json:"open_total"`
		Tools     int            `json:"tools"`
		Manual    int            `json:"manual"`
		Counts    map[string]int `json:"counts"`
		Sources   []src          `json:"sources"`
		Generated string         `json:"generated"`
	}{
		Title: rep.Title, Tier: s.Tier, Notice: s.Notice,
		Score: rep.Score, Rating: rep.Rating, OpenTotal: rep.OpenTotal,
		Tools: len(rep.Sources), Manual: len(rep.ManualReviews),
		Counts:    map[string]int{},
		Generated: rep.GeneratedAt.Format("2 Jan 2006 15:04"),
	}
	for k, v := range rep.Counts {
		out.Counts[k] = v
	}
	for _, x := range rep.Sources {
		out.Sources = append(out.Sources, src{Product: x.Product, Open: x.Open})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(report.HTML(s.Load())))
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="posture-report.html"`)
	w.Write([]byte(report.HTML(s.Load())))
}
