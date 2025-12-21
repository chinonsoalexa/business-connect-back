package webhooks

import (
	Dataa "business-connect/models"
	"time"
)

// working with webhook
type WebhookData struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}

type Data struct {
	ID              int                   `json:"id"`
	Domain          string                `json:"domain"`
	Status          string                `json:"status"`
	Reference       string                `json:"reference"`
	Amount          int                   `json:"amount"`
	Message         interface{}           `json:"message"` // Null values can be handled with interface{}
	GatewayResponse string                `json:"gateway_response"`
	PaidAt          time.Time             `json:"paid_at"`
	CreatedAt       time.Time             `json:"created_at"`
	Channel         string                `json:"channel"`
	Currency        string                `json:"currency"`
	IPAddress       string                `json:"ip_address"`
	Metadata        Dataa.ServiceMetaData `json:"metadata"` // Adjust according to the actual type
	Log             Log                   `json:"log"`
	Fees            interface{}           `json:"fees"` // Null values can be handled with interface{}
	Customer        Customer              `json:"customer"`
	Authorization   Authorization         `json:"authorization"`
	Plan            interface{}           `json:"plan"` // Adjust according to the actual type
}

type Log struct {
	TimeSpent      int         `json:"time_spent"`
	Attempts       int         `json:"attempts"`
	Authentication string      `json:"authentication"`
	Errors         int         `json:"errors"`
	Success        bool        `json:"success"`
	Mobile         bool        `json:"mobile"`
	Input          []string    `json:"input"`   // Adjust according to the actual type
	Channel        interface{} `json:"channel"` // Null values can be handled with interface{}
	History        []History   `json:"history"`
}

type History struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    int    `json:"time"`
}

type Customer struct {
	ID           int         `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Email        string      `json:"email"`
	CustomerCode string      `json:"customer_code"`
	Phone        interface{} `json:"phone"`    // Null values can be handled with interface{}
	Metadata     interface{} `json:"metadata"` // Adjust according to the actual type
	RiskAction   string      `json:"risk_action"`
}

type Authorization struct {
	AuthorizationCode string `json:"authorization_code"`
	Bin               string `json:"bin"`
	Last4             string `json:"last4"`
	ExpMonth          string `json:"exp_month"`
	ExpYear           string `json:"exp_year"`
	CardType          string `json:"card_type"`
	Bank              string `json:"bank"`
	CountryCode       string `json:"country_code"`
	Brand             string `json:"brand"`
	AccountName       string `json:"account_name"`
}

// Root structure for the webhook response
type PaystackWebhook struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    TFData `json:"data"`
}

// Data structure for the main transfer details
type TFData struct {
	Amount          int           `json:"amount"`
	CreatedAt       time.Time     `json:"createdAt"`
	Currency        string        `json:"currency"`
	Domain          string        `json:"domain"`
	Failures        interface{}   `json:"failures"` // Handle null with interface{}
	ID              int           `json:"id"`
	Integration     TFIntegration `json:"integration"`
	Reason          string        `json:"reason"`
	Reference       string        `json:"reference"`
	Source          string        `json:"source"`
	SourceDetails   interface{}   `json:"source_details"` // Handle null with interface{}
	Status          string        `json:"status"`
	TitanCode       interface{}   `json:"titan_code"` // Handle null with interface{}
	TransferCode    string        `json:"transfer_code"`
	TransferredAt   *time.Time    `json:"transferred_at"` // Nullable field, use a pointer
	UpdatedAt       time.Time     `json:"updatedAt"`
	Recipient       TFRecipient   `json:"recipient"`
	Session         TFSession     `json:"session"`
	FeeCharged      int           `json:"fee_charged"`
	GatewayResponse interface{}   `json:"gateway_response"` // Handle null with interface{}
}

// Integration structure for integration details
type TFIntegration struct {
	ID           int    `json:"id"`
	IsLive       bool   `json:"is_live"`
	BusinessName string `json:"business_name"`
	LogoPath     string `json:"logo_path"`
}

// Recipient structure for recipient details
type TFRecipient struct {
	Active        bool                  `json:"active"`
	CreatedAt     time.Time             `json:"createdAt"`
	Currency      string                `json:"currency"`
	Description   string                `json:"description"` // Can be empty string
	Domain        string                `json:"domain"`
	Email         *string               `json:"email"` // Nullable field, use a pointer
	ID            int                   `json:"id"`
	Integration   int                   `json:"integration"`
	Metadata      Dataa.ServiceMetaData `json:"metadata"` // Adjust according to the actual type
	Name          string                `json:"name"`
	RecipientCode string                `json:"recipient_code"`
	Type          string                `json:"type"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	IsDeleted     bool                  `json:"is_deleted"`
	Details       TFDetails             `json:"details"`
}

// Details structure for bank account details
type TFDetails struct {
	AuthorizationCode interface{} `json:"authorization_code"` // Handle null with interface{}
	AccountNumber     string      `json:"account_number"`
	AccountName       string      `json:"account_name"`
	BankCode          string      `json:"bank_code"`
	BankName          string      `json:"bank_name"`
}

// Session structure for session details
type TFSession struct {
	Provider *string `json:"provider"` // Nullable field, use a pointer
	ID       *string `json:"id"`       // Nullable field, use a pointer
}

type ChargeSuccessWebhook struct {
	Event string `json:"event"`
	Data  struct {
		ID              int         `json:"id"`
		Domain          string      `json:"domain"`
		Status          string      `json:"status"`
		Reference       string      `json:"reference"`
		Amount          int64       `json:"amount"`
		Message         string      `json:"message"`
		GatewayResponse string      `json:"gateway_response"`
		PaidAt          string      `json:"paid_at"`
		CreatedAt       string      `json:"created_at"`
		Channel         string      `json:"channel"`
		Currency        string      `json:"currency"`
		IPAddress       string      `json:"ip_address"`
		Metadata        interface{} `json:"metadata"` // Use a map if you expect specific metadata structure
		Log             struct {
			TimeSpent      int           `json:"time_spent"`
			Attempts       int           `json:"attempts"`
			Authentication string        `json:"authentication"`
			Errors         int           `json:"errors"`
			Success        bool          `json:"success"`
			Mobile         bool          `json:"mobile"`
			Input          []interface{} `json:"input"`
			Channel        interface{}   `json:"channel"`
			History        []struct {
				Type    string `json:"type"`
				Message string `json:"message"`
				Time    int    `json:"time"`
			} `json:"history"`
		} `json:"log"`
		Fees     interface{} `json:"fees"` // Use float64 or int64 if you expect numeric fees
		Customer struct {
			ID           int         `json:"id"`
			FirstName    string      `json:"first_name"`
			LastName     string      `json:"last_name"`
			Email        string      `json:"email"`
			CustomerCode string      `json:"customer_code"`
			Phone        string      `json:"phone"`
			Metadata     interface{} `json:"metadata"` // Use map if needed
			RiskAction   string      `json:"risk_action"`
		} `json:"customer"`
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			Bin               string `json:"bin"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			CardType          string `json:"card_type"`
			Bank              string `json:"bank"`
			CountryCode       string `json:"country_code"`
			Brand             string `json:"brand"`
			AccountName       string `json:"account_name"`
		} `json:"authorization"`
		Plan interface{} `json:"plan"` // Define specific struct if you expect detailed plan data
	} `json:"data"`
}
