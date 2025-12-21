package fundAccount

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func VerifyAccountNumberPayuee(c *fiber.Ctx) error {

	accountNumber := c.Params("accountNumber")
	bankCode := c.Params("bankCode")

	paystackResponse, err := verifyAccountNumberPaystack(accountNumber, bankCode)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Error verifying account number",
		})
	}

	return c.JSON(paystackResponse)
}

func verifyAccountNumberPaystack(accountNumber, bankCode string) (fiber.Map, error) {

	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Printf("Failed to load .env file: %v\n", envErr)
	}

	SECRET_KEY := os.Getenv("PAYSTACK_LIVE_SECRET_KEY")

	url := fmt.Sprintf("https://api.paystack.co/bank/resolve?account_number=%s&bank_code=%s", accountNumber, bankCode)
	authorization := "Bearer " + SECRET_KEY

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response body as JSON
	var paystackResponse fiber.Map
	err = json.NewDecoder(resp.Body).Decode(&paystackResponse)
	if err != nil {
		return nil, err
	}

	return paystackResponse, nil
}
