// Package cluster models the nodes of a distributed cluster and their health.
package cluster

// Status describes the health of a node as observed by a checker.
type Status int

// Node health states. StatusUnknown is the zero value: a node that has not
// been observed yet is never assumet to be healthy.
const (
	StatusUnknown Status = iota
	StatusHealthy
	StatusDegraded
	StatusUnreachable
)

// String returns the human-readable name of the status.
func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusDegraded:
		return "degraded"
	case StatusUnreachable:
		return "unreachable"
	default:
		return "unknown"
	}
}

// Node is a single machine in the cluster, addressable over the network.
type Node struct {
	ID   string
	Addr string
}
