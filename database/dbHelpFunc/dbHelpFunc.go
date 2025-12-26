package dbHelpFunc

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	random "business-connect/controllers/authentication/utils"
	conn "business-connect/database"
	Data "business-connect/models"
)

// Define the interface with function signatures
type DatabaseHelper interface {
	FindByEmail(email string) (user Data.User, err error)
	FindByEmailExcludingUser(userToSendFunds uint, SendToEmail string) (user Data.User, err error)
	FindByUuid(ID uint) (user Data.User, err error)
	FindByUuidFromLocal(ID interface{}) (user Data.User, err error)
	FindByPhoneNumber(number string) (user Data.User, err error)
	CheckIfTransactionIDExist(ID uint) (err error)
	CheckByUserByEmail(NewEmail, OldEmail string) (err error)
	CreatePasswordHash(password string) (string, error)
	ComparePasswordHash(hash string, password string) error
	CreateNewUser(NewUser Data.User) (CreatedUser Data.User, err error)
	DeleteUser(uuid uint) (err error)
	FindByJti(jti string) (string, error)
	StoreRefreshToken() (jti string, err error)
	GetAndDeleteRefreshToken(uuidStr interface{}) error
	DeleteRefreshToken(jti string) (err error)
	CheckRefreshToken(jti string) bool
	UpdateUser(NewData Data.User) error
	CheckSpecialCharacters(email string) bool
	EditProfileText(text Data.User, userID uint) (string, error)
	CreateOTP(OTP Data.OTP) error
	GetOTPByEmail(email string) (emailOTPBody Data.OTP, err error)
	GetOTPByNumber(number string) (numberOTPBody Data.OTP, err error)
	GetAndCheckOTPByEmail(Email, OTP string) (userOTP Data.OTP, err error)
	GetAndCheckOTPByNumber(number, OTP string) (userOTP Data.OTP, err error)
	UpdateExistingOTP(OTPBody Data.OTP, CustomID uint) (err error)
	DeleteExistingOTPByID(ID uint) (err error)
	UpdateMaxTry(Email string) (err error)
	UpdateMaxTryNumber(number string) (err error)
	UpdateMaxTryToZero(Email string) (err error)
	GetBusinessConnectProductsByLimit(limit, offset int) ([]Data.Post, bool, error)
	GetBusinessConnectProductsByLimit2( /*userID uint64, */ fingerprintHash string, limit, offset int) ([]Data.Post, int64, error)
	GetProductsAll(limit, offset int, sortField, sortOrder string) ([]Data.Post, int64, error)
	GetProductsByCategory(category string, limit, offset int, sortField, sortOrder string) ([]Data.Post, int64, error)
	GetBusinessConnectAdminProductsByLimit( /*userID uint64, */ limit, offset int) ([]Data.Post, int64, error)
	// GetBusinessConnectRecommendedProductsByLimit( /*userID uint64, */ category string, limit int) ([]Data.Post, int64, error)
	GetBusinessConnectRecommendedProductsByLimit(currentProductID uint64, category string, limit int) ([]Data.Post, int64, error)
	GetBusinessConnectBlogByLimit( /*userID uint64, */ limit, offset int) ([]Data.Blog, int64, error)
	GetBusinessConnectHomeAllProductsByLimit(limit int) ([]Data.Post, error)
	GetBusinessConnectHomeFeaturedProductsByLimit(limit int) ([]Data.Post, error)
	GetBusinessConnectHomeBestSellingProductsByLimit(limit int) ([]Data.Post, error)
	GetBusinessConnectHomeOnSaleProductsByLimit(limit int) ([]Data.Post, error)
	GetBusinessConnectProductsByIDs(productIDs []uint64) ([]Data.Post, error)
	GetBusinessConnectProductByID(productID uint64) (Data.Post, error)
	GetProductByID(productID uint) (Data.Post, error)
	GetNextProductID(currentID uint64) (uint64, error)
	GetPreviousProductID(currentID uint64) (uint64, error)
	SearchProductsByTitle(searchTerm string) ([]struct {
		Title        string `json:"title"`
		ProductUrlID string `json:"product_url_id"`
		Category     string `json:"category"`
	}, error)
	SearchProductsByTitleAndCategory(searchTerm string, categorySlug string) ([]Data.ProductSearchResponse, error)
	SearchAdminProductsByTitle(searchTerm string) ([]Data.Post, error)
	SearchAdminOrderByTitle(searchTerm string) ([]Data.OrderHistory, error)
	GetTransactionsByUserAndDateWithLimit(userID uint, dateString string) ([]Data.Post, error)
	GetTransactionHistoryForAi(userID uint64, limit int) ([]Data.Post, int64, error)
	AddEmailSubscriber(Email string) error
	UpdateSiteVisits() error
	GetLast12DaysSiteVisits() (map[string]int64, error)
	UpdateSubscriptionStatus(transactionID uint64) error
	AddProduct(post Data.Post, user Data.User) (Data.Post, error)
	AddProductImage(image Data.PostImage, postID uint) error
	AddBlog(post Data.Blog, user Data.User) (Data.Blog, error)
	AddBlogImage(image Data.BlogImage, postID uint) error
	// SaveCustomerReview(productID uint, email string, name string, reviewText string, rating int) (Data.CustomerReview, error)
	// GetCustomerReviewsByProduct(productID uint, limit int, offset int) ([]Data.CustomerReview, int64, error)
	SaveCustomerBlogReview(blogID uint, email string, name string, reviewText string, rating int) (Data.CustomerBlogReview, error)
	GetCustomerBlogReviewsByBlogPost(blogID uint, limit int, offset int) ([]Data.CustomerBlogReview, int64, error)
	GetBlogPostById(blogID uint) (*Data.Blog, error)
	AddOrder(orderHistoryBody Data.OrderHistoryBody, ordersBody []Data.ProductOrderBody) (uint, *Data.OrderHistory, []Data.ProductOrder, error)
	// UpsertShippingFee(fee int64, feesGreater, feesLess int64) (error)
	UpsertShippingFee(fee int64, feesGreater, feesLess int64, storeLatitude, storeLongitude float64, storeState,
		storeCity, stateISO string, calculateUsingKg bool) error
	GetShippingFee() (Data.ShippingFees, error)
	GetOrder(orderID uint) (*Data.OrderHistory, error)
	GetAndUpdateOrder(orderID uint, status string) (*Data.OrderHistory, error)
	GetBusinessConnectOrdersByLimit(limit, offset int) ([]Data.OrderHistory, int64, error)
	UpdateOrderStatus(orderID uint, newStatus string) error
	GetAnalyticsData() (*Data.Analytics, error)
	UpdateBusinessConnectProduct(Post Data.Post, ProductID uint) error
	UpdateBusinessConnectBlog(Post Data.Blog, BlogID uint) error
	DeleteBusinessConnectProduct(ProductID uint) error
	DeleteBusinessConnectBlog(BlogID uint) error
	GetBusinessConnectEmailSubscribers() ([]Data.BusinessConnectEmailSubscriber, error)
	SaveBusinessConnectSentEmail(sentEmail Data.Email) error
	GetBusinessConnectUniqueUserFingerPrintHash(fingerprintHash string) (Data.BusinessConnectDeviceFingerprint, error)
	CreateBusinessConnectDeviceFingerprint(fingerprintHash string) error
	RecommendProductsForUser(fingerprintHash string, limit, offset int) ([]Data.Post, error)
	LogUserClickData(fingerprintHash string, productID uint, ActivityType, Category, TitleOrSearchQuery string) error
}

// Define a struct that implements the interface
type DatabaseHelperImpl struct{}

// Create an instance of the struct to use as your database helper
var DBHelper DatabaseHelper = &DatabaseHelperImpl{}

const (
	PostTypePersonal = "personal"
	PostTypeBusiness = "business"
	PostTypeGroup    = "group"
	PostTypeEvent    = "event"
	PostTypeAd       = "ad"

	EntryFree = "free"
	EntryPaid = "paid"
)


// the findByEmail() function accepts an email as an argument
// and return a record of any found user else it returns an error
// and returns an empty struct of data
func (d *DatabaseHelperImpl) FindByEmail(email string) (user Data.User, err error) {
	result := conn.DB.First(&user, "email = ?", email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return Data.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Data.User{}, errors.New("error retrieving user")
		}
	}

	return
}

func (d *DatabaseHelperImpl) FindByEmailExcludingUser(userToSendFunds uint, SendToEmail string) (Data.User, error) {
	var user Data.User
	result := conn.DB.Where("email = ? AND id != ?", SendToEmail, userToSendFunds).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Data.User{}, errors.New("user not found")
		}
		return Data.User{}, errors.New("error retrieving user")
	}

	return user, nil
}

// the FindByUuid() function accepts an ID as an argument
// and return a record of any found user else it returns an error
// and returns an empty struct of data
func (d *DatabaseHelperImpl) FindByUuid(ID uint) (user Data.User, err error) {
	result := conn.DB.First(&user, "ID = ?", ID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return Data.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Data.User{}, errors.New("error retrieving user")
		}
	}

	return user, nil
}

func (d *DatabaseHelperImpl) FindByUuidFromLocal(ID interface{}) (user Data.User, err error) {

	if ID == "" || ID == nil {
		return Data.User{}, errors.New("error getting user id from request")
	}

	// Convert userId to a uint
	userIdUint, ok := ID.(uint)
	if !ok {
		return Data.User{}, errors.New("user id is not a valid string")
	}

	result := conn.DB.First(&user, "ID = ?", userIdUint)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return Data.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Data.User{}, errors.New("error retrieving user")
		}
	}

	return
}

func (d *DatabaseHelperImpl) FindByPhoneNumber(number string) (user Data.User, err error) {
	result := conn.DB.First(&user, "phone_number = ?", number)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return Data.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Data.User{}, errors.New("error retrieving user")
		}
	}

	return
}

func (d *DatabaseHelperImpl) CheckIfTransactionIDExist(ID uint) (err error) {
	var Details Data.Post
	result := conn.DB.First(&Details, "ID = ?", ID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("user not found")
		} else {
			// Some other error occurred
			return errors.New("error retrieving user")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) CheckByUserByEmail(NewEmail, OldEmail string) (err error) {
	var (
		email1 Data.User
		email2 Data.User
	)
	email1, err = d.FindByEmail(NewEmail)
	if err != nil {
		return
	}

	email2, err = d.FindByEmail(OldEmail)
	if err != nil {
		return
	}

	// checking if user already exist in the database
	if email1.Email == email2.Email {
		return
	}

	return
}

func (d *DatabaseHelperImpl) CreatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash[:]), err
}

func (d *DatabaseHelperImpl) ComparePasswordHash(hash string, password string) error {
	// Compare password hash to check if password is valid
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	// Check if there was an error comparing the hash and password
	if err != nil {
		// bcrypt.ErrMismatchedHashAndPassword indicates that the password hash and password do not match
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("password does not match")
		} else if errors.Is(err, bcrypt.ErrHashTooShort) {
			return errors.New("password too short to be a bcrypt password")
		}
		// Return the original error for any other bcrypt-related errors
		return err
	}

	// Password matches
	return nil
}

func (d *DatabaseHelperImpl) CreateNewUser(NewUser Data.User) (CreatedUser Data.User, err error) {
	var (
		NewPassword string
	)

	// ⚠️ Hash password here before save
	NewPassword, err = d.CreatePasswordHash(NewUser.Password)
	if err != nil {
		return Data.User{}, err
	}

	NewUser.EmailVerified = false
	// NewUser.GoogleAuth = false
	NewUser.Password = NewPassword

	// creating new users in the data base
	result := conn.DB.Create(&NewUser)

	// checks if an error occurred when creating a user
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return Data.User{}, errors.New("user not found")
		} else {
			// Some other error occurred
			return Data.User{}, errors.New("error creating user")
		}
	}

	// The created user is already updated in the NewUser object
	CreatedUser = NewUser

	return CreatedUser, nil
}

func (d *DatabaseHelperImpl) DeleteUser(uuid uint) (err error) {
	var (
		user Data.User
	)

	user, err = d.FindByUuid(uuid)
	if err != nil {
		return
	}

	result := conn.DB.Delete(&Data.User{}, user.ID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("user not found")
		} else {
			// Some other error occurred
			return errors.New("error retrieving user")
		}
	}

	return
}

func (d *DatabaseHelperImpl) FindByJti(jti string) (string, error) {
	var oldJti Data.JTI
	result := conn.DB.Where("jti = ?", jti).First(&oldJti)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("jti not found")
		}
		return "", errors.New("error retrieving user")
	}

	return oldJti.Jti, nil
}

func (d *DatabaseHelperImpl) StoreRefreshToken() (jti string, err error) {
	jti, err = random.RandomAlphanumericString(32)
	if err != nil {
		return "", fmt.Errorf("error generating jti: %v", err)
	}

	// Ensure jti is unique
	for {
		if _, err := d.FindByJti(jti); err != nil {
			break
		}
		jti, err = random.RandomAlphanumericString(32)
		if err != nil {
			return "", fmt.Errorf("error generating jti: %v", err)
		}
	}

	// uintValue, err := strconv.ParseUint(uuid, 10, 64)
	// if err != nil {
	// 	log.Printf("Error converting string to uint: %v", err)
	// 	return "", fmt.Errorf("error converting string to uint: %v", err)
	// }

	newJti := Data.JTI{
		Jti: jti,
		// UserID: uint(uintValue),
	}

	if err := conn.DB.Create(&newJti).Error; err != nil {
		return "", fmt.Errorf("error creating jti token: %v", err)
	}

	return jti, nil
}

func (d *DatabaseHelperImpl) GetAndDeleteRefreshToken(uuidStr interface{}) error {
	var userID uint

	// Check the type of uuidStr and convert if necessary
	switch v := uuidStr.(type) {
	case string:
		uintValue, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			log.Printf("Error converting string to uint: %v", err)
			return fmt.Errorf("error converting string to uint: %v", err)
		}
		userID = uint(uintValue)
	case uint:
		userID = v
	default:
		return fmt.Errorf("unsupported type for uuidStr: %T", v)
	}

	// Find all JTI records by UserID
	var jtiRecords []Data.JTI
	if err := conn.DB.Where("user_id = ?", userID).Find(&jtiRecords).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			log.Printf("Records not found for UUID: %d", userID)
			return fmt.Errorf("records not found for UUID: %d", userID)
		}
		log.Printf("Error finding records: %v", err)
		return fmt.Errorf("error finding records: %v", err)
	}

	if len(jtiRecords) == 0 {
		log.Printf("No records found for UUID: %d", userID)
		return fmt.Errorf("no records found for UUID: %d", userID)
	}

	// Hard delete the JTI records
	if err := conn.DB.Unscoped().Delete(&jtiRecords).Error; err != nil {
		log.Printf("Error deleting records: %v", err)
		return fmt.Errorf("error deleting records: %v", err)
	}

	return nil
}

func (d *DatabaseHelperImpl) DeleteRefreshToken(jti string) error {
	result := conn.DB.Unscoped().Where("jti = ?", jti).Delete(&Data.JTI{})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("refresh token not found")
		}
		return errors.New("error deleting refresh token")
	}

	return nil
}

func (d *DatabaseHelperImpl) CheckRefreshToken(jti string) bool {
	_, err := d.FindByJti(jti)
	return err == nil
}

func (d *DatabaseHelperImpl) CheckSpecialCharacters(email string) bool {
	// Define a regular expression pattern to match special characters
	pattern := `[!#$%&'*+/=?^_{|}~]`

	// Use regexp.MatchString to check if any special characters are present in the email
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		// errors.New("error while matching pattern")
		return false
	}
	return match
}

func (d *DatabaseHelperImpl) EditProfileText(text Data.User, userID uint) (string, error) {
	// save updated user
	result := conn.DB.Save(text)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return "", errors.New("user not found")
		} else {
			// Some other error occurred
			return "", errors.New("error retrieving user")
		}
	}

	return "user profile updated successfully", nil
}

func (d *DatabaseHelperImpl) UpdateUser(NewData Data.User) error {

	result := conn.DB.Save(NewData)

	// checks if an error occurred when creating a user
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("user to update not found")
		} else {
			// Some other error occurred
			return errors.New("error retrieving user to update")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) CreateOTP(OTP Data.OTP) error {
	hashOTP, otpErr := d.CreateOTPHash(OTP.OTP)
	if otpErr != nil {
		return errors.New("error creating otp hash")
	}
	// save hashed otp to the database
	OTP.OTP = hashOTP
	// creating new users in the data base
	result := conn.DB.Create(&OTP)

	// checks if an error occurred when creating a user
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("otp not found")
		} else {
			// Some other error occurred
			return errors.New("error retrieving otp")
		}
	}

	return nil
}

// get the users otp by an existing email Note: email ID's cannot be the same so each
// email ID is mapped to a otp value for verification purposes
func (d *DatabaseHelperImpl) GetOTPByEmail(email string) (emailOTPBody Data.OTP, err error) {

	result := conn.DB.First(&emailOTPBody, "email = ?", email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			log.Printf("OTP not found for email: %s", email)
			return Data.OTP{}, errors.New("otp not found by email")
		} else {
			// Some other error occurred
			log.Printf("Error retrieving OTP: %v", result.Error)
			return Data.OTP{}, errors.New("error retrieving otp")
		}
	}

	// log.Printf("OTP fetched from the database: %+v", emailOTPBody)

	return emailOTPBody, nil
}

func (d *DatabaseHelperImpl) GetOTPByNumber(number string) (numberOTPBody Data.OTP, err error) {

	result := conn.DB.First(&numberOTPBody, "phone_number = ?", number)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			log.Printf("OTP not found for number: %s", number)
			return Data.OTP{}, errors.New("otp not found by number")
		} else {
			// Some other error occurred
			log.Printf("Error retrieving OTP: %v", result.Error)
			return Data.OTP{}, errors.New("error retrieving otp")
		}
	}

	return numberOTPBody, nil
}

func (d *DatabaseHelperImpl) GetAndCheckOTPByEmail(Email, OTP string) (userOTP Data.OTP, err error) {
	userOTP, otpErr := d.GetOTPByEmail(Email)
	if otpErr != nil {
		if otpErr.Error() == "otp not found by email" {
			return Data.OTP{}, errors.New("otp not found")
		}
		return Data.OTP{}, errors.New("an error occurred while getting otp email")
	}

	// Create a copy of the original userOTP for logging purposes
	originalUserOTP := userOTP
	hashErr := d.CompareOTPHash(originalUserOTP.OTP, OTP)

	if hashErr != nil {
		log.Println("Comparison failed: ", hashErr)
		if hashErr.Error() == "incorrect OTP" {
			return Data.OTP{}, errors.New("incorrect otp value")
		} else if hashErr.Error() == "hashed OTP is too short to be a bcrypt hash" {
			return Data.OTP{}, errors.New("hashed OTP is too short to be a bcrypt hash")
		}
	}

	return userOTP, nil
}

func (d *DatabaseHelperImpl) GetAndCheckOTPByNumber(number, OTP string) (userOTP Data.OTP, err error) {
	userOTP, otpErr := d.GetOTPByNumber(number)
	if otpErr != nil {
		if otpErr.Error() == "otp not found by email" {
			return Data.OTP{}, errors.New("otp not found")
		}
		return Data.OTP{}, errors.New("an error occurred while getting otp email")
	}

	// Create a copy of the original userOTP for logging purposes
	originalUserOTP := userOTP
	hashErr := d.CompareOTPHash(originalUserOTP.OTP, OTP)

	if hashErr != nil {
		log.Println("Comparison failed: ", hashErr)
		if hashErr.Error() == "incorrect OTP" {
			return Data.OTP{}, errors.New("incorrect otp value")
		} else if hashErr.Error() == "hashed OTP is too short to be a bcrypt hash" {
			return Data.OTP{}, errors.New("hashed OTP is too short to be a bcrypt hash")
		}
	}

	return userOTP, nil
}

func (d *DatabaseHelperImpl) UpdateExistingOTP(OTPBody Data.OTP, CustomID uint) (err error) {
	// Find the existing record by CustomID
	var existingOTP Data.OTP
	if result := conn.DB.Where("custom_id = ?", CustomID).First(&existingOTP).Error; result != nil {
		if errors.Is(result, gorm.ErrRecordNotFound) {
			// The record with the specified id was not found
			return errors.New("otp not found")
		} else {
			// Some other error occurred
			return errors.New("error updating otp")
		}
	}

	oTPBodyHash, hashErr := d.CreateOTPHash(OTPBody.OTP)
	if hashErr != nil {
		return errors.New("error creating otp hash")
	}

	// Update only the OTP and CreatedAt fields
	existingOTP.OTP = oTPBodyHash
	existingOTP.CreatedAT = OTPBody.CreatedAT
	existingOTP.MaxTry = OTPBody.MaxTry

	// Save the updated record
	if err := conn.DB.Save(&existingOTP).Error; err != nil {
		return errors.New("error saving the updated otp")
	}

	return nil
}

func (d *DatabaseHelperImpl) UpdateMaxTry(Email string) (err error) {

	userOTP, otpErr := d.GetOTPByEmail(Email)
	if otpErr != nil {
		if otpErr.Error() == "otp not found by email" {
			return errors.New("otp not found")
		}
		return errors.New("an error occurred while getting otp email")
	}

	userOTP.MaxTry += 1
	maxOtpToUpdate := &Data.OTP{
		CustomID: userOTP.CustomID,
		MaxTry:   userOTP.MaxTry, // The new max_try value
	}

	// Update specific fields of the OTP record based on the email
	result := conn.DB.Model(&userOTP).Updates(maxOtpToUpdate)

	// Check if an error occurred
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("otp not found")
		} else {
			// Some other error occurred
			return errors.New("error updating max_try")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) UpdateMaxTryNumber(number string) (err error) {

	userOTP, otpErr := d.GetOTPByNumber(number)
	if otpErr != nil {
		if otpErr.Error() == "otp not found by email" {
			return errors.New("otp not found")
		}
		return errors.New("an error occurred while getting otp email")
	}

	userOTP.MaxTry += 1
	maxOtpToUpdate := &Data.OTP{
		CustomID: userOTP.CustomID,
		MaxTry:   userOTP.MaxTry, // The new max_try value
	}

	// Update specific fields of the OTP record based on the email
	result := conn.DB.Model(&userOTP).Updates(maxOtpToUpdate)

	// Check if an error occurred
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("otp not found")
		} else {
			// Some other error occurred
			return errors.New("error updating max_try")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) UpdateMaxTryToZero(Email string) (err error) {

	userOTP, otpErr := d.GetOTPByEmail(Email)
	if otpErr != nil {
		if otpErr.Error() == "otp not found by email" {
			return errors.New("otp not found")
		}
		return errors.New("an error occurred while getting otp email")
	}

	userOTP.MaxTry = 0
	maxOtpToUpdate := &Data.OTP{
		MaxTry: userOTP.MaxTry, // The new max_try value
	}

	// Update specific fields of the OTP record based on the id
	result := conn.DB.Model(&userOTP).Updates(maxOtpToUpdate)

	// Check if an error occurred
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("otp not found")
		} else {
			// Some other error occurred
			return errors.New("error updating max_try")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) DeleteExistingOTPByID(ID uint) (err error) {
	result := conn.DB.Where("custom_id = ?", ID).Delete(&Data.OTP{})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("otp not found to delete")
		} else {
			// Some other error occurred
			return errors.New("error deleting otp")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) CreateOTPHash(OTP string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(OTP), bcrypt.DefaultCost)
	return string(hash[:]), err
}

func (d *DatabaseHelperImpl) CompareOTPHash(hash string, plainOTP string) error {

	hash = strings.TrimSpace(hash)
	plainOTP = strings.TrimSpace(plainOTP)

	// Compare the hashed OTP to the plain OTP
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainOTP))

	// Check if there was an error comparing the hash and plain OTP
	if err != nil {
		// bcrypt.ErrMismatchedHashAndPassword indicates that the hashed OTP and plain OTP do not match
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("incorrect OTP")
		} else if errors.Is(err, bcrypt.ErrHashTooShort) {
			return errors.New("hashed OTP is too short to be a bcrypt hash")
		}
		// Return the original error for any other bcrypt-related errors
		return err
	}

	// OTP matches
	return nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectProductsByLimit(
	limit, offset int,
) ([]Data.Post, bool, error) {

	var posts []Data.Post

	allowedPostTypes := []string{
		PostTypePersonal,
		PostTypeBusiness,
		PostTypeGroup,
		PostTypeEvent,
		PostTypeAd,
	}

	result := conn.DB.
		Model(&Data.Post{}).
		Preload("Images").
		Where(`
			is_active = ? 
			AND approved = ? 
			AND post_type IN ?
		`, true, true, allowedPostTypes).
		Order("created_at DESC").
		Limit(limit + 1).
		Find(&posts)

	if result.Error != nil {
		return nil, false, result.Error
	}

	hasMore := false
	if len(posts) > limit {
		hasMore = true
		posts = posts[:limit]
	}

	return posts, hasMore, nil
}


// func (d *DatabaseHelperImpl) GetBusinessConnectProductsByLimit(
// 	limit,
// 	offset int,
// ) ([]Data.Post, int64, error) {

// 	var postRecordsCount int64
// 	var posts []Data.Post

// 	result := conn.DB.
// 		Preload("Images").
// 		// Where("is_active = ? AND approved = ?", true, true).
// 		Order("created_at DESC").
// 		Limit(limit).
// 		Offset(offset).
// 		Find(&posts)

// 	if result.Error != nil {
// 		return []Data.Post{}, 0, result.Error
// 	}

// 	conn.DB.
// 		Model(&Data.Post{}).
// 		// Where("is_active = ? AND approved = ?", true, true).
// 		Count(&postRecordsCount)

// 	return posts, postRecordsCount, nil
// }

func (d *DatabaseHelperImpl) GetBusinessConnectProductsByLimit2( /*userID uint64, */ fingerprintHash string, limit, offset int) ([]Data.Post, int64, error) {
	var productRecords []Data.Post
	var productRecordsCount int64

	// Get the count of transaction records for the user
	if err := conn.DB.Model(&Data.Post{}).Count(&productRecordsCount).Error; err != nil {
		return []Data.Post{}, 0, err
	}

	productRecords, productErr := d.RecommendProductsForUser(fingerprintHash, limit, offset)
	if productErr != nil {
		// fmt.Println("product returned error: ", productErr)
		// Retrieve transaction history for the user with pagination
		if productRecordsErr := conn.DB. /*Where("user_id = ?", userID).*/ Order("created_at desc").Limit(limit).Offset(offset).Find(&productRecords).Error; productRecordsErr != nil {
			if errors.Is(productRecordsErr, gorm.ErrRecordNotFound) {
				// The record with the specified UserID was not found
				return []Data.Post{}, 0, errors.New("no transaction record found")
			} else {
				// Some other error occurred
				return []Data.Post{}, 0, errors.New("error retrieving transaction record")
			}
		}
	}
	// fmt.Println("product did not return an error: ", productErr)

	return productRecords, productRecordsCount, nil
}

func (d *DatabaseHelperImpl) GetProductsAll(limit, offset int, sortField, sortOrder string) ([]Data.Post, int64, error) {
	var products []Data.Post
	var count int64

	query := conn.DB.Model(&Data.Post{})

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (d *DatabaseHelperImpl) GetProductsByCategory(category string, limit, offset int, sortField, sortOrder string) ([]Data.Post, int64, error) {
	var products []Data.Post
	var count int64

	query := conn.DB.Model(&Data.Post{}).Where("category = ?", category)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectAdminProductsByLimit( /*userID uint64, */ limit, offset int) ([]Data.Post, int64, error) {
	var productRecords []Data.Post
	var productRecordsCount int64

	// Get the count of transaction records for the user
	if err := conn.DB.Model(&Data.Post{}).Count(&productRecordsCount).Error; err != nil {
		return []Data.Post{}, 0, err
	}

	// Retrieve transaction history for the user with pagination
	if productRecordsErr := conn.DB. /*Where("user_id = ?", userID).*/ Order("created_at desc").Limit(limit).Offset(offset).Find(&productRecords).Error; productRecordsErr != nil {
		if errors.Is(productRecordsErr, gorm.ErrRecordNotFound) {
			// The record with the specified UserID was not found
			return []Data.Post{}, 0, errors.New("no transaction record found")
		} else {
			// Some other error occurred
			return []Data.Post{}, 0, errors.New("error retrieving transaction record")
		}
	}

	return productRecords, productRecordsCount, nil
}

// func (d *DatabaseHelperImpl) GetBusinessConnectRecommendedProductsByLimit( /*userID uint64, */ category string, limit int) ([]Data.Post, int64, error) {
// 	var productRecords []Data.Post
// 	var productRecordsCount int64

// 	// Get the count of transaction records for the user
// 	if err := conn.DB.Model(&Data.Post{}).Count(&productRecordsCount).Error; err != nil {
// 		return []Data.Post{}, 0, err
// 	}

// 	// Retrieve transaction history for the user with pagination
// 	if productRecordsErr := conn.DB. /*Where("user_id = ?", userID).*/ Where("category = ?", category).Order("created_at desc").Limit(limit).Find(&productRecords).Error; productRecordsErr != nil {
// 		if errors.Is(productRecordsErr, gorm.ErrRecordNotFound) {
// 			// The record with the specified UserID was not found
// 			return []Data.Post{}, 0, errors.New("no transaction record found")
// 		} else {
// 			// Some other error occurred
// 			return []Data.Post{}, 0, errors.New("error retrieving transaction record")
// 		}
// 	}

// 	return productRecords, productRecordsCount, nil
// }

func (d *DatabaseHelperImpl) GetBusinessConnectRecommendedProductsByLimit(currentProductID uint64, category string, limit int) ([]Data.Post, int64, error) {
	var productRecordsCount int64
	var productRecords []Data.Post
	// var fallbackProducts []Data.Post
	// var latestProducts []Data.Post

	// 1. Try to get related products in the same category, excluding current product
	err := conn.DB.Model(&Data.Post{}).
		Where("category = ? AND id != ? AND publish_status = ?", category, currentProductID, "publish").
		Order("created_at DESC").
		Limit(limit).
		Find(&productRecords).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}

	// 2. If less than the limit, try to fetch products from "related" categories
	if len(productRecords) < limit {
		remaining := limit - len(productRecords)

		var fallbackProducts []Data.Post
		// This is where you can add custom logic to find "related" categories — for now we do a fuzzy match
		err = conn.DB.Model(&Data.Post{}).
			Where("LOWER(category) LIKE ? AND id != ? AND publish_status = ?", "%"+category+"%", currentProductID, "publish").
			Where("id NOT IN (?)", getProductIDs(productRecords)).
			Order("created_at DESC").
			Limit(remaining).
			Find(&fallbackProducts).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, err
		}

		productRecords = append(productRecords, fallbackProducts...)
	}

	// 3. Final fallback: fill remaining with latest published products
	if len(productRecords) < limit {
		remaining := limit - len(productRecords)

		var latestProducts []Data.Post
		err = conn.DB.Model(&Data.Post{}).
			Where("publish_status = ? AND id != ?", "publish", currentProductID).
			Where("id NOT IN (?)", getProductIDs(productRecords)).
			Order("created_at DESC").
			Limit(remaining).
			Find(&latestProducts).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, err
		}

		productRecords = append(productRecords, latestProducts...)
	}

	// Get the count for reference (optional)
	if countErr := conn.DB.Model(&Data.Post{}).Where("publish_status = ?", "publish").Count(&productRecordsCount).Error; countErr != nil {
		return productRecords, 0, countErr
	}

	return productRecords, productRecordsCount, nil
}

// Helper to extract IDs from slice
func getProductIDs(products []Data.Post) []uint64 {
	ids := make([]uint64, 0, len(products))
	for _, p := range products {
		ids = append(ids, uint64(p.ID))
	}
	// Avoid empty slice in SQL "IN (?)" clause
	if len(ids) == 0 {
		return []uint64{0} // or -1 for signed IDs
	}
	return ids
}

func (d *DatabaseHelperImpl) GetBusinessConnectBlogByLimit( /*userID uint64, */ limit, offset int) ([]Data.Blog, int64, error) {
	var productRecords []Data.Blog
	var productRecordsCount int64

	// Get the count of transaction records for the user
	if err := conn.DB.Model(&Data.Blog{}).Count(&productRecordsCount).Error; err != nil {
		return []Data.Blog{}, 0, err
	}

	// Retrieve transaction history for the user with pagination
	if productRecordsErr := conn.DB. /*Where("user_id = ?", userID).*/ Order("created_at desc").Limit(limit).Offset(offset).Find(&productRecords).Error; productRecordsErr != nil {
		if errors.Is(productRecordsErr, gorm.ErrRecordNotFound) {
			// The record with the specified UserID was not found
			return []Data.Blog{}, 0, errors.New("no transaction record found")
		} else {
			// Some other error occurred
			return []Data.Blog{}, 0, errors.New("error retrieving transaction record")
		}
	}

	return productRecords, productRecordsCount, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectHomeAllProductsByLimit(limit int) ([]Data.Post, error) {
	var productRecords []Data.Post

	err := conn.DB.
		Preload("ProductImages").
		Order("created_at desc").
		Limit(limit).
		Find(&productRecords).Error

	if err != nil {
		return []Data.Post{}, fmt.Errorf("error retrieving latest products: %w", err)
	}

	return productRecords, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectHomeFeaturedProductsByLimit(limit int) ([]Data.Post, error) {
	var featured []Data.Post

	err := conn.DB.
		Preload("ProductImages").
		Where("featured = ?", true).
		Order("created_at desc").
		Limit(limit).
		Find(&featured).Error

	if err != nil {
		return []Data.Post{}, fmt.Errorf("error retrieving featured products: %w", err)
	}

	remaining := limit - len(featured)
	if remaining > 0 {
		existingIDs := make([]uint, len(featured))
		for i, p := range featured {
			existingIDs[i] = p.ID
		}

		var fallback []Data.Post
		fallbackErr := conn.DB.
			Preload("ProductImages").
			Where("featured = ?", false).
			Where("id NOT IN ?", existingIDs).
			Order("created_at desc").
			Limit(remaining).
			Find(&fallback).Error

		if fallbackErr == nil {
			featured = append(featured, fallback...)
		}
	}

	return featured, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectHomeBestSellingProductsByLimit(limit int) ([]Data.Post, error) {
	var bestSelling []Data.Post

	err := conn.DB.
		Preload("ProductImages").
		Where("best_seller = ?", true).
		Order("sales desc").
		Limit(limit).
		Find(&bestSelling).Error

	if err != nil {
		return []Data.Post{}, fmt.Errorf("error retrieving best-selling products: %w", err)
	}

	remaining := limit - len(bestSelling)
	if remaining > 0 {
		existingIDs := make([]uint, len(bestSelling))
		for i, p := range bestSelling {
			existingIDs[i] = p.ID
		}

		var fallback []Data.Post
		fallbackErr := conn.DB.
			Preload("ProductImages").
			Where("best_seller = ?", false).
			Where("id NOT IN ?", existingIDs).
			Order("sales desc").
			Limit(remaining).
			Find(&fallback).Error

		if fallbackErr == nil {
			bestSelling = append(bestSelling, fallback...)
		}
	}

	return bestSelling, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectHomeOnSaleProductsByLimit(limit int) ([]Data.Post, error) {
	var onSale []Data.Post

	err := conn.DB.
		Preload("ProductImages").
		Where("on_sale = ?", true).
		Order("created_at desc").
		Limit(limit).
		Find(&onSale).Error

	if err != nil {
		return []Data.Post{}, fmt.Errorf("error retrieving on-sale products: %w", err)
	}

	remaining := limit - len(onSale)
	if remaining > 0 {
		existingIDs := make([]uint, len(onSale))
		for i, p := range onSale {
			existingIDs[i] = p.ID
		}

		var fallback []Data.Post
		fallbackErr := conn.DB.
			Preload("ProductImages").
			Where("on_sale = ?", false).
			Where("id NOT IN ?", existingIDs).
			Order("created_at desc").
			Limit(remaining).
			Find(&fallback).Error

		if fallbackErr == nil {
			onSale = append(onSale, fallback...)
		}
	}

	return onSale, nil
}

// Function to retrieve multiple products by their IDs
func (d *DatabaseHelperImpl) GetBusinessConnectProductsByIDs(productIDs []uint64) ([]Data.Post, error) {
	var products []Data.Post

	// Fetch all products with IDs in the given slice
	if err := conn.DB.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("error retrieving products: %v", err)
	}

	// Check if no products were found
	if len(products) == 0 {
		return nil, errors.New("no products found for the given IDs")
	}

	return products, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectProductByID(productID uint64) (Data.Post, error) {
	var productRecord Data.Post
	var reviewLimit = 4

	// Try to retrieve the product by its ID and preload a limited number of customer reviews
	if err := conn.DB.Where("id = ?", productID).
		Preload("CustomerReviews", func(db *gorm.DB) *gorm.DB {
			return db.Limit(reviewLimit) // Limit the number of reviews
		}).
		First(&productRecord).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Post not found, find the next closest available product ID
			nextProductID, findErr := d.findClosestAvailableProductID(productID)
			if findErr != nil {
				return Data.Post{}, findErr
			}

			// Retrieve the product by the closest available ID and preload a limited number of customer reviews
			if err := conn.DB.Where("id = ?", nextProductID).
				Preload("CustomerReviews", func(db *gorm.DB) *gorm.DB {
					return db.Limit(reviewLimit) // Limit the number of reviews
				}).
				First(&productRecord).Error; err != nil {
				return Data.Post{}, errors.New("error retrieving product record")
			}
		} else {
			// Some other error occurred
			return Data.Post{}, errors.New("error retrieving product record")
		}
	}

	return productRecord, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectProductByIDd(productID uint64) (Data.Post, error) {
	var productRecord Data.Post
	var reviewLimit = 4

	// Try to retrieve the product by its ID and preload a limited number of customer reviews
	if err := conn.DB.Where("id = ?", productID).
		Preload("CustomerReviews", func(db *gorm.DB) *gorm.DB {
			return db.Limit(reviewLimit) // Limit the number of reviews
		}).
		First(&productRecord).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Post not found, find the next closest available product ID
			return Data.Post{}, errors.New("error retrieving product record")
		} else {
			// Some other error occurred
			return Data.Post{}, errors.New("error retrieving product record")
		}
	}

	return productRecord, nil
}

func (d *DatabaseHelperImpl) GetProductByID(productID uint) (Data.Post, error) {
	var product Data.Post
	err := conn.DB.First(&product, productID).Error
	return product, err
}

func (d *DatabaseHelperImpl) findClosestAvailableProductID(targetID uint64) (Data.Post, error) {
	var closestProduct Data.Post

	// Look for the next higher product ID and preload CustomerReviews
	if err := conn.DB.Where("id > ?", targetID).
		Order("id ASC").
		Preload("CustomerReviews"). // Preload customer reviews
		Limit(1).
		First(&closestProduct).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no higher product ID is found, look for the closest lower ID
			if err := conn.DB.Where("id < ?", targetID).
				Order("id DESC").
				Preload("CustomerReviews"). // Preload customer reviews
				Limit(1).
				First(&closestProduct).Error; err != nil {
				return Data.Post{}, errors.New("error finding closest lower product ID")
			}
		} else {
			// Some other error occurred
			return Data.Post{}, errors.New("error finding closest higher product ID")
		}
	}

	// If no IDs are found, return an error
	if closestProduct.ID == 0 {
		return Data.Post{}, errors.New("no available product found")
	}

	return closestProduct, nil
}

func (d *DatabaseHelperImpl) GetNextProductID(currentID uint64) (uint64, error) {
	var nextID uint64

	// Query for the next higher ID
	err := conn.DB.Where("id > ?", currentID).Order("id ASC").Limit(1).Pluck("id", &nextID).Error
	if err != nil {
		return 0, fmt.Errorf("error finding next product ID: %w", err)
	}

	// Return the next ID found, even if it's not a specific number
	return nextID, nil
}

func (d *DatabaseHelperImpl) GetPreviousProductID(currentID uint64) (uint64, error) {
	var previousID uint64

	// Query for the previous lower ID
	err := conn.DB.Where("id < ?", currentID).Order("id DESC").Limit(1).Pluck("id", &previousID).Error
	if err != nil {
		return 0, fmt.Errorf("error finding previous product ID: %w", err)
	}

	// Return the previous ID found, even if it's not a specific number
	return previousID, nil
}

// SearchProductsByTitleAndCategory searches products by title and optionally by category.
// It preloads product images and limits results to 7.
func (d *DatabaseHelperImpl) SearchProductsByTitleAndCategory(searchTerm string, categorySlug string) ([]Data.ProductSearchResponse, error) {
	var products []Data.Post
	var results []Data.ProductSearchResponse

	// If search term is empty, return an empty slice.
	// Consider if you want to show 'featured' or 'trending' products when search term is empty.
	if searchTerm == "" {
		return []Data.ProductSearchResponse{}, nil
	}

	// Initialize the GORM query builder
	tx := conn.DB.Model(&Data.Post{}).Preload("ProductImages")

	// Apply search by title (case-insensitive)
	searchTermLower := "%" + strings.ToLower(searchTerm) + "%"
	tx = tx.Where("LOWER(title) LIKE ?", searchTermLower)

	// Apply category filter ONLY if categorySlug is provided and not "all-categories"
	if categorySlug != "" && categorySlug != "all-categories" {
		// Assuming your Post model has a `Category` field that directly stores the category slug
		// If `Post.Category` stores the actual category name, you might need to adjust this
		// to use a JOIN with a `categories` table if you have a `Category` model as well.
		// Based on your `Post` model, `Category` is a string field, so we'll filter directly on it.
		tx = tx.Where("LOWER(category) = ?", strings.ToLower(categorySlug))
	}

	// Limit the results
	tx = tx.Limit(7)

	// Execute the query
	err := tx.Find(&products).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []Data.ProductSearchResponse{}, nil // No products found, return empty slice
		}
		return nil, fmt.Errorf("error searching for products: %w", err)
	}

	// Process products to fit the ProductSearchResponse format
	for _, product := range products {
		var imageUrl string
		// Check if there are any images and take the first one
		if len(product.Images) > 0 {
			imageUrl = product.Images[0].URL
		} else {
			// Provide a default or placeholder image if no images are found
			imageUrl = "https://via.placeholder.com/150x150?text=No+Image" // Consider a local asset path
		}

		results = append(results, Data.ProductSearchResponse{
			Title: product.Title,
			// ProductUrlID: product.ProductUrlID,
			Category: product.BusinessCategory, // Using the string category from Post model
			// SellingPrice: product.SellingPrice,
			ImageUrl: imageUrl,
		})
	}

	fmt.Printf("Search for '%s' in category '%s' returned %d results.\n", searchTerm, categorySlug, len(results))
	return results, nil
}

func (d *DatabaseHelperImpl) SearchProductsByTitle(searchTerm string) ([]struct {
	Title        string `json:"title"`
	ProductUrlID string `json:"product_url_id"`
	Category     string `json:"category"`
}, error) {
	var results []struct {
		Title        string `json:"title"`
		ProductUrlID string `json:"product_url_id"`
		Category     string `json:"category"`
	}

	if searchTerm == "" {
		return nil, fmt.Errorf("error product do not exist")
	}
	// Convert search term to lowercase
	// searchTermLower := strings.ToLower(searchTerm)
	fmt.Println("Search term:", searchTerm)
	// Perform a case-insensitive search by converting the title to lowercase
	err := conn.DB.Model(&Data.Post{}).
		Select("title, product_url_id", "category").
		Where("LOWER(title) LIKE ?", "%"+searchTerm+"%").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("error searching for products: %w", err)
	}

	return results, nil
}

func (d *DatabaseHelperImpl) SearchAdminProductsByTitle(searchTerm string) ([]Data.Post, error) {
	var results []Data.Post

	query := conn.DB.Model(&Data.Post{}).
		Where("publish_status = ?", "publish").
		Where("deleted_at IS NULL")

	if searchTerm != "" {
		fmt.Println("Search term:", searchTerm)
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(searchTerm)+"%").Limit(20)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("error searching for products: %w", err)
	}

	return results, nil
}

func (d *DatabaseHelperImpl) SearchAdminOrderByTitle(searchTerm string) ([]Data.OrderHistory, error) {
	var results []Data.OrderHistory

	query := conn.DB.Preload("ProductOrders").Model(&Data.OrderHistory{}).
		Where("deleted_at IS NULL")

	if searchTerm != "" {
		fmt.Println("Search term:", searchTerm)

		// Try to convert searchTerm to an integer to check if it's an ID
		if id, err := strconv.Atoi(searchTerm); err == nil {
			// Search by ID
			query = query.Where("id = ?", id)
		} else {
			// Search by title (case-insensitive)
			query = query.Where("LOWER(customer_email) LIKE ?", "%"+strings.ToLower(searchTerm)+"%")
		}

		query = query.Limit(20)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("error searching for products: %w", err)
	}

	return results, nil
}

func (d *DatabaseHelperImpl) GetTransactionsByUserAndDateWithLimit(userID uint, dateString string) ([]Data.Post, error) {
	// Parse the date string to a time.Time object
	parsedDate, err := time.Parse("01/02/2006", dateString)
	if err != nil {
		return nil, err
	}

	// Format the parsed date to match the database date format
	formattedDate := parsedDate.Format("2006-01-02")

	// Query the database for transactions for the specified user and date with a limit of 10
	var transactions []Data.Post
	if err := conn.DB.Where("user_id = ? AND DATE(created_at) = ?", userID, formattedDate).Limit(10).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (d *DatabaseHelperImpl) GetTransactionHistoryForAi(userID uint64, limit int) ([]Data.Post, int64, error) {
	var history []Data.Post
	var transactionCount int64

	// Count total number of transactions for the user
	if err := conn.DB.Model(&Data.Post{}).Where("user_id = ?", userID).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve transaction history for the user with pagination
	if transactionErr := conn.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Find(&history).Error; transactionErr != nil {
		if errors.Is(transactionErr, gorm.ErrRecordNotFound) {
			// The record with the specified UserID was not found
			return []Data.Post{}, 0, errors.New("no transaction record found")
		} else {
			// Some other error occurred
			return []Data.Post{}, 0, errors.New("error retrieving transaction record")
		}
	}

	return history, transactionCount, nil
}

// Helper function to get ordinal suffix
func getOrdinal(n int) string {
	if n%10 == 1 && n%100 != 11 {
		return "st"
	}
	if n%10 == 2 && n%100 != 12 {
		return "nd"
	}
	if n%10 == 3 && n%100 != 13 {
		return "rd"
	}
	return "th"
}

func (d *DatabaseHelperImpl) AddEmailSubscriber(email string) error {
	var user Data.SubscribeToEmail
	result1 := conn.DB.First(&user, "email = ?", email)

	if result1.Error != nil {
		if errors.Is(result1.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			// Create a new email subscriber
			var emailSubscriber Data.SubscribeToEmail
			emailSubscriber.Email = email
			result := conn.DB.Create(&emailSubscriber)

			// Checking if an error occurred when creating the email subscriber
			if result.Error != nil {
				return errors.New("failed to create email subscriber")
			}
			return nil // Return nil to indicate success
		} else {
			// Some other error occurred
			return errors.New("error retrieving user")
		}
	}

	// User already exists, return specific error message
	return errors.New("user with email already exists")
}

func (d *DatabaseHelperImpl) UpdateSiteVisits() error {

	// Get today's date
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Check if there's an existing record for today
	var visit Data.SiteVisit
	result := conn.DB.Where("DATE(created_at) = ?", today.Format("2006-01-02")).First(&visit)
	// if result != nil {
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// If no record found for today, create a new record
			visit = Data.SiteVisit{
				SiteVisitNumber: 1,
			}
			conn.DB.Create(&visit)
		} else {
			return errors.New("failed to update site visit")
		}
	} else {
		// If record found, update the visit number
		visit.SiteVisitNumber++
		conn.DB.Save(&visit)
	}

	return nil
}

func (d *DatabaseHelperImpl) GetLast12DaysSiteVisits() (map[string]int64, error) {
	// Get today's date
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Map to hold site visits for each of the last 12 days
	visitCounts := make(map[string]int64)

	// Iterate over the last 12 days
	for i := 0; i < 12; i++ {
		day := today.AddDate(0, 0, -i).Format("2006-01-02")
		dayName := today.AddDate(0, 0, -i).Weekday().String() // Get day of the week
		ordinal := getOrdinal(i + 1)                          // Get ordinal suffix
		dayLabel := fmt.Sprintf("%d%s %s", i+1, ordinal, dayName)

		var visit Data.SiteVisit
		result := conn.DB.Where("DATE(created_at) = ?", day).First(&visit)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// No visits for this day, continue to the next day
				visitCounts[dayLabel] = 0
			} else {
				return nil, fmt.Errorf("error retrieving site visits for %s: %v", day, result.Error)
			}
		} else {
			// Convert uint to int64 for consistency
			visitCounts[dayLabel] = int64(visit.SiteVisitNumber)
		}
	}

	return visitCounts, nil
}

type SubscriptionService struct {
	ServiceID     uint
	ServiceName   string
	ServiceNumber string
}

func (d *DatabaseHelperImpl) UpdateSubscriptionStatus(transactionID uint64) error {
	// Update auto renewal status to false for all Post records with the specified user ID
	if err := conn.DB.Model(&Data.Post{}).Where("id = ?", transactionID).Update("auto_renew", false).Error; err != nil {
		return err
	}

	return nil
}

// GenerateProductURL generates a product URL from the product name and ID
func GenerateProductURL(productName string, productID int) string {
	// Convert product name to lowercase
	productName = strings.ToLower(productName)

	// Replace spaces and underscores with hyphens
	productName = strings.ReplaceAll(productName, " ", "-")
	productName = strings.ReplaceAll(productName, "_", "-")

	// Remove any non-alphanumeric characters (except hyphens)
	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	productName = reg.ReplaceAllString(productName, "")

	// Construct the final URL with the product ID at the end
	productURL := fmt.Sprintf("%s-%d", productName, productID)

	return productURL
}

// AddProduct adds a product and updates its ProductUrlID
func (d *DatabaseHelperImpl) AddProduct(post Data.Post, user Data.User) (Data.Post, error) {
	// Create the product in the database
	post.UserName = user.FullName
	post.ProfilePhotoURL = user.ProfilePhotoURL
	post.PhoneNumber = user.PhoneNumber
	post.Verified = user.Verified
	post.IsActive = true
	post.Approved = true
	if post.Location == nil || *post.Location == "" {
		post.Location = &user.State
	}
	result := conn.DB.Create(&post)

	// Check if an error occurred when creating the product
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Data.Post{}, fmt.Errorf("failed to create product: product not found, %w", result.Error)
		} else if errors.Is(result.Error, gorm.ErrInvalidData) {
			return Data.Post{}, fmt.Errorf("failed to create product: invalid data provided, %w", result.Error)
		} else {
			return Data.Post{}, fmt.Errorf("failed to create product: %w", result.Error)
		}
	}

	// Generate the ProductUrlID
	post.ProductUrlID = GenerateProductURL(post.Title, int(post.ID))

	// Update the product with the new ProductUrlID
	updateResult := conn.DB.Model(&post).Update("product_url_id", post.ProductUrlID)
	if updateResult.Error != nil {
		return Data.Post{}, fmt.Errorf("failed to update product URL ID: %w", updateResult.Error)
	}

	// Update user with posts amount
	updateUserResult := conn.DB.Model(&user).Update("total_product", user.TotalProduct+1)
	if updateUserResult.Error != nil {
		return Data.Post{}, fmt.Errorf("failed to update user's total products: %w", updateUserResult.Error)
	}

	// Return the updated product
	return post, nil
}

func (d *DatabaseHelperImpl) AddProductImage(image Data.PostImage, postID uint) error {
	// Set the post_id for the image
	image.PostID = postID

	result := conn.DB.Create(&image)

	// Check if an error occurred when creating the post
	if result.Error != nil {
		// Some other error occurred
		return errors.New("an unknown error occurred")
	}

	// Return the ID of the newly created post
	return nil
}

// AddProduct adds a product and updates its ProductUrlID
func (d *DatabaseHelperImpl) AddBlog(post Data.Blog, user Data.User) (Data.Blog, error) {
	// Create the product in the database
	result := conn.DB.Create(&post)

	// Check if an error occurred when creating the product
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Data.Blog{}, fmt.Errorf("failed to create product: product not found, %w", result.Error)
		} else if errors.Is(result.Error, gorm.ErrInvalidData) {
			return Data.Blog{}, fmt.Errorf("failed to create product: invalid data provided, %w", result.Error)
		} else {
			return Data.Blog{}, fmt.Errorf("failed to create product: %w", result.Error)
		}
	}

	// Generate the ProductUrlID
	post.BlogUrlID = GenerateProductURL(post.Title, int(post.ID))

	// Update the product with the new ProductUrlID
	updateResult := conn.DB.Model(&post).Update("blog_url_id", post.BlogUrlID)
	if updateResult.Error != nil {
		return Data.Blog{}, fmt.Errorf("failed to update product URL ID: %w", updateResult.Error)
	}

	// Update user with posts amount
	updateUserResult := conn.DB.Model(&user).Update("total_product", user.TotalProduct+1)
	if updateUserResult.Error != nil {
		return Data.Blog{}, fmt.Errorf("failed to update user's total products: %w", updateUserResult.Error)
	}

	// Return the updated product
	return post, nil
}

// SaveCustomerReview saves a customer review for a product
// func (d *DatabaseHelperImpl) SaveCustomerReview(productID uint, email string, name string, reviewText string, rating int) (Data.CustomerReview, error) {
// 	// Check if the product exists in the database
// 	var product Data.Post
// 	if err := conn.DB.First(&product, productID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return Data.CustomerReview{}, fmt.Errorf("product not found: %w", err)
// 		}
// 		return Data.CustomerReview{}, fmt.Errorf("failed to retrieve product: %w", err)
// 	}

// 	// Create the new customer review
// 	customerReview := Data.CustomerReview{
// 		ProducttID: productID,
// 		Email:      email,
// 		Name:       name,
// 		Review:     reviewText,
// 		Rating:     rating,
// 	}

// 	// Save the review to the database
// 	if err := conn.DB.Create(&customerReview).Error; err != nil {
// 		return Data.CustomerReview{}, fmt.Errorf("failed to save customer review: %w", err)
// 	}

// 	// Update the product's review count
// 	if err := conn.DB.Model(&product).Update("product_reviews_count", product.ProductReviewsCount+1).Error; err != nil {
// 		return Data.CustomerReview{}, fmt.Errorf("failed to update product review count: %w", err)
// 	}

// 	// Return the saved customer review
// 	return customerReview, nil
// }

// GetCustomerReviewsByProduct retrieves a list of customer reviews for a given product with pagination
// func (d *DatabaseHelperImpl) GetCustomerReviewsByProduct(productID uint, limit int, offset int) ([]Data.CustomerReview, int64, error) {
// 	var reviews []Data.CustomerReview

// 	// Check if the product exists in the database
// 	var product Data.Post
// 	if err := conn.DB.First(&product, productID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, 0, fmt.Errorf("product not found: %w", err)
// 		}
// 		return nil, 0, fmt.Errorf("failed to retrieve product: %w", err)
// 	}

// 	// Retrieve reviews with limit and offset for pagination
// 	if err := conn.DB.Where("productt_id = ?", productID).
// 		Limit(limit).
// 		Offset(offset).
// 		Find(&reviews).Error; err != nil {
// 		return nil, 0, fmt.Errorf("failed to retrieve customer reviews: %w", err)
// 	}

// 	// Return the list of reviews
// 	return reviews, product.ProductReviewsCount, nil
// }

// SaveCustomerReview saves a customer review for a product
func (d *DatabaseHelperImpl) SaveCustomerBlogReview(blogID uint, email string, name string, reviewText string, rating int) (Data.CustomerBlogReview, error) {
	// Check if the blog post exists in the database
	var blog Data.Blog
	if err := conn.DB.First(&blog, blogID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Data.CustomerBlogReview{}, fmt.Errorf("blog not found: %w", err)
		}
		return Data.CustomerBlogReview{}, fmt.Errorf("failed to retrieve product: %w", err)
	}

	// Create the new customer review
	customerReview := Data.CustomerBlogReview{
		BlogID: blogID,
		Email:  email,
		Name:   name,
		Review: reviewText,
		Rating: rating,
	}

	// Save the review to the database
	if err := conn.DB.Create(&customerReview).Error; err != nil {
		return Data.CustomerBlogReview{}, fmt.Errorf("failed to save customer review: %w", err)
	}

	// Update the product's review count
	if err := conn.DB.Model(&blog).Update("blog_reviews_count", blog.BlogReviewsCount+1).Error; err != nil {
		return Data.CustomerBlogReview{}, fmt.Errorf("failed to update product review count: %w", err)
	}

	// Return the saved customer review
	return customerReview, nil
}

// GetCustomerReviewsByProduct retrieves a list of customer reviews for a given product with pagination
func (d *DatabaseHelperImpl) GetCustomerBlogReviewsByBlogPost(blogID uint, limit int, offset int) ([]Data.CustomerBlogReview, int64, error) {
	var reviews []Data.CustomerBlogReview

	// Check if the blog post exists in the database
	var blog Data.Blog
	if err := conn.DB.First(&blog, blogID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("blog not found: %w", err)
		}
		return nil, 0, fmt.Errorf("failed to retrieve blog: %w", err)
	}

	// Retrieve reviews with limit and offset for pagination
	if err := conn.DB.Where("blog_id = ?", blogID).
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve customer reviews: %w", err)
	}

	// Return the list of reviews
	return reviews, blog.BlogReviewsCount, nil
}

func (d *DatabaseHelperImpl) AddBlogImage(image Data.BlogImage, postID uint) error {
	// Set the post_id for the image
	image.BlogID = postID

	result := conn.DB.Create(&image)

	// Check if an error occurred when creating the post
	if result.Error != nil {
		// Some other error occurred
		return errors.New("an unknown error occurred")
	}

	// Return the ID of the newly created post
	return nil
}

func ConvertToOrderHistory(orderHistoryBody Data.OrderHistoryBody) *Data.OrderHistory {
	return &Data.OrderHistory{
		OrderStatus:            orderHistoryBody.OrderStatus,
		PaymentStatus:          "pending", // Default to pending
		Quantity:               orderHistoryBody.Quantity,
		OrderCost:              orderHistoryBody.OrderCost,
		OrderNote:              orderHistoryBody.OrderNote,
		OrderSubTotalCost:      orderHistoryBody.OrderSubTotalCost,
		ShippingCost:           orderHistoryBody.ShippingCost,
		OrderDiscount:          orderHistoryBody.OrderDiscount,
		CustomerEmail:          orderHistoryBody.CustomerEmail,
		CustomerFName:          orderHistoryBody.CustomerFName,
		CustomerSName:          orderHistoryBody.CustomerSName,
		CustomerCompanyName:    orderHistoryBody.CustomerCompanyName,
		CustomerState:          orderHistoryBody.CustomerState,
		CustomerCity:           orderHistoryBody.CustomerCity,
		CustomerStreetAddress1: orderHistoryBody.CustomerStreetAddress1,
		CustomerStreetAddress2: orderHistoryBody.CustomerStreetAddress2,
		CustomerZipCode:        orderHistoryBody.CustomerZipCode,
		CustomerProvince:       orderHistoryBody.CustomerProvince,
		CustomerPhoneNumber:    orderHistoryBody.CustomerPhoneNumber,
	}
}

func ConvertToProductOrders(orderBodies []Data.ProductOrderBody) []Data.ProductOrder {
	var productOrders []Data.ProductOrder
	for _, body := range orderBodies {
		productOrders = append(productOrders, Data.ProductOrder{
			ProducttID:   body.ID,
			ProductUrlID: body.ProductUrlID,
			Title:        body.Title,
			Description:  body.Description,
			NetWeight:    body.NetWeight,
			OrderCost:    body.OrderCost,
			Currency:     body.Currency,
			Quantity:     body.Quantity,
			Category:     body.Category,
			Image1:       body.Image1,
			Image2:       body.Image2,
		})
	}
	return productOrders
}

func (d *DatabaseHelperImpl) AddOrder(orderHistoryBody Data.OrderHistoryBody, ordersBody []Data.ProductOrderBody) (uint, *Data.OrderHistory, []Data.ProductOrder, error) {
	orderHistory := ConvertToOrderHistory(orderHistoryBody)
	productOrders := ConvertToProductOrders(ordersBody)

	err := conn.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the order history record
		if err := tx.Create(&orderHistory).Error; err != nil {
			return err
		}

		// 2. Track totals
		var totalProducts int64
		for i := range productOrders {
			productOrders[i].OrderHistoryID = orderHistory.ID

			var product Data.Post
			if err := tx.Where("product_url_id = ?", ordersBody[i].ProductUrlID).First(&product).Error; err != nil {
				return err
			}

			// if product.StockRemaining < ordersBody[i].Quantity {
			// 	return fmt.Errorf("insufficient stock for product: %s", product.Title)
			// }

			// product.StockRemaining -= ordersBody[i].Quantity
			// product.Sales += ordersBody[i].Quantity

			totalProducts += int64(ordersBody[i].Quantity)

			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}

		// 3. Save product orders
		if err := tx.Create(&productOrders).Error; err != nil {
			return err
		}

		// 4. Update analytics
		now := time.Now()
		currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		var analytics Data.Analytics
		var prevAnalytics Data.Analytics

		// Find current month's analytics
		if err := tx.Where("month = ?", currentMonth).First(&analytics).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				analytics = Data.Analytics{
					Month:          currentMonth,
					TotalRevenue:   0,
					TotalSales:     0,
					TotalCustomers: 0,
					TotalProducts:  0,
				}
			} else {
				return err
			}
		}

		// Get previous month
		prevMonth := currentMonth.AddDate(0, -1, 0)
		tx.Where("month = ?", prevMonth).First(&prevAnalytics)

		// Calculate analytics
		analytics.TotalRevenue += orderHistory.OrderCost
		analytics.TotalSales += 1
		analytics.TotalCustomers += 1 // You can deduplicate by customer email if needed
		analytics.TotalProducts += totalProducts

		if prevAnalytics.TotalRevenue > 0 {
			analytics.RevenueChange = ((analytics.TotalRevenue - prevAnalytics.TotalRevenue) / prevAnalytics.TotalRevenue) * 100
			analytics.IsRevenueBetter = analytics.TotalRevenue > prevAnalytics.TotalRevenue
		}
		if prevAnalytics.TotalSales > 0 {
			analytics.SalesChange = float64(analytics.TotalSales-prevAnalytics.TotalSales) / float64(prevAnalytics.TotalSales) * 100
			analytics.IsSalesBetter = analytics.TotalSales > prevAnalytics.TotalSales
		}
		if prevAnalytics.TotalCustomers > 0 {
			analytics.CustomerChange = float64(analytics.TotalCustomers-prevAnalytics.TotalCustomers) / float64(prevAnalytics.TotalCustomers) * 100
			analytics.IsCustomersBetter = analytics.TotalCustomers > prevAnalytics.TotalCustomers
		}
		if prevAnalytics.TotalProducts > 0 {
			analytics.ProductChange = float64(analytics.TotalProducts-prevAnalytics.TotalProducts) / float64(prevAnalytics.TotalProducts) * 100
			analytics.IsProductsBetter = analytics.TotalProducts > prevAnalytics.TotalProducts
		}

		// Upsert analytics
		if analytics.ID == 0 {
			if err := tx.Create(&analytics).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Save(&analytics).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, nil, nil, errors.New("failed to create order and update analytics: " + err.Error())
	}

	return orderHistory.ID, orderHistory, productOrders, nil
}

func (d *DatabaseHelperImpl) UpsertShippingFee(fee int64, feesGreater, feesLess int64, storeLatitude, storeLongitude float64,
	storeState, storeCity, stateISO string, calculateUsingKg bool) error {
	var shippingFee Data.ShippingFees

	// Try to find the existing shipping fee
	result := conn.DB.First(&shippingFee)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		// Record not found, create a new one
		shippingFee := Data.ShippingFees{
			ShippingFeePerKm:   fee,
			ShippingFeeGreater: feesGreater,
			ShippingFeeLess:    feesLess,
			StoreLatitude:      storeLatitude,
			StoreLongitude:     storeLongitude,
			StoreState:         storeState,
			StoreCity:          storeCity,
			StateISO:           stateISO,
			CalculateUsingKg:   calculateUsingKg,
		}
		result = conn.DB.Create(&shippingFee)
	} else if result.Error != nil {
		return result.Error
	} else {
		// Record found, update the existing one
		shippingFee.ShippingFeePerKm = fee
		shippingFee.ShippingFeePerKm = fee
		shippingFee.ShippingFeeGreater = feesGreater
		shippingFee.ShippingFeeLess = feesLess
		shippingFee.StoreLatitude = storeLatitude
		shippingFee.StoreLongitude = storeLongitude
		shippingFee.StoreState = storeState
		shippingFee.StoreCity = storeCity
		shippingFee.StateISO = stateISO
		shippingFee.CalculateUsingKg = calculateUsingKg
		result = conn.DB.Save(&shippingFee)
	}

	return result.Error
}

func (d *DatabaseHelperImpl) GetShippingFee() (Data.ShippingFees, error) {
	var shippingFee Data.ShippingFees
	result := conn.DB.First(&shippingFee) // Retrieves the first entry
	if result.Error != nil {
		return Data.ShippingFees{}, result.Error
	}
	return shippingFee, nil
}

func (d *DatabaseHelperImpl) GetOrder(orderID uint) (*Data.OrderHistory, error) {
	// Initialize variables
	var orderHistory Data.OrderHistory

	// Find the order history by ID and preload the associated ProductOrders
	if err := conn.DB.Preload("ProductOrders").First(&orderHistory, orderID).Error; err != nil {
		return nil, errors.New("failed to find order history: " + err.Error())
	}

	// Return the orderHistory and the associated productOrders
	return &orderHistory, nil
}

func (d *DatabaseHelperImpl) GetAndUpdateOrder(orderID uint, status string) (*Data.OrderHistory, error) {
	// Initialize variable
	var orderHistory Data.OrderHistory

	// Find the order with preloaded ProductOrders
	if err := conn.DB.Preload("ProductOrders").First(&orderHistory, orderID).Error; err != nil {
		return nil, errors.New("failed to find order history: " + err.Error())
	}

	// Update the payment status to "Successful"
	orderHistory.PaymentStatus = status

	// Save the update to the database
	if err := conn.DB.Save(&orderHistory).Error; err != nil {
		return nil, errors.New("failed to update payment status: " + err.Error())
	}

	return &orderHistory, nil
}

func (d *DatabaseHelperImpl) GetBlogPostById(blogID uint) (*Data.Blog, error) {
	// Initialize variables
	var blog Data.Blog

	// Preload related data with a limit on CustomerBlogReviews
	err := conn.DB.Preload("CustomerBlogReviews", func(db *gorm.DB) *gorm.DB {
		return db.Limit(4)
	}).First(&blog, blogID).Error

	if err != nil {
		return nil, errors.New("failed to find blog post: " + err.Error())
	}

	// Return the blog post and the associated customer reviews
	return &blog, nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectOrdersByLimit(limit, offset int) ([]Data.OrderHistory, int64, error) {
	var orderRecords []Data.OrderHistory
	var orderRecordsCount int64

	// Get the count of order records
	if err := conn.DB.Model(&Data.OrderHistory{}).Count(&orderRecordsCount).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve order history with pagination
	if err := conn.DB.Preload("ProductOrders").Order("created_at desc").Limit(limit).Offset(offset).Find(&orderRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found
			return nil, 0, errors.New("no order records found")
		}
		// Some other error occurred
		return nil, 0, errors.New("error retrieving order records: " + err.Error())
	}

	return orderRecords, orderRecordsCount, nil
}

func CalculateMonthlyAnalytics() (*Data.Analytics, error) {
	// Current month and year
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	var analytics Data.Analytics
	var previousAnalytics Data.Analytics
	var currentAnalytics Data.Analytics
	var orderHistories []Data.OrderHistory
	var topProducts []struct {
		ProducttID uint
		Sales      int64
	}

	// Retrieve current month's order histories
	if err := conn.DB.
		Where("EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?", currentMonth, currentYear).
		Preload("ProductOrders"). // Preload the ProductOrders relation
		Find(&orderHistories).Error; err != nil {
		return nil, fmt.Errorf("error retrieving order histories with product orders: %v", err)
	}

	// Check if analytics for the current month already exist
	if err := conn.DB.Where("EXTRACT(MONTH FROM month) = ? AND EXTRACT(YEAR FROM month) = ?", currentMonth, currentYear).
		First(&currentAnalytics).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error retrieving current month's analytics: %v", err)
	}

	// Initialize totals
	totalRevenue := float64(0)
	totalSales := int64(0)
	totalCustomers := int64(0)
	totalProducts := int64(0)
	dailyVisitors := int64(0) // Depends on visitor tracking mechanism

	customerEmails := make(map[string]struct{})
	productSales := make(map[uint]int64)

	// Calculate totals
	for _, order := range orderHistories {
		totalRevenue += order.OrderCost
		totalSales += order.Quantity

		// Track unique customers
		if _, exists := customerEmails[order.CustomerEmail]; !exists {
			customerEmails[order.CustomerEmail] = struct{}{}
			totalCustomers++
		}

		// Track product sales
		for _, productOrder := range order.ProductOrders {
			productSales[productOrder.ProducttID] += productOrder.Quantity
			totalProducts += productOrder.Quantity // Add to total products sold
		}
	}

	// Get top two products by sales
	for productID, sales := range productSales {
		topProducts = append(topProducts, struct {
			ProducttID uint
			Sales      int64
		}{ProducttID: productID, Sales: sales})
	}

	sort.Slice(topProducts, func(i, j int) bool {
		return topProducts[i].Sales > topProducts[j].Sales
	})

	topProduct1ID := uint(0)
	topProduct2ID := uint(0)
	if len(topProducts) > 0 {
		topProduct1ID = topProducts[0].ProducttID
	}
	if len(topProducts) > 1 {
		topProduct2ID = topProducts[1].ProducttID
	}

	// Calculate previous month and year
	previousMonth := currentMonth - 1
	previousYear := currentYear
	if previousMonth == 0 {
		previousMonth = 12
		previousYear--
	}

	// Retrieve previous month's analytics
	if err := conn.DB.Where("EXTRACT(MONTH FROM month) = ? AND EXTRACT(YEAR FROM month) = ?", previousMonth, previousYear).First(&previousAnalytics).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error retrieving previous month's analytics: %v", err)
	}

	// Calculate percentage changes
	revenueChange, salesChange, customerChange, productChange := float64(0), float64(0), float64(0), float64(0)
	isRevenueBetter, isSalesBetter, isCustomersBetter, isProductsBetter := false, false, false, false

	if previousAnalytics.ID != 0 { // If previous analytics exist
		revenueChange = ((totalRevenue - previousAnalytics.TotalRevenue) / previousAnalytics.TotalRevenue) * 100
		salesChange = float64(totalSales-previousAnalytics.TotalSales) / float64(previousAnalytics.TotalSales) * 100
		customerChange = float64(totalCustomers-previousAnalytics.TotalCustomers) / float64(previousAnalytics.TotalCustomers) * 100
		productChange = float64(totalProducts-previousAnalytics.TotalProducts) / float64(previousAnalytics.TotalProducts) * 100

		// Compare each metric to determine if this month is better
		isRevenueBetter = totalRevenue > previousAnalytics.TotalRevenue
		isSalesBetter = totalSales > previousAnalytics.TotalSales
		isCustomersBetter = totalCustomers > previousAnalytics.TotalCustomers
		isProductsBetter = totalProducts > previousAnalytics.TotalProducts
	}

	// If current analytics already exist, update them
	if currentAnalytics.ID != 0 {
		currentAnalytics.TotalRevenue = totalRevenue
		currentAnalytics.RevenueChange = revenueChange
		currentAnalytics.TotalSales = totalSales
		currentAnalytics.SalesChange = salesChange
		currentAnalytics.TotalCustomers = totalCustomers
		currentAnalytics.CustomerChange = customerChange
		currentAnalytics.TotalProducts = totalProducts
		currentAnalytics.ProductChange = productChange
		currentAnalytics.DailyVisitors = dailyVisitors
		currentAnalytics.TopProduct1ID = topProduct1ID
		currentAnalytics.TopProduct2ID = topProduct2ID
		currentAnalytics.IsRevenueBetter = isRevenueBetter
		currentAnalytics.IsSalesBetter = isSalesBetter
		currentAnalytics.IsCustomersBetter = isCustomersBetter
		currentAnalytics.IsProductsBetter = isProductsBetter

		// Update existing analytics
		if err := conn.DB.Save(&currentAnalytics).Error; err != nil {
			return nil, fmt.Errorf("error updating analytics: %v", err)
		}
		return &currentAnalytics, nil
	}

	// Create new Analytics entry if not found
	newMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	analytics = Data.Analytics{
		Month:             newMonth,
		TotalRevenue:      totalRevenue,
		RevenueChange:     revenueChange,
		TotalSales:        totalSales,
		SalesChange:       salesChange,
		TotalCustomers:    totalCustomers,
		CustomerChange:    customerChange,
		TotalProducts:     totalProducts,
		ProductChange:     productChange,
		DailyVisitors:     dailyVisitors,
		TopProduct1ID:     topProduct1ID,
		TopProduct2ID:     topProduct2ID,
		IsRevenueBetter:   isRevenueBetter,
		IsSalesBetter:     isSalesBetter,
		IsCustomersBetter: isCustomersBetter,
		IsProductsBetter:  isProductsBetter,
	}

	// Save new analytics
	if err := conn.DB.Save(&analytics).Error; err != nil {
		return nil, fmt.Errorf("error saving analytics: %v", err)
	}

	return &analytics, nil
}

func UpdateProductSalesByOrderHistory(orderHistoryID uint) error {
	var orderHistory Data.OrderHistory

	// Step 1: Retrieve OrderHistory with associated ProductOrders
	if err := conn.DB.Preload("ProductOrders").First(&orderHistory, orderHistoryID).Error; err != nil {
		return fmt.Errorf("error retrieving order history: %v", err)
	}

	// Step 2: Iterate over each ProductOrder
	for _, productOrder := range orderHistory.ProductOrders {
		var product Data.Post

		// Step 3: Retrieve the Post by ProducttID from ProductOrder
		if err := conn.DB.First(&product, productOrder.ProducttID).Error; err != nil {
			return fmt.Errorf("error retrieving product: %v", err)
		}

		// Step 4: Update the Sales count (increment by ProductOrder.Quantity)
		// product.Sales += productOrder.Quantity

		// Save the updated Post
		if err := conn.DB.Save(&product).Error; err != nil {
			return fmt.Errorf("error updating product sales: %v", err)
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) UpdateOrderStatus(orderID uint, newStatus string) error {
	// Assuming your GORM model for the order history is named OrderHistory
	// Update only the OrderStatus field where the Order ID matches
	result := conn.DB.Model(&Data.OrderHistory{}).Where("id = ?", orderID).Updates(map[string]interface{}{
		"order_status": newStatus,
	})

	// Check for errors
	if result.Error != nil {
		return result.Error
	}

	if newStatus == "shipped" {
		// Update product sales
		updateErr := UpdateProductSalesByOrderHistory(orderID)
		if updateErr != nil {
			return updateErr
		}

		// Update monthly analytics
		if _, err := CalculateMonthlyAnalytics(); err != nil {
			return fmt.Errorf("failed to update monthly analytics: %v", err)
		}
	}

	// Check if any rows were affected (i.e., updated)
	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found with the given order ID")
	}

	return nil
}

func (d *DatabaseHelperImpl) GetAnalyticsData() (*Data.Analytics, error) {
	var analytics Data.Analytics
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	// Retrieve current month's analytics data using EXTRACT for PostgreSQL
	if err := conn.DB.Where("EXTRACT(MONTH FROM month) = ? AND EXTRACT(YEAR FROM month) = ?", currentMonth, currentYear).First(&analytics).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No data found for the current month
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving analytics data: %v", err)
	}

	return &analytics, nil
}

func (d *DatabaseHelperImpl) UpdateBusinessConnectProduct(Post Data.Post, ProductID uint) error {
	// Find the product and update it
	returnedProduct, productErr := d.GetBusinessConnectProductByIDd(uint64(ProductID))

	if productErr != nil {
		return errors.New("error retrieving product")
	}

	returnedProduct.Title = Post.Title
	returnedProduct.Description = Post.Description
	// returnedProduct.InitialCost = Post.InitialCost
	// returnedProduct.SellingPrice = Post.SellingPrice
	// returnedProduct.ProductStock = Post.ProductStock
	// returnedProduct.StockRemaining = Post.StockRemaining
	// returnedProduct.NetWeight = Post.NetWeight
	returnedProduct.BusinessCategory = Post.BusinessCategory
	// returnedProduct.ProductRank = Post.ProductRank
	// returnedProduct.Tags = Post.Tags
	// returnedProduct.PublishStatus = Post.PublishStatus

	// Update the product
	// save updated user
	result := conn.DB.Save(returnedProduct)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("product not found")
		} else {
			// Some other error occurred
			return errors.New("error retrieving product")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) UpdateBusinessConnectBlog(Blog Data.Blog, BlogID uint) error {
	// Find the product and update it
	returnedBlog, productErr := d.GetBlogPostById(uint(BlogID))

	if productErr != nil {
		return errors.New("error retrieving product")
	}

	returnedBlog.Title = Blog.Title
	returnedBlog.Description1 = Blog.Description1
	returnedBlog.Description2 = Blog.Description2
	returnedBlog.BlogCategory = Blog.BlogCategory

	// Update the blog
	// save updated blog
	result := conn.DB.Save(returnedBlog)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// The record with the specified email was not found
			return errors.New("blog not found")
		} else {
			// Some other error occurred
			return errors.New("error retrieving blog")
		}
	}

	return nil
}

func (d *DatabaseHelperImpl) DeleteBusinessConnectProduct(ProductID uint) error {
	var returnedProduct Data.Post

	if err := conn.DB.First(&returnedProduct, "id = ?", uint64(ProductID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Post not found, find the next closest available product ID
			// nextProductID, findErr := d.findClosestAvailableProductID(productID)
			// if findErr != nil {
			return errors.New("product not found")
			// }
		}
	}

	// Delete the product
	result := conn.DB.Delete(&returnedProduct)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return result.Error // Return the actual error
	}

	// Optional: Check if a row was actually deleted
	if result.RowsAffected == 0 {
		return errors.New("product not found or already deleted")
	}

	return nil
}

func (d *DatabaseHelperImpl) DeleteBusinessConnectBlog(BlogID uint) error {
	var returnedBlog Data.Blog

	if err := conn.DB.First(&returnedBlog, "id = ?", uint64(BlogID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Post not found, find the next closest available product ID
			// nextProductID, findErr := d.findClosestAvailableProductID(productID)
			// if findErr != nil {
			return errors.New("blog not found")
			// }
		}
	}

	// Delete the product
	result := conn.DB.Delete(&returnedBlog)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("blog not found")
		}
		return result.Error // Return the actual error
	}

	// Optional: Check if a row was actually deleted
	if result.RowsAffected == 0 {
		return errors.New("blog not found or already deleted")
	}

	return nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectEmailSubscribers() ([]Data.BusinessConnectEmailSubscriber, error) {
	var emailSubscribers []Data.BusinessConnectEmailSubscriber

	// Retrieve order history with pagination
	if err := conn.DB.Find(&emailSubscribers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No records found
			return nil, errors.New("no email records found")
		}
		// Some other error occurred
		return nil, errors.New("error retrieving email records: " + err.Error())
	}

	return emailSubscribers, nil
}

func (d *DatabaseHelperImpl) SaveBusinessConnectSentEmail(sentEmail Data.Email) error {

	// Retrieve order history with pagination
	if err := conn.DB.Save(&sentEmail).Error; err != nil {
		// Some other error occurred
		return errors.New("error saving email records: " + err.Error())
	}

	return nil
}

func (d *DatabaseHelperImpl) GetBusinessConnectUniqueUserFingerPrintHash(fingerprintHash string) (Data.BusinessConnectDeviceFingerprint, error) {
	var deviceFingerprint Data.BusinessConnectDeviceFingerprint

	err := conn.DB.Where("fingerprint_hash = ?", fingerprintHash).First(&deviceFingerprint).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Data.BusinessConnectDeviceFingerprint{}, errors.New("no fingerprint record found")
		}
		return Data.BusinessConnectDeviceFingerprint{}, fmt.Errorf("error retrieving fingerprint record: %w", err)
	}

	return deviceFingerprint, nil
}

func (d *DatabaseHelperImpl) CreateBusinessConnectDeviceFingerprint(fingerprintHash string) error {
	// Check if the fingerprint already exists
	var existing Data.BusinessConnectDeviceFingerprint
	err := conn.DB.Where("fingerprint_hash = ?", fingerprintHash).First(&existing).Error
	if err == nil {
		// Fingerprint already exists
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other error occurred while checking
		return fmt.Errorf("error checking for existing fingerprint: %w", err)
	}

	// Fingerprint doesn't exist, proceed to create
	fp := Data.BusinessConnectDeviceFingerprint{
		FingerprintHash: fingerprintHash,
	}

	// Create the new fingerprint
	if err := conn.DB.Create(&fp).Error; err != nil {
		return fmt.Errorf("error creating fingerprint: %w", err)
	}
	return nil
}

func (d *DatabaseHelperImpl) RecommendProductsForUser2(fingerprintHash string, limit, offset int) ([]Data.Post, error) {
	var products []Data.Post

	// Check if fingerprintHash is empty for anonymous users
	if fingerprintHash == "" {
		fmt.Println("User is anonymous, fetching anonymous recommendations")
		return d.RecommendForAnonymous(limit)
	}

	// Fetch viewed product IDs by the user
	var viewedProductIDs []uint
	err := conn.DB.Model(&Data.BusinessConnectUserActivity{}).
		Where("fingerprint_hash = ? AND product_id IS NOT NULL", fingerprintHash).
		Pluck("DISTINCT product_id", &viewedProductIDs).Error
	if err != nil {
		fmt.Println("Error fetching viewed product IDs:", err)
		return nil, err
	}

	fmt.Printf("User has viewed the following products: %v\n", viewedProductIDs)

	// Check if there are any viewed products
	if len(viewedProductIDs) > 0 {
		var similarFingerprints []string
		// Fetch fingerprints of users who viewed similar products
		err = conn.DB.Model(&Data.BusinessConnectUserActivity{}).
			Where("product_id IN ?", viewedProductIDs).
			Where("fingerprint_hash != ?", fingerprintHash).
			Pluck("DISTINCT fingerprint_hash", &similarFingerprints).Error
		if err != nil {
			fmt.Println("Error fetching similar fingerprints:", err)
			return nil, err
		}

		fmt.Printf("Found similar fingerprints: %v\n", similarFingerprints)

		var recommendedIDs []uint
		// Fetch recommended product IDs based on similar fingerprints
		err = conn.DB.Model(&Data.BusinessConnectUserActivity{}).
			Where("fingerprint_hash IN ?", similarFingerprints).
			Where("product_id NOT IN ?", viewedProductIDs).
			Pluck("DISTINCT product_id", &recommendedIDs).Error
		if err != nil {
			fmt.Println("Error fetching recommended product IDs:", err)
			return nil, err
		}

		fmt.Printf("Recommended product IDs: %v\n", recommendedIDs)

		// Fetch the recommended products
		if len(recommendedIDs) > 0 {
			err = conn.DB.Where("id IN ? AND publish_status = ?", recommendedIDs, "publish").
				Order("product_rank * 1000 + sales DESC").
				Limit(limit).
				Offset(offset).
				Find(&products).Error
			if err != nil {
				fmt.Println("Error fetching recommended products:", err)
				return nil, err
			}
		}
	}

	// If no products found, fallback to category-based recommendations
	if len(products) == 0 {
		fmt.Println("No recommended products found, trying category-based recommendations")
		var topCategories []string
		err = conn.DB.
			Model(&Data.BusinessConnectUserActivity{}).
			Where("fingerprint_hash = ? AND category IS NOT NULL", fingerprintHash).
			Select("DISTINCT category").
			Limit(5).
			Pluck("category", &topCategories).Error
		if err != nil {
			fmt.Println("Error fetching top categories:", err)
			return nil, err
		}

		fmt.Printf("Top categories based on user activity: %v\n", topCategories)

		query := conn.DB.Model(&Data.Post{}).Where("publish_status = ?", "publish")
		if len(topCategories) > 0 {
			query = query.Where("category IN ?", topCategories)
		}

		err = query.
			Order("product_rank * 1000 + sales DESC").
			Limit(limit).
			Offset(offset).
			Find(&products).Error
		if err != nil {
			fmt.Println("Error fetching category-based products:", err)
			return nil, err
		}
	}

	// If still no products found, fallback to best sellers
	if len(products) == 0 {
		fmt.Println("No products found, trying best-seller fallback")
		err = conn.DB.Where("publish_status = ? AND best_seller = ?", "publish", true).
			Order("product_rank * 1000 + sales DESC").
			Limit(limit).
			Offset(offset).
			Find(&products).Error
		if err != nil {
			fmt.Println("Error fetching best-seller products:", err)
			return nil, err
		}
	}

	// Fill up the result with anonymous products if fewer than 'limit'
	if len(products) < limit {
		fmt.Println("Fewer products found than limit, fetching more anonymous products")
		missing := limit - len(products)
		var topUp []Data.Post
		err = conn.DB.Where("publish_status = ?", "publish").
			Order("product_rank * 1000 + sales DESC, RAND()").
			Limit(missing).
			Find(&topUp).Error
		if err != nil {
			fmt.Println("Error fetching anonymous products:", err)
			return nil, err
		}

		// Build a set of existing product IDs
		existing := make(map[uint]bool)
		for _, p := range products {
			existing[p.ID] = true
		}

		// Filter out duplicates before appending
		for _, p := range topUp {
			if !existing[p.ID] {
				products = append(products, p)
				existing[p.ID] = true
			}
		}
	}

	// Debug: Output the final product list
	// fmt.Printf("Returning %d products: %v\n", len(products), products)

	return products, nil
}

func (d *DatabaseHelperImpl) RecommendProductsForUser(fingerprintHash string, limit, offset int) ([]Data.Post, error) {
	var finalProducts []Data.Post
	existing := make(map[uint]bool)
	target := offset + limit

	if fingerprintHash == "" {
		fmt.Println("Anonymous user — fallback to anonymous recommendations")
		all, err := d.RecommendForAnonymous(target)
		if err != nil {
			return nil, err
		}
		// Deduplicate
		for _, p := range all {
			if !existing[p.ID] {
				finalProducts = append(finalProducts, p)
				existing[p.ID] = true
				if len(finalProducts) >= target {
					break
				}
			}
		}
		return finalProducts[offset:min(offset+limit, len(finalProducts))], nil
	}

	// Step 1: Weighted scoring
	type WeightedProduct struct {
		ProductID uint
		Score     int
	}
	var weightedProducts []WeightedProduct
	err := conn.DB.Raw(`
		SELECT product_id,
		       SUM(CASE 
		            WHEN activity_type = 'search' THEN 5
		            WHEN activity_type = 'click' THEN 3
		            WHEN activity_type = 'view' THEN 1
		            ELSE 0 END) AS score
		FROM dorng_user_activities
		WHERE fingerprint_hash = ?
		  AND product_id IS NOT NULL
		  AND last_updated >= ?
		GROUP BY product_id
		ORDER BY score DESC
		LIMIT 50
	`, fingerprintHash, time.Now().AddDate(0, 0, -30)).Scan(&weightedProducts).Error
	if err != nil {
		return nil, err
	}

	var viewedProductIDs []uint
	for _, wp := range weightedProducts {
		viewedProductIDs = append(viewedProductIDs, wp.ProductID)
		existing[wp.ProductID] = true
	}

	// Step 2: Collaborative filtering
	if len(viewedProductIDs) > 0 {
		var similarFingerprints []string
		err = conn.DB.Model(&Data.BusinessConnectUserActivity{}).
			Where("product_id IN ?", viewedProductIDs).
			Where("fingerprint_hash != ?", fingerprintHash).
			Pluck("DISTINCT fingerprint_hash", &similarFingerprints).Error
		if err != nil {
			return nil, err
		}

		var recommendedIDs []uint
		err = conn.DB.Model(&Data.BusinessConnectUserActivity{}).
			Where("fingerprint_hash IN ?", similarFingerprints).
			Where("product_id NOT IN ?", viewedProductIDs).
			Pluck("DISTINCT product_id", &recommendedIDs).Error
		if err != nil {
			return nil, err
		}

		if len(recommendedIDs) > 0 {
			var collabProducts []Data.Post
			err = conn.DB.
				Where("id IN ? AND publish_status = ? AND stock_remaining > 0", recommendedIDs, "publish").
				Order("product_rank * 1000 + sales DESC").
				Find(&collabProducts).Error
			if err != nil {
				return nil, err
			}

			for _, p := range collabProducts {
				if !existing[p.ID] {
					finalProducts = append(finalProducts, p)
					existing[p.ID] = true
					if len(finalProducts) >= target {
						goto RETURN_PRODUCTS
					}
				}
			}
		}
	}

	// Step 3: Category-based
	if len(finalProducts) < target {
		var topCategories []string
		err = conn.DB.
			Model(&Data.BusinessConnectUserActivity{}).
			Where("fingerprint_hash = ? AND category IS NOT NULL", fingerprintHash).
			Where("last_updated >= ?", time.Now().AddDate(0, 0, -30)).
			Select("DISTINCT category").
			Limit(5).
			Pluck("category", &topCategories).Error
		if err != nil {
			return nil, err
		}

		if len(topCategories) > 0 {
			var categoryProducts []Data.Post
			err = conn.DB.
				Where("category IN ? AND publish_status = ? AND stock_remaining > 0", topCategories, "publish").
				Order("product_rank * 1000 + sales DESC").
				Find(&categoryProducts).Error
			if err != nil {
				return nil, err
			}

			for _, p := range categoryProducts {
				if !existing[p.ID] {
					finalProducts = append(finalProducts, p)
					existing[p.ID] = true
					if len(finalProducts) >= target {
						goto RETURN_PRODUCTS
					}
				}
			}
		}
	}

	// Step 4: Search-based
	if len(finalProducts) < target {
		var searchQueries []string
		err = conn.DB.Model(&Data.BusinessConnectUserActivity{}).
			Where("fingerprint_hash = ? AND activity_type = ?", fingerprintHash, "search").
			Where("last_updated >= ?", time.Now().AddDate(0, 0, -30)).
			Select("DISTINCT title_or_search_query").
			Limit(5).
			Pluck("title_or_search_query", &searchQueries).Error
		if err != nil {
			return nil, err
		}

		if len(searchQueries) > 0 {
			var searchProducts []Data.Post
			query := conn.DB.
				Where("publish_status = ? AND stock_remaining > 0", "publish")
			for _, q := range searchQueries {
				qLower := strings.ToLower(q)
				query = query.Or("LOWER(title) LIKE ?", "%"+qLower+"%").
					Or("LOWER(tags) LIKE ?", "%"+qLower+"%")
			}
			err = query.Order("product_rank * 1000 + sales DESC").
				Find(&searchProducts).Error
			if err != nil {
				return nil, err
			}

			for _, p := range searchProducts {
				if !existing[p.ID] {
					finalProducts = append(finalProducts, p)
					existing[p.ID] = true
					if len(finalProducts) >= target {
						goto RETURN_PRODUCTS
					}
				}
			}
		}
	}

	// Step 5: Best sellers
	if len(finalProducts) < target {
		var bestSellers []Data.Post
		err = conn.DB.
			Where("publish_status = ? AND best_seller = ? AND stock_remaining > 0", "publish", true).
			Order("product_rank * 1000 + sales DESC").
			Find(&bestSellers).Error
		if err != nil {
			return nil, err
		}
		for _, p := range bestSellers {
			if !existing[p.ID] {
				finalProducts = append(finalProducts, p)
				existing[p.ID] = true
				if len(finalProducts) >= target {
					goto RETURN_PRODUCTS
				}
			}
		}
	}

	// Step 6: Random top-up
	if len(finalProducts) < target {
		var topUp []Data.Post
		err = conn.DB.
			Where("publish_status = ? AND stock_remaining > 0", "publish").
			Order("product_rank * 1000 + sales DESC, RAND()").
			Find(&topUp).Error
		if err != nil {
			return nil, err
		}
		for _, p := range topUp {
			if !existing[p.ID] {
				finalProducts = append(finalProducts, p)
				existing[p.ID] = true
				if len(finalProducts) >= target {
					break
				}
			}
		}
	}

RETURN_PRODUCTS:
	if offset >= len(finalProducts) {
		return []Data.Post{}, nil
	}
	return finalProducts[offset:min(offset+limit, len(finalProducts))], nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *DatabaseHelperImpl) RecommendForAnonymous(limit int) ([]Data.Post, error) {
	var products []Data.Post

	err := conn.DB.
		Where("publish_status = ?", "publish").
		Order("product_rank * 1000 + sales DESC, RAND()").
		Limit(limit).
		Find(&products).Error

	return products, err
}

func (d *DatabaseHelperImpl) LogUserClickData(fingerprintHash string, productID uint, ActivityType, Category, TitleOrSearchQuery string) error {
	// Check if the activity (click) already exists for this user and product
	var activity Data.BusinessConnectUserActivity
	err := conn.DB.Where("fingerprint_hash = ? AND activity_type = ? AND product_id = ?", fingerprintHash, ActivityType, productID).First(&activity).Error
	if err == nil {
		// If the record exists, increment the ClickCount
		activity.ClickCount++
		activity.LastUpdated = time.Now() // Update the timestamp
		return conn.DB.Save(&activity).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// If the record does not exist, create a new one
		activity = Data.BusinessConnectUserActivity{
			FingerprintHash:    fingerprintHash,
			ActivityType:       ActivityType,
			ClickCount:         1, // Initial click count
			ProductID:          productID,
			Category:           Category,           // Optional: You could add the category here if available
			TitleOrSearchQuery: TitleOrSearchQuery, // Optional: You could add the search query if relevant
			LastUpdated:        time.Now(),
		}
		return conn.DB.Create(&activity).Error
	}
	return err
}
