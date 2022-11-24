package detector

import (
	"fmt"
	"os"

	"gitee.com/cpds/cpds-detector/config"
	"github.com/sirupsen/logrus"
)

func RunDetector(opts *config.Config) error {
	if err := configureLogLevel(opts); err != nil {
		return err
	}

	if opts.Debug {
		enableDebug()
		logrus.Debugf("Enable debug mode")
	}

	logrus.Infof("Starting cpds-detector......")
	logrus.Infof("Using config: database address: %s, database port: %s", opts.DatabaseAddress, opts.DatabasePort)
	logrus.Infof("Using config: bind address: %s, listening port: %s", opts.BindAddress, opts.Port)
	// TODO: complete this function
	return nil
}

// enableDebug sets the DEBUG env var to true
// and makes the logger to log at debug level.
func enableDebug() {
	os.Setenv("DEBUG", "1")
	logrus.SetLevel(logrus.DebugLevel)
}

// disableDebug sets the DEBUG env var to false
// and makes the logger to log at info level.
func disableDebug() {
	os.Setenv("DEBUG", "")
	logrus.SetLevel(logrus.InfoLevel)
}

// isDebugEnabled checks whether the debug flag is set or not.
func isDebugEnabled() bool {
	return os.Getenv("DEBUG") != ""
}

// configureLogLevel "debug"|"info"|"warn"|"error"|"fatal", default: "info"
func configureLogLevel(opts *config.Config) error {
	if opts.LogLevel != "" {
		lvl, err := logrus.ParseLevel(opts.LogLevel)
		if err != nil {
			return fmt.Errorf("unable to parse logging level: %s", opts.LogLevel)
		}
		logrus.SetLevel(lvl)
	} else {
		// Set InfoLevel as default logLevel
		// Only log the info severity or above.
		logrus.SetLevel(logrus.InfoLevel)
	}
	return nil
}