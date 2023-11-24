package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	transceivercollector "github.com/wobcom/transceiver-exporter/transceiver-collector"
)

const version string = "1.5.0"

var (
	showVersion              = flag.Bool("version", false, "Print version and exit")
	listenAddress            = flag.String("web.listen-address", "[::]:9458", "Address to listen on")
	metricsPath              = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	collectInterfaceFeatures = flag.Bool("collector.interface-features.enable", true, "Collect interface features")
	excludeInterfaces        = flag.String("exclude.interfaces", "", "Comma seperated list of interfaces to exclude")
	includeInterfaces        = flag.String("include.interfaces", "", "Comma seperated list of interfaces to include")
	excludeInterfacesRegex   = flag.String("exclude.interfaces-regex", "", "Regex of interfaces to exclude")
	includeInterfacesRegex   = flag.String("include.interfaces-regex", "", "Regex of interfaces to include")
	excludeInterfacesDown    = flag.Bool("exclude.interfaces-down", false, "Don't report on interfaces being management DOWN")
	powerUnitdBm             = flag.Bool("collector.optical-power-in-dbm", false, "Report optical powers in dBm instead of mW (default false -> mW)")

	excludeInterfacesRegexCompiled = &regexp.Regexp{}
	includeInterfacesRegexCompiled = &regexp.Regexp{}
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if err := compileRegexFlags(); err != nil {
		log.Fatalf(err.Error())
	}

	startServer()
}

func printVersion() {
	fmt.Println("transceiver-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): @fluepke, @BarbarossaTM, @vidister")
	fmt.Println("Metrics Exporter for pluggable transceivers on Linux based hosts / switches")
}

// compileRegexFlags compiles the cli regex flags into the global variables
// and returns an error if the regex is invalid
func compileRegexFlags() error {
	var err error
	excludeInterfacesRegexCompiled, err = regexp.Compile(*excludeInterfacesRegex)
	if err != nil {
		return fmt.Errorf("error compiling exclude.interfaces-regex: %v", err)
	}
	includeInterfacesRegexCompiled, err = regexp.Compile(*includeInterfacesRegex)
	if err != nil {
		return fmt.Errorf("error compiling include.interfaces-regex: %v", err)
	}
	return nil
}

func startServer() {
	log.Infof("Starting transceiver-exporter (version: %s)\n", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
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

func (t transceiverCollectorWrapper) Collect(ch chan<- prometheus.Metric) {
	errs := make(chan error)
	done := make(chan struct{})
	go t.collector.Collect(ch, errs, done)
	for {
		select {
		case err := <-errs:
			log.Errorf("Error while collecting metrics: %v", err)
		case <-done:
			return
		}
	}
}

func (t transceiverCollectorWrapper) Describe(ch chan<- *prometheus.Desc) {
	t.collector.Describe(ch)
}

func handleMetricsRequest(w http.ResponseWriter, request *http.Request) {
	var excludedIfaceNames []string
	var includedIfaceNames []string
	registry := prometheus.NewRegistry()

	if len(*excludeInterfaces) > 0 {
		excludedIfaceNames = strings.Split(*excludeInterfaces, ",")
		for index, excludedIfaceName := range excludedIfaceNames {
			excludedIfaceNames[index] = strings.Trim(excludedIfaceName, " ")
		}
	}
	if len(*includeInterfaces) > 0 {
		includedIfaceNames = strings.Split(*includeInterfaces, ",")
		for index, includedIfaceName := range includedIfaceNames {
			includedIfaceNames[index] = strings.Trim(includedIfaceName, " ")
		}
	}

	transceiverCollector := transceivercollector.NewCollector(excludedIfaceNames, includedIfaceNames, excludeInterfacesRegexCompiled, includeInterfacesRegexCompiled, *excludeInterfacesDown, *collectInterfaceFeatures, *powerUnitdBm)
	wrapper := &transceiverCollectorWrapper{
		collector: transceiverCollector,
	}

	registry.MustRegister(wrapper)
	l := log.New()
	l.Level = log.ErrorLevel

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      l,
		ErrorHandling: promhttp.ContinueOnError,
	}).ServeHTTP(w, request)
}
