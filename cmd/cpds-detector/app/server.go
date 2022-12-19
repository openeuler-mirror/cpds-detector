package app

import (
	"gitee.com/cpds/cpds-detector/config"
	"gitee.com/cpds/cpds-detector/detector"
	"github.com/spf13/cobra"
)

func NewDetectorCommand() (*cobra.Command, error) {
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
			return Run(conf)
		},
	}
	flags := cmd.Flags()

	conf.LoadConfig(flags)

	return cmd, nil
}

func Run(opts *config.Config) error {
	detector := detector.NewDetector()
	return detector.Run(opts)
}
