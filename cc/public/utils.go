package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) checkInitialized(ctx contractapi.TransactionContextInterface) (bool, error) {
	initStatus, err := ctx.GetStub().GetState("initRan")
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if initStatus == nil {
		return false, nil
	}
	return true, nil
}

func (s *SmartContract) readCounter(ctx contractapi.TransactionContextInterface, counterName string) (int, error) {
	initialized, err := s.checkInitialized(ctx)
	if err != nil {
		return 0, err
	}
	if !initialized {
		return 0, fmt.Errorf("need to run init first")
	}

	counterValue, err := ctx.GetStub().GetState(counterName)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}

	return strconv.Atoi(string(counterValue))
}

func (s *SmartContract) incrementCounter(ctx contractapi.TransactionContextInterface, counterName string) error {
	counter, err := s.readCounter(ctx, counterName)
	if err != nil {
		return err
	}
	counter++

	counterStr := strconv.Itoa(counter)

	return ctx.GetStub().PutState(counterName, []byte(counterStr))
}
