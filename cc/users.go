package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type User struct {
	DocType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name    string `json:"name"`
	PubKey  string `json:"pubkey"`
	Balance string `json:"balance"`
}

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
	if (user != User{}) {
		return fmt.Errorf("public key already exists")
	}

	user = User{
		DocType: "user",
		Name:    name,
		PubKey:  pubkey,
		Balance: strconv.Itoa(initialBalance),
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

// ReadUser returns the User stored in the world state with given name.
func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, id string) (User, error) {
	userBytes, err := ctx.GetStub().GetState(fmt.Sprintf("user_%v", id))
	if err != nil {
		return User{}, fmt.Errorf("user with id %v could not be read from world state: %v", id, err)
	}
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return User{}, fmt.Errorf("failed to unmarshal file: %v", err)
	}
	return user, nil
}

func (s *SmartContract) QueryUsersByPublicKey(ctx contractapi.TransactionContextInterface,
	publicKey string) (User, error) {
	publicKey = strings.ReplaceAll(publicKey, "\n", "")

	username, err := ctx.GetStub().GetState(publicKey)
	if err != nil {
		return User{}, fmt.Errorf("query failed: %v", err)
	}

	user, err := s.ReadUser(ctx, string(username))
	if err != nil {
		return User{}, err
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
