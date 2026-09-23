package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/vmagueta/clusterwatch/api"
	"github.com/vmagueta/clusterwatch/cluster"
)

func main() {
	checker := cluster.TCPChecker{
		Timeout:       2 * time.Second,
		DegradedAfter: 200 * time.Millisecond,
	}

	nodes := []cluster.Node{
		{ID: "google", Addr: "google.com:80"},
		{ID: "cloudflare", Addr: "1.1.1.1:443"},
		{ID: "dead", Addr: "127.0.0.1:9999"},
	}

	mux := http.NewServeMux()
	mux.Handle("GET /report", api.ReportHandler(checker, nodes))

	addr := ":8080"
	slog.Info("listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
