package main

import (
	"flag"

	"bitbucket.org/hornetdefiant/core-engine/server"
	"golang.org/x/exp/slog"
)

func main() {
	// Setup the configuration management
	envType := flag.String("env", "dev", "set the env type to dev or prod or staging")
	flag.Parse()

	slog.Info("Running in", "env", *envType)
	server.Run(envType)
}
