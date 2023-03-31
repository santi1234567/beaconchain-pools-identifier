package poolHistory

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/santi1234567/eth-pools-identifier/config"
	"github.com/santi1234567/eth-pools-identifier/postgresql"
	"github.com/santi1234567/eth-pools-identifier/utils"

	log "github.com/sirupsen/logrus"

	progressbar "github.com/schollz/progressbar/v3"
)

type poolHistory struct {
	postgresql       *postgresql.Postgresql
	validatorPoolMap *map[string]string
	config           *config.Config // TODO: Remove repeated parameters
}

func NewpoolHistory(
	ctx context.Context,
	config *config.Config, validatorPoolMap map[string]string) (*poolHistory, error) {

	var pg *postgresql.Postgresql
	var err error
	if config.Postgres != "" {
		pg, err = postgresql.New(config.Postgres)
		if err != nil {
			return nil, errors.Wrap(err, "could not create postgresql")
		}

		if config.WriteMode == "database" {
			if config.History {
				err = pg.CreateValidatorPoolHistoryTable()
				if err != nil {
					return nil, errors.Wrap(err, "error creating validator pool history table")
				}
			}
		}
	}

	return &poolHistory{postgresql: pg, config: config, validatorPoolMap: &validatorPoolMap}, nil

}
func (a *poolHistory) Run() {
	defer postgresql.Close(a.postgresql)

	err := GetPooHistory(a)

	if err != nil {
		log.Fatal(err)
	}
}

func GetPooHistory(a *poolHistory) error {
	log.Info("Getting pool history")
	var latestEpoch, err = a.postgresql.GetLatestEpoch()
	log.Info("Latest epoch recorded: ", latestEpoch)
	if err != nil {
		return errors.Wrap(err, "could not get latest epoch from postgresql")
	}
	history := make([]map[string]int, latestEpoch+1)
	for i := range history {
		history[i] = make(map[string]int)
	}
	validators, err := a.postgresql.GetValidators()
	if err != nil {
		return errors.Wrap(err, "could not get validators from postgresql")
	}
	for validator, data := range validators {
		var pool string = (*a.validatorPoolMap)[validator]
		if pool == "" {
			pool = "unknown"
		}
		activationEpoch := data[0]
		exitEpoch := data[1]
		history[activationEpoch][pool]++
		if exitEpoch != -1 {
			history[exitEpoch][pool]--
		}
	}
	for i := range history[1:] {
		for pool := range history[i] {
			history[i+1][pool] += history[i][pool]
		}
	}
	if a.config.WriteMode == "file" {
		var rows []string
		// write header
		var header string = "epoch,"
		for pool := range history[len(history)-1] {
			header += pool + ","
		}
		header = header[:len(header)-1]
		rows = append(rows, header)
		for epoch := range history {
			var row string = fmt.Sprint(epoch) + ","
			for _, pool := range strings.Split(header, ",")[1:] {
				row += fmt.Sprint(history[epoch][pool]) + ","
			}
			rows = append(rows, row[:len(row)-1])
		}
		err = utils.WriteTextFile("./poolHistory/poolHistory.csv", rows)
		if err != nil {
			return errors.Wrap(err, "could not write pool history file")
		}
	} else {
		bar := progressbar.Default(int64(len(history)))
		for epoch := range history {
			bar.Add(1)
			err = a.postgresql.InsertValidatorPoolHistory(epoch, history[epoch])
			if err != nil {
				return errors.Wrap(err, "could not insert pool history")
			}

		}
	}

	log.Info("Done getting pool history")
	return nil
}
