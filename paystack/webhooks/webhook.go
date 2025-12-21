package webhooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	// "strings"

	"net/http"
	// EmailsVer "business-connect/controllers/authentication/emails"
	SendEmail "business-connect/controllers/authentication/emails"
	dbFunc "business-connect/database/dbHelpFunc"

	// conn "business-connect/database"
	// Dataa "business-connect/models"
	// helperFunc "business-connect/paystack"

	// paystackBuyServices "business-connect/paystack/buyServicesPaystack"

	// intiTrasfer "business-connect/paystack/transferFundsToUser"
	// PayueeHelper "business-connect/payueeTrans"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Webhook for getting transaction status
func WebHookStatus(ctx *fiber.Ctx) error {
	// Read the request body
	fmt.Println(" this is the response from paystack web hook origin: ", ctx.Get("Origin"))
	responseBody := ctx.Body()

	// Key to be converted to an int
	key := "price"
	// key2 := "Amount"
	// key3 := "TranCharge"
	// key4 := "AutoRenew"
	// key5 := "TransNumb"

	fmt.Println(" this is the response from paystack: ", string(responseBody))

	PaystackJsonStr1, err := ConvertKeyToIntSafe(responseBody, key)
	// PaystackJsonStr1, err := ConvertKeyToInt(string(responseBody), key)
	if err != nil {
		fmt.Println("Error converting key value to int:", err)
		// return err
	}

	// PaystackJsonStr2, err := ConvertKeyToInt(PaystackJsonStr1, key2)
	// if err != nil {
	// return err
	// }

	// PaystackJsonStr3, err := ConvertKeyToInt(PaystackJsonStr2, key3)
	// if err != nil {
	// 	// return err
	// }

	// PaystackJsonStr4, err := ConvertToBool(PaystackJsonStr3, key4)
	// if err != nil {
	// 	// return err
	// }

	// PaystackJsonStr5, err := ConvertKeyToUint(PaystackJsonStr4, key5)
	// if err != nil {
	// 	// return err
	// }

	transactionType := "empty"

	// Attempt to unmarshal into WebhookData
	var webhookData WebhookData
	err1 := json.Unmarshal([]byte(PaystackJsonStr1), &webhookData)
	if err1 == nil {
		transactionType = "charge.success"
		// Successfully unmarshaled into WebhookData
		fmt.Println("Parsed using WebhookData:", webhookData)
		// return nil
	}

	if err1 != nil {
		fmt.Println("transaction type err 1 :-----------------------------------------:", err1)
	}

	fmt.Println("transaction type 0 :-----------------------------------------:", transactionType)

	// Attempt to unmarshal into helperFunc.TransferEventPayload
	// var webhookDataTF helperFunc.TransferEventPayload
	// err2 := json.Unmarshal(responseBody, &webhookDataTF)
	// if err2 == nil {
	// 	transactionType = "transfer.success"
	// 	// Successfully unmarshaled into TransferEventPayload
	// 	fmt.Println("Parsed using TransferEventPayload:", webhookDataTF)
	// 	// return nil
	// }

	// if err2 != nil {
	// 	fmt.Println("transaction type err 2 :-----------------------------------------:", err2)
	// }

	fmt.Println("transaction type 1 :-----------------------------------------:", transactionType)
	// Print the raw JSON data
	fmt.Println("Raw WebHook Data:", PaystackJsonStr1)

	fmt.Println("Parsed WebHook Data:", string(responseBody))

	// if transactionType == "charge.success" {
	fmt.Println("transaction type 01 :-----------------------------------------:", transactionType)
	// Check the type of event
	switch webhookData.Event {
	case "charge.success":
		fmt.Println("Charge successful event")
		// Access data specific to charge.success event
		// fmt.Println("WebHook Data:", webhookData)
		// fmt.Println("Amount:", webhookData.Data.Amount)
		// fmt.Println("Recipient Reference:", webhookData.Data.Reference)
		// user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocalPaystack(webhookData.Data.Metadata.TransactionID)
		// // fmt.Println("this is the user id:", webhookData.Data.Metadata.UserId)
		// if uuidErr != nil {
		// 	fmt.Println("response error: ", uuidErr)
		// 	return ctx.SendStatus(http.StatusNotAcceptable)
		// }
		PaystackWebHookSaveToDbCallbackHandler(webhookData)
		// ... (access other fields as needed)
	default:
		fmt.Println("Unknown event")
	}
	// } else if transactionType == "transfer.success" {
	// 	// Check the type of event
	// 	switch webhookDataTF.EventName {
	// 	case "transfer.success":
	// 		// fmt.Println("Started the transfer")

	// 		// fmt.Println("User ID:", webhookDataTF.EventData.RecipientData.Metadata.UserId)

	// 		user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocalPaystack(webhookDataTF.EventData.RecipientData.Metadata.UserId)
	// 		if uuidErr != nil {
	// 			fmt.Println("Response error:", uuidErr)
	// 			return ctx.SendStatus(http.StatusNotAcceptable)
	// 		}

	// 		// fmt.Println("Transfer successful event")
	// 		// fmt.Println("Webhook Data:", webhookDataTF)
	// 		// fmt.Println("Amount:", webhookDataTF.EventData.RecipientData.Metadata.Amount)
	// 		// fmt.Println("Recipient Name:", webhookDataTF.EventData.TransactionRef)

	// 		addSendFundErr := PayueeHelper.PaystackHelper.AddSendFundsSuccessTransactionHistoryP(user, webhookDataTF, webhookDataTF.EventData.RecipientData.Metadata.TransNumb)
	// 		if addSendFundErr != nil {
	// 			return errors.New("an error occurred while sending funds")
	// 		}

	// 		emailErr := EmailsVer.SendFundsConfirmationEmail(
	// 			user.FirstName+" "+user.LastName,
	// 			user.Email,
	// 			webhookDataTF.EventData.RecipientData.Metadata.AccountName,
	// 			webhookDataTF.EventData.RecipientData.Metadata.Bank,
	// 			webhookDataTF.EventData.RecipientData.Metadata.AccountNumber,
	// 			"₦"+strconv.Itoa(webhookDataTF.EventData.RecipientData.Metadata.Amount),
	// 			"₦"+strconv.Itoa(int(user.WalletBalance)),
	// 			"₦"+strconv.Itoa(int(user.WalletBalance)+webhookDataTF.EventData.RecipientData.Metadata.Amount),
	// 		)

	// 		if emailErr != nil {
	// 			return errors.New("an error occurred while sending verification email for fund wallet")
	// 		}

	// 	case "transfer.failed":
	// 		// fmt.Println("Started the transfer")

	// 		// fmt.Println("User ID:", webhookDataTF.EventData.RecipientData.Metadata.UserId)

	// 		user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocalPaystack(webhookDataTF.EventData.RecipientData.Metadata.UserId)
	// 		if uuidErr != nil {
	// 			fmt.Println("Response error:", uuidErr)
	// 			return ctx.SendStatus(http.StatusNotAcceptable)
	// 		}

	// 		addSendFundErr := PayueeHelper.PaystackHelper.AddSendFundsFailedTransactionHistoryWP(user, webhookDataTF, webhookDataTF.EventData.RecipientData.Metadata.TransNumb)
	// 		if addSendFundErr != nil {
	// 			return errors.New("an error occurred while sending funds")
	// 		}
	// 	case "transfer.reversed":
	// 		fmt.Println("Transfer reversed event")
	// 		// Access data specific to transfer.reversed event
	// 		fmt.Println("Amount:", webhookData.Data.Amount)
	// 		// fmt.Println("Started the transfer")

	// 		// fmt.Println("User ID:", webhookDataTF.EventData.RecipientData.Metadata.UserId)

	// 		user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocalPaystack(webhookDataTF.EventData.RecipientData.Metadata.UserId)
	// 		if uuidErr != nil {
	// 			fmt.Println("Response error:", uuidErr)
	// 			return ctx.SendStatus(http.StatusNotAcceptable)
	// 		}

	// 		addSendFundErr := PayueeHelper.PaystackHelper.AddSendFundsFailedTransactionHistoryWP(user, webhookDataTF, webhookDataTF.EventData.RecipientData.Metadata.TransNumb)
	// 		if addSendFundErr != nil {
	// 			return errors.New("an error occurred while sending funds")
	// 		}
	// 	default:
	// 		fmt.Println("Unknown event")
	// 	}
	// }

	fmt.Println("transaction type 2 :-----------------------------------------:", transactionType)

	return ctx.SendStatus(http.StatusOK)
}

// ConvertKeyToInt converts a specified key's value to an integer if it's a string.
func ConvertKeyToInt2(responseBody string, key string) (string, error) {
	// Unmarshal JSON into a generic map
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &jsonData); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	// Check if the key exists
	value, exists := jsonData[key]
	if !exists {
		fmt.Printf("Key %s not found in the response, skipping conversion.\n", key)
		return responseBody, nil // Return unchanged JSON
	}

	// If the value is already an integer, skip conversion
	if _, ok := value.(float64); ok {
		fmt.Printf("Key %s is already an integer, skipping conversion.\n", key)
		return responseBody, nil // Return unchanged JSON
	}

	// If the value is a string, attempt to convert to integer
	if strValue, ok := value.(string); ok {
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			return "", fmt.Errorf("error converting key %s to int: %w", key, err)
		}

		// Update the value in the map
		jsonData[key] = intValue
	}

	// Marshal the map back into JSON format
	updatedJSON, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	return string(updatedJSON), nil
}

// ConvertNestedKeyToInt converts a nested key's value to an integer if it's a string.
func ConvertNestedKeyToInt(responseBody string, nestedKeys ...string) (string, error) {
	// Unmarshal JSON into a generic map
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &jsonData); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	// Handle case where there are no nested keys
	if len(nestedKeys) == 0 {
		return responseBody, fmt.Errorf("no keys provided for conversion")
	}

	// Navigate through the nested keys
	current := jsonData
	for i, key := range nestedKeys {
		if i == len(nestedKeys)-1 { // Final key
			// Check if the key exists
			value, exists := current[key]
			if !exists {
				fmt.Printf("Key %s not found in the response, skipping conversion.\n", key)
				return responseBody, nil // Return unchanged JSON
			}

			// If the value is already an integer, skip conversion
			if _, ok := value.(float64); ok {
				fmt.Printf("Key %s is already an integer, skipping conversion.\n", key)
				return responseBody, nil // Return unchanged JSON
			}

			// If the value is a string, attempt to convert to integer
			if strValue, ok := value.(string); ok {
				intValue, err := strconv.Atoi(strValue)
				if err != nil {
					return "", fmt.Errorf("error converting key %s to int: %w", key, err)
				}

				// Update the value in the map
				current[key] = intValue
			}
		} else { // Intermediate keys
			// Drill down into the nested object
			if next, ok := current[key].(map[string]interface{}); ok {
				current = next
			} else {
				return "", fmt.Errorf("key %s not found or not a nested object", key)
			}
		}
	}

	// Marshal the map back into JSON format
	updatedJSON, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	return string(updatedJSON), nil
}

// ConvertToBool converts a specified key's string value to a boolean if necessary.
func ConvertToBool2(responseBody, boolKey string) (string, error) {
	// Unmarshal JSON into a generic map
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &jsonData); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	// Check if the key exists
	value, exists := jsonData[boolKey]
	if !exists {
		fmt.Printf("Key %s not found in the response, skipping conversion.\n", boolKey)
		return responseBody, nil // Return unchanged JSON
	}

	// If the value is already a boolean, skip conversion
	if _, ok := value.(bool); ok {
		fmt.Printf("Key %s is already a boolean, skipping conversion.\n", boolKey)
		return responseBody, nil // Return unchanged JSON
	}

	// If the value is a string, attempt to convert to boolean
	if strValue, ok := value.(string); ok {
		boolValue, err := strconv.ParseBool(strValue)
		if err != nil {
			return "", fmt.Errorf("error converting key %s to bool: %w", boolKey, err)
		}

		// Update the value in the map
		jsonData[boolKey] = boolValue
	}

	// Marshal the map back into JSON format
	updatedJSON, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	return string(updatedJSON), nil
}

// ConvertNestedKeyToBool converts a nested key's string value to a boolean if necessary.
func ConvertNestedKeyToBool(responseBody string, nestedKeys ...string) (string, error) {
	// Unmarshal JSON into a generic map
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &jsonData); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	// Navigate through the nested keys
	current := jsonData
	for i, key := range nestedKeys {
		if i == len(nestedKeys)-1 { // Final key
			// Check if the key exists
			value, exists := current[key]
			if !exists {
				fmt.Printf("Key %s not found in the response, skipping conversion.\n", key)
				return responseBody, nil // Return unchanged JSON
			}

			// If the value is already a boolean, skip conversion
			if _, ok := value.(bool); ok {
				fmt.Printf("Key %s is already a boolean, skipping conversion.\n", key)
				return responseBody, nil // Return unchanged JSON
			}

			// If the value is a string, attempt to convert to boolean
			if strValue, ok := value.(string); ok {
				boolValue, err := strconv.ParseBool(strValue)
				if err != nil {
					return "", fmt.Errorf("error converting key %s to bool: %w", key, err)
				}

				// Update the value in the map
				current[key] = boolValue
			}
		} else { // Intermediate keys
			// Drill down into the nested object
			if next, ok := current[key].(map[string]interface{}); ok {
				current = next
			} else {
				return "", fmt.Errorf("key %s not found or not a nested object", key)
			}
		}
	}

	// Marshal the map back into JSON format
	updatedJSON, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	return string(updatedJSON), nil
}

func ConvertKeyToInt(responseBody string, key string) (string, error) {
	// Regular expression to match the key and its associated value
	re := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, key))

	// Find the key-value pair in the JSON string
	match := re.FindStringSubmatch(responseBody)
	if match == nil {
		// Key not found, handle the error
		fmt.Println("Key not found in JSON.")
		return "", fmt.Errorf("key not found: %s", key)
	}

	// Extract the value associated with the key
	valueStr := match[1]
	fmt.Println("Original value:", valueStr)

	// Convert the value to an integer
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		// Conversion error, handle the error
		fmt.Println("Error converting value to int:", err)
		return "", err
	}

	// Replace the original value with the new integer value in the JSON string
	paystackJsonStr := re.ReplaceAllString(string(responseBody), fmt.Sprintf(`"%s":%d`, key, valueInt))

	return paystackJsonStr, nil
}

func ConvertKeyToIntSafe(responseBody []byte, key string) ([]byte, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(responseBody, &jsonData); err != nil {
		return nil, err
	}

	if val, ok := jsonData[key]; ok {
		if strVal, ok := val.(string); ok {
			if intVal, err := strconv.Atoi(strVal); err == nil {
				jsonData[key] = intVal
			}
		}
	}

	return json.Marshal(jsonData)
}

func ConvertKeyToUint(responseBody string, key string) (string, error) {
	// Regular expression to match the key and its associated value (as a string)
	re := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, key))

	// Find the key-value pair in the JSON string
	match := re.FindStringSubmatch(responseBody)
	if match == nil {
		// Key not found, handle the error
		return "", fmt.Errorf("key not found: %s", key)
	}

	// Extract the value associated with the key
	valueStr := match[1]

	// Convert the value to a uint
	valueUint, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("error converting value to uint: %w", err)
	}

	// Replace the original value (string) with the new uint (number) value in the JSON string
	modifiedJson := re.ReplaceAllString(responseBody, fmt.Sprintf(`"%s":%d`, key, valueUint))

	return modifiedJson, nil
}

func ConvertToBool(responseBody, boolValue string) (string, error) {
	// Regular expression to match the "AutoRenew" key and its associated value
	re := regexp.MustCompile(fmt.Sprintf(`"%s":"([^"]+)"`, boolValue))

	// Find the "AutoRenew" key-value pair in the JSON string
	match := re.FindStringSubmatch(responseBody)
	if match == nil {
		// Key not found, handle the error
		fmt.Println("AutoRenew key not found in JSON.")
		return "", fmt.Errorf("AutoRenew key not found")
	}

	// Extract the value associated with the "AutoRenew" key
	valueStr := match[1]
	fmt.Println("Original value:", valueStr)

	// Convert the value to a bool
	valueBool, err := strconv.ParseBool(valueStr)
	if err != nil {
		// Conversion error, handle the error
		fmt.Println("Error converting value to bool:", err)
		// return "", err
	}

	// Replace the original value with the new bool value in the JSON string
	paystackJsonStr := re.ReplaceAllString(responseBody, fmt.Sprintf(`"%s":%t`, boolValue, valueBool))

	return paystackJsonStr, nil
}

func PaystackWebHookSaveToDbCallbackHandler(data WebhookData) error {

	metadataString := data.Data.Metadata
	// fmt.Println("this is the service id: ", data.Data.Metadata.ServiceID)

	num, err := strconv.ParseUint(metadataString.TransactionID, 10, 64) // base 10, 64-bit unsigned
	if err != nil {
		return errors.New("Conversion error:" + err.Error())
	}

	// Call the database helper function to retrieve the order and update payment status
	orderHistory, err := dbFunc.DBHelper.GetAndUpdateOrder(uint(num), data.Data.Status)
	if err != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return errors.New("failed to retrieve order")
	}

	// send a confirmation email
	emailErr := SendEmail.ShopsphereConfirmationEmail(*orderHistory, orderHistory.ProductOrders)
	if emailErr != nil {
		return errors.New("failed to send order email")
	}

	return nil
}
