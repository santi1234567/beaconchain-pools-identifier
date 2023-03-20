# eth-pools-identifier

This project allows the identification of validators which are part of pools in the consenseus layer of Ethereum. After running it, files will be created containing the public keys of validators from each pool. You can also get the amount of validators for each pool for every epoch.

# Requirements

The program in this repository checks the ETH deposits for each validator and then links it to a list of depositor addresses corresponding to each pool. The deposits are stored in a database created by [`chaind`](https://github.com/wealdtech/chaind). The lists of depositors were sourced from [this repository](https://github.com/alrevuelta/eth-metrics) and can be modified/added by the user (note that they could be wrong). 

# Chaind

The database is created running chaind to obtain the table `t_eth1_deposits`and `t_validators` (in case you want to create the epoch history table). They should take less than a day to be fully synchronized and weight less than 1GB. 

# Usage

There are two main flags for using `eth-pools-identifier`. `--postgres`is the address of the database created by `chaind` and `pool-history` is used to activate the pool history mode, which creates a table containing the number of validators for each pool on every epoch.

# Acknowledgements

This repository was created using these as references: 

* [https://github.com/alrevuelta/eth-deposits](https://github.com/alrevuelta/eth-deposits)
* [https://github.com/alrevuelta/eth-metrics](https://github.com/alrevuelta/eth-metrics)
