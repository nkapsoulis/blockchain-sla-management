package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
)

// SLAViolated changes the number of violations that have happened.
func (s *SmartContract) CheckIfSLAViolated(ctx contractapi.TransactionContextInterface, metricJson string) (bool, error) {
	metr, err := iso19086.ReadMetric([]byte(metricJson))
	if err != nil {
		return false, err
	}

	contract, err := s.GetContract(ctx, metr.SLAID)
	if err != nil {
		return false, err
	}
	if contract.SLA.State == "stopped" {
		return false, fmt.Errorf("the contract %s is completed, no violations can happen", metr.SLAID)
	}

	if !contract.ConsumerApproved || !contract.ProviderApproved {
		return false, fmt.Errorf("the contract %s has not been validated by the provider or the consumer", metr.SLAID)
	}

	slaData, err := json.Marshal(contract.SLA)
	if err != nil {
		return false, err
	}

	violated, err := iso19086.Parse(slaData, []byte(metricJson))
	if err != nil {
		return false, err
	}

	if violated {

		contract.DailyViolations[0] += 1
	}
	ContractJSON, err := json.Marshal(contract)
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, fmt.Errorf("could not transfer tokens from violation: %v", err)
	}

	return violated, ctx.GetStub().PutState(fmt.Sprintf("contract_%v", contract.ID), ContractJSON)
}
