package order

import (
	"errors"
	"fmt"
	"reflect"

	// "fmt"
	"math"
	"net/http"
	"strconv"

	dbFunc "business-connect/database/dbHelpFunc"
	Data "business-connect/models"
	initTrans "business-connect/paystack/initTransactionForPaystack"

	// SendEmail "business-connect/controllers/authentication/emails"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PaginationData struct {
	NextPage     int
	PreviousPage int
	CurrentPage  int
	TotalPages   int
	TwoAfter     int
	TwoBefore    int
	ThreeAfter   int
	Offset       int
	AllRecords   int64
}

type OrderBody struct {
	OrderHistoryBody Data.OrderHistoryBody   `json:"order_history_body"`
	ProductOrderBody []Data.ProductOrderBody `json:"product_order_body"`
}

func AddOrder(ctx *fiber.Ctx) error {

	var NewOrder OrderBody

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&NewOrder); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// fmt.Println("this is the order body: ", NewOrder)
	NewOrder.OrderHistoryBody.OrderStatus = "processing"

	// adding new order history
	orderResultID, _, _, orderErr := dbFunc.DBHelper.AddOrder(NewOrder.OrderHistoryBody, NewOrder.ProductOrderBody)

	// checking if there was an error comparing the orders
	if orderErr != nil {
		fmt.Println("order error: ", orderErr)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to add new order",
		})
	}

	transIdStr := strconv.FormatUint(uint64(orderResultID), 10)

	var details Data.ServiceMetaData = Data.ServiceMetaData{
		TransactionID: transIdStr,
		Price:         int(NewOrder.OrderHistoryBody.OrderCost),
		Status:        "pending",
		PhoneNumber:   NewOrder.OrderHistoryBody.CustomerPhoneNumber,
		EmailID:       NewOrder.OrderHistoryBody.CustomerEmail,
	}

	// initialize paystack payment
	// let's convert details to an interface in order to add it to the metadata field in paystack so that
	// it can be used in the callback function to know which transaction took place and update the database accordingly
	metadata, metadataErr := ConvertToMetadata(details)
	if metadataErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "an error occurred",
		})
	}
	// perform transaction using paystack
	initTransURL, airtimeErr := initTrans.InitializePaystackTransaction(NewOrder.OrderHistoryBody.CustomerEmail, transIdStr, int(NewOrder.OrderHistoryBody.OrderCost), metadata)
	if airtimeErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "an error occurred",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  initTransURL,
		"order_id": orderResultID,
	})
}

func GetOrder(ctx *fiber.Ctx) error {
	// Extract order ID from path or query parameter (adjust based on your implementation)
	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	// Call the database helper function to retrieve the order
	orderHistory, err := dbFunc.DBHelper.GetOrder(uint(orderID))
	if err != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "order not found",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve order",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": *orderHistory,
	})
}

func GetBusinessConnectOrdersByLimit(ctx *fiber.Ctx) error {

	// // Get stored user id from request timeline
	// userId := ctx.Locals("user-id")

	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// _, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)

	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": uuidErr.Error(),
	// 	})
	// }

	// this is to get the page number to route to
	pageNumber := 1
	limitNumber := ctx.Params("orderLimit")
	eachPage := 15

	if limitNumber != "" {
		pageNumber, _ = strconv.Atoi(limitNumber)
	}

	offset := (pageNumber - 1) * eachPage

	productRecords, totalRecords, productRecordsErr := dbFunc.DBHelper.GetBusinessConnectOrdersByLimit(eachPage, offset)
	if productRecordsErr != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product history",
		})
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(eachPage)))

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": productRecords,
		"pagination": PaginationData{
			NextPage:     pageNumber + 1,
			PreviousPage: pageNumber - 1,
			CurrentPage:  pageNumber,
			TotalPages:   totalPages,
			TwoAfter:     pageNumber + 2,
			TwoBefore:    pageNumber - 2,
			ThreeAfter:   pageNumber + 4,
			AllRecords:   totalRecords,
		},
	})
}

func GetBusinessConnectOrder(ctx *fiber.Ctx) error {

	// // Get stored user id from request timeline
	// userId := ctx.Locals("user-id")

	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// _, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)

	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": uuidErr.Error(),
	// 	})
	// }

	// Extract order ID from path or query parameter (adjust based on your implementation)
	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	// Call the database helper function to retrieve the order
	orderHistory, err := dbFunc.DBHelper.GetOrder(uint(orderID))
	if err != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "order not found",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve order",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": *orderHistory,
	})
}

func UpdateBusinessConnectOrderStatus(ctx *fiber.Ctx) error {

	// // Get stored user id from request timeline
	// userId := ctx.Locals("user-id")

	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// _, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)

	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": uuidErr.Error(),
	// 	})
	// }

	type OrderStatus struct {
		OrderID int    `json:"orderID"`
		Status  string `json:"status"`
	}

	var StatusUpdate OrderStatus

	// Check if there is an error binding the request
	if bindErr := ctx.BodyParser(&StatusUpdate); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// Call the database helper function to retrieve the order
	err := dbFunc.DBHelper.UpdateOrderStatus(uint(StatusUpdate.OrderID), StatusUpdate.Status)
	if err != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "order not found",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve order",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "successfully updated status",
	})
}

func UpdateBusinessConnectProduct(ctx *fiber.Ctx) error {
	type ProductToUpdate struct {
		ProductID          int     `json:"product_id"`
		ProductTitle       string  `json:"product_title"`
		ProductDescription string  `json:"product_description"`
		InitialCost        float64 `json:"initial_cost"`
		SellingPrice       float64 `json:"selling_price"`
		ProductStock       int64   `json:"product_stock"`
		NetWeight          int64   `json:"net_weight"`
		Category           string  `json:"category"`
		ProductRank        int     `json:"product_rank"`
		Tags               string  `json:"tags"`
		PublishStatus      string  `json:"publish_status"`
		FeaturedStatus     string  `json:"featured_status"`
	}

	var ProductUpdate ProductToUpdate

	if bindErr := ctx.BodyParser(&ProductUpdate); bindErr != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	// Fetch the existing product first
	existingProduct, err := dbFunc.DBHelper.GetProductByID(uint(ProductUpdate.ProductID))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Start with existing product
	UpdatedProduct := existingProduct

	// Overwrite updated fields
	UpdatedProduct.Title = ProductUpdate.ProductTitle
	UpdatedProduct.Description = ProductUpdate.ProductDescription
	UpdatedProduct.InitialCost = ProductUpdate.InitialCost
	UpdatedProduct.SellingPrice = ProductUpdate.SellingPrice
	UpdatedProduct.ProductStock = ProductUpdate.ProductStock
	UpdatedProduct.NetWeight = ProductUpdate.NetWeight
	UpdatedProduct.Category = ProductUpdate.Category
	UpdatedProduct.ProductRank = ProductUpdate.ProductRank
	UpdatedProduct.Tags = ProductUpdate.Tags
	UpdatedProduct.PublishStatus = ProductUpdate.PublishStatus

	// Make sure product rank do not go below 1 or above 5
	if UpdatedProduct.ProductRank < 1 {
		UpdatedProduct.ProductRank = 1
	} else if UpdatedProduct.ProductRank > 5 {
		UpdatedProduct.ProductRank = 5
	}

	// Handle OnSale
	if UpdatedProduct.InitialCost != UpdatedProduct.SellingPrice {
		UpdatedProduct.OnSale = true
	} else {
		UpdatedProduct.OnSale = false
	}

	// Handle Featured
	if ProductUpdate.FeaturedStatus == "featured" {
		UpdatedProduct.Featured = true
	} else {
		UpdatedProduct.Featured = false
	}

	// Now intelligently update StockRemaining:
	stockDifference := ProductUpdate.ProductStock - existingProduct.ProductStock
	UpdatedProduct.StockRemaining = existingProduct.StockRemaining + stockDifference

	// Make sure StockRemaining doesn't go negative
	if UpdatedProduct.StockRemaining < 0 {
		UpdatedProduct.StockRemaining = 0
	}

	// Save update
	updateErr := dbFunc.DBHelper.UpdateBusinessConnectProduct(UpdatedProduct, uint(ProductUpdate.ProductID))
	if updateErr != nil {
		if errors.Is(updateErr, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "error retrieving product",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update product",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "successfully updated product",
	})
}

func DeleteBusinessConnectProduct(ctx *fiber.Ctx) error {

	// // Get stored user id from request timeline
	// userId := ctx.Locals("user-id")

	// if userId == nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "failed to get user",
	// 	})
	// }

	// _, uuidErr := dbFunc.DBHelper.FindByUuidFromLocal(userId)

	// if uuidErr != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": uuidErr.Error(),
	// 	})
	// }

	// Extract order ID from path or query parameter (adjust based on your implementation)
	productID, err := strconv.Atoi(ctx.Params("productID"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid order ID",
		})
	}

	// Call the database helper function to retrieve the order
	err2 := dbFunc.DBHelper.DeleteBusinessConnectProduct(uint(productID))
	if err2 != nil {
		// Handle potential errors based on your GetOrder function implementation
		if errors.Is(err2, gorm.ErrRecordNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "product not found",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve product",
		})
	}

	// Return successful response with order details
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": "successfully deleted product",
	})
}

func ConvertToMetadataFromType(responseData interface{}) (interface{}, error) {
	// Check if the input is a struct
	val := reflect.ValueOf(responseData)
	if val.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct")
	}

	// Create a new instance of Metadata
	responseDataBody := reflect.New(reflect.TypeOf(Data.ServiceMetaData{})).Elem()

	// Iterate over the fields of the input struct and copy them to Metadata
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		responseDataBody.FieldByName(fieldName).Set(val.Field(i))
	}

	return responseDataBody.Interface(), nil
}

func ConvertToMetadata(responseDataType interface{}) (Data.ServiceMetaData, error) {

	var responseDataBody Data.ServiceMetaData

	// Use reflection to get the type and value of the struct
	dataType := reflect.TypeOf(responseDataType)
	dataValue := reflect.ValueOf(responseDataType)

	fieldNamesInterface := GetFieldNamesFromInterface(responseDataType)

	// Iterate through the fields of the struct
	for i := 0; i < dataType.NumField(); i++ {
		// field := dataType.Field(i)
		fieldValue := dataValue.Field(i).Interface()

		// Check if the field name exists in the list of field names to check
		for _, nameToCheck := range Data.FieldNames {
			if fieldNamesInterface[i] == nameToCheck {
				// Field name found, print the field name and value
				// fmt.Printf("field name: %s fields value: %v\n", field.Name, fieldValue)
				// all need to have field values
				switch fieldNamesInterface[i] {
				case "TransactionID":
					responseDataBody.TransactionID = fieldValue.(string)
				case "Price":
					responseDataBody.Price = fieldValue.(int)
				case "Status":
					responseDataBody.Status = fieldValue.(string)
				case "PhoneNumber":
					responseDataBody.PhoneNumber = fieldValue.(string)
				case "EmailID":
					responseDataBody.EmailID = fieldValue.(string)
				default:
					fmt.Println("field name do not exist: ", fieldNamesInterface[i], "this is the field value", fieldValue)
				}
				// break
			}
		}
	}

	return responseDataBody, nil
}

func GetFieldNamesFromInterface(data interface{}) []string {
	var fieldNames []string

	// Use reflection to get the type of the interface
	dataType := reflect.TypeOf(data)

	// If the type is a pointer, get the underlying element type
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}

	// Iterate through the fields of the type
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)

		// Append the field name to the slice
		fieldNames = append(fieldNames, field.Name)
	}

	return fieldNames
}
