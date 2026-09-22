package cluster

import "time"

// Report is the aggregated outcome of checking a fleet of nodes.
type Report struct {
	CheckedAt time.Time
	Results   []Result
	Counts    map[Status]int
}

// CheckAll probes every node with the given checker and aggregates the results.
//
// Nodes are probed sequentially, so the total time is the sum of all checks.
func CheckAll(c Checker, nodes []Node) Report {
	results := make([]Result, 0, len(nodes))
	counts := make(map[Status]int)

	for _, n := range nodes {
		r := c.Check(n)
		results = append(results, r)
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
