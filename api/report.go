// Package api exposes cluster health over HTTP as JSON.
package api

import (
	"time"

	"github.com/vmagueta/clusterwatch/cluster"
)

// reportResponse is the JSON shape of a cluster.Report. It is kept apart from
// the domain type so the wire format can stay stable while cluster evolves.
type reportResponse struct {
	CheckedAt time.Time              `json:"checked_at"`
	Counts    map[cluster.Status]int `json:"counts"`
	Results   []resultResponse       `json:"results"`
}

// resultResponse is the JSON shape of a single cluster.Result.
type resultResponse struct {
	Node      string         `json:"node"`
	Status    cluster.Status `json:"status"`
	LatencyMS int64          `json:"latency_ms"`
	Error     string         `json:"error,omitempty"`
}

// newReportResponse converts a domain report into its JSON shape.
func newReportResponse(r cluster.Report) reportResponse {
	results := make([]resultResponse, len(r.Results))
	for i, res := range r.Results {
		results[i] = resultResponse{
			Node:      res.Node.ID,
			Status:    res.Status,
			LatencyMS: res.Latency.Milliseconds(),
		}
		if res.Err != nil {
			results[i].Error = res.Err.Error()
		}
	}

	return reportResponse{
		CheckedAt: r.CheckedAt,
		Counts:    r.Counts,
		Results:   results,
	}
}
