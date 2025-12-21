package authentication

import (
	Data "business-connect/models"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func SendTransactionalSMS(
	payload Data.SendSMSRequest,
) (*Data.SendSMSResponse, error) {

	url := "https://api.brevo.com/v3/transactionalSMS/send"

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if os.Getenv("RENDER") == "" { // assume in local dev
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env file: %v\n", err)
		}
	}

	// Then just use os.Getenv("SMS_KEY") as usual
	apiKey := os.Getenv("SMS_KEY")

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to send SMS: status " + resp.Status)
	}

	var response Data.SendSMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
