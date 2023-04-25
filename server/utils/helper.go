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

func PromptAddress(fn func(string) error) func(...string) error {
	return func(args ...string) error {
		reader := bufio.NewReader(os.Stdin)
		var address string
		for {
			fmt.Print("Enter address: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			input = strings.TrimSpace(input)
			if !common.IsHexAddress(input) {
				fmt.Println("Invalid address. Please enter a valid Ethereum address.")
				continue
			}
			address = input
			break
		}
		return fn(address)
	}
}

func EnvHelper(key string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
