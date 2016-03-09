// Command rtorrent_exporter provides a Prometheus exporter for rTorrent.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/mdlayher/rtorrent"
	"github.com/mdlayher/rtorrent_exporter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9135", "host:port for rTorrent exporter")
	metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")

	rtorrentAddr = flag.String("rtorrent.addr", "", "address of rTorrent XML-RPC server")
)

func main() {
	flag.Parse()

	if *rtorrentAddr == "" {
		log.Fatal("address of rTorrent XML-RPC server must be specified with '-rtorrent.addr' flag")
	}

	c, err := rtorrent.New(*rtorrentAddr, nil)
	if err != nil {
		log.Fatalf("cannot create rTorrent client: %v", err)
	}

	prometheus.MustRegister(rtorrentexporter.New(c))

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting rTorrent exporter on %q for server %q", *telemetryAddr, *rtorrentAddr)

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("cannot start rTorrent exporter: %s", err)
	}
}
