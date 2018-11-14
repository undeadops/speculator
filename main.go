package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const appName = "speculator"

var (
	// Info Logger
	Info *log.Logger
	// Warning Logger
	Warning *log.Logger
)

func setupLogging(
	infoHandle io.Writer,
	warningHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// Configuration Details
type settings struct {
	ASGroupName     string `required:"true"`
	Enabled         bool
	DisableCoolDown bool
	MemWatermark    int `default:"25"`
	CPUWatermark    int `default:"40"`
	MemThreshold    int `default:"60"`
	CPUThreshold    int `default:"56"`
}

// ClusterResources - Latest status of Cluster Resources
type ClusterResources struct {
	ClusterCPU     float64
	ClusterMemory  float64
	CPULimits      float64
	CPURequests    float64
	MemoryLimits   float64
	MemoryRequests float64
	ASGDesired     int
}

// Settings for application
var Settings settings

func loadconfig() {
	// Load Configuration from Args or Env Vars
	err := envconfig.Process("", &Settings)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(podCount)
	prometheus.MustRegister(promClusterCPU)
	prometheus.MustRegister(promClusterMemory)
	prometheus.MustRegister(promCPULimits)
	prometheus.MustRegister(promCPULimitsPerc)
	prometheus.MustRegister(promCPURequests)
	prometheus.MustRegister(promCPURequestsPerc)
	prometheus.MustRegister(promMemLimits)
	prometheus.MustRegister(promMemLimitsPerc)
	prometheus.MustRegister(promMemRequests)
	prometheus.MustRegister(promMemRequestsPerc)
	prometheus.MustRegister(promASGcount)
	prometheus.MustRegister(promASGscale)
}

func serveHTTP() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.Handle("/metrics", promhttp.Handler())
	Warning.Fatal(http.ListenAndServe(":5000", router))
}

func main() {
	setupLogging(os.Stdout, os.Stdout)

	// Configure App
	loadconfig()

	// Create Channel for Resources
	resources := make(chan *ClusterResources)

	// Start Watching for changes in Kubernetes
	go watchK8s(resources)
	go scaleCluster(resources)

	// Start Metrics/Status Http Server
	serveHTTP()
}
