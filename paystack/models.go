package paystack

import (
	Dataa "business-connect/models"
	"time"
)

var (
	// all url to make request to each transaction endpoints
	InitializeTransactions  = "https://api.paystack.co/transaction/initialize"
	VerifyTransactions      = "https://api.paystack.co/transfer/verify/%s"
	FetchTransactions       = "https://api.paystack.co/transfer/%s"
	CreateTransferRecipient = "https://api.paystack.co/transferrecipient"
	InitiateTransfer        = "https://api.paystack.co/transfer"
	VerifyTransfer          = "https://api.paystack.co/transfer/verify/"
)

// initialize transaction body
type (
	TransactionRequestBody struct {
		Amount            string   `json:"amount"`
		Email             string   `json:"email"`
		Currency          string   `json:"currency,omitempty"`
		Reference         string   `json:"reference"`
		CallbackURL       string   `json:"callback_url,omitempty"`
		Plan              string   `json:"plan,omitempty"`
		InvoiceLimit      int      `json:"invoice_limit,omitempty"`
		Metadata          string   `json:"metadata,omitempty"`
		Channels          []string `json:"channels,omitempty"`
		SplitCode         string   `json:"split_code,omitempty"`
		Subaccount        string   `json:"subaccount,omitempty"`
		TransactionCharge int      `json:"transaction_charge,omitempty"`
		Bearer            string   `json:"bearer,omitempty"`
	}

	ChargeRequestBody struct {
		Amount            string   `json:"amount"`
		Email             string   `json:"email"`
		AuthorizationCode string   `json:"authorization_code"`
		Reference         string   `json:"reference"`
		Currency          string   `json:"currency,omitempty"`
		Metadata          string   `json:"metadata,omitempty"`
		Channels          []string `json:"channels,omitempty"`
		Subaccount        string   `json:"subaccount,omitempty"`
		TransactionCharge int      `json:"transaction_charge,omitempty"`
		Bearer            string   `json:"bearer,omitempty"`
		CustomFields      []struct {
			DisplayName  string `json:"display_name"`
			VariableName string `json:"variable_name"`
			Value        string `json:"value"`
		} `json:"custom_fields,omitempty"`
		Queue bool `json:"queue,omitempty"`
	}

	DebitAuthorizationRequestBody struct {
		AuthorizationCode string `json:"authorization_code"`
		Currency          string `json:"currency"`
		Amount            string `json:"amount"`
		Email             string `json:"email"`
		Reference         string `json:"reference"`
		AtLeast           string `json:"at_least,omitempty"`
	}
)

// visual account creation structs
type (
	VirtualAccountResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Bank struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
				Slug string `json:"slug"`
			} `json:"bank"`
			AccountName   string      `json:"account_name"`
			AccountNumber string      `json:"account_number"`
			Assigned      bool        `json:"assigned"`
			Currency      string      `json:"currency"`
			Metadata      interface{} `json:"metadata"` // You may replace this with the actual type if needed
			Active        bool        `json:"active"`
			ID            int         `json:"id"`
			CreatedAt     time.Time   `json:"created_at"`
			UpdatedAt     time.Time   `json:"updated_at"`
			Assignment    struct {
				Integration  int       `json:"integration"`
				AssigneeID   int       `json:"assignee_id"`
				AssigneeType string    `json:"assignee_type"`
				Expired      bool      `json:"expired"`
				AccountType  string    `json:"account_type"`
				AssignedAt   time.Time `json:"assigned_at"`
			} `json:"assignment"`
			Customer struct {
				ID           int    `json:"id"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				Email        string `json:"email"`
				CustomerCode string `json:"customer_code"`
				Phone        string `json:"phone"`
				RiskAction   string `json:"risk_action"`
			} `json:"customer"`
		} `json:"data"`
	}
)

// Paystack Response represents the Paystack response data structure
type (
	PaystackVerificationResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}

	// Data represents the data field in the Paystack response
	Data struct {
		ID                 uint                  `json:"id"`
		Domain             string                `json:"domain"`
		Status             string                `json:"status"`
		Reference          string                `json:"reference"`
		Amount             int                   `json:"amount"`
		GatewayResponse    string                `json:"gateway_response"`
		PaidAt             time.Time             `json:"paid_at"`
		CreatedAt          time.Time             `json:"created_at"`
		Channel            string                `json:"channel"`
		Currency           string                `json:"currency"`
		IPAddress          string                `json:"ip_address"`
		Metadata           Dataa.ServiceMetaData `json:"metadata"`
		Log                Log                   `json:"log"`
		Fees               int                   `json:"fees"`
		Authorization      Authorization         `json:"authorization"`
		Customer           Customer              `json:"customer"`
		Plan               interface{}           `json:"plan"`
		Split              Split                 `json:"split"`
		OrderID            interface{}           `json:"order_id"`
		PaidAtAlt          time.Time             `json:"paidAt"`
		CreatedAtAlt       time.Time             `json:"createdAt"`
		RequestedAmount    int                   `json:"requested_amount"`
		PosTransactionData interface{}           `json:"pos_transaction_data"`
		Source             interface{}           `json:"source"`
		FeesBreakdown      interface{}           `json:"fees_breakdown"`
		TransactionDate    time.Time             `json:"transaction_date"`
		PlanObject         interface{}           `json:"plan_object"`
		Subaccount         Subaccount            `json:"subaccount"`
	}

	// Log represents the log field in the Paystack response
	Log struct {
		StartTime int       `json:"start_time"`
		TimeSpent int       `json:"time_spent"`
		Attempts  int       `json:"attempts"`
		Errors    int       `json:"errors"`
		Success   bool      `json:"success"`
		Mobile    bool      `json:"mobile"`
		Input     []int     `json:"input"`
		History   []History `json:"history"`
	}

	// History represents the history field in the Paystack response
	History struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Time    int    `json:"time"`
	}

	// Authorization represents the authorization field in the Paystack response
	Authorization struct {
		AuthorizationCode string `json:"authorization_code"`
		Bin               string `json:"bin"`
		Last4             string `json:"last4"`
		ExpMonth          string `json:"exp_month"`
		ExpYear           string `json:"exp_year"`
		Channel           string `json:"channel"`
		CardType          string `json:"card_type"`
		Bank              string `json:"bank"`
		CountryCode       string `json:"country_code"`
		Brand             string `json:"brand"`
		Reusable          bool   `json:"reusable"`
		Signature         string `json:"signature"`
		AccountName       string `json:"account_name"`
	}

	// Customer represents the customer field in the Paystack response
	Customer struct {
		ID                       int               `json:"id"`
		FirstName                string            `json:"first_name"`
		LastName                 string            `json:"last_name"`
		Email                    string            `json:"email"`
		CustomerCode             string            `json:"customer_code"`
		Phone                    string            `json:"phone"`
		Metadata                 map[string]string `json:"metadata"`
		RiskAction               string            `json:"risk_action"`
		InternationalFormatPhone string            `json:"international_format_phone"`
	}

	// Split represents the split field in the Paystack response
	Split struct {
	}

	// Subaccount represents the subaccount field in the Paystack response
	Subaccount struct {
	}
)

// Paystack Create Transfer Recipient Transaction
// Details struct represents the nested details in the JSON structure.
type (
	Details struct {
		AuthorizationCode string `json:"authorization_code"`
		AccountNumber     string `json:"account_number"`
		AccountName       string `json:"account_name"`
		BankCode          string `json:"bank_code"`
		BankName          string `json:"bank_name"`
	}

	// Data struct represents the "data" field in the JSON structure.
	TansData struct {
		Active        bool      `json:"active"`
		CreatedAt     time.Time `json:"createdAt"`
		Currency      string    `json:"currency"`
		Domain        string    `json:"domain"`
		ID            int       `json:"id"`
		Integration   int       `json:"integration"`
		Name          string    `json:"name"`
		RecipientCode string    `json:"recipient_code"`
		Type          string    `json:"type"`
		UpdatedAt     time.Time `json:"updatedAt"`
		IsDeleted     bool      `json:"is_deleted"`
		Details       Details   `json:"details"`
	}

	// Response struct represents the overall JSON structure.
	CreateTransferRecipientResponse struct {
		Status  bool     `json:"status"`
		Message string   `json:"message"`
		Data    TansData `json:"data"`
	}
)

// Paystack Initiate Transfer
type (
	TransferData struct {
		Integration  int       `json:"integration"`
		Domain       string    `json:"domain"`
		Amount       int       `json:"amount"`
		Currency     string    `json:"currency"`
		Reference    string    `json:"reference"`
		Source       string    `json:"source"`
		Reason       string    `json:"reason"`
		Recipient    int       `json:"recipient"`
		Status       string    `json:"status"`
		TransferCode string    `json:"transfer_code"`
		ID           int       `json:"id"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
	}

	// TransferResponse struct represents the overall JSON structure.
	InitiateTransferResponse struct {
		Status  bool         `json:"status"`
		Message string       `json:"message"`
		Data    TransferData `json:"data"`
	}
)

// TransferRecipientDetails represents details of the transfer recipient.
type (
	TransferRecipientDetails struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankCode      string `json:"bank_code"`
		BankName      string `json:"bank_name"`
	}

	// TransferRecipient represents the recipient information.
	TransferRecipient struct {
		Domain        string                   `json:"domain"`
		Type          string                   `json:"type"`
		Currency      string                   `json:"currency"`
		Name          string                   `json:"name"`
		Details       TransferRecipientDetails `json:"details"`
		Description   string                   `json:"description"`
		Metadata      string                   `json:"metadata"`
		RecipientCode string                   `json:"recipient_code"`
		Active        bool                     `json:"active"`
		Email         string                   `json:"email"`
		ID            int                      `json:"id"`
		Integration   int                      `json:"integration"`
		CreatedAt     time.Time                `json:"createdAt"`
		UpdatedAt     time.Time                `json:"updatedAt"`
	}

	// TransferData represents the data field in the JSON response.
	VerifyTransferData struct {
		Integration   int               `json:"integration"`
		Recipient     TransferRecipient `json:"recipient"`
		Domain        string            `json:"domain"`
		Amount        int               `json:"amount"`
		Currency      string            `json:"currency"`
		Reference     string            `json:"reference"`
		Source        string            `json:"source"`
		SourceDetails interface{}       `json:"source_details"` // Assuming it can be of any type, change it as needed.
		Reason        string            `json:"reason"`
		Status        string            `json:"status"`
		Failures      interface{}       `json:"failures"` // Assuming it can be of any type, change it as needed.
		TransferCode  string            `json:"transfer_code"`
		TitanCode     interface{}       `json:"titan_code"`     // Assuming it can be of any type, change it as needed.
		TransferredAt interface{}       `json:"transferred_at"` // Assuming it can be of any type, change it as needed.
		ID            int               `json:"id"`
		CreatedAt     time.Time         `json:"createdAt"`
		UpdatedAt     time.Time         `json:"updatedAt"`
	}

	// VerifyTransferResponse represents the overall JSON response structure.
	VerifyTransferResponse struct {
		Status  bool               `json:"status"`
		Message string             `json:"message"`
		Data    VerifyTransferData `json:"data"`
	}
)

// TRANSFER RESPOSE DATA
// RecipientDetails represents the bank account details of the transfer recipient.
type RecipientDetails struct {
	AuthCode     interface{} `json:"authorization_code"` // Can be null
	AccountNo    string      `json:"account_number"`
	AccountOwner string      `json:"account_name"`
	BankCode     string      `json:"bank_code"`
	Bank         string      `json:"bank_name"`
}

// RecipientInfo represents the details of the transfer recipient.
type RecipientInfo struct {
	IsActive     bool                  `json:"active"`
	CreatedOn    time.Time             `json:"createdAt"`
	TransferUnit string                `json:"currency"`
	DetailsDesc  string                `json:"description"`
	DomainType   string                `json:"domain"`
	UserEmail    interface{}           `json:"email"` // Can be null
	RecipientID  int                   `json:"id"`
	Integration  int                   `json:"integration"`
	Metadata     Dataa.ServiceMetaData `json:"metadata"`
	FullName     string                `json:"name"`
	Code         string                `json:"recipient_code"`
	AccountType  string                `json:"type"`
	UpdatedOn    time.Time             `json:"updatedAt"`
	DeletedFlag  bool                  `json:"is_deleted"`
	Details      RecipientDetails      `json:"details"`
}

// IntegrationDetails provides information about the integration.
type IntegrationDetails struct {
	IntegrationID int    `json:"id"`
	IsLive        bool   `json:"is_live"`
	BusinessTitle string `json:"business_name"`
	LogoURL       string `json:"logo_path"`
}

// SessionDetails contains session-related information for the transfer.
type SessionDetails struct {
	ServiceProvider string `json:"provider"`
	SessionID       string `json:"id"`
}

// TransferPayload contains the details of the transfer event.
type TransferPayload struct {
	TransferAmount   int                `json:"amount"`
	CreatedTimestamp time.Time          `json:"createdAt"`
	TransferCurrency string             `json:"currency"`
	DomainCategory   string             `json:"domain"`
	ErrorDetails     interface{}        `json:"failures"` // Can be null
	TransferID       int                `json:"id"`
	IntegrationData  IntegrationDetails `json:"integration"`
	ReasonText       string             `json:"reason"`
	TransactionRef   string             `json:"reference"`
	FundingSource    string             `json:"source"`
	ExtraSourceInfo  interface{}        `json:"source_details"` // Can be null
	TransferStatus   string             `json:"status"`
	ExtraCode        interface{}        `json:"titan_code"` // Can be null
	CodeTransfer     string             `json:"transfer_code"`
	TransferredOn    time.Time          `json:"transferred_at"`
	UpdatedTimestamp time.Time          `json:"updatedAt"`
	RecipientData    RecipientInfo      `json:"recipient"`
	SessionData      SessionDetails     `json:"session"`
	TransactionFee   int                `json:"fee_charged"`
	GatewayResponse  interface{}        `json:"gateway_response"` // Can be null
}

// TransferEventPayload represents the overall structure of the JSON event.
type TransferEventPayload struct {
	EventName string          `json:"event"`
	EventData TransferPayload `json:"data"`
}
