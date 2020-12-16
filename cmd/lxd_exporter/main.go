package main

import (

	"log"
	"net/http"
	"os"

	lxd "github.com/lxc/lxd/client"
	"github.com/nieltg/lxd_exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "staging-UNVERSIONED"

	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default("127.0.0.1:8000").String()

)

func main() {
	logger := log.New(os.Stderr, "lxd_exporter: ", log.LstdFlags)

	kingpin.Version(version)
	kingpin.Parse()

	server, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		logger.Fatalf("Unable to contact LXD server: %s", err)
		return
	}

	prometheus.MustRegister(metrics.NewCollector(logger, server))
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
