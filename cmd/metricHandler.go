package cmd

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/spf13/viper"
)

var (
	metricSpecs  map[string]metricSpec
	storageUnits map[string]metricUnit
)

type metricSpec struct {
	Name    string
	handler func() ([]metricHandlerResponse, error)
}

type metricHandlerResponse struct {
	Dimension dimension
	Value     float64
	Unit      string
}

type metricUnit struct {
	Name       string
	Multiplier float64
}

type dimension struct {
	Name string
	Value string
}

func (ms *metricSpec) Run() error {
	responses, err := ms.handler()

	if err != nil {
		return err
	}

	for _, resp := range responses {
		if dryRun {
			if resp.Dimension.Name != "" {
				fmt.Printf("%s(%s: %s) %s: %v %s\n", systemId, resp.Dimension.Name ,resp.Dimension.Value, ms.Name, resp.Value, resp.Unit)
			} else {
				fmt.Printf("%s %s: %v %s\n", systemId, ms.Name, resp.Value, resp.Unit)
			}
			continue
		}

		_, err = SendCwMetric(ms.Name, resp.Unit, resp.Value, resp.Dimension)
		if err != nil {
			log.Errorf("Error writing metric %s to CloudWatch.\n", ms.Name)
			return err
		}
	}

	return nil
}

func SendCwMetric(metricName, metricUnit string, metricValue float64, d dimension) (*cloudwatch.PutMetricDataOutput, error) {
	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("InstanceID"),
			Value: aws.String(systemId),
		}}

	if d.Name != "" {
		dimensions = append(dimensions,
		&cloudwatch.Dimension{
			Name: aws.String(d.Name),
			Value: aws.String(d.Value),
		})
	}

	params := &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Dimensions: dimensions,
				Timestamp: aws.Time(time.Now()),
				Unit:      aws.String(metricUnit),
				Value:     aws.Float64(metricValue),
			},
		},
		Namespace: aws.String(viper.GetString("namespace")),
	}
	cw := cloudwatch.New(sess)
	return cw.PutMetricData(params)
}
