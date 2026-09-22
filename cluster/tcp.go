package cluster

import (
	"net"
	"time"
)

// TCPChecker probes a node by opening a TCP connection to its address.
type TCPChecker struct {
	// Timeout bounds how long a single dial may take.
	Timeout time.Duration
	// DegradedAfter is the latency above which a reachable node is
	// considerer degraded rather than healthy.
	DegradedAfter time.Duration
}

// Check dials the node's address and classifies the ouctome by latency.
func (c TCPChecker) Check(n Node) Result {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", n.Addr, c.Timeout)
	latency := time.Since(start)

	if err != nil {
		return Result{Node: n, Status: StatusUnreachable, Latency: latency, Err: err}
	}
	defer conn.Close()

	status := StatusHealthy
	if latency > c.DegradedAfter {
		status = StatusDegraded
	}

	return Result{Node: n, Status: status, Latency: latency}
}
