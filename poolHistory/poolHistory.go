package poolHistory

import (
	"context"

	"github.com/pkg/errors"
	"github.com/santi1234567/eth-pools-identifier/config"
	"github.com/santi1234567/eth-pools-identifier/postgresql"

	log "github.com/sirupsen/logrus"
)


type poolHistory struct {
	postgresql     *postgresql.Postgresql	
	validatorPoolMap *map[string]string
	config     *config.Config // TODO: Remove repeated parameters
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

		
		// err = pg.CreateValidatorPoolTable()
		// if err != nil {
		// 	return nil, errors.Wrap(err, "error creating pool table to store data")
		// }
	}

	return &poolHistory{postgresql:  pg,config:      config, validatorPoolMap: &validatorPoolMap}, nil
	
	
}
func (a *poolHistory) Run() {
	err := GetPooHistory(a)

	if err != nil {
		log.Fatal(err)
	}
}




func GetPooHistory(a *poolHistory) (error) {
	log.Info("Getting pool history")
	history := make([]map[string]int, 600000)
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
		history[data[0]][pool] ++ 
		if data[1] != -1 {
			history[data[1]][pool] --
		}
	}		
	// for i := range history[1:] {		
	// 	fmt.Println(history[i])
	// 	for pool := range history[i] {
	// 		history[i+1][pool] += history[i][pool]
	// 	}
	// 	fmt.Println(history[i+1])
	// }
	return nil
}

