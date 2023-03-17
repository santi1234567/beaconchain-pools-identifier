package config

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

// By default the release is a custom build. CI takes care of upgrading it with
// go build -v -ldflags="-X 'github.com/alrevuelta/eth-pools-metrics/config.ReleaseVersion=x.y.z'"
var ReleaseVersion = "custom-build"

type Config struct {
	Postgres              string
	Verbosity             string
	History               bool
}


func NewCliConfig() (*Config, error) {
	var version = flag.Bool("version", false, "Prints the release version and exits")
	var postgres = flag.String("postgres", "", "Postgres db endpoit: postgresql://user:password@netloc:port")
	var verbosity = flag.String("verbosity", "info", "Logging verbosity (trace, debug, info=default, warn, error, fatal, panic)")
	var poolHistory = flag.Bool("pool-history", false, "If true, it will create a file with daily pool data")
	flag.Parse()

	if *version {
		log.Info("Version: ", ReleaseVersion)
		os.Exit(0)
	}



	conf := &Config{
		Postgres:              *postgres,
		Verbosity:             *verbosity,
		History:               *poolHistory,
	}
	logConfig(conf)
	return conf, nil
}

func logConfig(cfg *Config) {
	log.WithFields(log.Fields{
		"Postgres":              cfg.Postgres,
		"Verbosity":             cfg.Verbosity,
		"Pool-History":               cfg.History,
	}).Info("Cli Config:")
}