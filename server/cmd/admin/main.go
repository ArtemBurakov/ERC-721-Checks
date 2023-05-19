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
		fmt.Println(err)
		return
	}
}

func grantRole(address string) error {
	return smartContract.GrantRole(address, minterRepository)
}

func revokeRole(address string) error {
	return smartContract.RevokeRole(address, minterRepository)
}

func printMinters(args ...string) error {
	minters, err := smartContract.GetMinters()
	if err != nil {
		return fmt.Errorf("failed to get minters: %v", err)
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
		return fmt.Errorf("failed to fetch existing minters from database: %v", err)
	}

	fmt.Println("Syncing minters with local database...")
	for _, minter := range minters {
		err = smartContract.SyncMinterRole(minter.Address)
		if err != nil {
			return fmt.Errorf("failed to sync minter: %v", err)
		}
	}

	return nil
}

func fetchMinters(args ...string) error {
	minters, err := smartContract.GetMinters()
	if err != nil {
		log.Fatalf("failed to get minters: %v", err)
	}

	if err := minterRepository.InitializeMintersTable(utils.ToMinters(minters)); err != nil {
		log.Fatalf("failed to initialize minters table: %v", err)
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
