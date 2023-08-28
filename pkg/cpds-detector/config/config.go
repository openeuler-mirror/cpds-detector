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

package config

import (
	"cpds/cpds-detector/pkg/cpds-detector/config/database"
	"cpds/cpds-detector/pkg/cpds-detector/config/generic"
	"cpds/cpds-detector/pkg/cpds-detector/config/logger"
	"cpds/cpds-detector/pkg/cpds-detector/config/prometheus"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// DefaultConfigurationName is the default name of configuration
	DefaultConfigurationName = "config"

	// DefaultConfigurationPath the default location of the configuration file
	DefaultConfigurationPath = "/etc/cpds/detector"
)

// Config defines everything needed for cpds-detector to deal with external services
type Config struct {
	GenericOptions    *generic.Options    `json:"generic,omitempty" yaml:"generic,omitempty" mapstructure:"generic"`
	DatabaseOptions   *database.Options   `json:"database,omitempty" yaml:"database,omitempty" mapstructure:"database"`
	PrometheusOptions *prometheus.Options `json:"prometheus,omitempty" yaml:"prometheus,omitempty" mapstructure:"prometheus"`
	LoggerOptions     *logger.Options     `json:"log,omitempty" yaml:"log,omitempty" mapstructure:"log"`
}

func New() *Config {
	return &Config{
		GenericOptions:    generic.NewGenericOptions(),
		DatabaseOptions:   database.NewDatabaseOptions(),
		PrometheusOptions: prometheus.NewPrometheusOptions(),
		LoggerOptions:     logger.NewLoggerOptions(),
	}
}

func TryLoadFromDisk(path string, debug bool) (*Config, error) {
	viper.SetConfigName(DefaultConfigurationName)

	// Config flag not set, using default configuration file
	if path != DefaultConfigurationPath {
		viper.AddConfigPath(path)
	} else {
		viper.AddConfigPath(DefaultConfigurationPath)
	}

	// Load from current working directory, only used for debugging
	if debug {
		viper.AddConfigPath(".")
	}

	// Load from Environment variables
	viper.SetEnvPrefix("cpds-detector")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		} else {
			return nil, fmt.Errorf("error parsing configuration file %s", err)
		}
	}

	conf := New()

	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
