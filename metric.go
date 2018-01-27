package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	podCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "pod_count",
		Help:      "Number of running Pods",
	},
	)

	promCPURequests = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_cpu_requests",
		Help:      "Current CPU Requests defined by Deployments/Pods",
	})

	promCPURequestsPerc = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_cpu_requests_percent",
		Help:      "Current CPU Requests as Percent of Cluster CPU",
	})

	promCPULimits = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluter_cpu_limits",
		Help:      "Current CPU Limits defined by Deployments/Pods",
	})

	promCPULimitsPerc = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_cpu_limits_percent",
		Help:      "Current CPU Limits as Percent of Cluster CPU, Can be over 100%",
	})

	promMemRequests = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_mem_requests",
		Help:      "Current Memory Requests defined by Deployments/Pods",
	})

	promMemRequestsPerc = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_mem_requests_percent",
		Help:      "Current Memory Requests as Percent of Cluster Memory",
	})

	promMemLimits = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_mem_limits",
		Help:      "Current Memory Limits defined by Deployments/Pods",
	})

	promMemLimitsPerc = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_mem_limits_percent",
		Help:      "Current Memory Limits as Percent of Cluster Memory, Can be over 100%? maybe...",
	})

	promClusterCPU = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_cpu",
		Help:      "Current CPU available on the Cluster",
	})

	promClusterMemory = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: appName,
		Name:      "cluster_memory",
		Help:      "Current Memory available on the Cluster",
	})

	promASGcount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: appName,
			Name:      "asg_count",
			Help:      "AutoScaler Group Current Desired Count",
		},
		[]string{"name"},
	)

	promASGscale = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: appName,
			Name:      "asg_scale_event",
			Help:      "AutoScaler Group Scale Event",
		},
		[]string{"name"},
	)
)
