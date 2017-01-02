package cmd

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/viper"
)

func memAvail() (resp []metricHandlerResponse, err error) {
	v, err := mem.VirtualMemory()

	if err != nil {
		return nil, err
	}

	return []metricHandlerResponse{metricHandlerResponse{Value: 100.0 - v.UsedPercent, Unit: "Percent"}}, nil
}

func memUtil() (resp []metricHandlerResponse, err error) {
	v, err := mem.VirtualMemory()

	if err != nil {
		return nil, err
	}

	return []metricHandlerResponse{metricHandlerResponse{Value: v.UsedPercent, Unit: "Percent"}}, nil
}

func memFree() (resp []metricHandlerResponse, err error) {
	v, err := mem.VirtualMemory()

	if err != nil {
		return nil, err
	}

	return []metricHandlerResponse{
		metricHandlerResponse{Value: float64(v.Available) / storageUnits[viper.GetString("memoryUnit")].Multiplier,
			Unit: storageUnits[viper.GetString("memoryUnit")].Name},
	}, nil
}

func memUsed() (resp []metricHandlerResponse, err error) {
	v, err := mem.VirtualMemory()

	if err != nil {
		return nil, err
	}

	return []metricHandlerResponse{
		metricHandlerResponse{Value: float64(v.Used) / storageUnits[viper.GetString("memoryUnit")].Multiplier,
			Unit: storageUnits[viper.GetString("memoryUnit")].Name},
	}, nil
}
