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

package database

import (
	"cpds/cpds-detector/pkg/utils/net"
	timeutils "cpds/cpds-detector/pkg/utils/time"
	"fmt"

	"github.com/spf13/pflag"
)

type Options struct {
	Host               string `json:"host,omitempty" yaml:"host,omitempty"`
	Port               int    `json:"port,omitempty" yaml:"port,omitempty"`
	Username           string `json:"username,omitempty" yaml:"username,omitempty"`
	Password           string `json:"password,omitempty" yaml:"password,omitempty"`
	MaxOpenConnections int    `json:"maxOpenConnections,omitempty" yaml:"maxOpenConnections,omitempty"`
	MaxIdleConnections int    `json:"maxIdleConnections,omitempty" yaml:"maxIdleConnections,omitempty"`
	MaxLifetime        string `json:"maxLifeTime,omitempty" yaml:"maxLifeTime,omitempty"`
}

func NewDatabaseOptions() *Options {
	return &Options{
		Host:               "127.0.0.1",
		Port:               3306,
		MaxOpenConnections: 100,
		MaxIdleConnections: 100,
		MaxLifetime:        "60m",
	}
}

func (s *Options) Validate() []error {
	errs := []error{}

	if !net.IsValidIPAdress(s.Host) {
		errs = append(errs, fmt.Errorf("wrong IP Address format: %s", s.Host))
	}

	if !net.IsValidPort(s.Port) {
		errs = append(errs, fmt.Errorf("invalid port number range: %d, should be 0 - 65535", s.Port))
	}

	if !timeutils.IsValidDuration(s.MaxLifetime) {
		errs = append(errs, fmt.Errorf("invalid time duration format: %s", s.MaxLifetime))
	}
	return errs
}

func (s *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	fs.StringVar(&s.Host, "database-host", c.Host, "database service host address. If left blank, the following related mysql options will be ignored.")
	fs.IntVar(&s.Port, "database-port", c.Port, "database service port.")
	fs.StringVar(&s.Username, "database-username", c.Username, "Username for access to database service.")
	fs.StringVar(&s.Password, "database-password", c.Password, "Password for access to database, should be used pair with password.")
	fs.IntVar(&s.MaxOpenConnections, "database-max-open-connections", c.MaxOpenConnections, "Maximum open connections allowed to connect to database.")
}
