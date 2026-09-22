package cluster

import "time"

// Result is the outcome of checking a single node.
type Result struct {
	Node    Node
	Status  Status
	Latency time.Duration
	Err     error
}

// Checker probes a single node and reports what it found.
//
// Implementations decide how a node is probed (TCP dial, HTTP request, ...);
// callers only care about the Result.
type Checker interface {
	Check(n Node) Result
}
