package main

import (
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func watchK8s(resources chan *ClusterResources) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		Warning.Panic(err.Error())
	}
	// creates the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		Warning.Panic(err.Error())
	}

	// Maybe not good to put for here...
	for {
		pods, err := client.CoreV1().Pods(apiv1.NamespaceAll).List(
			metav1.ListOptions{})
		if err != nil {
			Warning.Panic(err)
		}
		podCount.Set(float64(len(pods.Items)))

		nodes, err := client.CoreV1().Nodes().List(
			metav1.ListOptions{})
		if err != nil {
			Warning.Panic(err)
		}

		clusterCPU := calculateNodeResources(nodes, apiv1.ResourceCPU)
		clusterMem := calculateNodeResources(nodes, apiv1.ResourceMemory)

		promClusterCPU.Set(float64(clusterCPU))
		promClusterMemory.Set(float64(clusterMem))

		cpuRequests := calculateReqResources(pods, apiv1.ResourceCPU)
		cpuLimits := calculateLimitResources(pods, apiv1.ResourceCPU)

		memoryRequests := calculateReqResources(pods, apiv1.ResourceMemory)
		memoryLimits := calculateLimitResources(pods, apiv1.ResourceMemory)

		cores := resource.NewMilliQuantity(cpuRequests, resource.DecimalSI)
		memory := resource.NewQuantity(memoryRequests, resource.BinarySI)

		r := &ClusterResources{
			ClusterCPU:     float64(clusterCPU),
			ClusterMemory:  float64(clusterMem),
			CPULimits:      percentOf(cpuLimits, clusterCPU),
			CPURequests:    percentOf(cpuRequests, clusterCPU),
			MemoryLimits:   percentOf(memoryLimits, clusterMem),
			MemoryRequests: percentOf(memoryRequests, clusterMem),
		}
		resources <- r

		Info.Printf("Requested cores = %v (%v%%)\n", cores, percentOf(cpuRequests, clusterCPU))
		Info.Printf("Requested memory = %v(%v%%)\n", memory, percentOf(memoryRequests, clusterMem))

		promCPURequests.Set(float64(cpuRequests))
		promCPULimits.Set(float64(cpuLimits))

		promCPURequestsPerc.Set(percentOf(cpuRequests, clusterCPU))
		promCPULimitsPerc.Set(percentOf(cpuLimits, clusterCPU))

		promMemRequests.Set(float64(memoryRequests))
		promMemLimits.Set(float64(memoryLimits))

		promMemRequestsPerc.Set(percentOf(memoryRequests, clusterMem))
		promMemLimitsPerc.Set(percentOf(memoryLimits, clusterMem))

		time.Sleep(120 * time.Second)
	}
}

func calculateReqResources(pods *apiv1.PodList, resourceName apiv1.ResourceName) int64 {
	var resources int64
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if resourceValue, found := container.Resources.Requests[resourceName]; found {
				resources = resources + resourceValue.MilliValue()
			}
		}
	}
	return resources
}

func calculateLimitResources(pods *apiv1.PodList, resourceName apiv1.ResourceName) int64 {
	var resources int64
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if resourceValue, found := container.Resources.Limits[resourceName]; found {
				resources = resources + resourceValue.MilliValue()
			}
		}
	}
	return resources
}

func calculateNodeResources(nodes *apiv1.NodeList, resourceName apiv1.ResourceName) int64 {
	var resources int64
	for _, node := range nodes.Items {
		if !node.Spec.Unschedulable {
			//fmt.Printf("Node=%s, Resource=%s, Capacity=%v\n", node.ObjectMeta.Name, resourceName, node.Status.Capacity[resourceName])
			cap := node.Status.Capacity[resourceName]
			resources = resources + cap.MilliValue()
		}
	}
	return resources
}

func percentOf(current int64, all int64) float64 {
	percent := (float64(current) * float64(100)) / float64(all)
	return round(percent, 0.05)
}

func round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}
