package cmd

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/viper"
)

func volHandler(metric string) (resp []metricHandlerResponse, err error) {
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

		value := 0.0
		unit := storageUnits[viper.GetString("volumeUnit")].Name
		multiplier := storageUnits[viper.GetString("volumeUnit")].Multiplier
		metDim := dimension{Name: "Volume", Value: v}

		switch metric {
		case "avail":
			value = 100.0 - d.UsedPercent
			unit = "Percent"
		case "free":
			value = float64(d.Free) / multiplier
		case "total":
			value = float64(d.Total) / multiplier
		case "used":
			value = float64(d.Used) / multiplier
		case "util":
			value = d.UsedPercent
			unit = "Percent"
		}
		resp = append(resp,
			metricHandlerResponse{Dimension: metDim, Value: value, Unit: unit})
	}

	return resp, nil
}

func volumeAvailable() (resp []metricHandlerResponse, err error) {
	return volHandler("avail")
}

func volumeFree() (resp []metricHandlerResponse, err error) {
	return volHandler("free")
}

func volumeTotal() (resp []metricHandlerResponse, err error) {
	return volHandler("total")
}

func volumeUsed() (resp []metricHandlerResponse, err error) {
	return volHandler("used")
}

func volumeUtil() (resp []metricHandlerResponse, err error) {
	return volHandler("util")
}

func getVolumesConfigured() (vols []string) {
	viperVols := viper.GetStringSlice("volumes")

	//if it's 1 it could be from cobraflags and comma delimited
	if len(viperVols) == 1 {
		viperVols = strings.Split(viperVols[0], ",")
	}

	volMap := make(map[string]struct{})

	for _, vol := range viperVols {
		//uppercase windows drives and remove \
		match, _ := regexp.MatchString("^[[:alpha:]]:\\\\*$", vol)
		if match {
			vol = strings.ToUpper(vol[0:2])
		}
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

	for k := range volMap {
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
		return curdir[0:2]
	}
	return "/"
}
