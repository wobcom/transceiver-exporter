package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"os"
)

const version string = "0.1"

const prefix = "transceiver_exporter_"

var (
	showVersion              = flag.Bool("version", false, "Print version and exit")
	listenAddress            = flag.String("web.listen-address", "[::]:9458", "Address to listen on")
	metricsPath              = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	collectInterfaceFeatures = flag.Bool("collector.interface-features.enable", true, "Collect interface features")
	excludeInterfaces        = flag.String("exclude.interfaces", "", "Comma seperated list of interfaces to exclude")
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("transceiver-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): @fluepke")
	fmt.Println("Metrics Exporter for pluggable transceivers on Linux based hosts / switches")
}

func startServer() {
	log.Infof("Starting transceiver-exporter (version: %s)\n", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>transceiver-exporter (Version ` + version + `)</title></head>
            <body>
            <h1>transceiver-exporter</h1>
            </body>
            </html>`))
	})
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening on %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, request *http.Request) {
	registry := prometheus.NewRegistry()
	transceiverCollector := NewTransceiverCollector()

	registry.MustRegister(transceiverCollector)
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError,
	}).ServeHTTP(w, request)
}
