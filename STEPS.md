# Steps to start the system

1. Go to `utils/docker` and build the development container with `make build-dev`
2. Run the development container with `make run-dev`
3. Compile the chaincodes and clients by running `make` on the following folders
    * `cc/public`
    * `cc/private`
    * `clients/api`
    * `clients/ipfs_client`
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

