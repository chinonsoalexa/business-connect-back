package paystack

import (
	conn "business-connect/database"
	Dataa "business-connect/models"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// Define the interface with function signatures
type PaystackImpl interface {
	FindByUuidFromLocal(ID interface{}) (user Dataa.User, err error)
	FindByUuidFromLocalPaystack(ID interface{}) (user Dataa.User, err error)
	FindByEmail(Email string) (user Dataa.User, err error)
}

// Define a struct that implements the interface
type PaystackHelperImpl struct{}

// Create an instance of the struct to use as your Database helper
var PaystackHelper PaystackImpl = &PaystackHelperImpl{}

func (d *PaystackHelperImpl) FindByUuidFromLocal(ID interface{}) (user Dataa.User, err error) {

	if ID == "" || ID == nil {
		return Dataa.User{}, errors.New("error getting user id from request")
	}

	// Convert userId to a uint
	userIdUint, ok := ID.(uint)
	if !ok {
		return Dataa.User{}, errors.New("user id is not a valid string")
	}

	result := conn.DB.First(&user, "ID = ?", userIdUint)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return Dataa.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Dataa.User{}, errors.New("error retrieving user")
		}
	}

	return
}

func (d *PaystackHelperImpl) FindByUuidFromLocalPaystack(ID interface{}) (user Dataa.User, err error) {
	// Check if ID is empty or nil
	if ID == "" || ID == nil {
		return Dataa.User{}, errors.New("error getting user id from request")
	}

	// Convert ID to a string directly
	userIdStr := fmt.Sprint(ID)

	// Convert string userIdStr to a uint directly
	userIdUint, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return Dataa.User{}, errors.New("error converting string to uint")
	}

	// Proceed with the uint userIdUint
	result := conn.DB.First(&user, "ID = ?", userIdUint)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified ID was not found
			return Dataa.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Dataa.User{}, errors.New("error retrieving user")
		}
	}

	return user, nil
}

func (d *PaystackHelperImpl) FindByEmail(Email string) (user Dataa.User, err error) {
	result := conn.DB.First(&user, "email = ?", Email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return Dataa.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Dataa.User{}, errors.New("error retrieving user")
		}
	}

	return
}
