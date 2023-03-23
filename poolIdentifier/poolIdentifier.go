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
	ValidatorPoolMap *map[string]string
	config     *config.Config // TODO: Remove repeated parameters
}

func NewPoolIdentifier(
	ctx context.Context,
	config *config.Config) (*PoolIdentifier, error) {

	var validatorPoolMap = make(map[string]string)
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

	return &PoolIdentifier{postgresql:  pg,config:      config,ValidatorPoolMap: &validatorPoolMap}, nil	
}


func (a *PoolIdentifier) Run() () {
	err := ReadCoinbaseValidators(a)
	if err != nil {
		log.Fatal(err)
	}

	err = ReadDepositorAddresses(a)

	if err != nil {
		log.Fatal(err)
	}
}

func ReadCoinbaseValidators(a *PoolIdentifier) (error) {
	filePath := "./poolValidators/coinbase.txt"
	log.Info("Getting validators for pool: coinbase")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
        log.Info("No coinbase validators file found")
		return nil
	} else {
		f, err := os.Open(filePath)
		if err != nil {
			return errors.Wrap(err, "could not read file coinbase.txt")
		}
	
		defer f.Close()
	
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			validator := scanner.Text()
			(*a.ValidatorPoolMap)[validator] = "coinbase"
		}
		if err := scanner.Err(); err != nil {
			return errors.Wrap(err, "could not get coinbase validators from file coinbase.txt")
		}
		log.Info("Done getting pool validators for pool: coinbase. Found ", len(*a.ValidatorPoolMap), " validators")
		return nil
	}
}

func ReadDepositorAddresses(a *PoolIdentifier) (error) {
	var dir string = "./poolDepositors/"
	files, err := os.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, "could not read files in directory "+ dir)
	}
	var poolSummary []string

	if len(*a.ValidatorPoolMap) >0 {
		poolSummary = append(poolSummary, "coinbase," + fmt.Sprint(len(*a.ValidatorPoolMap)))
	}
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
			return errors.Wrap(err, "could not get depositor addresses corresponding to file "+ fileName)
		}

		validators, err := a.postgresql.GetPoolValidators(poolName, depositors)
		if a.config.History {
			for _, validator := range validators {
				(*a.ValidatorPoolMap)[validator] = poolName
			}
		}
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