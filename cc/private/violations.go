package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// SLAViolated changes the number of violations that have happened.
func (s *SmartContract) CheckIfSLAViolated(ctx contractapi.TransactionContextInterface, metricJson, contractJson string) (bool, error) {
	metr, err := iso19086.ReadMetric([]byte(metricJson))
	if err != nil {
		return false, err
	}

	contract, err := iso19086.ReadSLA([]byte(contractJson))
	if err != nil {
		return false, err
	}

	if contract.State == "stopped" {
		return false, fmt.Errorf("the contract %s is completed, no violations can happen", metr.SLAID)
	}

	violated, err := iso19086.Parse([]byte(contractJson), []byte(metricJson))
	if err != nil {
		return false, err
	}
	return violated, nil
}
