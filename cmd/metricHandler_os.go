package cmd

import (
	"github.com/shirou/gopsutil/host"
)

func osHandler(metric string) (resp []metricHandlerResponse, err error) {
	h, err := host.Info()

	if err != nil {
		return nil, err
	}

	value := 0.0
	unit := "None"

	switch metric {
	case "uptime":
		value = float64(h.Uptime)
		unit = "Seconds"
	case "procs":
		value = float64(h.Procs)
	}
	return []metricHandlerResponse{metricHandlerResponse{Value: value, Unit: unit}}, nil
}

func uptime() (resp []metricHandlerResponse, err error) {
	return osHandler("uptime")
}

func procs() (resp []metricHandlerResponse, err error) {
	return osHandler("procs")
}
