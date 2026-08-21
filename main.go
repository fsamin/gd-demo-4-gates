// gd-demo-4-gates is step 4 of the git-deploy demo: the guard rails a commit
// must clear before it is released. The page prints an order total computed by
// price.go, so a regression in the business rule is visible on the page — and
// the point of the demo is that a commit carrying that regression never reaches
// this page at all.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "v1"

// basket is the sample order shown on the page. It is worth exactly the
// discount threshold, so it sits on the boundary the tests pin — the one a
// regression is most likely to move.
var basket = []Item{
	{Label: "widget", Cents: 4_000, Qty: 2},
	{Label: "gadget", Cents: 2_000, Qty: 1},
}

func main() {
	log.Printf("gd-demo-4-gates %s starting on :8080", version)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		hostname, _ := os.Hostname()
		printf(w, "gd-demo-4-gates %s (pod %s)\n", version, hostname)

		printf(w, "\n== Order ==\n")
		for _, item := range basket {
			printf(w, "%-8s %2d × %9s = %9s\n", item.Label, item.Qty,
				euros(item.Cents), euros(item.Cents*item.Qty))
		}
		subtotal := Subtotal(basket)
		total := Total(basket)
		printf(w, "\nsubtotal %22s\n", euros(subtotal))
		printf(w, "discount %22s   (%d%% from %s on)\n", euros(total-subtotal),
			discountPercent, euros(discountThresholdCents))
		printf(w, "TOTAL    %22s\n", euros(total))

		printf(w, "\n== Why this total can be trusted ==\n")
		printf(w, "Every commit clears four gates before it is released:\n")
		printf(w, "  test       go test ./...        the discount rule and its boundary\n")
		printf(w, "  quality    golangci-lint run    dead code, unchecked errors\n")
		printf(w, "  vulncheck  govulncheck ./...    reachable vulnerabilities in the deps\n")
		printf(w, "  imageScan  trivy image          HIGH and CRITICAL CVEs in the image\n")
		printf(w, "\nA commit failing a blocking gate is never released: this page keeps\n")
		printf(w, "serving the previous one. Each branch of this repository breaks exactly\n")
		printf(w, "one gate — see the README.\n")

		log.Printf("%s %s from %s (%s)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start).Round(time.Millisecond))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// printf writes one line of the page, discarding the write error on purpose: a
// client hanging up mid-response is this server's normal weather, and there is
// nothing left to send it. Discarding it explicitly is what keeps the quality
// gate green without teaching the linter to look away.
func printf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

// euros renders a cent amount the way an invoice would.
func euros(cents int) string {
	sign := ""
	if cents < 0 {
		sign, cents = "-", -cents
	}
	return fmt.Sprintf("%s%d.%02d €", sign, cents/100, cents%100)
}
