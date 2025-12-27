package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"

	// "gorm.io/driver/postgres"
	"gorm.io/gorm"

	Data "business-connect/models"
)

var DB *gorm.DB

func init() {
	if os.Getenv("RENDER") == "" {
		// Local development only
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env file: %v\n", err)
		}
	}


	dsn := os.Getenv("DATABASE_URL")

	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect to the database")
	}

	DbMigration()

	// Get the underlying sql.DB object
	sqlDB, err := DB.DB()

	// Set the maximum number of idle connections.
	sqlDB.SetMaxIdleConns(30)

	// Set the maximum number of open connections.
	sqlDB.SetMaxOpenConns(200)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		// Handle the error
		panic("failed to create multiple connections")
	}

	fmt.Println("Successfully connected to the data base")
}

func DbMigration() {
	// AutoMigrate the user model
	err := DB.AutoMigrate(&Data.User{})
	if err != nil {
		panic("failed to migrate the User's database")
	}

	// AutoMigrate the otp model
	otpErr := DB.AutoMigrate(&Data.OTP{})
	if otpErr != nil {
		panic("failed to migrate the otp database")
	}

	// AutoMigrate the transaction history model
	// err = DB.AutoMigrate(&Data.Post{})
	// if err != nil {
	// 	panic("failed to migrate the Post database")
	// }

	// // AutoMigrate the jti model
	// err = DB.AutoMigrate(&Data.OrderHistory{})
	// if err != nil {
	// 	panic("failed to migrate the OrderHistory to the database")
	// }

	// err = DB.AutoMigrate(&Data.UserProfile{})
	// if err != nil {
	// 	panic("failed to migrate the UserProfile to the database")
	// }

	err = DB.AutoMigrate(&Data.Post{})
	if err != nil {
		panic("failed to migrate the Post database")
	}

	err = DB.AutoMigrate(&Data.PostImage{})
	if err != nil {
		panic("failed to migrate the PostImage database")
	}

	err = DB.AutoMigrate(&Data.ProfileImage{})
	if err != nil {
		panic("failed to migrate the ProfileImage database")
	}

	err = DB.AutoMigrate(&Data.GroupParticipant{})
	if err != nil {
		panic("failed to migrate the GroupParticipant database")
	}

	err = DB.AutoMigrate(&Data.Connection{})
	if err != nil {
		panic("failed to migrate the Connection database")
	}

	// err = DB.AutoMigrate(&Data.SubscribeToEmail{})
	// if err != nil {
	// 	panic("failed to migrate the SubscribeToEmail database")
	// }

	// err = DB.AutoMigrate(&Data.SiteVisit{})
	// if err != nil {
	// 	panic("failed to migrate the SiteVisit database")
	// }

	// err = DB.AutoMigrate(&Data.BusinessConnectEmailSubscriber{})
	// if err != nil {
	// 	panic("failed to migrate the BusinessConnectEmailSubscriber database")
	// }

	// err = DB.AutoMigrate(&Data.Email{})
	// if err != nil {
	// 	panic("failed to migrate the Email database")
	// }

	// err = DB.AutoMigrate(&Data.ShippingFees{})
	// if err != nil {
	// 	panic("failed to migrate the ShippingFees database")
	// }

	// err = DB.AutoMigrate(&Data.Analytics{})
	// if err != nil {
	// 	panic("failed to migrate the Analytics database")
	// }

	err = DB.AutoMigrate(&Data.JTI{})
	if err != nil {
		panic("failed to migrate the JTI database")
	}

	// err = DB.AutoMigrate(&Data.BusinessConnectDeviceFingerprint{})
	// if err != nil {
	// 	panic("failed to migrate the BusinessConnectDeviceFingerprint database")
	// }

	// err = DB.AutoMigrate(&Data.BusinessConnectUserActivity{})
	// if err != nil {
	// 	panic("failed to migrate the BusinessConnectUserActivity database")
	// }
}
