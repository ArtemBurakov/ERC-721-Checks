package main

import (
	"context"
	"erc-721-checks/contract"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/joho/godotenv"
	"github.com/turret-io/go-menu/menu"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	contractClient *ethclient.Client
	instance       *contract.Contract
	auth           *bind.TransactOpts
)

func envHelper(key string) string {
	// load .env file
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func init() {
	// Create instance of the contract client
	var err error
	contractClient, err = ethclient.Dial(envHelper("TESTNET_PROVIDER"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Create instance of the contract
	contractAddress := common.HexToAddress(envHelper("DEPLOYED_CONTRACT_ADDRESS"))
	instance, err = contract.NewContract(contractAddress, contractClient)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	// Set up the transaction options
	privateKey, err := crypto.HexToECDSA(envHelper("SUPER_USER_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}
	auth = bind.NewKeyedTransactor(privateKey)
	nonce, err := contractClient.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Fatalf("Failed to retrieve account nonce: %v", err)
	}
	gasPrice, err := contractClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to retrieve suggested gas price: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(30000)
	auth.GasPrice = gasPrice
}

func promptAddress(fn func(string) error) func(...string) error {
	return func(args ...string) error {
		var address string
		fmt.Print("Enter wallet address: ")
		fmt.Scanln(&address)
		return fn(address)
	}
}

func grantRole(address string) error {
	// Grant the MINTER_ROLE permission to the recipient address
	roleHash := [32]byte{}
	copy(roleHash[:], []byte("MINTER_ROLE"))
	tx, err := instance.GrantRole(auth, roleHash, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to grant role: %v", err)
	}
	fmt.Printf("Role granted successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
}

func revokeRole(address string) error {
	// Revoke the MINTER_ROLE permission to the recipient address
	roleHash := [32]byte{}
	copy(roleHash[:], []byte("MINTER_ROLE"))
	tx, err := instance.RevokeRole(auth, roleHash, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to revoke role: %v", err)
	}
	fmt.Printf("Role revoked successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
}

func main() {
	commandOptions := []menu.CommandOption{
		{Command: "grantRole", Description: "Grant user minter role", Function: promptAddress(grantRole)},
		{Command: "revokeRole", Description: "Revoke user minter role", Function: promptAddress(revokeRole)},
	}
	menuOptions := menu.NewMenuOptions("> ", 0)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
