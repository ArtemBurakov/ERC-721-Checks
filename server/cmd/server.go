package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"erc-721-checks/contract"
	"erc-721-checks/db"
	"erc-721-checks/models"
	"erc-721-checks/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/turret-io/go-menu/menu"
)

var (
	auth             *bind.TransactOpts
	instance         *contract.Contract
	contractClient   *ethclient.Client
	minterRepository models.MinterRepository
	minterRoleHash   = crypto.Keccak256Hash([]byte("MINTER_ROLE"))
)

func init() {
	var err error
	fmt.Println("Initializing database connection...")
	if err = db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize the database connection pool: %v", err)
	}
	minterRepository = models.NewMinterRepository(db.GetDB())

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
}

func grantRole(address string) error {
	tx, err := instance.GrantRole(auth, minterRoleHash, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to grant role: %v", err)
	}
	fmt.Printf("Role granted successfully. Transaction hash: %v\n", tx.Hash().Hex())

	err = minterRepository.CreateMinter(address)
	if err != nil {
		log.Fatalf("Failed to add minter to the database: %v", err)
	}

	return nil
}

func revokeRole(address string) error {
	tx, err := instance.RevokeRole(auth, minterRoleHash, common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to revoke role: %v", err)
	}
	fmt.Printf("Role revoked successfully. Transaction hash: %v\n", tx.Hash().Hex())

	err = minterRepository.DeleteMinter(address)
	if err != nil {
		log.Fatalf("Failed to remove minter from the database: %v", err)
	}

	return nil
}

func printMinters(args ...string) error {
	minters, err := getMinters()
	if err != nil {
		return fmt.Errorf("failed to get minters: %v", err)
	}

	for _, minter := range minters {
		fmt.Println(minter)
	}

	return nil
}

func getMinters() ([]string, error) {
	opts := &bind.CallOpts{
		Context: context.Background(),
	}

	minterCount, err := instance.GetRoleMemberCount(opts, minterRoleHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get minter count: %v", err)
	}
	fmt.Printf("Minters count: %s\n", minterCount)

	var minters []string
	for i := uint64(0); i < minterCount.Uint64(); i++ {
		minter, err := instance.GetRoleMember(opts, minterRoleHash, big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to get minter at index %d: %v", i, err)
		}
		minters = append(minters, minter.Hex())
	}

	return minters, nil
}

func fetchMinters(args ...string) error {
	minters, err := getMinters()
	if err != nil {
		log.Fatalf("failed to get minters: %v", err)
	}

	if err := minterRepository.InitializeMintersTable(utils.ToMinters(minters)); err != nil {
		log.Fatalf("failed to initialize minters table: %v", err)
	}
	fmt.Println("Minters inserted to the database")

	return nil
}

func syncMinters(args ...string) error {
	fmt.Println("Fetching minters from the database...")
	minters, err := minterRepository.GetAllMinters()
	if err != nil {
		return fmt.Errorf("failed to fetch existing minters from database: %v", err)
	}

	fmt.Println("Syncing minters with local database...")
	for _, m := range minters {
		minter := common.HexToAddress(m.Address)
		hasRole, err := instance.HasRole(nil, minterRoleHash, minter)
		if err != nil {
			return fmt.Errorf("failed to check if minter has role: %v", err)
		}

		if !hasRole {
			tx, err := instance.GrantRole(auth, minterRoleHash, minter)
			if err != nil {
				return fmt.Errorf("failed to grant role to minter: %v", err)
			}

			_, err = bind.WaitMined(context.Background(), contractClient, tx)
			if err != nil {
				return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
			}

			fmt.Printf("Granted access to minter: %s\n", minter.Hex())
		}
	}

	return nil
}

func main() {
	commandOptions := []menu.CommandOption{
		{Command: "grantRole", Description: "Grant user minter role", Function: utils.PromptAddress(grantRole)},
		{Command: "revokeRole", Description: "Revoke user minter role", Function: utils.PromptAddress(revokeRole)},
		{Command: "printMinters", Description: "Get all users with minter role", Function: printMinters},
		{Command: "syncMinters", Description: "Sync local minters with contract", Function: syncMinters},
		{Command: "fetchMinters", Description: "Save all users with minter role to local db", Function: fetchMinters},
	}
	menuOptions := menu.NewMenuOptions("> ", 0)

	menu := menu.NewMenu(commandOptions, menuOptions)
	menu.Start()
}
