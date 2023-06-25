package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"erc-721-checks/internal/checks"
	"erc-721-checks/internal/contract"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/turret-io/go-menu/menu"
)

type Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

var (
	contractAbi   abi.ABI
	smartContract *contract.SmartContract
)

func processEvents(logs chan types.Log) {
	for vLog := range logs {
		go func(vLog types.Log) {
			fmt.Println("Log Name: Transfer")
			fmt.Printf("Transaction hash: %s\n", vLog.TxHash.Hex())
			fmt.Printf("Block Number: %d\n", vLog.BlockNumber)
			fmt.Printf("Block Hash: %s\n", vLog.BlockHash.Hex())

			var transferEvent Transfer
			if err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data); err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			tokenId := new(big.Int).SetBytes(vLog.Topics[3][:])
			transferEvent.TokenId = tokenId

			fmt.Printf("Sender Address: %s\n", transferEvent.From.Hex())
			fmt.Printf("Recipient Address: %s\n", transferEvent.To.Hex())
			fmt.Printf("Token ID: %s\n\n", transferEvent.TokenId.String())
		}(vLog)
	}
}

func listen(args ...string) error {
	var err error
	smartContract, err = contract.InitContract()
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(smartContract.ContractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))}},
	}

	contractAbi, err = abi.JSON(strings.NewReader(string(checks.ChecksABI)))
	if err != nil {
		log.Fatal(err)
	}

	logs := make(chan types.Log)
	sub, err := smartContract.ContractClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening to the smart contract events. Waiting for new events...\n\n")
	go processEvents(logs)
	for err := range sub.Err() {
		log.Fatal(err)
	}

	return nil
}

func main() {
	commandOptions := []menu.CommandOption{
		{Command: "listen", Description: "Start listening to the smart contract events", Function: listen},
	}
	menuOptions := menu.NewMenuOptions("> ", 0)
	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
