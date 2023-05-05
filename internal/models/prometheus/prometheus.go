package prometheus

import (
	"cpds/cpds-detector/pkg/prometheus"
	"errors"
	"time"
)

type Operator interface {
	Query(expr string, timestamp int64) (*prometheus.MetricData, error)
	QueryRange(expr string, startTime, endTime int64, step time.Duration) (*prometheus.MetricData, error)
}

type operator struct {
	prometheusConfig *prometheusConfig
}

type prometheusConfig struct {
	host string
	port int
}

func NewOperator(prometheusHost string, prometheusPort int) Operator {
	return &operator{
		prometheusConfig: &prometheusConfig{
			host: prometheusHost,
			port: prometheusPort,
		},
	}
}

func (o *operator) Query(expr string, timestamp int64) (*prometheus.MetricData, error) {
	client, err := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
	if err != nil {
		return nil, err
	}

	metric := client.GetSingleMetric(expr, time.Unix(timestamp, 0))
	if metric.Error != "" {
		return nil, errors.New(metric.Error)
	}

	return &metric.MetricData, nil
}

func (o *operator) QueryRange(expr string, startTime, endTime int64, step time.Duration) (*prometheus.MetricData, error) {
	client, err := prometheus.NewPrometheus(o.prometheusConfig.host, o.prometheusConfig.port)
	if err != nil {
		return nil, err
	}

	metric := client.GetSingleMetricOverTime(expr, time.Unix(startTime, 0), time.Unix(endTime, 0), step)
	if metric.Error != "" {
		return nil, errors.New(metric.Error)
	}

	return &metric.MetricData, nil
}
