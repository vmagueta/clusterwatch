package api

import (
      "encoding/json"
      "errors"
      "net/http"
      "net/http/httptest"
      "testing"

      "github.com/vmagueta/clusterwatch/cluster"
)

// stubChecker reports the node with ID "down" as unreachable and every other
// node as healthy.
type stubChecker struct{}

func (stubChecker) Check(n cluster.Node) cluster.Result {
      if n.ID == "down" {
              return cluster.Result{Node: n, Status: cluster.StatusUnreachable, Err: errors.New("connection refused")}
      }
      return cluster.Result{Node: n, Status: cluster.StatusHealthy}
}

func TestReportHandlerServesJSON(t *testing.T) {
      handler := ReportHandler(stubChecker{}, []cluster.Node{{ID: "up"}, {ID: "down"}})

      req := httptest.NewRequest(http.MethodGet, "/report", nil)
      rec := httptest.NewRecorder()
      handler.ServeHTTP(rec, req)

      if rec.Code != http.StatusOK {
              t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
      }
      if got := rec.Header().Get("Content-Type"); got != "application/json" {
              t.Errorf("Content-Type = %q, want %q", got, "application/json")
      }

      // Decode into plain strings, the way a client would read the wire format.
      var body struct {
              Counts  map[string]int `json:"counts"`
              Results []struct {
                      Node   string `json:"node"`
                      Status string `json:"status"`
                      Error  string `json:"error"`
              } `json:"results"`
      }
      if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
              t.Fatalf("decode body: %v", err)
      }

      if got := body.Counts["unreachable"]; got != 1 {
              t.Errorf(`counts["unreachable"] = %d, want 1`, got)
      }
      if got := body.Results[1].Error; got != "connection refused" {
              t.Errorf("results[1].error = %q, want %q", got, "connection refused")
      }
      if got := body.Results[0].Error; got != "" {
              t.Errorf("results[0].error = %q, want it omitted", got)
      }
}
