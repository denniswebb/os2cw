package cmd

import (
	"github.com/shirou/gopsutil/host"
)

func uptime() (resp []metricHandlerResponse, err error) {
	h, err := host.Info()

	if err != nil {
		return nil, err
	}

	return []metricHandlerResponse{
		metricHandlerResponse{Value: float64(h.Uptime),
			Unit: "Seconds"},
	}, nil
}
