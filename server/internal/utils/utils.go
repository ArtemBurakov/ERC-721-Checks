package utils

import (
	"bufio"
	"erc-721-checks/internal/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

const (
	DotEnvPath          = "../../.env"
	ProviderKey         = "TESTNET_PROVIDER"
	SuperUserPrivateKey = "SUPER_USER_PRIVATE_KEY"
)

func PromptAddress(fn func(string) error) func(...string) error {
	return func(args ...string) error {
		address, err := promptAddress("Enter user address: ")
		if err != nil {
			return err
		}
		return fn(address)
	}
}

func PromptContractAddress() (string, error) {
	return promptAddress("Enter contract address: ")
}

func promptAddress(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)
		if common.IsHexAddress(input) {
			return input, nil
		}
		fmt.Println("Invalid address. Please enter a valid Ethereum address.")
	}
}

func ToMinters(minters []string) []models.Minter {
	var m []models.Minter
	for _, address := range minters {
		m = append(m, models.Minter{Address: address})
	}
	return m
}

func EnvHelper(key string) string {
	err := godotenv.Load(DotEnvPath)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
