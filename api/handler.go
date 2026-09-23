package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/vmagueta/clusterwatch/cluster"
)

// ReportHandler returns a handler that probes the given nodes on every request
// and responds with the resulting report as JSON.
func ReportHandler(c cluster.Checker, nodes []cluster.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		report := cluster.CheckAll(c, nodes)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(newReportResponse(report)); err != nil {
			// The status line is already sent once the body starts, so the
			// failure can only be logged not reported to the client.
			slog.Error("encode report", "err", err)
		}
	}
}
