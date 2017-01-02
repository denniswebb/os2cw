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
	Dimension string
	Value     float64
	Unit      string
}

type metricUnit struct {
	Name       string
	Multiplier float64
}

func (ms *metricSpec) Run() error {
	responses, err := ms.handler()

	if err != nil {
		return err
	}

	for _, resp := range responses {
		if dryRun {
			if resp.Dimension != "" {
				fmt.Printf("%s(%s) %s: %v %s\n", systemId, resp.Dimension, ms.Name, resp.Value, resp.Unit)
			} else {
				fmt.Printf("%s %s: %v %s\n", systemId, ms.Name, resp.Value, resp.Unit)
			}
			continue
		}

		_, err = SendCwMetric(ms.Name, resp.Unit, resp.Value)
		if err != nil {
			log.Errorf("Error writing metric %s to CloudWatch.\n", ms.Name)
			return err
		}
	}

	return nil
}

func SendCwMetric(metricName string, metricUnit string, metricValue float64) (*cloudwatch.PutMetricDataOutput, error) {
	params := &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("InstanceID"),
						Value: aws.String(systemId),
					},
				},
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
