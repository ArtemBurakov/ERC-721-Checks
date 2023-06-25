package main

import (
	"fmt"
	"log"

	"erc-721-checks/internal/contract"
	"erc-721-checks/internal/database"
	"erc-721-checks/internal/models"
	"erc-721-checks/internal/utils"

	"github.com/turret-io/go-menu/menu"
)

var (
	smartContract    *contract.SmartContract
	minterRepository models.MinterRepository
)

func init() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize the database connection pool: %v", err)
	}
	minterRepository = models.NewMinterRepository(database.GetDB())

	var err error
	smartContract, err = contract.InitContract()
	if err != nil {
		log.Fatalf("Failed to initialize the smart contract: %v", err)
	}
}

func grantRole(address string) error {
	if err := minterRepository.CreateMinter(address); err != nil {
		fmt.Printf("failed to add minter to the database: %v\n", err)
	}

	if err := smartContract.GrantRole(address); err != nil {
		_ = minterRepository.DeleteMinter(address)
		fmt.Printf("failed to grant role: %v\n", err)
	}

	return nil
}

func revokeRole(address string) error {
	if err := minterRepository.DeleteMinter(address); err != nil {
		fmt.Printf("failed to remove minter from the database: %v\n", err)
	}

	if err := smartContract.RevokeRole(address); err != nil {
		_ = minterRepository.CreateMinter(address)
		fmt.Printf("failed to revoke role: %v\n", err)
	}

	return nil
}

func printMinters(args ...string) error {
	minters, err := smartContract.GetMinters()
	if err != nil {
		fmt.Printf("failed to get minters: %v\n", err)
	}

	for _, minter := range minters {
		fmt.Println(minter)
	}

	return nil
}

func syncMinters(args ...string) error {
	fmt.Println("Fetching minters from the database...")
	minters, err := minterRepository.GetAllMinters()
	if err != nil {
		fmt.Printf("failed to fetch existing minters from database: %v\n", err)
	}

	fmt.Println("Syncing minters with local database...")
	for _, minter := range minters {
		err = smartContract.SyncMinterRole(minter.Address)
		if err != nil {
			fmt.Printf("failed to sync minter: %v\n", err)
		}
	}

	return nil
}

func fetchMinters(args ...string) error {
	minters, err := smartContract.GetMinters()
	if err != nil {
		fmt.Printf("failed to get minters: %v\n", err)
	}

	if err := minterRepository.InitializeMintersTable(utils.ToMinters(minters)); err != nil {
		fmt.Printf("failed to initialize minters table: %v\n", err)
	}

	fmt.Println("Minters inserted to the database")

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
