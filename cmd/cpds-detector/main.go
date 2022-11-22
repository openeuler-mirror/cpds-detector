package main

import (
	"os"

	"gitee.com/cpds/cpds-detector/config"
	"gitee.com/cpds/cpds-detector/detector"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newDetector() (*cobra.Command, error) {
	conf := config.New()
	cmd := &cobra.Command{
		Use:                   "cpds-detector [OPTIONS]",
		Short:                 "Detect exceptions for Container Problem Detect System",
		Version:               "undefined",
		SilenceUsage:          true,
		SilenceErrors:         true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return detector.RunDetector(conf)
		},
	}
	flags := cmd.Flags()

	conf.LoadConfig(flags)

	return cmd, nil
}

func initLogging() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)
}

func main() {
	initLogging()

	cmd, err := newDetector()
	if err != nil {
		logrus.Error(err)
		// if cannot create new Detector, just exit
		os.Exit(1)
	}
	if err := cmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
