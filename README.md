# SLA management with FPC an IPFS

This repository creates a network where SLAs are handled between a provider and a client and where data coming from IPFS get compared inside in a Fabric Private Chaincode to calculate if a violation happened.

This provides security, since it guarantees that even the computation is secured and untempered.

## Using the repository

Make sure to recursively clone this repository by running

```bash
git clone --recursive https://github.com/nkapsoulis/blockchain-sla-management.git
```
## Set-Up

0. Make sure you have a new version of Docker installed, as well as make.
1. Go to `utils/docker` and build the development container with `make build-dev`
2. Run the development container with `make run-dev`
3. Compile the protocol buffers definitions by running `make -C protos`
3. Compile the chaincodes and clients by running `make` on the following folders
    * `cc/public`
    * `cc/private`
    * `clients/api`
    * `clients/ipfs_client`
    * `clients/violation_producer`
4. Run `export FPC_PATH=$(pwd)` on the top level folder.
5. In `cc/public/details.env` add the prefix `PUBLIC` to all the variables (make sure you have the correct permissions to do that).

The next steps could be run in virtual terminals so they can survive a disconnect.

To start the network

6. Add the environmental variables of the ccs by running `source` on `cc/public/details.env` and `cc/private/details.env`
7. `$FPC_PATH/samples/deployment/fabric-smart-client/the-simple-testing-network/env.sh Org1` to get credentials for the gateway
8. In `samples/deployment/fabric-smart-client/the-simple-testing-network` run `make run`.

To run the API

9. On another terminal do step 6 and then also source the `Org1.env` file created in step 7.
10. Run `clients/api/api`


To run the IPFS client

11. On another terminal do step 6 and then also source the `Org1.env` file created in step 7.
12. Start IPFS by running `docker compose up -d` in the `clients/ipfs_client` folder.
13. Run `clients/ipfs_client/ipfsclient`

## Testing

To test the deployment, you will need to have Postman (or a similar tool) installed to interact with the API. A Postman collection is found on the `postman` folder.

The flow that needs to be followed is the following:

* Login as provider-user with the mnemonic created during api startup
* Create a new SLA, with examples from `test/slas` (Only Incident Resolution and Incident Response currently supported)
* Approve given SLA
* Logout
* Login as client-user
* Approve SLA
* Now data from IPFS can start flowing in and violations to happen.

You can put data on IPFS using the `producer` found in `clients/violation_producer`.

## Shutdown

When finished working with the network the following processes need to be stopped (mostly a Ctrl+C would suffice):

1) IPFS: With `docker compose down` in the `clients/ipfs_client` folder.
2) The ipfs_client
3) The api
4) The whole Hyperledger Fabric network