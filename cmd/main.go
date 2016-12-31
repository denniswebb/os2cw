package cmd

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"time"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

func main() {
	flag.Parse()

	sess := session.New()
	cloudwatchService := cloudwatch.New(sess)

	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")
	h, _ := host.Info()

	var memUnitDivisor, diskUnitDivisor, uptimeUnitDivisor float64
	var memUnitText, diskUnitText, uptimeUnitText string
	memUnitDivisor = 1
	diskUnitDivisor = 1
	uptimeUnitDivisor = 1

	if config.InstanceId=="" {
		//try to read instanceid from metadata
		metadataService := ec2metadata.New(sess)
		instanceid, err := metadataService.GetMetadata("instance-id")
		if err != nil {
			fmt.Println(err)
		} else {
			config.InstanceId = instanceid
		}
	}

	switch {
	case config.MemoryUnit == "kb":
		memUnitText = "Kilobytes"
		memUnitDivisor = 1024
	case config.MemoryUnit == "mb":
		memUnitText = "Megabytes"
		memUnitDivisor = 1024 * 1024
	case config.MemoryUnit == "gb":
		memUnitText = "Gigabytes"
		memUnitDivisor = 1024 * 1024 * 1024
	default:
		config.MemoryUnit = "b"
		memUnitText = "Bytes"
		memUnitDivisor = 1
	}

	if config.MemoryFreePercentage {
		memFreePercentageValue := 100.0 - v.UsedPercent

		if config.Verbose {
			fmt.Printf("Memory Free Percent: %f%%\n", memFreePercentageValue)
		}

		_, err := SendMetrics(cloudwatchService, "MemoryFreePercentage", "Percent", memFreePercentageValue)

		if err != nil {
			fmt.Println(err)
		}
	}

	if config.MemoryUsedPercentage {
		if config.Verbose {
			fmt.Printf("Memory Used Percent: %f%%\n", v.UsedPercent)
		}

		_, err := SendMetrics(cloudwatchService, "MemoryUtilization", "Percent", v.UsedPercent)

		if err != nil {
			fmt.Println(err)
		}
	}

	if config.MemoryFreeAmount {
		memFreeValue := float64(v.Available) / memUnitDivisor

		if config.Verbose {
			fmt.Printf("Memory Free(%s): %f\n", memUnitText, memFreeValue)
		}

		_, err := SendMetrics(cloudwatchService, "Memory Free", memUnitText, memFreeValue)

		if err != nil {
			fmt.Println(err)
		}
	}

	if config.MemoryUsedAmount {
		memUsedValue := float64(v.Used) / memUnitDivisor
		if config.Verbose {
			fmt.Printf("Memory Used(%s): %f\n", memUnitText, memUsedValue)
		}

		_, err := SendMetrics(cloudwatchService, "Memory Used", memUnitText, memUsedValue)

		if err != nil {
			fmt.Println(err)
		}
	}

	switch {
	case config.DiskUnit == "kb":
		diskUnitText = "Kilobytes"
		diskUnitDivisor = 1024
	case config.DiskUnit == "mb":
		diskUnitText = "Megabytes"
		diskUnitDivisor = 1024 * 1024
	case config.DiskUnit == "gb":
		diskUnitText = "Gigabytes"
		diskUnitDivisor = 1024 * 1024 * 1024
	case config.DiskUnit == "tb":
		diskUnitText = "Terabytes"
		diskUnitDivisor = 1024 * 1024 * 1024 * 1024
	default:
		config.DiskUnit = "b"
		diskUnitText = "Bytes"
		diskUnitDivisor = 1024 * 1024
	}

	if config.DiskFreePercentage {
		diskFreePercentageValue := 100.0 - d.UsedPercent

		if config.Verbose {
			fmt.Printf("Disk Free Percent: %f%%\n", diskFreePercentageValue)
		}

		_, err := SendMetrics(cloudwatchService, "DiskFreePercentage", "Percent", diskFreePercentageValue)

		if err != nil {
			fmt.Println(err)
		}
	}

	if config.DiskUsedPercentage {
		if config.Verbose {
			fmt.Printf("Disk Used Percent: %f%%\n", d.UsedPercent)
		}

		_, err := SendMetrics(cloudwatchService, "DiskUtilization", "Percent", d.UsedPercent)

		if err != nil {
			fmt.Println(err)
		}
	}

	if config.DiskFreeAmount {
		diskFreeValue := float64(d.Free) / diskUnitDivisor

		if config.Verbose {
			fmt.Printf("Disk Free(%s): %f\n", diskUnitText, diskFreeValue)
		}

		_, err := SendMetrics(cloudwatchService, "Disk Free", diskUnitText, diskFreeValue)

		if err != nil {
			fmt.Println(err)
		}
	}

	if config.DiskUsedAmount {
		diskUsedValue := float64(d.Used) / diskUnitDivisor

		if config.Verbose {
			fmt.Printf("Disk Used(%s): %f\n", diskUnitText, diskUsedValue)
		}

		_, err := SendMetrics(cloudwatchService, "Disk Used", diskUnitText, diskUsedValue)

		if err != nil {
			fmt.Println(err)
		}
	}

	switch {
	case config.UptimeUnit == "s":
		uptimeUnitDivisor = 1
		uptimeUnitText = "seconds"
	case config.UptimeUnit == "m":
		uptimeUnitDivisor = 60
		uptimeUnitText = "minutes"
	case config.UptimeUnit == "h":
		uptimeUnitDivisor = 60 * 60
		uptimeUnitText = "hours"
	case config.UptimeUnit == "d":
		uptimeUnitDivisor = 24 * 60 * 60
		uptimeUnitText = "days"
	default:
		config.UptimeUnit = "s"
		uptimeUnitDivisor = 1
		uptimeUnitText = "seconds"
	}

	if config.Uptime {
		if config.Verbose {
			fmt.Printf("System Uptime(%s): %v\n", uptimeUnitText, float64(h.Uptime) / uptimeUnitDivisor)
		}

		_, err := SendMetrics(cloudwatchService, "Uptime", "Seconds", float64(h.Uptime))

		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendMetrics(cloudwatchService *cloudwatch.CloudWatch, metricName string, metricUnit string, metricValue float64) (*cloudwatch.PutMetricDataOutput, error) {
	params := &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("InstanceID"),
						Value: aws.String(config.InstanceId),
					},
				},
				Timestamp: aws.Time(time.Now()),
				Unit:      aws.String(metricUnit),
				Value:     aws.Float64(metricValue),
			},
		},
		Namespace: aws.String(config.NameSpace),
	}
	return cloudwatchService.PutMetricData(params)
}