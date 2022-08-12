package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"gitlab.com/wobcom/transceiver-exporter/transceiver-collector"
	"net/http"
	"os"
	"strings"
)

const version string = "1.0"

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
	fmt.Println("Author(s): @fluepke, @BarbarossaTM")
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

type transceiverCollectorWrapper struct {
    collector *transceivercollector.TransceiverCollector
}

func (t transceiverCollectorWrapper) Collect (ch chan<- prometheus.Metric) {
    errs := make(chan error)
    done := make(chan struct{})
    go t.collector.Collect(ch, errs, done)
    for {
        select {
        case err := <-errs:
            log.Errorf("Error while collecting metrics: %v", err)
        case <- done:
            return
        }
    }
}

func (t transceiverCollectorWrapper) Describe(ch chan<- *prometheus.Desc) {
    t.collector.Describe(ch)
}

func handleMetricsRequest(w http.ResponseWriter, request *http.Request) {
	registry := prometheus.NewRegistry()

	excludedIfaceNames := strings.Split(*excludeInterfaces, ",")
	for index, excludedIfaceName := range excludedIfaceNames {
		excludedIfaceNames[index] = strings.Trim(excludedIfaceName, " ")
	}
	transceiverCollector := transceivercollector.NewCollector(excludedIfaceNames, *collectInterfaceFeatures)
    wrapper := &transceiverCollectorWrapper{
        collector: transceiverCollector,
    }

	registry.MustRegister(wrapper)
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError,
	}).ServeHTTP(w, request)
}
