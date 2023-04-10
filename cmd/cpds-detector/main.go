package main

import (
	detector "cpds/cpds-detector/pkg/cpds-detector/server"
)

func main() {
	cmd := detector.NewCommand()

	if err := cmd.Execute(); err != nil {

	}
}
