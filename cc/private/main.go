package main

import (
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	fpc "github.com/hyperledger/fabric-private-chaincode/ecc_go/chaincode"
)

func main() {

	ccid := os.Getenv("CHAINCODE_PKG_ID")
	addr := os.Getenv("CHAINCODE_SERVER_ADDRESS")

	assetChaincode, _ := contractapi.NewChaincode(&SmartContract{})
	chaincode := fpc.NewPrivateChaincode(assetChaincode)

	// start chaincode as a service
	server := &shim.ChaincodeServer{
		CCID:    ccid,
		Address: addr,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true, // just for testing good enough
		},
	}

	if err := server.Start(); err != nil {
		panic(err)
	}
}
