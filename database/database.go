package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"

	Data "business-connect/models"
)

var DB *gorm.DB

func init() {
	dsn := os.Getenv("DATABASE_URL")

	var err error
	var db *gorm.DB

	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 1; i <= maxRetries; i++ {
		fmt.Printf("⏳ Attempt %d: Connecting to database...\n", i)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})

		if err == nil {
			sqlDB, err2 := db.DB()
			if err2 == nil && sqlDB.Ping() == nil {
				fmt.Println("✅ Connected to database successfully!")
				DB = db
				break
			}
		}

		fmt.Printf("❌ DB connection failed: %v\n", err)

		if i == maxRetries {
			fmt.Println("❌ Could not connect to database after multiple attempts")
		}

		time.Sleep(retryDelay)
	}

	// Run migrations AFTER successful connection
	DbMigration()

	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("❌ failed to get DB connection: " + err.Error())
	}

	// Pooling config
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	fmt.Println("🚀 Database fully initialized and optimized")
}

func DbMigration() {
	// AutoMigrate the user model
	err := DB.AutoMigrate(&Data.User{})
	if err != nil {
		fmt.Println("failed to migrate the User's database")
	}

	// AutoMigrate the otp model
	otpErr := DB.AutoMigrate(&Data.OTP{})
	if otpErr != nil {
		fmt.Println("failed to migrate the otp database")
	}

	// AutoMigrate the transaction history model
	err = DB.AutoMigrate(&Data.SiteVisit{})
	if err != nil {
		fmt.Println("failed to migrate the SiteVisit database")
	}

	// // AutoMigrate the jti model
	err = DB.AutoMigrate(&Data.UserConnectLimit{})
	if err != nil {
		fmt.Println("failed to migrate the UserConnectLimit to the database")
	}

	err = DB.AutoMigrate(&Data.Notification{})
	if err != nil {
		fmt.Println("failed to migrate the Notification database")
	}

	err = DB.AutoMigrate(&Data.Post{})
	if err != nil {
		fmt.Println("failed to migrate the Post database")
	}

	err = DB.AutoMigrate(&Data.PostImage{})
	if err != nil {
		fmt.Println("failed to migrate the PostImage database")
	}

	err = DB.AutoMigrate(&Data.ProfileImage{})
	if err != nil {
		fmt.Println("failed to migrate the ProfileImage database")
	}

	err = DB.AutoMigrate(&Data.GroupParticipant{})
	if err != nil {
		fmt.Println("failed to migrate the GroupParticipant database")
	}

	err = DB.AutoMigrate(&Data.Connection{})
	if err != nil {
		fmt.Println("failed to migrate the Connection database")
	}

	err = DB.AutoMigrate(&Data.Status{})
	if err != nil {
		fmt.Println("failed to migrate the Status database")
	}

	err = DB.AutoMigrate(&Data.StatusImage{})
	if err != nil {
		fmt.Println("failed to migrate the StatusImage database")
	}

	// err = DB.AutoMigrate(&Data.BusinessConnectEmailSubscriber{})
	// if err != nil {
	// 	fmt.Println("failed to migrate the BusinessConnectEmailSubscriber database")
	// }

	// err = DB.AutoMigrate(&Data.Email{})
	// if err != nil {
	// 	fmt.Println("failed to migrate the Email database")
	// }

	// err = DB.AutoMigrate(&Data.ShippingFees{})
	// if err != nil {
	// 	fmt.Println("failed to migrate the ShippingFees database")
	// }

	// err = DB.AutoMigrate(&Data.Analytics{})
	// if err != nil {
	// 	fmt.Println("failed to migrate the Analytics database")
	// }

	err = DB.AutoMigrate(&Data.JTI{})
	if err != nil {
		fmt.Println("failed to migrate the JTI database")
	}

	// err = DB.AutoMigrate(&Data.BusinessConnectDeviceFingerprint{})
	// if err != nil {
	// 	fmt.Println("failed to migrate the BusinessConnectDeviceFingerprint database")
	// }

	// err = DB.AutoMigrate(&Data.BusinessConnectUserActivity{})
	// if err != nil {
	// 	fmt.Println("failed to migrate the BusinessConnectUserActivity database")
	// }
}
