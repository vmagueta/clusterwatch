package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/vmagueta/clusterwatch/cluster"
)

func main() {
	checker := cluster.TCPChecker{
		Timeout:       2 * time.Second,
		DegradedAfter: 200 * time.Millisecond,
	}

	nodes := []cluster.Node{
		{ID: "google", Addr: "google.com:80"},
		{ID: "cloudflare", Addr: "1.1.1.1:443"},
		{ID: "dead", Addr: "127.0.0.1:9999"},
	}

	report := cluster.CheckAll(checker, nodes)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "NODE\tSTATUS\tLATENCY\tERROR")
	for _, r := range report.Results {
		fmt.Fprintf(w, "%s\t%s\t%v\t%v\n", r.Node.ID, r.Status, r.Latency, r.Err)
	}

	fmt.Fprintf(w, "\n%d healthy\t%d degraded\t%d unreachable\t\n",
		report.Counts[cluster.StatusHealthy],
		report.Counts[cluster.StatusDegraded],
		report.Counts[cluster.StatusUnreachable],
	)
}
