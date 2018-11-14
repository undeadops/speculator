package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	scaleUpBy   = 4
	scaleDownBy = 1
	asgDesired  int
	asgMin      int
	asgMax      int
	scaleTo     int
)

func scaleCluster(resources chan *ClusterResources) {

	for {
		asgDesired, asgMin, asgMax = getAutoScaleDesired()
		// Report current count to prometheus exporter
		promASGcount.With(prometheus.Labels{"name": Settings.ASGroupName}).Set(float64(asgDesired))

		// Set scaleTo to equal Desired
		scaleTo = asgDesired
		r := <-resources

		// Not scaling down for low memory useage
		Info.Printf("ASG: %s, Current: %d, Min: %d, Max: %d\n", Settings.ASGroupName, asgDesired, asgMin, asgMax)
		// Check CPU Resources
		if r.CPURequests >= float64(Settings.CPUThreshold) {
			scaleUp()
		} else if r.CPURequests <= float64(Settings.CPUWatermark) {
			scaleDown()
		} else if r.MemoryRequests > float64(Settings.MemThreshold) {
			scaleUp()
		}
	}
}

func scaleUp() {
	var d int
	d = asgDesired + scaleUpBy
	if d >= asgMax && asgDesired < asgMax {
		scaleASG(asgMax)
	} else if d < asgMax {
		scaleASG(d)
	} else {
		Info.Printf("Workers Already Scaled up to Max\n")
	}
}

func scaleDown() {
	var d int
	d = asgDesired - scaleDownBy
	if d < asgMin && asgDesired > asgMin {
		scaleASG(asgMin)
	} else if d >= asgMin {
		scaleASG(d)
	} else {
		Info.Printf("Workers Already Scaled to Min\n")
	}
}

func getAutoScaleDesired() (int, int, int) {
	var desired int
	var min int
	var max int

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		Warning.Panic("failed to load config, " + err.Error())
	}

	svc := autoscaling.New(cfg)
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{
			Settings.ASGroupName,
		},
	}

	req := svc.DescribeAutoScalingGroupsRequest(input)
	result, err := req.Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case autoscaling.ErrCodeInvalidNextToken:
				Warning.Println(autoscaling.ErrCodeInvalidNextToken, aerr.Error())
			case autoscaling.ErrCodeResourceContentionFault:
				Warning.Println(autoscaling.ErrCodeResourceContentionFault, aerr.Error())
			default:
				Warning.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			Warning.Println(err.Error())
		}
		return 0, 0, 0
	}

	for _, g := range result.AutoScalingGroups {
		desired = int(*g.DesiredCapacity)
		min = int(*g.MinSize)
		max = int(*g.MaxSize)
	}
	return desired, min, max
}

func scaleASG(desired int) {
	Info.Printf("I will scale ASG to: %d\n", desired)

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		Warning.Panic("failed to load config, " + err.Error())
	}

	svc := autoscaling.New(cfg)

	var coolDown = true
	if Settings.DisableCoolDown {
		coolDown = false
	}

	desiredCapacity := int64(desired)

	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: &Settings.ASGroupName,
		DesiredCapacity:      &desiredCapacity,
		HonorCooldown:        &coolDown,
	}

	if Settings.Enabled {
		promASGscale.With(prometheus.Labels{"name": Settings.ASGroupName}).Inc()
		req := svc.SetDesiredCapacityRequest(input)
		resp, err := req.Send()
		if err != nil {
			Info.Println("Scaling Failed, Cooldown window may be active")
			Info.Println(resp)
		}

	} else {
		Info.Println("Scaling is disabled: envvar Enabled: True to enable")
	}
}
