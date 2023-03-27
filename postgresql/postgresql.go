package postgresql

import (
	"context"
	"database/sql"
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

func Close(db *Postgresql) {
	db.postgresql.Close(context.Background())
}

// Table containing validator list with it's corresponding pool
func (db *Postgresql) CreateValidatorPoolTable() error {

	var removeTableQuery = `DROP TABLE IF EXISTS t_validators_pools;`
	if _, err := db.postgresql.Exec(
		context.Background(),
		removeTableQuery); err != nil {
		return err
	}
	var createTableQuery = `
	CREATE TABLE IF NOT EXISTS t_validators_pools
	(
		f_public_key bytea NOT NULL,
		f_pool_name text NOT NULL,
		CONSTRAINT t_validators_pools_f_public_key_key UNIQUE (f_public_key)
	)
	
	TABLESPACE pg_default;

	ALTER TABLE IF EXISTS public.t_validators_pools
    OWNER to chain;
	ALTER TABLE t_validators_pools ADD PRIMARY KEY (f_public_key);
	
	`

	if _, err := db.postgresql.Exec(
		context.Background(),
		createTableQuery); err != nil {
		return err
	}
	return nil
}

func (db *Postgresql) CreateValidatorPoolHistoryTable() error {

	var removeTableQuery = `DROP TABLE IF EXISTS t_validators_pools_history;`
	if _, err := db.postgresql.Exec(
		context.Background(),
		removeTableQuery); err != nil {
		return err
	}
	var createHistoryQuery = `
	CREATE TABLE IF NOT EXISTS t_validators_pools_history
	(
		f_epoch integer NOT NULL,
		f_pool_name text NOT NULL,
		f_active_validators integer NOT NULL
	)

	TABLESPACE pg_default;

	ALTER TABLE IF EXISTS public.t_validators_pools_history
	OWNER to chain;
	ALTER TABLE t_validators_pools_history ADD PRIMARY KEY (f_epoch, f_pool_name);

	`
	if _, err := db.postgresql.Exec(
		context.Background(),
		createHistoryQuery); err != nil {
		return err
	}
	return nil
}

func (db *Postgresql) InsertValidatorPoolHistory(epoch int, pool string, activeValidators int) error {
	var insertQuery = `INSERT INTO t_validators_pools_history(f_epoch,f_pool_name,f_active_validators) VALUES ($1,$2,$3);`
	if _, err := db.postgresql.Exec(
		context.Background(),
		insertQuery,
		epoch,
		pool,
		activeValidators); err != nil {
		return errors.Wrap(err, "could not insert validator pool history")
	}
	return nil
}

func (db *Postgresql) InsertValidatorPool(validator string, pool string) error {
	var insertQuery = `
	INSERT INTO t_validators_pools(f_public_key,f_pool_name)
	VALUES ($1,$2)
	ON CONFLICT (f_public_key) DO UPDATE SET f_pool_name = $2;
	`
	if _, err := db.postgresql.Exec(
		context.Background(),
		insertQuery,
		validator,
		pool); err != nil {
		return errors.Wrap(err, "could not insert validator pool")
	}
	return nil
}

func (db *Postgresql) GetLatestEpoch() (int, error) {
	var query = `SELECT max(f_activation_epoch) FROM t_validators;`
	var latestEpoch int
	err := db.postgresql.QueryRow(context.Background(), query).Scan(&latestEpoch)
	if err != nil {
		return 0, errors.Wrap(err, "could not get latest epoch")
	}
	return latestEpoch, nil
}

func (db *Postgresql) GetValidators() (map[string][]int64, error) {
	var query = `SELECT f_public_key,f_activation_epoch,f_exit_epoch FROM t_validators;`
	rows, err := db.postgresql.Query(context.Background(), query)
	if err != nil {
		return nil, errors.Wrap(err, "could not get validators")
	}

	defer rows.Close()

	validators := map[string][]int64{}
	for rows.Next() {
		var publicKey []byte
		var activationEpochTemp sql.NullInt64
		var activationEpoch int64
		var exitEpochTemp sql.NullInt64
		var exitEpoch int64
		err := rows.Scan(&publicKey, &activationEpochTemp, &exitEpochTemp)
		if err != nil {
			return nil, errors.Wrap(err, "could not get values from row")
		}
		if activationEpochTemp.Valid { // If not, the validator isn't active
			activationEpoch = activationEpochTemp.Int64
			if exitEpochTemp.Valid {
				exitEpoch = exitEpochTemp.Int64
			} else {
				exitEpoch = -1
			}
			validators["\\x"+hex.EncodeToString(publicKey)] = append(validators["\\x"+hex.EncodeToString(publicKey)], activationEpoch)
			validators["\\x"+hex.EncodeToString(publicKey)] = append(validators["\\x"+hex.EncodeToString(publicKey)], exitEpoch)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "could not get values from row")
	}
	return validators, nil
}

func (db *Postgresql) GetPoolValidators(pool string, depositors []string) ([]string, error) {
	var query = `SELECT f_validator_pubkey FROM t_eth1_deposits
	where f_eth1_sender in
	(`
	for _, depositor := range depositors {
		query += "'" + depositor + "',"
	}
	query = query[:len(query)-1] + ");"
	rows, err := db.postgresql.Query(context.Background(), query)
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
