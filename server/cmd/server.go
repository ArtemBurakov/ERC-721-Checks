package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"erc-721-checks/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type NFTData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func createNFT(recipientAddress string, ipfsHash string) error {
	// Create a new instance of the contract client
	contractClient, err := ethclient.Dial(goDotEnvVariable("TESTNET_PROVIDER"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer contractClient.Close()

	// Create a new instance of the NFT contract
	contractAddress := common.HexToAddress(goDotEnvVariable("DEPLOYED_CONTRACT_ADDRESS"))
	instance, err := contract.NewContract(contractAddress, contractClient)
	if err != nil {
		log.Fatalf("Failed to instantiate the NFT contract: %v", err)
	}

	// Set up the transaction options
	privateKey, err := crypto.HexToECDSA(goDotEnvVariable("RECIPIENT_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
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
	auth.GasLimit = uint64(300000) // Or any other gas limit you want to set
	auth.GasPrice = gasPrice

	// Call the MintTo function to create a new NFT
	recipient := common.HexToAddress(recipientAddress)
	tokenURI := "ipfs://" + ipfsHash
	tx, err := instance.MintTo(auth, recipient, tokenURI)
	if err != nil {
		log.Fatalf("Failed to create NFT: %v", err)
	}
	fmt.Printf("NFT created successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
}

func pinJSONToIPFS(nftData *NFTData) (string, error) {
	pinataJWT := goDotEnvVariable("PINATA_JWT_TOKEN")
	pinataAPI := "https://api.pinata.cloud/pinning/pinJSONToIPFS"

	jsonBytes, err := json.Marshal(nftData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, pinataAPI, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+pinataJWT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData struct {
		IpfsHash string `json:"IpfsHash"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	return respData.IpfsHash, nil
}

func createNFTEndpoint(c *fiber.Ctx) error {
	var nftData NFTData
	if err := c.BodyParser(&nftData); err != nil {
		return err
	}

	ipfsHash, err := pinJSONToIPFS(&nftData)
	if err != nil {
		log.Println(err)
		return c.SendString("Error!")
	}

	err = createNFT(goDotEnvVariable("RECIPIENT_ADDRESS"), ipfsHash)
	if err != nil {
		log.Println(err)
		return c.SendString("Error!")
	}

	return c.SendString("Success!")
}

func setupRoutes(app *fiber.App) {
	app.Get("api/v1/create", createNFTEndpoint)
}

func main() {
	app := fiber.New()
	setupRoutes(app)
	app.Listen(":3000")
}
