package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/mmohamed/wmbusmeters-prometheus-metric/pkg/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsPath = "/metrics"
	pushPath = "/push"
	healthzPath = "/healthz"
)

var (
	inMemoriesData = map[string]string {}
)

func metricsServer(registry prometheus.Gatherer, port int) {
	// Address to listen on for web interface and telemetry
	listenAddress := fmt.Sprintf(":%d", port)

	glog.Infof("Starting metrics server: %s", listenAddress)
	// Add metricsPath
	http.Handle(metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	// Add pushPath
	http.HandleFunc(pushPath, func(w http.ResponseWriter, r *http.Request) {
		var parsedData map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&parsedData)
		if err != nil {
			panic(err)
		}
		body, err := json.Marshal(parsedData)
		if err != nil {
			panic(err)
		}
		inMemoriesData[parsedData["id"].(string)+"-"+parsedData["name"].(string)+"-"+parsedData["meter"].(string)] = string(body[:])
		w.WriteHeader(200)
		w.Write(body)
	})
	// Add healthzPath
	http.HandleFunc(healthzPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	// Add index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
	<head>
		<title>Wmbusmeters Metrics Server</title>
	</head>
	<body>
		<h1>Kube Metrics</h1>
		<ul>
			<li><a href='` + metricsPath + `'>metrics</a></li>
			<li><a href='` + pushPath + `'>push</a></li>
			<li><a href='` + healthzPath + `'>healthz</a></li>
		</ul>
	</body>
</html>`))
	})
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

var (
	optHelp           bool
	optPort           int
	optKubeconfigPath string
)

func init() {
	flag.BoolVar(&optHelp, "help", false, "print help info and exit")
	flag.IntVar(&optPort, "port", 9001, "port to expose metrics on")
}

func main() {
	// We log to stderr because glog will default to logging to a file.
	flag.Set("logtostderr", "true")
	flag.Parse()

	if optHelp {
		flag.Usage()
		return
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewWMBusmetersCollector(inMemoriesData))

	metricsServer(registry, optPort)
}
