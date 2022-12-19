package main

import (
	"os"

	"gitee.com/cpds/cpds-detector/cmd/cpds-detector/app"
	"github.com/sirupsen/logrus"
)

func initLogging() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)
}

func main() {
	initLogging()

	cmd, err := app.NewDetectorCommand()
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
