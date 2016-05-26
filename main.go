package main

import (
	//	"flag"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Config struct {
	MemoryFreePercentage bool
	MemoryFreeAmount     bool
	MemoryUsedPercentage bool
	MemoryUsedAmount     bool
	MemoryUnit           string

	DiskFreePercentage bool
	DiskFreeAmount     bool
	DiskUsedPercentage bool
	DiskUsedAmount     bool
	DiskUnit           string

	Uptime     bool
	UptimeUnit string
}

var config *Config

func init() {
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

		diskFreePercentageDefault = false
		diskFreePercentageDesc    = "Push disk Free Percentage Custom Metric"
		diskFreeAmountDefault     = false
		diskFreeAmountDesc        = "Push disk Free Amount Custom Metric"
		diskUsedPercentageDefault = false
		diskUsedPercentageDesc    = "Push disk Used Percentage Custom Metric"
		diskUsedAmountDefault     = false
		diskUsedAmountDesc        = "Push disk Used Amount Custom Metric"
		diskUnitDefault           = "kb"
		diskUnitDesc              = "Unit for disk metrics. Allowed: tb, gb, mb, kb, b. Default: mb"

		uptimeDefault     = false
		uptimeDesc        = "System uptime"
		uptimeUnitDefault = "s"
		uptimeUnitDesc    = "Unit for uptime metric. Allowed: s, m, h, d. Default: s"
	)

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
}

func main() {
	flag.Parse()

	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")
	h, _ := host.Info()

	var memUnitDivisor, diskUnitDivisor, uptimeUnitDivisor float64
	var uptimeUnitText string
	memUnitDivisor = 1
	diskUnitDivisor = 1
	uptimeUnitDivisor = 1

	switch {
	case config.MemoryUnit == "kb":
		memUnitDivisor = 1024
	case config.MemoryUnit == "mb":
		memUnitDivisor = 1024 * 1024
	case config.MemoryUnit == "gb":
		memUnitDivisor = 1024 * 1024 * 1024
	default:
		config.MemoryUnit = "b"
		memUnitDivisor = 1
	}

	if config.MemoryFreePercentage {
		fmt.Printf("Memory Free Percent: %f%%\n", (100.0 - v.UsedPercent))
	}

	if config.MemoryUsedPercentage {
		fmt.Printf("Memory Used Percent: %f%%\n", v.UsedPercent)
	}

	if config.MemoryFreeAmount {
		fmt.Printf("Memory Free(%s): %f\n", config.MemoryUnit, float64(v.Available)/memUnitDivisor)
	}

	if config.MemoryUsedAmount {
		fmt.Printf("Memory Used(%s): %f\n", config.MemoryUnit, float64(v.Used)/memUnitDivisor)
	}

	switch {
	case config.DiskUnit == "kb":
		diskUnitDivisor = 1024
	case config.DiskUnit == "mb":
		diskUnitDivisor = 1024 * 1024
	case config.DiskUnit == "gb":
		diskUnitDivisor = 1024 * 1024 * 1024
	case config.DiskUnit == "tb":
		diskUnitDivisor = 1024 * 1024 * 1024 * 1024
	default:
		config.DiskUnit = "b"
		diskUnitDivisor = 1024 * 1024
	}

	if config.DiskFreePercentage {
		fmt.Printf("Disk Free Percent: %f%%\n", (100.0 - d.UsedPercent))
	}

	if config.DiskUsedPercentage {
		fmt.Printf("Disk Used Percent: %f%%\n", d.UsedPercent)
	}

	if config.DiskFreeAmount {
		fmt.Printf("Disk Free(%s): %f\n", config.DiskUnit, float64(d.Free)/diskUnitDivisor)
	}

	if config.DiskUsedAmount {
		fmt.Printf("Disk Used(%s): %f\n", config.DiskUnit, float64(d.Used)/diskUnitDivisor)
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
		fmt.Printf("System Uptime(%s): %v\n", uptimeUnitText, float64(h.Uptime)/uptimeUnitDivisor)
	}
}
