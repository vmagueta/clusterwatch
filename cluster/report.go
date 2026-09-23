package cluster

import (
	"sync"
	"time"
)

// Report is the aggregated outcome of checking a fleet of nodes.
type Report struct {
	CheckedAt time.Time
	Results   []Result
	Counts    map[Status]int
}

// CheckAll probes every node with the given checker and aggregates the results.
//
// Nodes are probed concurrently, one goroutine per node, so the total time is
// roughly that of the slowest check. Results keep the order of nodes.
func CheckAll(c Checker, nodes []Node) Report {
	results := make([]Result, len(nodes))

	var wg sync.WaitGroup
	for i, n := range nodes {
		wg.Go(func() {
			results[i] = c.Check(n)
		})
	}
	wg.Wait()

	counts := make(map[Status]int)
	for _, r := range results {
		counts[r.Status]++
	}

	return Report{
		CheckedAt: time.Now(),
		Results:   results,
		Counts:    counts,
	}
}

// Unhealthy returns the results when the status are not “StatusHealthy“
func (rpt Report) Unhealthy() []Result {
	var results []Result
	for _, r := range rpt.Results {
		if r.Status != StatusHealthy {
			results = append(results, r)
		}
	}
	return results
}
