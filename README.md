# clusterwatch

A small service that probes the health of cluster nodes and reports their status.

> **Note:** this is a learning project. I started writing Go recently, coming
> from Python and Rust. The commit history is the honest record of that.

## What it does today

Serves a fleet health report over HTTP. On every `GET /report` it probes a list of
nodes over TCP, concurrently, measures dial latency, and classifies each node as
healthy, degraded, or unreachable.

```json
{
  "checked_at": "2026-09-23T16:40:06.4853538+01:00",
  "counts": { "degraded": 2, "unreachable": 1 },
  "results": [
    { "node": "google", "status": "degraded", "latency_ms": 426 },
    { "node": "cloudflare", "status": "degraded", "latency_ms": 324 },
    {
      "node": "dead",
      "status": "unreachable",
      "latency_ms": 0,
      "error": "dial tcp 127.0.0.1:9999: connect: connection refused"
    }
  ]
}
```

## Running it

```bash
go run .
```

The server listens on `:8080`. From another terminal:

```bash
curl -s localhost:8080/report
```

Only `GET` is routed; any other method on `/report` gets `405 Method Not Allowed`.

Run the tests with the race detector:

```bash
go test -race ./...
```

Requires Go 1.26 or later. No external dependencies: the standard library covers
networking, HTTP, JSON, logging, and tests.

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

**Nodes are probed concurrently.** `CheckAll` starts one goroutine per node and
each goroutine writes only to its own slot of a preallocated slice, so there is no
shared mutable state to lock. Counting happens after every probe has finished. The
total time is that of the slowest node, not the sum of all of them. As a
consequence, `Checker` implementations must be safe for concurrent use.

**The JSON shape is its own type.** The `api` package converts `cluster.Report`
into response types before encoding. The domain types can change freely without
breaking clients, and fields that do not serialize well (`error`, `time.Duration`)
get an explicit wire representation.
