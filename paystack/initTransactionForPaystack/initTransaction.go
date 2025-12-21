package initTransactionForPaystack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	// EmailsVer "business-connect/controllers/authentication/emails"
	helperFunc "business-connect/paystack"
	// paystackBuyServices "business-connect/paystack/buyServicesPaystack"
	// PayueeHelper "business-connect/payueeTrans"
	"strconv"

	// "strings"

	Data "business-connect/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Replace with the actual Paystack API endpoint for verification
const paystackVerifyURL = "https://api.paystack.co/transaction/verify/%s"

func InitializePaystackTransaction(Email string, TransactionID string, Amount int, Metadata Data.ServiceMetaData) (map[string]interface{}, error) {

	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Printf("Failed to load .env file: %v\n", envErr)
	}

	SECRET_KEY := os.Getenv("PAYSTACK_LIVE_SECRET_KEY")
	// SECRET_KEY := os.Getenv("PAYSTACK_TEST_SECRET_KEY")

	// let's send the amount in the currency's sub unit to paystack
	amount := Amount * 100
	Metadata.CancelAction = "https://shopsphereafrica.com/cancel-transaction.html"
	Metadata.TransactionID = TransactionID

	// fmt.Println("this is the type for the metadata price 1: ", reflect.TypeOf(Metadata.Price))

	method := "POST"
	// Add a callback URL to the payload
	callbackURL := "https://shopsphereafrica.com/track-order.html"
	payload := map[string]interface{}{
		"email":        Email,
		"amount":       strconv.Itoa(amount),
		"callback_url": callbackURL,
		"metadata":     Metadata,
	}

	// fmt.Println("this is the type for the metadata price 2: ", reflect.TypeOf(Metadata.Price))
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return make(map[string]interface{}), errors.New("error encoding JSON")
	}

	req, err := http.NewRequest(method, helperFunc.InitializeTransactions, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return make(map[string]interface{}), errors.New("error creating request")
	}

	req.Header.Set("Authorization", "Bearer "+SECRET_KEY)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return make(map[string]interface{}), errors.New("error making request")
	}
	defer res.Body.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return make(map[string]interface{}), errors.New("error decoding JSON")
	}

	return data, nil
}

func PaystackCallbackHandler(ctx *fiber.Ctx) error {
	// reference := ctx.Query("reference")

	// referenceErr := helperFunc.PaystackHelper.FindByReference(reference)
	// if referenceErr == nil {
	// 	// fmt.Println("response error:", referenceErr)
	// 	return ctx.Redirect("https://shopsphereafrica.com/successful.html")
	// }

	// Verify the transaction
	// verificationResponse, err := VerifyPaystackTransaction(reference)
	// fmt.Println("response data 1 this is the error message: ", err)
	// if err != nil {
	// 	// fmt.Println("response data 2: ")
	// 	// Handle the error
	// 	if strings.Contains(err.Error(), "response data is empty") {
	// 		// Response data is empty, handle accordingly
	// 		// fmt.Println("response data 3: ")
	// 		fmt.Println("response error 1:")
	// 		return ctx.Redirect("https://shopsphereafrica.com/halfSuccessful.html")
	// 	} else {
	// 		// Some other error occurred, handle accordingly
	// 		// fmt.Println("response data 4: ")
	// 		fmt.Println("response error 2:")
	// 		return ctx.Redirect("https://shopsphereafrica.com/halfSuccessful.html")
	// 	}
	// }

	// get stored user id from request time line
	// userId := ctx.Locals("user-id")
	// StringConvertedToUint, stringToUintErr := StringToUint(verificationResponse.Data.Metadata.UserId)
	// if stringToUintErr != nil {
	// 	// Handle error
	// 	fmt.Println("response error string to uint:", stringToUintErr)
	// 	return ctx.Redirect("https://shopsphereafrica.com/halfSuccessful.html")
	// }

	// user, uuidErr := helperFunc.PaystackHelper.FindByUuidFromLocalPaystack(verificationResponse.Data.Metadata.UserId)
	// fmt.Println("this is the user id:", verificationResponse.Data.Metadata.UserId)
	// if uuidErr != nil {
	// 	fmt.Println("response error:", uuidErr)
	// 	return ctx.Redirect("https://shopsphereafrica.com/halfSuccessful.html")
	// }

	// if !verificationResponse.Status {
	// 	fmt.Println("response error 4:")
	// 	return ctx.Redirect("https://shopsphereafrica.com/halfSuccessful.html")
	// }

	// fmt.Println("response data2: ", verificationResponse)
	// Transaction successful let's register the user to the database
	// addTransErr := PaystackSaveToDbCallbackHandler(user, verificationResponse)
	// if addTransErr != nil {
	// 	// fmt.Println("error occurring: ", addTransErr)
	// 	// this only runs when there is an error updating the database
	// 	fmt.Println("response error 5:")
	// 	return ctx.Redirect("https://shopsphereafrica.com/halfSuccessful.html")
	// }

	// Redirect to success page
	return ctx.Redirect("https://shopsphereafrica.com/track-order.html")
}

// Function to verify Paystack transaction
func VerifyPaystackTransaction(reference string) (helperFunc.PaystackVerificationResponse, error) {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Printf("Failed to load .env file: %v\n", envErr)
	}

	SECRET_KEY := os.Getenv("PAYSTACK_LIVE_SECRET_KEY")

	url := fmt.Sprintf(paystackVerifyURL, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return helperFunc.PaystackVerificationResponse{}, err
	}

	// Set Paystack API secret key in the Authorization header
	req.Header.Set("Authorization", "Bearer "+SECRET_KEY)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return helperFunc.PaystackVerificationResponse{}, err
	}
	defer res.Body.Close()

	// Read the response body into a buffer
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return helperFunc.PaystackVerificationResponse{}, err
	}

	// Key to be converted to an int
	key := "price"
	key2 := "Amount"
	key3 := "TranCharge"

	fmt.Println(" this is the response from paystack: ", string(responseBody))

	PaystackJsonStr1, err := ConvertKeyToInt(string(responseBody), key)
	if err != nil {
		fmt.Println("Error converting key value to int:", err)
		return helperFunc.PaystackVerificationResponse{}, err
	}

	PaystackJsonStr2, err := ConvertKeyToInt(PaystackJsonStr1, key2)
	if err != nil {
		return helperFunc.PaystackVerificationResponse{}, err
	}

	PaystackJsonStr3, err := ConvertKeyToInt(PaystackJsonStr2, key3)
	if err != nil {
		return helperFunc.PaystackVerificationResponse{}, err
	}

	PaystackJsonStr4, err := ConvertToBool(PaystackJsonStr3, "AutoRenew")
	if err != nil {
		return helperFunc.PaystackVerificationResponse{}, err
	}

	// Unmarshal Paystack response into your struct
	var verificationResponse helperFunc.PaystackVerificationResponse
	err = json.Unmarshal([]byte(PaystackJsonStr4), &verificationResponse)

	// fmt.Println("response data here 7: ")
	if err != nil {
		fmt.Println("Error un-marshaling Paystack response:", err)
		return helperFunc.PaystackVerificationResponse{}, err
	}

	return verificationResponse, nil
}

// func PaystackSaveToDbCallbackHandler(user Data.User, data helperFunc.PaystackVerificationResponse) error {

// 	metadataString := data.Data.Metadata
// 	// fmt.Println("this is the service id: ", data.Data.Metadata.ServiceID)

// 	switch data.Data.Metadata.ServiceID {
// 	case "airtime":
// 		// fmt.Println("this is the price to update to the database: ", metadataString.Price)
// 		// perform transaction using paystack verified payment update the user balance if there was an error buying any service
// 		_, airtimeErr := paystackBuyServices.PaystackBuyAirtime(int64(metadataString.Price), metadataString.Network, metadataString.PhoneNumber)
// 		// fmt.Println("6")
// 		if airtimeErr != nil {
// 			// if there was an error updating the buying airtime let's update the user's balance
// 			// cause the user is making payment with paystack and the payment was successful
// 			_, err := PayueeHelper.PaystackHelper.UpdateUserBalance(user, metadataString.Price)
// 			// fmt.Println("8")
// 			if err != nil {
// 				// fmt.Println("9")
// 				return errors.New("an error occurred while updating users balance after paying with paystack")
// 			}
// 			// Let's update the database based on the failed transaction that occurred
// 			data.Data.Status = "failed"
// 			metadataString.Status = "failed"
// 			addDataErr := PayueeHelper.PaystackHelper.AddAirtimeTransactionHistoryP(user, data, metadataString)
// 			if addDataErr != nil {
// 				return errors.New("an error occurred while saving to the database")
// 			}
// 			return errors.New("an error occurred while trying to buy airtime")
// 		}
// 		// Let's update the database based on the transaction that occurred
// 		metadataString.Status = "success"
// 		addAirtimeErr := PayueeHelper.PaystackHelper.AddAirtimeTransactionHistoryP(user, data, metadataString)
// 		if addAirtimeErr != nil {
// 			return errors.New("an error occurred")
// 		}

// 		emailErr := EmailsVer.AirtimeConfirmationEmail(user.FirstName+" "+user.LastName, user.Email, "₦"+strconv.Itoa(metadataString.Price), metadataString.PhoneNumber, "₦"+strconv.Itoa(int(user.WalletBalance)), "External Bank")
// 		if emailErr != nil {
// 			return errors.New("an error while sending verification email")
// 		}
// 	case "rechargePin":
// 		// perform transaction using paystack verified payment update the user balance if there was an error buying any service
// 		rechargePinResponse, dataErr := paystackBuyServices.PaystackBuyRechargePin(metadataString.Network, metadataString.NumberOfPin, metadataString.Value)
// 		if dataErr != nil {
// 			// if there was an error updating the buying recharge pin let's update the user's balance
// 			// cause the user is making payment with paystack and the payment was successful
// 			_, err := PayueeHelper.PaystackHelper.UpdateUserBalance(user, metadataString.Price)
// 			if err != nil {
// 				return errors.New("an error occurred while updating users balance after paying with paystack")
// 			}
// 			// Let's update the database based on the failed transaction that occurred
// 			data.Data.Status = "failed"
// 			metadataString.Status = "failed"
// 			addDataErr := PayueeHelper.PaystackHelper.AddRechargePinTransactionHistoryP(user, data, metadataString, nil)
// 			if addDataErr != nil {
// 				return errors.New("an error occurred while saving to the database")
// 			}
// 			return errors.New("an error occurred while trying to buy airtime")
// 		}
// 		// Let's update the database based on the transaction that occurred
// 		metadataString.Status = "success"
// 		addAirtimeErr := PayueeHelper.PaystackHelper.AddRechargePinTransactionHistoryP(user, data, metadataString, rechargePinResponse.PinNumbers)
// 		if addAirtimeErr != nil {
// 			return errors.New("an error occurred")
// 		}
// 	case "data":
// 		// perform transaction using paystack verified payment update the user balance if there was an error buying any service
// 		_, dataErr := paystackBuyServices.PaystackBuyData(metadataString.PlanID, metadataString.NetworkPlan, metadataString.PhoneNumber)
// 		if dataErr != nil {
// 			// if there was an error updating the buying data let's update the user's balance
// 			// cause the user is making payment with paystack and the payment was successful
// 			_, err := PayueeHelper.PaystackHelper.UpdateUserBalance(user, metadataString.Price)
// 			if err != nil {
// 				return errors.New("an error occurred while updating users balance after paying with paystack")
// 			}
// 			// Let's update the database based on the failed transaction that occurred
// 			data.Data.Status = "failed"
// 			metadataString.Status = "failed"
// 			addDataErr := PayueeHelper.PaystackHelper.AddDataTransactionHistoryP(user, data, metadataString)
// 			if addDataErr != nil {
// 				return errors.New("an error occurred while saving to the database")
// 			}
// 			return errors.New("an error occurred while trying to buy data")
// 		}
// 		// Let's update the database based on the transaction that occurred
// 		metadataString.Status = "success"
// 		addDataErr := PayueeHelper.PaystackHelper.AddDataTransactionHistoryP(user, data, metadataString)
// 		if addDataErr != nil {
// 			return errors.New("an error occurred while saving to the database")
// 		}
// 		emailErr := EmailsVer.DataConfirmationEmail(user.FirstName+" "+user.LastName, user.Email, "₦"+strconv.Itoa(metadataString.Price), metadataString.Bundle, metadataString.PhoneNumber, "₦"+strconv.Itoa(int(user.WalletBalance)), "External Bank")
// 		if emailErr != nil {
// 			return errors.New("an error while sending verification email")
// 		}
// 	case "educationalPayment":
// 		// perform transaction using paystack verified payment update the user balance if there was an error buying any service
// 		educationErr := paystackBuyServices.PaystackBuyWaecPin(metadataString.PhoneNumber)
// 		if educationErr != nil {
// 			// if there was an error updating the buying waec pin let's update the user's balance
// 			// cause the user is making payment with paystack and the payment was successful
// 			_, err := PayueeHelper.PaystackHelper.UpdateUserBalance(user, metadataString.Price)
// 			if err != nil {
// 				return errors.New("an error occurred while updating users balance after paying with paystack")
// 			}
// 			// Let's update the database based on the failed transaction that occurred
// 			data.Data.Status = "failed"
// 			metadataString.Status = "failed"
// 			addDataErr := PayueeHelper.PaystackHelper.AddEducationalPaymentsTransactionHistoryP(user, data, metadataString)
// 			if addDataErr != nil {
// 				return errors.New("an error occurred while saving to the database")
// 			}
// 			return errors.New("an error occurred while trying to buy waec pin")
// 		}
// 		// Let's update the database based on the transaction that occurred
// 		metadataString.Status = "success"
// 		addEducationErr := PayueeHelper.PaystackHelper.AddEducationalPaymentsTransactionHistoryP(user, data, metadataString)
// 		if addEducationErr != nil {
// 			return errors.New("an error occurred")
// 		}
// 	case "decoder":
// 		// perform transaction using paystack verified payment update the user balance if there was an error buying any service
// 		_, decoderErr := paystackBuyServices.PaystackTvSubscription(metadataString.Operator, metadataString.DecoderNumber, metadataString.VariationID, metadataString.PhoneNumber)
// 		if decoderErr != nil {
// 			// if there was an error updating the buying decoder let's update the user's balance
// 			fmt.Println("this is the decoder error: " + decoderErr.Error())
// 			// cause the user is making payment with paystack and the payment was successful
// 			// _, err := PayueeHelper.PaystackHelper.UpdateUserBalance(user, metadataString.Price)
// 			// if err != nil {
// 			// 	return errors.New("an error occurred while updating users balance after paying with paystack")
// 			// }
// 			// Let's update the database based on the failed transaction that occurred
// 			data.Data.Status = "failed"
// 			metadataString.Status = "failed"
// 			addDataErr := PayueeHelper.PaystackHelper.AddDecoderTransactionHistoryP(user, data, metadataString)
// 			if addDataErr != nil {
// 				return errors.New("an error occurred while saving to the database")
// 			}
// 			return errors.New("an error occurred while trying to buy decoder")
// 		}
// 		// Let's update the database based on the transaction that occurred
// 		metadataString.Status = "success"
// 		addDecodeErr := PayueeHelper.PaystackHelper.AddDecoderTransactionHistoryP(user, data, metadataString)
// 		if addDecodeErr != nil {
// 			return errors.New("an error occurred")
// 		}
// 		emailErr := EmailsVer.TvSubscriptionConfirmationEmail(user.FirstName+" "+user.LastName, user.Email, metadataString.Operator, metadataString.Plan, metadataString.DecoderNumber, metadataString.PhoneNumber, "₦"+strconv.Itoa(int(user.WalletBalance)-metadataString.Price), "External Bank")
// 		if emailErr != nil {
// 			return errors.New("an error while sending an airtime verification email")
// 		}
// 	case "electricity":
// 		// perform transaction using paystack verified payment update the user balance if there was an error buying any service
// 		electricityResponse, electricErr := paystackBuyServices.PaystackBuyElectricToken(metadataString.RegionID, strconv.Itoa(metadataString.Price), metadataString.MeterNumber, metadataString.VariationID, metadataString.PhoneNumber)
// 		if electricErr != nil {
// 			// if there was an error updating the buying electric token let's update the user's balance
// 			// cause the user is making payment with paystack and the payment was successful
// 			_, err := PayueeHelper.PaystackHelper.UpdateUserBalance(user, metadataString.Price)
// 			if err != nil {
// 				return errors.New("an error occurred while updating users balance after paying with paystack")
// 			}
// 			// Let's update the database based on the failed transaction that occurred
// 			data.Data.Status = "failed"
// 			metadataString.Status = "failed"
// 			addDataErr := PayueeHelper.PaystackHelper.AddElectricityTransactionHistoryP(user, data, metadataString, Data.ElectricityResponse{})
// 			if addDataErr != nil {
// 				return errors.New("an error occurred while saving to the database")
// 			}
// 			return errors.New("an error occurred while trying to buy electric token")
// 		}
// 		// Let's update the database based on the transaction that occurred
// 		metadataString.Status = "success"
// 		addElectricErr := PayueeHelper.PaystackHelper.AddElectricityTransactionHistoryP(user, data, metadataString, electricityResponse)
// 		if addElectricErr != nil {
// 			return errors.New("an error occurred")
// 		}
// 		emailErr := EmailsVer.ElectricityTransactionConfirmationEmail(user.FirstName+" "+user.LastName, user.Email, electricityResponse.Data.Electricity, metadataString.MeterNumber, electricityResponse.Data.Token, electricityResponse.Data.Units, metadataString.PhoneNumber, "₦"+strconv.Itoa(int(metadataString.Price)), "₦"+strconv.Itoa(int(user.WalletBalance)-metadataString.Price), "External Bank")
// 		if emailErr != nil {
// 			return errors.New("an error while sending an electricity verification email")
// 		}
// 	case "fundWallet":
// 		// perform transaction using paystack verified payment update the user balance only if there was no error making payments
// 		// Let's update the database based on the transaction that occurred
// 		addFundsErr := PayueeHelper.PaystackHelper.AddFundWalletTransactionHistory(user, data, metadataString)
// 		if addFundsErr != nil {
// 			return errors.New("an error occurred while updating funding wallet database record")
// 		}

// 		emailErr := EmailsVer.WalletFundingConfirmationEmail(user.FirstName+" "+user.LastName, user.Email, data.Data.Channel, "₦"+strconv.Itoa(metadataString.Amount), "₦"+strconv.Itoa(int(user.WalletBalance)), "₦"+strconv.Itoa(int(user.WalletBalance)+metadataString.Amount))
// 		if emailErr != nil {
// 			return errors.New("an error while sending verification email for fund wallet")
// 		}
// 	case "sendFunds":
// 		// perform transaction using paystack verified payment update the user balance only if there was no error making payments
// 		// Let's update the database based on the transaction that occurred
// 		// addSendFundErr := PayueeHelper.PaystackHelper.AddSendFundsTransactionHistoryP(user, data, metadataString)
// 		// if addSendFundErr != nil {
// 		// 	return errors.New("an error occurred while sending funds")
// 		// }
// 	default:
// 		return errors.New("invalid request body")
// 	}

// 	return nil
// }

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

func StringToUint(str string) (uint, error) {
	// Attempt to parse the string to a uint
	u, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(u), nil
}
