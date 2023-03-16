package postgresql

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)


type Postgresql struct {
	postgresql *pgx.Conn
	PoolName   string
}

// postgresql://user:password@netloc:port/dbname
func New(postgresEndpoint string) (*Postgresql, error) {
	conn, err := pgx.Connect(context.Background(), postgresEndpoint)

	if err != nil {
		return nil, err
	}

	return &Postgresql{
		postgresql: conn,
	}, nil
}

// Table containing validator list with it's corresponding pool
func (db *Postgresql) CreateValidatorPoolTable() error {

	var createTableQuery = `
	CREATE TABLE IF NOT EXISTS t_validators_pools
	(
		f_public_key bytea NOT NULL,
		f_pool_name text NOT NULL,
	);
	`

	if _, err := db.postgresql.Exec(
		context.Background(),
		createTableQuery); err != nil {
		return err
	}
	return nil
}


func (db *Postgresql) GetPoolValidators(pool string, depositors []string) ([]string, error) {
	var query = `SELECT f_validator_pubkey FROM t_eth1_deposits
	where f_eth1_sender in
	(`
	for _, depositor := range depositors {
		query += "'"+depositor+"',"
	}
	query = query[:len(query)-1] + ");"
	rows, err := db.postgresql.Query(context.Background(),query)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s: %s", "could not get validators for pool", pool))
	}

	defer rows.Close()
	
		
	var validators []string
	for rows.Next() {
		var data []byte
		err := rows.Scan(&data)
		if err != nil {
			return nil, errors.Wrap(err, "could not get values from row for pool "+pool)
		}
		validators = append(validators, "\\x"+hex.EncodeToString(data))
	}
    if err := rows.Err(); err != nil {
        return nil, errors.Wrap(err, "could not get values from row for pool "+pool)
    }
	return validators, nil
}