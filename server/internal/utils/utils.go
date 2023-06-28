package utils

import (
	"bufio"
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
	DBHost              = "DATABASE_HOST"
	DBPort              = "DATABASE_PORT"
	DBName              = "DATABASE_NAME"
	DBUser              = "DATABASE_USER"
	DBPassword          = "DATABASE_USER_PASSWORD"
)

func PromptAddress(fn func(string) error) func(...string) error {
	return func(args ...string) error {
		address, err := handleAddressPrompt("Enter user address: ")
		if err != nil {
			return err
		}
		return fn(address)
	}
}

func PromptContractAddress() (string, error) {
	return handleAddressPrompt("Enter contract address: ")
}

func handleAddressPrompt(prompt string) (string, error) {
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

func EnvHelper(key string) string {
	if err := godotenv.Load(DotEnvPath); err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
