package main

import (
	"fmt"
	"log"
	"sync"

	"erc-721-checks/internal/contract"
	"erc-721-checks/internal/database"
	"erc-721-checks/internal/models"
	"erc-721-checks/internal/utils"

	"github.com/turret-io/go-menu/menu"
)

var (
	smartContract    *contract.SmartContract
	minterRepository *models.MinterRepository
)

func init() {
	var err error
	if err = database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize the database connection pool: %v", err)
	}
	minterRepository = models.NewMinterRepository(database.DBInstance)

	smartContract, err = contract.InitContract()
	if err != nil {
		log.Fatalf("Failed to initialize the smart contract: %v", err)
	}
}

func grantRole(address string) error {
	if err := minterRepository.CreateMinter(address, models.ActiveMinterStatus); err != nil {
		fmt.Printf("failed to add minter to the database: %v\n", err)
		return nil
	}

	if err := smartContract.GrantRole(address); err != nil {
		fmt.Printf("failed to grant role: %v\n", err)

		if err := minterRepository.UpdateMinter(address, models.ArchivedMinterStatus); err != nil {
			fmt.Printf("failed to delete minter: %v\n", err)
		}
	}

	return nil
}

func revokeRole(address string) error {
	if err := minterRepository.UpdateMinter(address, models.ArchivedMinterStatus); err != nil {
		fmt.Printf("failed to remove minter from the database: %v\n", err)
		return nil
	}

	if err := smartContract.RevokeRole(address); err != nil {
		fmt.Printf("failed to revoke role: %v\n", err)

		if err := minterRepository.UpdateMinter(address, models.ActiveMinterStatus); err != nil {
			fmt.Printf("failed to create minter: %v\n", err)
		}
	}

	return nil
}

func printMinters(args ...string) error {
	minters, err := smartContract.GetMinters()
	if err != nil {
		fmt.Printf("failed to get minters: %v\n", err)
		return nil
	}

	for _, minter := range minters {
		fmt.Println(minter)
	}

	return nil
}

func syncMinters(args ...string) error {
	minters, err := minterRepository.GetAllMinters()
	if err != nil {
		fmt.Printf("failed to fetch existing minters from database: %v\n", err)
		return nil
	}

	var wg sync.WaitGroup
	for _, minter := range minters {
		wg.Add(1)

		go func(minter models.Minter) {
			defer wg.Done()
			err := smartContract.SyncMinterRole(minter)
			if err != nil {
				fmt.Printf("failed to sync minter: %v\n", err)
			}
		}(minter)
	}
	wg.Wait()

	fmt.Print("Sync completed successfully\n")
	return nil
}

func fetchMinters(args ...string) error {
	minters, err := smartContract.GetMinters()
	if err != nil {
		fmt.Printf("failed to get minters: %v\n", err)
		return nil
	}

	var mintersArray []models.Minter
	for _, address := range minters {
		mintersArray = append(mintersArray, models.Minter{Address: address})
	}

	if err := minterRepository.InitializeMintersTable(mintersArray); err != nil {
		fmt.Printf("failed to initialize minters table: %v\n", err)
		return nil
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
