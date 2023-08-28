/* 
 *  Copyright 2023 CPDS Author
 *  
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *  
 *       https://www.apache.org/licenses/LICENSE-2.0
 *  
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
