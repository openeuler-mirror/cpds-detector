package detector

import (
	"fmt"
	"net/http"
	"os"

	"gitee.com/cpds/cpds-detector/config"
	rulesv1 "gitee.com/cpds/cpds-detector/pkgs/apis/rules/v1"
	"gitee.com/cpds/cpds-detector/pkgs/rules"
	restful "github.com/emicklei/go-restful"

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

	wsContainer := restful.NewContainer()
	installAPIs(wsContainer)

	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	server := &http.Server{
		Addr:    ":" + opts.Port,
		Handler: wsContainer,
	}
	if err := server.ListenAndServeTLS(opts.CertFile, opts.KeyFile); err != nil {
		logrus.Infof("Failed to listen https://%s:%s: %w", opts.BindAddress, opts.Port, err)
	}
	defer server.Close()

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

func installAPIs(c *restful.Container) {
	r := rules.New()
	rulesv1.AddToContainer(c, r)
}
