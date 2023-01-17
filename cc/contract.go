package main

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-private-chaincode/structs"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type cc_SLA struct {
	structs.SLA
	ProviderApproved bool    `json:"providerApproved"`
	ConsumerApproved bool    `json:"consumerApproved"`
	RefundValue      int     `json:"RefundValue"` // compensation amount
	TotalViolations  []int   `json:"TotalViolations"`
	DailyValue       float64 `json:"DailyValue"`
	DailyViolations  []int   `json:"DailyViolations"`
}

// InitLedger is just a template for now.
// Used to test the connection and verify that applications can connect to the chaincode.
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	return nil
}

func (s *SmartContract) transferTokens(ctx contractapi.TransactionContextInterface,
	from, to string, amount float64) error {
	if from == to {
		return fmt.Errorf("cannot transfer from and to the same account")
	}

	fromBalance, err := s.UserBalance(ctx, from)
	if err != nil {
		return fmt.Errorf("could not get balance of transferer during token transfer: %v", err)
	}
	if fromBalance < amount {
		return fmt.Errorf("transferer does not have enough tokens to complete transfer")
	}

	toBalance, err := s.UserBalance(ctx, to)
	if err != nil {
		return fmt.Errorf("could not get balance of transferee during token transfer: %v", err)
	}

	updatedFromBalance := fromBalance - amount
	updatedToBalance := toBalance + amount

	err = s.updateUserBalance(ctx, from, updatedFromBalance)
	if err != nil {
		return fmt.Errorf("could not update sender's balance: %v", err)
	}

	err = s.updateUserBalance(ctx, to, updatedToBalance)
	if err != nil {
		return fmt.Errorf("could not update receiver's balance: %v", err)
	}
	return nil
}

// Approve gets the signature of the user and verifies they have signed the contract
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, slaId, userName string, signature []byte) error {
	contract, err := s.ReadContract(ctx, slaId)
	if err != nil {
		return err
	}

	user, err := s.ReadUser(ctx, userName)
	if err != nil {
		return err
	}

	// If the user is not a client nor a provider then return
	if contract.SLA.Details.Provider.Name != userName &&
		contract.SLA.Details.Client.Name != userName {
		return fmt.Errorf("the contract does not include the provided user")
	}

	// Now we know that the user is either a client or a provider
	// so we check one of them.
	client := false
	if contract.SLA.Details.Client.Name == userName {
		client = true
	}

	slaJSON, err := json.Marshal(contract.SLA)
	if err != nil {
		return err
	}

	if ed25519.Verify(ed25519.PublicKey(user.PubKey), slaJSON, signature) {
		if client {
			contract.ConsumerApproved = true
		} else {
			contract.ProviderApproved = true
		}
	}

	slaContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", contract.SLA.ID), slaContractJSON)
}

// CreateOrUpdateContract issues a new Contract to the world state with given details.
func (s *SmartContract) CreateOrUpdateContract(ctx contractapi.TransactionContextInterface, contractJSON string) error {
	var sla structs.SLA
	err := json.Unmarshal([]byte(contractJSON), &sla)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	exists, err := s.UserExists(ctx, sla.Details.Provider.Name)
	if err != nil {
		return fmt.Errorf("provider account %s could not be read: %v", sla.Details.Provider.ID, err)
	}
	if !exists {
		return fmt.Errorf("provider does not exist")
	}

	exists, err = s.UserExists(ctx, sla.Details.Client.Name)
	if err != nil {
		return fmt.Errorf("client account %s could not be read: %v", sla.Details.Client.ID, err)
	}
	if !exists {
		return fmt.Errorf("client does not exist")
	}

	exists, err = s.ContractExists(ctx, sla.ID)
	if err != nil {
		return err
	}

	value := rand.Intn(20) + 10
	totalViolations := make([]int, 1)
	dailyViolations := make([]int, 1)
	dailyValue := 0.0
	if exists {
		contract, err := s.ReadContract(ctx, sla.ID)
		if err != nil {
			return err
		}
		value = contract.RefundValue
		totalViolations = contract.TotalViolations
		dailyViolations = contract.DailyViolations
		dailyValue = contract.DailyValue
	}

	contract := cc_SLA{
		SLA:              sla,
		RefundValue:      value,
		TotalViolations:  totalViolations,
		DailyViolations:  dailyViolations,
		DailyValue:       dailyValue,
		ProviderApproved: false,
		ConsumerApproved: false,
	}

	slaContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", contract.SLA.ID), slaContractJSON)
}

// ReadContract returns the Contract stored in the world state with given id.
func (s *SmartContract) ReadContract(ctx contractapi.TransactionContextInterface, id string) (*cc_SLA, error) {
	ContractJSON, err := ctx.GetStub().GetState(fmt.Sprintf("contract_%v", id))
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if ContractJSON == nil {
		return nil, fmt.Errorf("the Contract %s does not exist", id)
	}
	var contract cc_SLA
	err = json.Unmarshal(ContractJSON, &contract)
	if err != nil {
		return nil, err
	}

	return &contract, nil
}

// DeleteContract deletes an given Contract from the world state.
func (s *SmartContract) DeleteContract(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.ContractExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the Contract %s does not exist", id)
	}

	return ctx.GetStub().DelState(fmt.Sprintf("contract_%v", id))
}

// ContractExists returns true when Contract with given ID exists in world state
func (s *SmartContract) ContractExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	ContractJSON, err := ctx.GetStub().GetState(fmt.Sprintf("contract_%v", id))
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return ContractJSON != nil, nil
}

// SLAViolated changes the number of violations that have happened.
func (s *SmartContract) SLAViolated(ctx contractapi.TransactionContextInterface, violation string) error {
	var vio structs.Violation
	err := json.Unmarshal([]byte(violation), &vio)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	contract, err := s.ReadContract(ctx, vio.SLAID)
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

func (s *SmartContract) RefundSLA(ctx contractapi.TransactionContextInterface, id string) error {
	contract, err := s.ReadContract(ctx, id)
	if err != nil {
		return err
	}

	if contract.SLA.State == "stopped" {
		return fmt.Errorf("the contract %s is completed, no violations can happen", id)
	}

	if !contract.ConsumerApproved || !contract.ProviderApproved {
		return fmt.Errorf("the contract %s has not been validated by the provider or the consumer", contract.ID)
	}

	err = s.transferTokens(ctx, contract.SLA.Details.Provider.Name, contract.SLA.Details.Client.Name, contract.DailyValue)
	if err != nil {
		return err
	}

	for i := 0; i < len(contract.DailyViolations); i++ {
		contract.TotalViolations[i] += contract.DailyViolations[i]
		contract.DailyViolations[i] = 0
	}
	contract.DailyValue = 0.0

	ContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", id), ContractJSON)
}

func (s *SmartContract) RefundAllSLAs(ctx contractapi.TransactionContextInterface) error {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all kv pairs in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
		}

		if strings.HasPrefix(queryResponse.Key, "user_") {
			continue
		}

		var contract cc_SLA
		err = json.Unmarshal(queryResponse.Value, &contract)
		if err != nil {
			return err
		}

		err = s.RefundSLA(ctx, contract.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
