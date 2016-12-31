package cmd

import (
	"flag"
)

const (
	memoryFreePercentageDefault = false
	memoryFreePercentageDesc    = "Push memory Free Percentage Custom Metric"
	memoryFreeAmountDefault     = false
	memoryFreeAmountDesc        = "Push memory Free Amount Custom Metric"
	memoryUsedPercentageDefault = false
	memoryUsedPercentageDesc    = "Push memory Used Percentage Custom Metric"
	memoryUsedAmountDefault     = false
	memoryUsedAmountDesc        = "Push memory Used Amount Custom Metric"
	memoryUnitDefault           = "kb"
	memoryUnitDesc              = "Unit for memory metrics. Allowed: gb, mb, kb, b. Default: kb"

	diskFreePercentageDefault 	= false
	diskFreePercentageDesc    	= "Push disk Free Percentage Custom Metric"
	diskFreeAmountDefault     	= false
	diskFreeAmountDesc        	= "Push disk Free Amount Custom Metric"
	diskUsedPercentageDefault 	= false
	diskUsedPercentageDesc    	= "Push disk Used Percentage Custom Metric"
	diskUsedAmountDefault     	= false
	diskUsedAmountDesc        	= "Push disk Used Amount Custom Metric"
	diskUnitDefault           	= "kb"
	diskUnitDesc              	= "Unit for disk metrics. Allowed: tb, gb, mb, kb, b. Default: mb"

	uptimeDefault     			= false
	uptimeDesc        			= "System uptime"
	uptimeUnitDefault 			= "s"
	uptimeUnitDesc    			= "Unit for uptime metric. Allowed: s, m, h, d. Default: s"

	instanceIdDesc 				= "InstanceId used as metric dimension. Auto-detected by default."
	nameSpaceDefault			= "Instance"
	nameSpaceDesc				= "Cloudwatch namespace for metrics. Defaults to Custom"

	verboseDefault				= false
	verboseDesc					= "Outputs values to stdout."
)

var config *Config

type Config struct {
	MemoryFreePercentage 	bool
	MemoryFreeAmount     	bool
	MemoryUsedPercentage 	bool
	MemoryUsedAmount     	bool
	MemoryUnit           	string

	DiskFreePercentage		bool
	DiskFreeAmount    		bool
	DiskUsedPercentage		bool
	DiskUsedAmount    		bool
	DiskUnit          		string

	Uptime     				bool
	UptimeUnit 				string

	InstanceId				string
	NameSpace 				string

	Verbose					bool
}

func init() {
	config = &Config{}

	flag.BoolVar(&config.MemoryFreePercentage, "mem-free-percent", memoryFreePercentageDefault, memoryFreePercentageDesc)
	flag.BoolVar(&config.MemoryFreeAmount, "mem-free", memoryFreeAmountDefault, memoryFreeAmountDesc)
	flag.BoolVar(&config.MemoryUsedPercentage, "mem-used-percent", memoryUsedPercentageDefault, memoryUsedPercentageDesc)
	flag.BoolVar(&config.MemoryUsedAmount, "mem-used", memoryUsedAmountDefault, memoryUsedAmountDesc)
	flag.StringVar(&config.MemoryUnit, "mem-unit", memoryUnitDefault, memoryUnitDesc)

	flag.BoolVar(&config.DiskFreePercentage, "disk-free-percent", diskFreePercentageDefault, diskFreePercentageDesc)
	flag.BoolVar(&config.DiskFreeAmount, "disk-free", diskFreeAmountDefault, diskFreeAmountDesc)
	flag.BoolVar(&config.DiskUsedPercentage, "disk-used-percent", diskUsedPercentageDefault, diskUsedPercentageDesc)
	flag.BoolVar(&config.DiskUsedAmount, "disk-used", diskUsedAmountDefault, diskUsedAmountDesc)
	flag.StringVar(&config.DiskUnit, "disk-unit", diskUnitDefault, diskUnitDesc)

	flag.BoolVar(&config.Uptime, "uptime", uptimeDefault, uptimeDesc)
	flag.StringVar(&config.UptimeUnit, "uptime-unit", uptimeUnitDefault, uptimeUnitDesc)

	flag.StringVar(&config.InstanceId, "instanceid", "", instanceIdDesc)
	flag.StringVar(&config.NameSpace, "namespace", nameSpaceDefault, nameSpaceDesc)

	flag.BoolVar(&config.Verbose,"verbose", verboseDefault, verboseDesc)
}
