package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/turret-io/go-menu/menu"
)

func envHelper(key string) string {
	// load .env file
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func promptAddress(fn func(string) error) func(...string) error {
	return func(args ...string) error {
		var address string
		fmt.Print("Enter wallet address: ")
		fmt.Scanln(&address)
		return fn(address)
	}
}

func listen(address string) error {
	client, err := ethclient.Dial(envHelper("TESTNET_PROVIDER"))
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening for the smart contract events started...\n")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog)
		}
	}
}

func main() {
	commandOptions := []menu.CommandOption{
		{Command: "listen", Description: "Start listening to the smart contract events", Function: promptAddress(listen)},
	}
	menuOptions := menu.NewMenuOptions("> ", 0)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
