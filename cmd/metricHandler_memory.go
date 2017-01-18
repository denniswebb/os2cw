package cmd

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/viper"
)

func memHandler(metric string) (resp []metricHandlerResponse, err error) {
	v, err := mem.VirtualMemory()

	if err != nil {
		return nil, err
	}

	value := 0.0
	unit := storageUnits[viper.GetString("memoryUnit")].Name
	multiplier := storageUnits[viper.GetString("memoryUnit")].Multiplier

	switch metric{
	case "avail":
		value = 100.0 - v.UsedPercent
		unit = "Percent"
	case "free":
		value = float64(v.Available) / multiplier
	case "total":
		value = float64(v.Total) / multiplier
	case "used":
		value = float64(v.Used) / multiplier
	case "util":
		value = v.UsedPercent
		unit = "Percent"
	}

	return []metricHandlerResponse{metricHandlerResponse{Value: value, Unit: unit}}, nil
}

func memAvail() (resp []metricHandlerResponse, err error) {
	return memHandler("avail")
}

func memFree() (resp []metricHandlerResponse, err error) {
	return memHandler("free")
}

func memTotal() (resp []metricHandlerResponse, err error) {
	return memHandler("total")
}

func memUsed() (resp []metricHandlerResponse, err error) {
	return memHandler("used")
}

func memUtil() (resp []metricHandlerResponse, err error) {
	return memHandler("util")
}