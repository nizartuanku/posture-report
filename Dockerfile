FROM golang:1.24 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=1 go build -ldflags "-s -w" -o /posturereport ./cmd/posturereport
FROM debian:bookworm-slim
COPY --from=build /posturereport /usr/local/bin/posturereport
EXPOSE 8432
ENTRYPOINT ["posturereport","-listen","0.0.0.0:8432"]
