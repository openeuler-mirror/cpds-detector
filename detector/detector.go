package detector

import (
	"fmt"
	"net/http"
	"time"

	"gitee.com/cpds/cpds-detector/config"
	"gitee.com/cpds/cpds-detector/detector/debug"
	commonv1 "gitee.com/cpds/cpds-detector/pkgs/apis/common/v1"
	rulesv1 "gitee.com/cpds/cpds-detector/pkgs/apis/rules/v1"
	"gitee.com/cpds/cpds-detector/pkgs/rules"
	restful "github.com/emicklei/go-restful"

	"github.com/sirupsen/logrus"
)

var (
	serverTimeout = 5000 * time.Millisecond
)

type Detector struct {
	*config.Config
}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) Run(opts *config.Config) error {
	if err := opts.CheckConfig(); err != nil {
		return err
	}

	if err := configureLogLevel(opts); err != nil {
		return err
	}

	if opts.Debug {
		debug.Enable()
		logrus.Debugf("enable debug mode")
	}

	logrus.Infof("Starting cpds-detector......")
	logrus.Infof("Using config: database address: %s, database port: %s", opts.DatabaseAddress, opts.DatabasePort)
	logrus.Infof("Using config: bind address: %s, listening port: %s", opts.BindAddress, opts.Port)

	wsContainer := restful.NewContainer()
	installAPIs(wsContainer)
	setRestfulConf(wsContainer)
	opts.RegisterSwagger(wsContainer)

	tlsconf := config.GetTlsConf()
	server := &http.Server{
		Addr:        ":" + opts.Port,
		Handler:     wsContainer,
		TLSConfig:   tlsconf,
		ReadTimeout: serverTimeout,
	}
	if err := server.ListenAndServeTLS(opts.CertFile, opts.KeyFile); err != nil {
		logrus.Infof("Failed to listen https://%s:%s: %w", opts.BindAddress, opts.Port, err)
	}
	defer server.Close()

	return nil
}

// configureLogLevel "debug"|"info"|"warn"|"error"|"fatal", default: "info"
func configureLogLevel(opts *config.Config) error {
	logrus.Debug("Configure Log Level: %s", opts.LogLevel)
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
	logrus.Debug("Installing APIs")
	r := rules.New()
	rulesv1.AddToContainer(c, r)
	commonv1.AddToContainer(c)
}

func setRestfulConf(c *restful.Container) {
	logrus.Debug("Setting restful configuration")
	// Add cross origin filter
	cors := config.GetCors(c)
	c.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	c.Filter(c.OPTIONSFilter)
}
