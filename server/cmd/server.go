package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type NFTData struct {
	Name        string `json:"name" xml:"name" form:"name"`
	Description string `json:"description" xml:"description" form:"description"`
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func setupRoutes(app *fiber.App) {
	app.Get("/create", func(c *fiber.Ctx) error {
		nftData := new(NFTData)

		if err := c.BodyParser(nftData); err != nil {
			return err
		}

		// client, err := ethclient.Dial("https://sepolia.infura.io/v3/98fec64d921f4387a663bb11cd26f845")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// address := common.HexToAddress("0x748E36d986c08222b69711726858e5f8737f5ae8")
		// instance, err := contract.NewContract(address, client)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// version, err := instance.Version()
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Println(version)

		method := "POST"
		token := goDotEnvVariable("PINATA_JWT_TOKEN")
		url := "https://api.pinata.cloud/pinning/pinJSONToIPFS"

		// Marshal the struct to JSON format
		payloadBytes, err := json.Marshal(nftData)
		if err != nil {
			return err
		}
		payload := strings.NewReader(string(payloadBytes))

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return c.SendString("Error!")
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		res, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			return c.SendString("Error!")
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return c.SendString("Error!")
		}
		fmt.Println(string(body))

		return c.SendString(string(body))
	})
}

func main() {
	app := fiber.New()
	setupRoutes(app)
	app.Listen(":3000")
}
