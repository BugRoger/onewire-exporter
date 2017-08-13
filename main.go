package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "onewire"
)

type Metrics struct {
	Temperatures []TempSensor
}

type TempSensor struct {
	Address     string
	Temperature float64
}

type Exporter struct {
	client       *OneWireClient
	up           *prometheus.Desc
	info         *prometheus.Desc
	temperatures *prometheus.Desc
}

type Client interface {
	Collect() (*Metrics, error)
}

func NewExporter(owserver string) *Exporter {
	client, err := NewOneWireClient(owserver)
	if err != nil {
		log.Fatalf("Couldn't initialize owfsclient to %s: %s", owserver, err)
	}

	return &Exporter{
		client: client,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Could the 1-Wire network be reached.",
			nil,
			nil,
		),
		temperatures: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "temperatures"),
			"Temperatures",
			[]string{"address"},
			nil,
		),
	}
}

func main() {
	var (
		listenAddress   = flag.String("web.listen-address", ":9401", "Address to listen on for web interface and telemetry.")
		metricsPath     = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		owserverAddress = flag.String("owserver", "localhost:4304", "Connection to the owserver")
	)
	flag.Parse()

	prometheus.MustRegister(NewExporter(*owserverAddress))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>1-Wire Exporter</title></head>
             <body>
             <h1>1-Wire Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	fmt.Println("Starting HTTP server on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.temperatures
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	metrics, err := e.client.Collect()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		log.Printf("Failed to collect stats from miner: %s\n", err)
		return
	} else {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
	}

	for _, temp := range metrics.Temperatures {
		ch <- prometheus.MustNewConstMetric(e.temperatures, prometheus.GaugeValue, temp.Temperature, temp.Address)
	}
}
