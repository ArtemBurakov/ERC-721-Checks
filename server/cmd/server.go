package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"

	"erc-721-checks/contract"
	"erc-721-checks/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/go-sql-driver/mysql"
	"github.com/turret-io/go-menu/menu"
)

var (
	minterRoleHash [32]byte
	contractClient *ethclient.Client
	instance       *contract.Contract
	auth           *bind.TransactOpts
	db             *sql.DB
)

func init() {
	var err error

	fmt.Println("Connecting to the database...")
	db, err = sql.Open("mysql", "artem:3Wxc84@jK5hp6az*@/erc721-checks")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	fmt.Println("Connecting to the smart contract...")
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
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	copy(minterRoleHash[:], []byte("MINTER_ROLE"))
}

func getMinters(args ...string) error {
	opts := &bind.CallOpts{
		Context: context.Background(),
	}

	roleName := "MINTER_ROLE"
	roleHash := crypto.Keccak256Hash([]byte(roleName))

	minterCount, err := instance.GetRoleMemberCount(opts, roleHash)
	if err != nil {
		log.Fatalf("Failed to get minter count: %v", err)
	}
	fmt.Printf("Minters count: %s\n", minterCount)

	var minters []common.Address
	for i := uint64(0); i < minterCount.Uint64(); i++ {
		minter, err := instance.GetRoleMember(opts, roleHash, big.NewInt(int64(i)))
		if err != nil {
			log.Fatalf("Failed to get minter at index %d: %v", i, err)
		}
		minters = append(minters, minter)
	}

	for _, minter := range minters {
		fmt.Printf("Minter: %s\n", minter.Hex())
	}

	return nil
}

func grantRole(address string) error {
	fmt.Println("Adding minter to the database...")
	query := "INSERT INTO `minters` (`wallet_address`) VALUES (?)"
	insertResult, err := db.ExecContext(context.Background(), query, address)
	if err != nil {
		log.Fatalf("Error while insert minter: %s", err)
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("Error while retrieving last inserted id: %s", err)
	}
	log.Printf("Inserted minter id: %d", id)

	tx, err := instance.GrantRole(auth, minterRoleHash, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to grant role: %v", err)
	}
	fmt.Printf("Role granted successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
}

func revokeRole(address string) error {
	fmt.Println("Deleting minter from the database...")
	query := "DELETE FROM `minters` WHERE `wallet_address` = ?"
	_, err := db.ExecContext(context.Background(), query, address)
	if err != nil {
		return fmt.Errorf("error while deleting minter: %s", err)
	}
	fmt.Println("Minter removed from the database")

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
		{Command: "getMinters", Description: "Get all users with minter role", Function: getMinters},
	}
	menuOptions := menu.NewMenuOptions("> ", 0)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
