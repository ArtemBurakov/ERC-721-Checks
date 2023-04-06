package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		method := "POST"
		token := goDotEnvVariable("PINATA_JWT_TOKEN")
		url := "https://api.pinata.cloud/pinning/pinJSONToIPFS"

		payload := strings.NewReader(`{"message": "Test message"}`)

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

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return c.SendString("Error!")
		}
		fmt.Println(string(body))

		return c.SendString(string(body))
	})

	app.Listen(":3000")
}
