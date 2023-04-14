package main

import (
	"bytes"
	"context"
	"encoding/json"
	"erc-721-checks/contract"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
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

func main() {
	app := fiber.New()
	setupRoutes(app)
	app.Listen(":3000")
}

func setupRoutes(app *fiber.App) {
	// Public routes
	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	// Restricted Routes
	app.Get("/whoAmI", whoAmI)
	app.Post("/grandRole", grandRole)
	app.Post("/revokeRole", revokeRole)
	app.Post("/create", createNFTEndpoint)
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func whoAmI(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func grandRole(c *fiber.Ctx) error {
	userAddress := c.FormValue("userAddress")

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
	privateKey, err := crypto.HexToECDSA(goDotEnvVariable("SUPER_USER_PRIVATE_KEY"))
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

	// Grant the MINTER_ROLE permission to the recipient address
	roleHash := [32]byte{}
	copy(roleHash[:], []byte("MINTER_ROLE"))
	tx, err := instance.GrantRole(auth, roleHash, common.HexToAddress(userAddress))
	if err != nil {
		log.Fatalf("Failed to grant role: %v", err)
	}
	fmt.Printf("Role granted successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
}

func revokeRole(c *fiber.Ctx) error {
	userAddress := c.FormValue("userAddress")

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
	privateKey, err := crypto.HexToECDSA(goDotEnvVariable("SUPER_USER_PRIVATE_KEY"))
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

	// Revoke the MINTER_ROLE permission to the recipient address
	roleHash := [32]byte{}
	copy(roleHash[:], []byte("MINTER_ROLE"))
	tx, err := instance.RevokeRole(auth, roleHash, common.HexToAddress(userAddress))
	if err != nil {
		log.Fatalf("Failed to revoke role: %v", err)
	}
	fmt.Printf("Role revoked successfully. Transaction hash: %v\n", tx.Hash().Hex())

	return nil
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
	privateKey, err := crypto.HexToECDSA(goDotEnvVariable("SUPER_USER_PRIVATE_KEY"))
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

	// // Call the MintTo function to create a new NFT
	recipient := common.HexToAddress(recipientAddress)
	tokenURI := "ipfs://" + ipfsHash
	tx, err := instance.SafeMint(auth, recipient, tokenURI)
	if err != nil {
		log.Fatalf("Failed to create NFT: %v", err)
	}
	fmt.Printf("NFT created successfully. Transaction hash: %v\n", tx.Hash().Hex())

	// Grant the MINTER_ROLE permission to the recipient address
	// roleHash := [32]byte{}
	// copy(roleHash[:], []byte("MINTER_ROLE"))
	// tx, err := instance.GrantRole(auth, roleHash, recipient)
	// if err != nil {
	// 	log.Fatalf("Failed to grant role: %v", err)
	// }
	// fmt.Printf("Role granted successfully. Transaction hash: %v\n", tx.Hash().Hex())

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
