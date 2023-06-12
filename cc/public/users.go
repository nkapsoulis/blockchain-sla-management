package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-private-chaincode/lib"
)

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface,
	name, pubkey string, initialBalance int) error {

	if initialBalance < 0 {
		return fmt.Errorf("initial amount must be zero or positive")
	}

	exists, err := s.UserExists(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get user info")
	}
	if exists {
		return fmt.Errorf("user already exists")
	}

	user, err := s.QueryUsersByPublicKey(ctx, pubkey)
	if err != nil {
		return fmt.Errorf("querying for public key failed: %v", err)
	}
	if (user != lib.User{}) {
		return fmt.Errorf("public key already exists")
	}

	user = lib.User{
		Name:       name,
		PubKey:     pubkey,
		Balance:    strconv.Itoa(initialBalance),
		ProviderOf: "",
		ClientOf:   "",
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("unable to marshal json: %v", err)
	}
	err = ctx.GetStub().PutState(pubkey, []byte(name))

	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(fmt.Sprintf("user_%v", name), userBytes)
}

// Returns the users balance.
func (s *SmartContract) UserBalance(ctx contractapi.TransactionContextInterface, id string) (float64, error) {
	user, err := s.ReadUser(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("could not read user: %v", err)
	}

	var currentBalance float64

	currentBalance, err = strconv.ParseFloat(string(user.Balance), 64)
	if err != nil {
		return 0, fmt.Errorf("could not convert balance: %v", err)
	}

	return currentBalance, nil
}

func (s *SmartContract) updateUserBalance(ctx contractapi.TransactionContextInterface,
	id string, newBalance float64) error {

	user, err := s.ReadUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to read user %v", err)
	}
	user.Balance = fmt.Sprintf("%f", newBalance)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user: %v", err)
	}
	return ctx.GetStub().PutState(fmt.Sprintf("user_%v", id), userBytes)
}

func (s *SmartContract) addProvidedSLA(ctx contractapi.TransactionContextInterface, userId, slaId string) error {
	user, err := s.ReadUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to read user %v", err)
	}

	if user.ProviderOf != "" {
		user.ProviderOf += ","
	}
	user.ProviderOf += slaId

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user: %v", err)
	}

	return ctx.GetStub().PutState(fmt.Sprintf("user_%v", userId), userBytes)
}

func (s *SmartContract) addConsumedSLA(ctx contractapi.TransactionContextInterface, userId, slaId string) error {
	user, err := s.ReadUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to read user %v", err)
	}

	if user.ClientOf != "" {
		user.ClientOf += ","
	}
	user.ClientOf += slaId

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user: %v", err)
	}

	return ctx.GetStub().PutState(fmt.Sprintf("user_%v", userId), userBytes)
}

func (s *SmartContract) SlaInUserContracts(ctx contractapi.TransactionContextInterface, userId, slaId string) (bool, error) {
	user, err := s.ReadUser(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("failed to read user %v", err)
	}

	clientList := strings.Split(user.ClientOf, ",")
	// Slice of size 1 means that the delimiter was not found in the string
	if len(clientList) != 1 {
		for _, sla := range clientList {
			if sla == slaId {
				return true, nil
			}
		}
	}

	providerList := strings.Split(user.ProviderOf, ",")
	// Slice of size 1 means that the delimiter was not found in the string
	if len(providerList) != 1 {
		for _, sla := range providerList {
			if sla == slaId {
				return true, nil
			}
		}
	}
	return false, nil
}

// ReadUser returns the User stored in the world state with given name.
func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, id string) (lib.User, error) {
	userBytes, err := ctx.GetStub().GetState(fmt.Sprintf("user_%v", id))
	if err != nil {
		return lib.User{}, fmt.Errorf("user with id %v could not be read from world state: %v", id, err)
	}
	if len(userBytes) == 0 || userBytes == nil {
		return lib.User{}, nil
	}

	var user lib.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return lib.User{}, fmt.Errorf("failed to unmarshal file: %v", err)
	}
	return user, nil
}

func (s *SmartContract) QueryUsersByPublicKey(ctx contractapi.TransactionContextInterface,
	publicKey string) (lib.User, error) {
	publicKey = strings.ReplaceAll(publicKey, "\n", "")

	username, err := ctx.GetStub().GetState(publicKey)
	if err != nil {
		return lib.User{}, fmt.Errorf("query failed: %v", err)
	}

	if len(username) == 0 || username == nil {
		return lib.User{}, nil
	}

	user, err := s.ReadUser(ctx, string(username))
	if err != nil {
		return lib.User{}, err
	}

	return user, nil

}

// UserExists returns true when a User with given name or public key exists in world state
func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	UserJSON, err := ctx.GetStub().GetState(fmt.Sprintf("user_%v", id))
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return UserJSON != nil, nil
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
