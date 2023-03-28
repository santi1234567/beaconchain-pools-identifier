# eth-pools-identifier

This project allows the identification of validators which are part of pools in the consenseus layer of Ethereum. After running it, files will be created containing the public keys of validators from each pool. You can also get the amount of validators for each pool for every epoch, taking into account activation and exit epochs.

# Requirements

The program in this repository checks the ETH deposits for each validator and then links it to a list of depositor addresses corresponding to each pool. The deposits are stored in a database created by [`chaind`](https://github.com/wealdtech/chaind). The lists of depositors were sourced from [this repository](https://github.com/alrevuelta/eth-metrics) and [https://beaconcha.in/](https://beaconcha.in/) and can be modified/added by the user (note that they could be wrong).

# Chaind

The database is created running chaind to obtain the table `t_eth1_deposits`and `t_validators` (in case you want to create the epoch history table). They should take less than a day to be fully synchronized and weight less than 1GB.

# Usage

The flags for using `eth-pools-identifier` are:

- `--postgres` is the address of the database created by `chaind`
- `--pool-history` is used to activate the pool history mode, which creates a table in the directory [poolHistory](https://github.com/santi1234567/eth-pools-identifier/tree/main/poolHistory) containing the number of validators for each pool on every epoch.
- `--write-mode` defines if results will be stored in files or in tables created on the database from `--postgres`. Possible values are `file` and `database`.
- `--read-from` indicates if depositors are read from files in the directory [poolDepositors](https://github.com/santi1234567/eth-pools-identifier/tree/main/poolDepositors) or if they are read from the database given in `--postgres`. In the case that flag `--read-from=database` is intented to be used, the database should contain table `t_depositors` with columns `f_pool_name` (text type) and `f_depositor_address` (bytea type). Possible values are `file` and `database`.

# Note on coinbase validators

Validators on the Coinbase pool make deposits from their own wallets so the method of identifying the pool in which a validator corresponds to doesn't work in that case. The list of validators should be sourced externally and entered manually in the [poolValidators directory](https://github.com/santi1234567/eth-pools-identifier/tree/main/poolValidators). The list in this repository was obtained utilizing the same method from [`eth-deposits`](https://github.com/alrevuelta/eth-deposits) but instead of using Dune Analytics, an API was used to obtain transactions.

# Acknowledgements

This repository was created using these as references:

- [https://github.com/alrevuelta/eth-deposits](https://github.com/alrevuelta/eth-deposits)
- [https://github.com/alrevuelta/eth-metrics](https://github.com/alrevuelta/eth-metrics)
