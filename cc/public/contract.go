package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"

	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-private-chaincode/lib"
	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
	"github.com/tyler-smith/go-bip32"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type SLA struct {
	iso19086.SLA
}

type Metric struct {
	iso19086.Metrics
}

type Approval struct {
	lib.Approval
}

type cc_SLA struct {
	SLA
	Approval
	State           string  `json:"state"`
	RefundValue     float64 `json:"RefundValue"` // compensation amount
	TotalViolations []int   `json:"TotalViolations"`
	DailyValue      float64 `json:"DailyValue"`
	DailyViolations []int   `json:"DailyViolations"`
}

// InitLedger is just a template for now.
// Used to test the connection and verify that applications can connect to the chaincode.
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	initStatus, err := ctx.GetStub().GetState("initRan")
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if initStatus != nil {
		return fmt.Errorf("init has already ran")
	}

	return ctx.GetStub().PutState("initRan", []byte("true"))

}

// CreateOrUpdateContract issues a new Contract to the world state with given details.
func (s *SmartContract) CreateOrUpdateContract(ctx contractapi.TransactionContextInterface, contractJSON string) error {
	var sla SLA
	err := json.Unmarshal([]byte(contractJSON), &sla)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	exists, err := s.UserExists(ctx, sla.Provider.Name)
	if err != nil {
		return fmt.Errorf("provider account %s could not be read: %v", sla.Provider.ID, err)
	}
	if !exists {
		return fmt.Errorf("provider does not exist")
	}

	exists, err = s.UserExists(ctx, sla.Client.Name)
	if err != nil {
		return fmt.Errorf("client account %s could not be read: %v", sla.Client.ID, err)
	}
	if !exists {
		return fmt.Errorf("client does not exist")
	}

	exists, err = s.ContractExists(ctx, sla.ID)
	if err != nil {
		return err
	}

	value := rand.Float64()
	totalViolations := make([]int, 1)
	dailyViolations := make([]int, 1)
	dailyValue := 0.0
	if exists {
		contract, err := s.GetContract(ctx, sla.ID)
		if err != nil {
			return err
		}
		value = contract.RefundValue
		totalViolations = contract.TotalViolations
		dailyViolations = contract.DailyViolations
		dailyValue = contract.DailyValue
	} else {
		s.addProvidedSLA(ctx, sla.Provider.Name, sla.ID)
		s.addConsumedSLA(ctx, sla.Client.Name, sla.ID)
	}

	approval := new(Approval)
	approval.ProviderApproved = false
	approval.ConsumerApproved = false

	contract := cc_SLA{
		SLA:             sla,
		Approval:        *approval,
		RefundValue:     value,
		State:           "unapproved",
		TotalViolations: totalViolations,
		DailyViolations: dailyViolations,
		DailyValue:      dailyValue,
	}

	slaContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", contract.SLA.ID), slaContractJSON)
}

func (s *SmartContract) GetApprovals(ctx contractapi.TransactionContextInterface, slaId string) (string, error) {
	contract, err := s.GetContract(ctx, slaId)
	if err != nil {
		return "", err
	}

	approvalJSON, err := json.Marshal(contract.Approval)
	if err != nil {
		return "", err
	}
	return string(approvalJSON), nil
}

func (s *SmartContract) GetState(ctx contractapi.TransactionContextInterface, slaId string) (string, error) {
	contract, err := s.GetContract(ctx, slaId)
	if err != nil {
		return "", err
	}
	return (*contract).State, nil
}

// Approve gets the signature of the user and verifies they have signed the contract
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, slaId, userName, signatureHex string) error {
	contract, err := s.GetContract(ctx, slaId)
	if err != nil {
		return err
	}

	user, err := s.ReadUser(ctx, userName)
	if err != nil {
		return err
	}

	// If the user is not a client nor a provider then return
	if contract.SLA.Provider.Name != userName &&
		contract.SLA.Client.Name != userName {
		return fmt.Errorf("the contract does not include the provided user")
	}

	if contract.State == "active" {
		return fmt.Errorf("the contract is already in active state")
	}

	// Now we know that the user is either a client or a provider
	// so we check one of them.
	client := false
	if contract.SLA.Client.Name == userName {
		client = true
	}

	slaJSON, err := json.Marshal(contract.SLA)
	if err != nil {
		return err
	}

	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return err
	}

	pubKey, err := bip32.B58Deserialize(user.PubKey)
	if err != nil {
		return err
	}

	pubKeyParsed, err := secp256k1.ParsePubKey(pubKey.PublicKey().Key)
	if err != nil {
		return err
	}

	signature, err := ecdsa.ParseDERSignature(signatureBytes)
	if err != nil {
		return err
	}

	// Create the hash of the data
	hash := sha256.New()
	hash.Write([]byte(slaJSON))
	hashedData := hash.Sum(nil)

	if signature.Verify(hashedData, pubKeyParsed) {
		if client {
			contract.Approval.ConsumerApproved = true
		} else {
			contract.Approval.ProviderApproved = true
		}
	} else {
		return fmt.Errorf("signature could not be verified")
	}

	if contract.Approval.ConsumerApproved && contract.Approval.ProviderApproved {
		contract.State = "active"
	}

	slaContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", contract.SLA.ID), slaContractJSON)
}

// ReadSLA returns the Contract stored in the world state with given id.
func (s *SmartContract) ReadSLA(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	contract, err := s.GetContract(ctx, id)
	if err != nil {
		return "", err
	}

	SLAJson, err := json.Marshal(contract.SLA)
	if err != nil {
		return "", err
	}

	return string(SLAJson), nil
}

// GetContract returns all the SLA information stored in the world state with given id.
func (s *SmartContract) GetContract(ctx contractapi.TransactionContextInterface, id string) (*cc_SLA, error) {
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

func (s *SmartContract) RefundSLA(ctx contractapi.TransactionContextInterface, id string) error {
	contract, err := s.GetContract(ctx, id)
	if err != nil {
		return err
	}

	if contract.State != "active" {
		return fmt.Errorf("the contract %s is not active, but in %s state, no violations can happen", contract.State, id)
	}

	if !contract.ConsumerApproved || !contract.ProviderApproved {
		return fmt.Errorf("the contract %s has not been validated by the provider or the consumer", contract.ID)
	}

	err = s.transferTokens(ctx, contract.SLA.Provider.Name, contract.SLA.Client.Name, contract.DailyValue)
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

func (s *SmartContract) SLAViolated(ctx contractapi.TransactionContextInterface, id string) error {
	contract, err := s.GetContract(ctx, id)
	if err != nil {
		return err
	}

	if !contract.ConsumerApproved || !contract.ProviderApproved {
		return fmt.Errorf("the contract %s has not been validated by the provider or the consumer", id)
	}

	contract.DailyViolations[0] += 1

	ContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("could not transfer tokens from violation: %v", err)
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", contract.ID), ContractJSON)
}

func (s *SmartContract) SLAViolatedAndRefunded(ctx contractapi.TransactionContextInterface, id string) error {
	contract, err := s.GetContract(ctx, id)
	if err != nil {
		return err
	}

	if contract.State != "active" {
		return fmt.Errorf("the contract %s is not active, but in %s state, no violations can happen", contract.State, id)
	}

	if !contract.ConsumerApproved || !contract.ProviderApproved {
		return fmt.Errorf("the contract %s has not been validated by the provider or the consumer", contract.ID)
	}

	err = s.transferTokens(ctx, contract.SLA.Provider.Name, contract.SLA.Client.Name, float64(contract.RefundValue))
	if err != nil {
		return err
	}

	ContractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("contract_%v", id), ContractJSON)
}
