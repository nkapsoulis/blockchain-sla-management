package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-private-chaincode/lib/errors"
	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
)

// MintNFT creates SLAs that don't have a client, so they are open to be bought from entities.
func (s *SmartContract) MintNFT(ctx contractapi.TransactionContextInterface, contractJSON string) error {
	var nft SLA
	err := json.Unmarshal([]byte(contractJSON), &nft)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}

	exists, err := s.UserExists(ctx, nft.Provider.Name)
	if err != nil {
		return fmt.Errorf("provider account %s could not be read: %v", nft.Provider.ID, err)
	}
	if !exists {
		return fmt.Errorf("provider does not exist")
	}

	if (nft.Client != iso19086.Entity{}) {
		return fmt.Errorf("the nft should not have a client defined")
	}

	nftCount, err := s.readCounter(ctx, nftCounterName)
	if err != nil {
		return err
	}

	nft.ID = strconv.Itoa(nftCount)

	err = s.incrementCounter(ctx, nftCounterName)
	if err != nil {
		return err
	}

	nftContractJSON, err := json.Marshal(nft)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(fmt.Sprintf("nft_%v", nft.ID), []byte(nftContractJSON))
	if err != nil {
		return err
	}

	err = s.addProvidedNFT(ctx, nft.Provider.Name, nft.ID)
	if err != nil {
		errs := make([]error, 0)
		errs = append(errs, err)
		errs = append(errs, ctx.GetStub().DelState(fmt.Sprintf("nft_%v", nft.ID)))
		return errors.ReportMultipleErrors(errs)
	}

	err = s.addToNTFList(ctx, nft.ID)
	if err != nil {
		errs := make([]error, 0)
		errs = append(errs, err)
		errs = append(errs, ctx.GetStub().DelState(fmt.Sprintf("nft_%v", nft.ID)))
		errs = append(errs, s.removeProvidedNFT(ctx, nft.Provider.Name, nft.ID))
		return errors.ReportMultipleErrors(errs)
	}
	return nil
}

func (s *SmartContract) readNFT(ctx contractapi.TransactionContextInterface, nftID string) (*SLA, error) {
	nftBytes, err := ctx.GetStub().GetState(fmt.Sprintf("nft_%v", nftID))
	if err != nil {
		return nil, err
	}

	var nft SLA
	err = json.Unmarshal(nftBytes, &nft)
	if err != nil {
		return nil, err
	}
	return &nft, nil
}

func (s *SmartContract) OwnerOf(ctx contractapi.TransactionContextInterface, nftID string) (string, error) {
	nft, err := s.readNFT(ctx, nftID)
	if err != nil {
		return "", err
	}

	return nft.Provider.Name, nil
}

func (s *SmartContract) readNTFList(ctx contractapi.TransactionContextInterface) ([]string, error) {
	initialized, err := s.checkInitialized(ctx)
	if err != nil {
		return nil, err
	}

	if !initialized {
		return nil, fmt.Errorf("chaincode is not initialized")
	}

	nftListBytes, err := ctx.GetStub().GetState(nftListName)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(nftListBytes), ","), nil
}

func (s *SmartContract) addToNTFList(ctx contractapi.TransactionContextInterface, nftId string) error {
	initialized, err := s.checkInitialized(ctx)
	if err != nil {
		return err
	}

	if !initialized {
		return fmt.Errorf("chaincode is not initialized")
	}

	nftListBytes, err := ctx.GetStub().GetState(nftListName)
	if err != nil {
		return err
	}

	nftListStr := string(nftListBytes)

	if nftListStr != "" {
		nftListStr += ","
	}
	nftListStr += nftId

	return ctx.GetStub().PutState(nftListName, []byte(nftListStr))
}
