package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ryanuber/columnize"
	"sort"
	"k8s.io/kops/_vendor/github.com/aws/aws-sdk-go/aws"
)

type SendCmd struct {
	cobraCommand *cobra.Command
}

var (
	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send OS metrics to CloudWatch",
	}

	sess *session.Session

	systemId string
	dryRun   bool
)

func init() {
	cmd := sendCmd
	cmd.Run = send

	rootCommand.AddCommand(cmd)

	cmd.PersistentFlags().StringP("mem-unit", "m", "", "memory size unit (b, kb, mb, gb)")
	cmd.PersistentFlags().StringP("vol-unit", "u", "", "volume size unit (b, kb, mb, gb, tb)")
	cmd.PersistentFlags().StringP("namespace", "n", "", "CloudWatch namespace")
	cmd.PersistentFlags().StringVarP(&systemId, "id", "i", "", "system id to store metrics")
	cmd.PersistentFlags().StringSliceP("volumes", "v", []string{}, "volumes to report (examples: /,/home,C:)")
	cmd.PersistentFlags().BoolVarP(&dryRun, "dryrun", "", false, "output metrics without sending to CloudWatch")

	viper.BindPFlag("memoryUnit", cmd.PersistentFlags().Lookup("mem-unit"))
	viper.BindPFlag("volumeUnit", cmd.PersistentFlags().Lookup("vol-unit"))
	viper.BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("systemId", cmd.PersistentFlags().Lookup("id"))
	viper.BindPFlag("volumes", cmd.PersistentFlags().Lookup("volumes"))
	viper.BindPFlag("dryrun", cmd.PersistentFlags().Lookup("dryrun"))

	viper.SetDefault("metrics", []string{})
	viper.SetDefault("memoryUnit", "kb")
	viper.SetDefault("volumeUnit", "mb")
	viper.SetDefault("namespace", "System")
	viper.SetDefault("volumes", []string{"all"})

	metricSpecs = make(map[string]metricSpec)
	metricSpecs["mem-avail"] = metricSpec{Name: "MemoryFreePercentage",
		handler: memAvail}
	metricSpecs["mem-free"] = metricSpec{Name: "MemoryFree",
		handler: memFree}
	metricSpecs["mem-used"] = metricSpec{Name: "MemoryUsed",
		handler: memUsed}
	metricSpecs["mem-util"] = metricSpec{Name: "MemoryUsedPercentage",
		handler: memUtil}
	metricSpecs["vol-avail"] = metricSpec{Name: "VolumeFreePercentage",
		handler: volumeAvailable}
	metricSpecs["vol-free"] = metricSpec{Name: "VolumeFree",
		handler: volumeFree}
	metricSpecs["vol-used"] = metricSpec{Name: "VolumeUsed",
		handler: volumeUsed}
	metricSpecs["vol-util"] = metricSpec{Name: "VolumeUsedPercentage",
		handler: volumeUtil}
	metricSpecs["uptime"] = metricSpec{Name: "Uptime",
		handler: uptime}

	updateUsageTemplate()

	storageUnits = make(map[string]metricUnit)
	storageUnits["b"] = metricUnit{Name: "Bytes", Multiplier: 1}
	storageUnits["kb"] = metricUnit{Name: "Kilobytes", Multiplier: 1024}
	storageUnits["mb"] = metricUnit{Name: "Megabytes", Multiplier: 1024 * 1024}
	storageUnits["gb"] = metricUnit{Name: "Gigabytes", Multiplier: 1024 * 1024 * 1024}
	storageUnits["tb"] = metricUnit{Name: "Terabytes", Multiplier: 1024 * 1024 * 1024 * 1024}
}

func updateUsageTemplate() {
	var metricArgHelp []string

	var args []string
	for k,_ := range metricSpecs {
		args = append(args,k)
	}
	sort.Strings(args)

	sendCmd.ValidArgs = args
	sendCmd.Example = "  os2cw send -u gb -m mb -v / -v /home mem-avail mem-used vol-free uptime"

	for _,arg := range args {
		metricArgHelp = append(metricArgHelp,fmt.Sprintf("%s | %s \n", arg,metricSpecs[arg].Name))
	}

	sendCmd.SetUsageTemplate(
		fmt.Sprintf("%s\nAvailable Metrics:\n%s\n\n",
			sendCmd.UsageTemplate(),
			columnize.Format(metricArgHelp,
				columnize.MergeConfig(columnize.DefaultConfig(),
					&columnize.Config{Prefix:"      "}))))
}

func send(cmd *cobra.Command, args []string) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	sess = session.New()

	if systemId == "" {
		systemId = generateId()
	}

	if systemId == "" {
		sendCmd.Usage()
		log.Errorf("Unable to generate system id.\n")
		code = 1
		return
	}

	setRegion(sess)

	metrics := viper.GetStringSlice("metrics")

	//override configured metrics when specified via CLI
	if len(args) > 0 {
		metrics = args
	}

	if len(metrics) == 0 {
		sendCmd.Usage()
		log.Errorf("No metrics specified.\n")
		code = 1
		return
	}

	if _, ok := storageUnits[viper.GetString("volumeUnit")]; !ok {
		sendCmd.Usage()
		log.Errorf("Invalid volume unit: %s\n\n", viper.GetString("volumeUnit"))
		code = 1
		return
	}

	if _, ok := storageUnits[viper.GetString("memoryUnit")]; !ok {
		sendCmd.Usage()
		log.Errorf("Invalid memory unit: %s\n\n", viper.GetString("memoryUnit"))
		code = 1
		return
	}

	//convert to map to remove dupes
	metricsDistinct := make(map[string]struct{})
	for _, metric := range metrics {
		metricsDistinct[metric] = struct{}{}
	}

	for metric, _ := range metricsDistinct {
		s, ok := metricSpecs[metric]
		if !ok {
			fmt.Printf("Error: Invalid metric %s provided.\n", metric)
			continue
		}

		err := s.Run()
		if err != nil {
			log.Errorf("An error occured during metric run.\n%s\n", err)
			code = 1
		}
	}
}

func generateId() string {
	//try to read instanceid from metadata
	metadataService := ec2metadata.New(sess)
	if instanceid, err := metadataService.GetMetadata("instance-id"); err == nil {
		return instanceid
	}

	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return ""
}

func setRegion(s *session.Session) {
	metadataService := ec2metadata.New(sess)
	if region, err := metadataService.Region(); err == nil {
		s.Config.Region = aws.String(region)
	}
}