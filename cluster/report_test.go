package cluster

import (
	"errors"
	"testing"
	"time"
)

// fakeChecker returns a canned result per node ID, so tests can describe any
// fleet state without touching the network.
type fakeChecker struct {
	results map[string]Result
}

func (f fakeChecker) Check(n Node) Result {
	if r, ok := f.results[n.ID]; ok {
		r.Node = n
		return r
	}
	return Result{Node: n, Status: StatusUnknown}
}

func TestCheckAllCountsEveryStatus(t *testing.T) {
	nodes := []Node{
		{ID: "a", Addr: "10.0.0.1:9000"},
		{ID: "b", Addr: "10.0.0.2:9000"},
		{ID: "c", Addr: "10.0.0.3:9000"},
	}

	checker := fakeChecker{results: map[string]Result{
		"a": {Status: StatusHealthy, Latency: 5 * time.Millisecond},
		"b": {Status: StatusDegraded, Latency: 900 * time.Millisecond},
		"c": {Status: StatusUnreachable, Err: errors.New("connection refused")},
	}}

	report := CheckAll(checker, nodes)

	if got, want := len(report.Results), len(nodes); got != want {
		t.Fatalf("len(Results) = %d, want %d", got, want)
	}

	wantCounts := map[Status]int{
		StatusHealthy:     1,
		StatusDegraded:    1,
		StatusUnreachable: 1,
	}
	for status, want := range wantCounts {
		if got := report.Counts[status]; got != want {
			t.Errorf("Counts[%s] = %d, want %d", status, got, want)
		}
	}
}

func TestUnhealthyExcludesHealthyResults(t *testing.T) {
	nodes := []Node{
		{ID: "a", Addr: "10.0.0.1:9000"},
		{ID: "b", Addr: "10.0.0.2:9000"},
		{ID: "c", Addr: "10.0.0.3:9000"},
	}

	checker := fakeChecker{results: map[string]Result{
		"a": {Status: StatusHealthy, Latency: 5 * time.Millisecond},
		"b": {Status: StatusDegraded, Latency: 900 * time.Millisecond},
		"c": {Status: StatusUnreachable, Err: errors.New("connection refused")},
	}}

	report := CheckAll(checker, nodes)
	unhealthy := report.Unhealthy()

	if got, want := len(unhealthy), 2; got != want {
		t.Fatalf("len(Unhealthy()) = %d, want %d", got, want)
	}

	for _, r := range unhealthy {
		if r.Status == StatusHealthy {
			t.Errorf("Unhealthy() returned node %q with status %s", r.Node.ID, r.Status)
		}
	}
}
