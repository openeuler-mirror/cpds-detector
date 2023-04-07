package prometheus

import (
	"cpds/cpds-detector/pkg/utils/net"
	"fmt"

	"github.com/spf13/pflag"
)

type Options struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
	Port int    `json:"port,omitempty" yaml:"port,omitempty"`
}

func NewPrometheusOptions() *Options {
	return &Options{
		Host: "127.0.0.1",
		Port: 9095,
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

	return errs
}

func (s *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	fs.StringVar(&s.Host, "prometheus-host", c.Host, "prometheus host IP address")
	fs.IntVar(&s.Port, "prometheus-port", c.Port, "prometheus port number")
}
