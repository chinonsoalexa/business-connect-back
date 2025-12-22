package models

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type OTP struct {
	CustomID                uint   `gorm:"primaryKey"`
	PasswordReset           bool   `json:"password_reset" gorm:"column:password_reset"`
	LinkWhatsapp            bool   `json:"link_whatsapp" gorm:"column:link_whatsapp"`
	EmailVerification       bool   `json:"email_verification" gorm:"column:email_verification"`
	PhoneNumberVerification bool   `json:"phone_number_verification" gorm:"column:phone_number_verification"`
	OTP                     string `json:"otp" gorm:"column:otp"`
	Email                   string `json:"email" gorm:"column:email"`
	PhoneNumber             string `json:"phone_number" gorm:"column:phone_number"`
	CreatedAT               int64  `json:"created_at" gorm:"column:created_at"`
	MaxTry                  int64  `json:"max_try" gorm:"column:max_try"`
}

type User struct {
	gorm.Model

	// Basic Info
	FullName        string `json:"full_name" gorm:"size:100;not null"`
	BusinessName    string `json:"business_name" gorm:"size:100;not null"`
	Email           string `json:"email" gorm:"uniqueIndex;not null"`
	Password        string `json:"password"` // store HASHED password only
	PhoneNumber     string `json:"phone_number" gorm:"size:15;index"`
	ProfilePhotoURL string `json:"profile_photo_url"`
	CoverPhotoURL   string `json:"cover_photo_url"`

	EmailVerified bool    `json:"email_verified" gorm:"default:false"`
	Suspended     bool    `json:"suspended" gorm:"default:false"`
	Address       string  `json:"address"`
	State         string  `json:"state"`
	City          string  `json:"city"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`

	// Business Metrics
	TotalRevenue  float64 `json:"total_revenue" gorm:"default:0"`
	TotalSales    int64   `json:"total_sales" gorm:"default:0"`
	TotalCustomer int64   `json:"total_customer" gorm:"default:0"`
	TotalProduct  int64   `json:"total_product" gorm:"default:0"`

	// Role
	UserType string `json:"user_type" gorm:"size:20;index"` // ADMIN | BUSINESS | USER

	// Auth
	RefreshToken string `json:"-"`

	// Relations
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// post struct
type (
	Post struct {
		gorm.Model

		// Ownership
		UserID uint `json:"user_id" gorm:"index"`

		// Core fields (COMMON)
		PostType     string `json:"post_type" gorm:"size:20;index"` // personal | business | group | event | ad
		Title        string `json:"title" gorm:"size:200;not null"`
		ProductUrlID string `json:"product_url_id" gorm:"size:200;not null"`
		Description  string `json:"description" gorm:"type:text;not null"`
		WhatsappURL  string `json:"whatsapp_url" gorm:"not null"`

		// Visibility
		IsSponsored bool  `json:"is_sponsored" gorm:"default:false"`
		IsActive    bool  `json:"is_active" gorm:"default:true"`
		Views       int64 `json:"views" gorm:"default:0"`
		Clicks      int64 `json:"clicks" gorm:"default:0"`

		// Location (used by business + event)
		Location string `json:"location,omitempty"`

		// BUSINESS FIELDS
		BusinessCategory string `json:"business_category,omitempty"`

		// GROUP / EVENT FIELDS
		EntryType  string `json:"entry_type,omitempty"`  // free | paid
		EntryPrice int64  `json:"entry_price,omitempty"` // store in kobo
		MaxMembers *int   `json:"max_members,omitempty"`

		// EVENT ONLY
		EventDate *time.Time `json:"event_date,omitempty"`

		// Moderation
		Approved bool `json:"approved" gorm:"default:true"`

		Images []PostImage `json:"images,omitempty" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	}
	PostImage struct {
		gorm.Model
		PostID           uint   `json:"post_id"`
		URL              string `json:"url" gorm:"column:url"`
		OriginalFilename string `json:"original_file_name" gorm:"column:original_file_name"`
	}
	ProductSearchResponse struct {
		Title        string  `json:"title"`
		ProductUrlID string  `json:"product_url_id"`
		Category     string  `json:"category"`
		SellingPrice float64 `json:"selling_price"` // Numeric price
		ImageUrl     string  `json:"image"`         // The first image URL
	}
)

type TokenClaims struct {
	gorm.Model
	jwt.RegisteredClaims
	Role string `json:"role"`
	Csrf string `json:"csrf"`
}

type JTI struct {
	// ID     uint   `gorm:"primaryKey;autoIncrement"`
	Jti    string `json:"jti" gorm:"type:varchar(255)"`
	UserID uint   `json:"user_id"`
}

type ChatMessages struct {
	gorm.Model
	ChatPhoneNumberID uint
	Text              string
	Role              string
}

// model for email subscribers
type SubscribeToEmail struct {
	gorm.Model
	Email string `json:"email" validate:"required,email"`
}

type SiteVisit struct {
	gorm.Model
	SiteVisitNumber uint
}

// Not all activities are equal. You can introduce weights:

// click = +3

// view = +1

// add_to_cart = +5

// purchase = +10

// Track actions with a weight column and use it to influence recommendations.

// blog types
type (
	Blog struct {
		gorm.Model
		UserID              uint                 `json:"user_id"`
		BlogUrlID           string               `json:"blog_url_id"`
		Title               string               `json:"title"`
		Description1        string               `json:"description1" gorm:"type:text"`
		Description2        string               `json:"description2" gorm:"type:text"`
		BlogCategory        string               `json:"blog_category"`
		Image1              string               `json:"Image1"`
		Image2              string               `json:"Image2"`
		Image3              string               `json:"Image3"`
		BlogReviewsCount    int64                `json:"blog_reviews_count"`
		CustomerBlogReviews []CustomerBlogReview `json:"customer_review,omitempty" gorm:"foreignKey:BlogID"`
		BlogImages          []BlogImage          `json:"product_image,omitempty" gorm:"foreignKey:BlogID"`
	}
	CustomerBlogReview struct {
		gorm.Model
		BlogID   uint   `json:"blog_id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Review   string `json:"review"`
		Rating   int    `json:"rating"`
		AddEmail bool   `json:"add_email"`
	}
	BlogImage struct {
		gorm.Model
		BlogID           uint   `json:"blog_id"`
		URL              string `json:"url" gorm:"column:url"`
		OriginalFilename string `json:"original_file_name" gorm:"column:original_file_name"`
	}
)

type (
	BusinessConnectEmailSubscriber struct {
		gorm.Model
		Email string `json:"email"`
	}
	Email struct {
		gorm.Model
		Subject string `json:"subject"`
		Content string `json:"content" gorm:"type:text"`
		SendTo  string `json:"send_to"`
	}
)

type (
	OrderHistory struct {
		gorm.Model
		ProducttID             uint           `json:"productt_id"`
		OrderStatus            string         `json:"order_status"`
		PaymentStatus          string         `json:"payment_status"`
		Quantity               int64          `json:"quantity"`
		OrderCost              float64        `json:"order_cost"`
		OrderNote              string         `json:"order_note"`
		OrderSubTotalCost      float64        `json:"order_sub_total_cost"`
		ShippingCost           float64        `json:"shipping_cost"`
		OrderDiscount          float64        `json:"order_discount"`
		CustomerEmail          string         `json:"customer_email"`
		CustomerFName          string         `json:"customer_fname"`
		CustomerSName          string         `json:"customer_user_sname"`
		CustomerCompanyName    string         `json:"customer_company_name"`
		CustomerState          string         `json:"customer_state"`
		CustomerCity           string         `json:"customer_city"`
		CustomerStreetAddress1 string         `json:"customer_street_address_1"`
		CustomerStreetAddress2 string         `json:"customer_street_address_2"`
		CustomerZipCode        string         `json:"customer_zip_code"`
		CustomerProvince       string         `json:"customer_province"`
		CustomerPhoneNumber    string         `json:"customer_phone_number"`
		ProductOrders          []ProductOrder `json:"product_orders,omitempty" gorm:"foreignKey:OrderHistoryID"`
	}
	ProductOrder struct {
		gorm.Model
		OrderHistoryID uint    `json:"order_history_id"`
		ProducttID     uint    `json:"productt_id"`
		ProductUrlID   string  `json:"product_url_id"`
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		NetWeight      int64   `json:"net_weight"`
		OrderCost      float64 `json:"order_cost"`
		Currency       string  `json:"currency"`
		Quantity       int64   `json:"quantity"`
		Category       string  `json:"category"`
		Image1         string  `json:"Image1"`
		Image2         string  `json:"Image2"`
	}
)

// Body to get request from client side
type (
	OrderHistoryBody struct {
		OrderStatus            string  `json:"order_status"`
		Quantity               int64   `json:"quantity"`
		OrderCost              float64 `json:"order_cost"`
		OrderNote              string  `json:"order_note"`
		OrderSubTotalCost      float64 `json:"order_sub_total_cost"`
		ShippingCost           float64 `json:"shipping_cost"`
		OrderDiscount          float64 `json:"order_discount"`
		CustomerEmail          string  `json:"customer_email"`
		CustomerFName          string  `json:"customer_fname"`
		CustomerSName          string  `json:"customer_user_sname"`
		CustomerCompanyName    string  `json:"customer_company_name"`
		CustomerState          string  `json:"customer_state"`
		CustomerCity           string  `json:"customer_city"`
		CustomerStreetAddress1 string  `json:"customer_street_address_1"`
		CustomerStreetAddress2 string  `json:"customer_street_address_2"`
		CustomerZipCode        string  `json:"customer_zip_code"`
		CustomerProvince       string  `json:"customer_province"`
		CustomerPhoneNumber    string  `json:"customer_phone_number"`
	}
	ProductOrderBody struct {
		ID           uint    `json:"ID"`
		ProductUrlID string  `json:"product_url_id"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		NetWeight    int64   `json:"net_weight"`
		OrderCost    float64 `json:"order_cost"`
		Currency     string  `json:"currency"`
		Quantity     int64   `json:"quantity"`
		Category     string  `json:"category"`
		Image1       string  `json:"Image1"`
		Image2       string  `json:"Image2"`
	}
)

type ShippingFees struct {
	gorm.Model
	EshopUserID        uint    `json:"eshop_user_id"`
	StoreName          string  `json:"store_name"`
	StoreEmail         string  `json:"store_email"`
	ShippingFeePerKm   int64   `json:"shipping_fee_per_km"`
	ShippingFeeGreater int64   `json:"shipping_fee_greater"`
	ShippingFeeLess    int64   `json:"shipping_fee_less"`
	StoreLatitude      float64 `json:"store_latitude"`
	StoreLongitude     float64 `json:"store_longitude"`
	StoreState         string  `json:"store_state"`
	StoreCity          string  `json:"store_city"`
	StateISO           string  `json:"state_iso"`
	CalculateUsingKg   bool    `json:"calculate_using_kg"`
}

// BusinessConnect Analytics
type Analytics struct {
	gorm.Model
	Month             time.Time `json:"month"`
	TotalRevenue      float64   `json:"total_revenue"`
	RevenueChange     float64   `json:"revenue_change"` // Percentage change from the previous month
	TotalSales        int64     `json:"total_sales"`
	SalesChange       float64   `json:"sales_change"` // Percentage change in sales
	TotalCustomers    int64     `json:"total_customers"`
	CustomerChange    float64   `json:"customer_change"` // Percentage change in customers
	TotalProducts     int64     `json:"total_products"`
	ProductChange     float64   `json:"product_change"` // Percentage change in products sold
	DailyVisitors     int64     `json:"daily_visitors"`
	TopProduct1ID     uint      `json:"top_product_1_id"`    // Top selling product of the month
	TopProduct2ID     uint      `json:"top_product_2_id"`    // Second top selling product
	IsRevenueBetter   bool      `json:"is_revenue_better"`   // Boolean to indicate if the month is better than the previous
	IsSalesBetter     bool      `json:"is_sales_better"`     // Boolean to indicate if the month is better than the previous
	IsCustomersBetter bool      `json:"is_customers_better"` // Boolean to indicate if the month is better than the previous
	IsProductsBetter  bool      `json:"is_products_better"`  // Boolean to indicate if the month is better than the previous
}

type BusinessConnectDeviceFingerprint struct {
	gorm.Model             // embeds ID, CreatedAt, UpdatedAt, DeletedAt
	FingerprintHash string `gorm:"size:64;uniqueIndex;not null"`
}

type BusinessConnectUserActivity struct {
	gorm.Model
	FingerprintHash    string    `gorm:"size:64;index" json:"fingerprint_hash"` // Link to the device
	ActivityType       string    `json:"activity_type"`                         // "search", "click", "view", etc.
	ClickCount         uint      `json:"click_count"`                           // Increment this field for clicks
	ProductID          uint      `json:"product_id"`                            // Optional: for product-specific activities
	Category           string    `json:"category"`                              // Optional: if activity is category-related
	TitleOrSearchQuery string    `json:"title_or_search_query"`                 // For search activity
	LastUpdated        time.Time `json:"last_updated"`                          // Time when the record was last updated
}

type ServiceMetaData struct {
	// Common fields for all types
	CancelAction  string `json:"cancel_action"`
	TransactionID string `json:"transaction_id"`
	Price         int    `json:"price"`
	Status        string `json:"status"`
	PhoneNumber   string `json:"PhoneNumber"`
	EmailID       string
}

// field name in string
var FieldNames = []string{
	"TransactionID",
	"Status",
	"Price",
	"PhoneNumber",
	"EmailID",
}

type SendSMSRequest struct {
	Sender             string `json:"sender"`
	Recipient          string `json:"recipient"`
	Content            string `json:"content"`
	Type               string `json:"type"` // "transactional" or "marketing"
	Tag                string `json:"tag,omitempty"`
	WebURL             string `json:"webUrl,omitempty"`
	UnicodeEnabled     bool   `json:"unicodeEnabled,omitempty"`
	OrganisationPrefix string `json:"organisationPrefix,omitempty"`
}

type SendSMSResponse struct {
	MessageID int64 `json:"messageId"`
}
