package detector

import (
	"gitee.com/cpds/cpds-detector/config"
	"github.com/sirupsen/logrus"
)

func RunDetector(opts *config.Options) error {
	logrus.Infof("Starting cpds-detector......")
	logrus.Infof("Using config: database address: %s, database port: %s", opts.DatabaseAddress, opts.DatabasePort)
	// TODO: complete this function
	return nil
}
