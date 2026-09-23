package cluster

import (
	"encoding/json"
	"testing"
)

func TestStatusString(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   string
	}{
		{"healthy", StatusHealthy, "healthy"},
		{"degraded", StatusDegraded, "degraded"},
		{"unreachable", StatusUnreachable, "unreachable"},
		{"unknown", StatusUnknown, "unknown"},
		{"value outside the known set falls back to unknown", Status(42), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("Status(%d).String() = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusMarshalsAsText(t *testing.T) {
	got, err := json.Marshal(map[Status]Status{StatusHealthy: StatusDegraded})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	if want := `{"healthy":"degraded"}`; string(got) != want {
		t.Errorf("json.Marshal = %s, want %s", got, want)
	}
}
