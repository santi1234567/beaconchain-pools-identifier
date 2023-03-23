package main

import (
	"context"

	"github.com/santi1234567/eth-pools-identifier/config"
	"github.com/santi1234567/eth-pools-identifier/poolHistory"
	"github.com/santi1234567/eth-pools-identifier/poolIdentifier"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := config.NewCliConfig()
	if err != nil {
		log.Fatal(err)
	}
	logLevel, err := log.ParseLevel(config.Verbosity)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	poolIdentifier, err := poolIdentifier.NewPoolIdentifier(context.Background(), config)

	if err != nil {
		log.Fatal(err)
	}

	poolIdentifier.Run()

	if config.History {
		poolHistory, err := poolHistory.NewpoolHistory(context.Background(), config, *poolIdentifier.ValidatorPoolMap)

		if err != nil {
			log.Fatal(err)
		}
		poolHistory.Run()
	}
}
