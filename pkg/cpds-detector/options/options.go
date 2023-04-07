package options

import (
	detector "cpds/cpds-detector/pkg/cpds-detector"
	"cpds/cpds-detector/pkg/cpds-detector/config"

	"github.com/spf13/pflag"
)

type ServerRunOptions struct {
	ConfigFile string
	*config.Config

	DebugMode bool
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Config: config.New(),
	}
}

func (s *ServerRunOptions) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("generic", pflag.ExitOnError)
	fs.StringVar(&s.ConfigFile, "config", config.DefaultConfigurationPath, "Directory where configuration files are stored")
	fs.BoolVar(&s.DebugMode, "debug", false, "Don't enable this if you don't know what it means.")

	return fs
}

func (s *ServerRunOptions) NewDetector() (*detector.Detector, error) {
	detector := &detector.Detector{
		Config: s.Config,
		Debug:  s.DebugMode,
	}
	return detector, nil
}
