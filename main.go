package main

import (
	"os"

	"github.com/hosseintrz/gaterun/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Info("hererer")
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Could not run command")
	}
}
