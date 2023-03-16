package poolIdentifier

import (
	"bufio"
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/santi1234567/eth-pools-identifier/config"
	"github.com/santi1234567/eth-pools-identifier/postgresql"

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
	_ = ReadDepositorAddresses(a)
}
func ReadDepositorAddresses(a *PoolIdentifier) (error) {
	var dir string = "./poolDepositors/"

	files, err := os.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, "could not read files in directory "+ dir)
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
			return errors.Wrap(err, "could not get validators from depositor addresses corresponding to file "+ fileName)
		}

		validators, err := a.postgresql.GetPoolValidators(poolName, depositors)
		if err != nil {
			return errors.Wrap(err, "could not get pool validators for pool"+poolName+" from postgresql")
		}		
		err = WritePoolValidatorsFile(poolName, validators)
		if err != nil {
			return errors.Wrap(err, "could not write validators from pool "+poolName +" to .txt file ")
		}
		
		log.Info("Done getting pool validators for pool: ", poolName, ". Found ", len(validators), " validators")
    }
	
	return nil
}
func WritePoolValidatorsFile(pool string, validators []string) (error) {
	var dir string = "./poolValidators/"
	var fileName string = dir+pool+".txt"
	f, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "could not create file "+ fileName)
	}
	defer f.Close()

	for _, validator := range validators {
		_, err := f.WriteString(validator + "\n")
		if err != nil {
			return errors.Wrap(err, "could not write validator "+ validator + " to file "+ fileName)
		}
	}
	return nil
}