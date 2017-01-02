package cmd

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/viper"
)

func volumeUsed() (resp []metricHandlerResponse, err error) {
	volumes := getVolumesConfigured()

	if len(volumes) == 0 {
		volumes = []string{getVolumeRoot()}
	}

	for _, v := range volumes {
		d, err := disk.Usage(v)

		if err != nil {
			log.Warn(err)
			continue
		}
		resp = append(resp,
			metricHandlerResponse{Dimension: v,
				Value: float64(d.Used) / storageUnits[viper.GetString("volumeUnit")].Multiplier,
				Unit:  storageUnits[viper.GetString("volumeUnit")].Name})
	}

	return resp, nil
}

func volumeAvailable() (resp []metricHandlerResponse, err error) {
	volumes := getVolumesConfigured()

	if len(volumes) == 0 {
		volumes = []string{getVolumeRoot()}
	}

	for _, v := range volumes {
		d, err := disk.Usage(v)

		if err != nil {
			log.Warn(err)
			continue
		}
		resp = append(resp,
			metricHandlerResponse{Dimension: v, Value: 100.0 - d.UsedPercent, Unit: "Percent"})
	}

	return resp, nil
}

func volumeUtil() (resp []metricHandlerResponse, err error) {
	volumes := getVolumesConfigured()

	if len(volumes) == 0 {
		volumes = []string{getVolumeRoot()}
	}

	for _, v := range volumes {
		d, err := disk.Usage(v)

		if err != nil {
			log.Warn(err)
			continue
		}
		resp = append(resp,
			metricHandlerResponse{Dimension: v, Value: d.UsedPercent, Unit: "Percent"})
	}

	return resp, nil
}

func volumeFree() (resp []metricHandlerResponse, err error) {
	volumes := getVolumesConfigured()

	if len(volumes) == 0 {
		volumes = []string{getVolumeRoot()}
	}

	for _, v := range volumes {
		d, err := disk.Usage(v)

		if err != nil {
			log.Warn(err)
			continue
		}
		resp = append(resp,
			metricHandlerResponse{Dimension: v, Value: float64(d.Free) / storageUnits[viper.GetString("volumeUnit")].Multiplier,
				Unit: storageUnits[viper.GetString("volumeUnit")].Name})
	}

	return resp, nil
}

func getVolumesConfigured() (vols []string) {
	viperVols := viper.GetStringSlice("volumes")

	//if it's 1 it could be from cobraflags and comma delimited
	if len(viperVols) == 1 {
		viperVols = strings.Split(viperVols[0], ",")
	}

	volMap := make(map[string]struct{})

	for _, vol := range viperVols {
		volMap[vol] = struct{}{}
	}

	//before we go further, if any one of these is "all", then lets get all we can and append to map
	if _, ok := volMap["all"]; ok {
		delete(volMap, "all")
		//todo:find all volumes
		allVols := getVolumesAll()
		for _, v := range allVols {
			volMap[v] = struct{}{}
		}
	}

	for k, _ := range volMap {
		//todo: test that each volume exists before adding - log.Warn
		if _, err := os.Stat(k); err != nil {
			log.Warn(fmt.Sprintf("Volume %s does not exist.", k))
			continue
		}
		vols = append(vols, k)
	}

	sort.Strings(vols)
	return
}

func getVolumesAll() (vols []string) {
	//todo:read all volumes
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Warnf("Error discovering all volumes. Returning root volume only.\n%s\n", err)
		return []string{getVolumeRoot()}
	}

	log.Debugf("Found %d partitions.\n", len(partitions))
	for _, p := range partitions {
		vols = append(vols, p.Mountpoint)
	}

	return vols
}

func getVolumeRoot() string {
	if runtime.GOOS == "windows" {
		curdir, _ := os.Getwd()
		return curdir[0:3]
	}
	return "/"
}
