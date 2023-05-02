package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-private-chaincode/lib"
)

// SLAViolated changes the number of violations that have happened.
func (s *SmartContract) SLAViolated(ctx contractapi.TransactionContextInterface, violation string) error {
	var vio lib.Violation
	err := json.Unmarshal([]byte(violation), &vio)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	contract, err := s.GetContract(ctx, vio.SLAID)
	if err != nil {
		return err
	}
	if contract.SLA.State == "stopped" {
		return fmt.Errorf("the contract %s is completed, no violations can happen", vio.SLAID)
	}

	if !contract.ConsumerApproved || !contract.ProviderApproved {
		return fmt.Errorf("the contract %s has not been validated by the provider or the consumer", vio.SLAID)
	}

	switch vio.GuaranteeID {
	case "40":
		// This should happen only the first time the SLA is violated, but it's the time
		// we actually have information about the violation itself.
		if len(contract.DailyViolations) < 3 {
			contract.DailyViolations = make([]int, 3)
			contract.TotalViolations = make([]int, 3)
		}
		switch vio.ImportanceName {
		case "Warning":
			contract.DailyValue += (1 - 0.985) * float64(contract.RefundValue)
			contract.DailyViolations[0] += 1
		case "Serious":
			contract.DailyValue += (1 - 0.965) * float64(contract.RefundValue)
			contract.DailyViolations[1] += 1
		case "Catastrophic":
			contract.DailyValue += (1 - 0.945) * float64(contract.RefundValue)
			contract.DailyViolations[2] += 1
		}
	// If we don't know the type of guarantee
	default:
		contract.DailyViolations[0] += 1
	}
	ContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("could not transfer tokens from violation: %v", err)
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", vio.SLAID), ContractJSON)
}
