package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// SLAViolated changes the number of violations that have happened.
func (s *SmartContract) CheckIfSLAViolated(ctx contractapi.TransactionContextInterface, metricJson, contractJson string) (bool, error) {
	violated, err := iso19086.Parse([]byte(contractJson), []byte(metricJson))
	if err != nil {
		return false, err
	}
	return violated, nil
}
