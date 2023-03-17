package poolIdentifier

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/santi1234567/eth-pools-identifier/config"
	"github.com/santi1234567/eth-pools-identifier/postgresql"

	"github.com/santi1234567/eth-pools-identifier/utils"

	log "github.com/sirupsen/logrus"
)


type PoolIdentifier struct {
	postgresql     *postgresql.Postgresql
	
	config     *config.Config // TODO: Remove repeated parameters
}

func NewPoolIdentifier(
	ctx context.Context,
	config *config.Config) (*PoolIdentifier, error) {


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

	return &PoolIdentifier{postgresql:  pg,config:      config,}, nil
	
	
}
func (a *PoolIdentifier) Run() {
	err := ReadDepositorAddresses(a)

	if err != nil {
		log.Fatal(err)
	}
}


func ReadDepositorAddresses(a *PoolIdentifier) (error) {
	var dir string = "./poolDepositors/"

	files, err := os.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, "could not read files in directory "+ dir)
	}
	var poolSummary []string
	for _, file := range files {
		var fileName string = file.Name()
		var filePath string = dir+fileName
		var poolName string = fileName[0:len(fileName)-4] // remove .txt extension
		f, err := os.Open(filePath)
		log.Info("Getting validators for pool: ", poolName)
		if err != nil {
			return errors.Wrap(err, "could not read file "+ fileName)
		}
	
		defer f.Close()
	
		scanner := bufio.NewScanner(f)
		var depositors []string
		for scanner.Scan() {
			depositors = append(depositors, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return errors.Wrap(err, "could not get validators from depositor addresses corresponding to file "+ fileName)
		}

		validators, err := a.postgresql.GetPoolValidators(poolName, depositors)
		if err != nil {
			return errors.Wrap(err, "could not get pool validators for pool"+poolName+" from postgresql")
		}		
		err = utils.WriteTextFile("./poolValidators/"+poolName+".txt", validators)
		if err != nil {
			return errors.Wrap(err, "could not write validators from pool "+poolName +" to .txt file ")
		}
		poolSummary = append(poolSummary, poolName + "," + fmt.Sprint(len(validators)))
		log.Info("Done getting pool validators for pool: ", poolName, ". Found ", len(validators), " validators")
    }
	
	log.Info("Writing summary file")
	err = utils.WriteTextFile("./poolValidators/poolSummary.txt", poolSummary)
	if err != nil {
		return errors.Wrap(err, "could not write file summary")
	}
	
	log.Info("Done writing summary file")
	return nil
}