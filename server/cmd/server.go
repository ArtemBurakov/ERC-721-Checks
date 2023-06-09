package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"erc-721-checks/contract"
	"erc-721-checks/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/turret-io/go-menu/menu"
)

var (
	minterRoleHash [32]byte
	contractClient *ethclient.Client
	instance       *contract.Contract
	auth           *bind.TransactOpts
)

func init() {
	fmt.Println("Connecting to the smart contract...")
	var err error
	contractClient, err = ethclient.Dial(utils.EnvHelper(utils.ProviderKey))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fmt.Println("Instantiating the smart contract...")
	contractAddress := common.HexToAddress(utils.EnvHelper(utils.ContractAddressKey))
	instance, err = contract.NewContract(contractAddress, contractClient)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	fmt.Printf("Setting up the transaction options...\n\n")
	privateKey, err := crypto.HexToECDSA(utils.EnvHelper(utils.SuperUserPrivateKey))
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

	copy(minterRoleHash[:], []byte("MINTER_ROLE"))
}

func grantRole(address string) error {
	tx, err := instance.GrantRole(auth, minterRoleHash, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to grant role: %v", err)
	}
	fmt.Printf("Role granted successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
}

func revokeRole(address string) error {
	tx, err := instance.RevokeRole(auth, minterRoleHash, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to revoke role: %v", err)
	}
	fmt.Printf("Role revoked successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
}

func main() {
	commandOptions := []menu.CommandOption{
		{Command: "grantRole", Description: "Grant user minter role", Function: utils.PromptAddress(grantRole)},
		{Command: "revokeRole", Description: "Revoke user minter role", Function: utils.PromptAddress(revokeRole)},
	}
	menuOptions := menu.NewMenuOptions("> ", 0)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
