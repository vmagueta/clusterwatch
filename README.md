# clusterwatch

A small service that probes the health of cluster nodes and reports their status.

> **Note:** this is a learning project. I started writing Go recently, coming
> from Python and Rust. The commit history is the honest record of that.

## What it does today

Probes a list of nodes over TCP, measures dial latency, and classifies each node as
healthy, degraded, or unreachable.

```
NODE        STATUS       LATENCY      ERROR
google      healthy      59.057616ms  <nil>
cloudflare  healthy      33.711840ms  <nil>
dead        unreachable  222.254µs    dial tcp 127.0.0.1:9999: connect: connection refused
```

## Running it

```bash
go run .
```

Requires Go 1.26 or later. No external dependencies: the standard library covers
networking, formatting, and tests.

## Design notes

**Probing is behind an interface.** `cluster.Checker` declares a single method,
`Check(Node) Result`. `TCPChecker` is one implementation; an HTTP or query-based
checker would be another, and tests inject a fake. Callers aggregate results without
knowing how a node was probed.

**Status is never serialized as a number.** `Status` is an `int` behind the scenes,
but it renders through `String()` and is emitted as text. Constants declared with
`iota` shift value whenever a line is inserted, which silently corrupts any data
that outlived the change.

**Degraded is distinct from unreachable.** A node that answers slowly is a different
operational problem from one that does not answer at all, and collapsing both into
"down" loses the signal that matters for capacity decisions.
