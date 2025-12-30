package router

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/gofiber/storage/redis"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	ai "business-connect/ai"
	"business-connect/controllers/authentication"
	"business-connect/controllers/blog"
	email "business-connect/controllers/emails"
	"business-connect/controllers/home"
	"business-connect/controllers/order"
	upload "business-connect/controllers/post"
	"business-connect/controllers/profile"
	initTrans "business-connect/paystack/initTransactionForPaystack"
	webHook "business-connect/paystack/webhooks"

	mid "business-connect/middleware"

	"github.com/joho/godotenv"
)

func NotAuthMiddleware(c *fiber.Ctx) error {

	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Printf("Failed to load .env file: %v\n", envErr)
	}

	allowedOrigins := map[string]bool{
		"https://business-connect-eta.vercel.app/": true,
		"https://payuee.shop":                      true,
	}

	allowedIPs := map[string]bool{
		"52.31.139.75":  true,
		"52.49.173.169": true,
		"52.214.14.220": true,
	}

	origin := c.Get("Origin")
	apiKey := c.Get("X-BUSCONNECT-APP-API-KEY")
	ip := c.IP() // Extract the IP address of the request

	// fmt.Println("this is the origin: ", origin)
	// fmt.Println("this is the apiKey: ", apiKey)
	// fmt.Println("this is the IP: ", ip)

	if origin != "" {
		// Web client request
		if !allowedOrigins[origin] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}
	} else if apiKey != "" {
		// App client request
		if apiKey != os.Getenv("PAYUEE_APP_API_KEY") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}
	} else if !allowedIPs[ip] {
		// Check if the IP address is in the allowed list
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	} else {
		// No valid Origin, API Key, or IP found
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Next()
}

func AuthMiddleware(c *fiber.Ctx) error {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Printf("Failed to load .env file: %v\n", envErr)
	}
	allowedOrigins := map[string]bool{
		"https://business-connect-eta.vercel.app/": true,
		"https://payuee.shop":                      true,
	}

	origin := c.Get("Origin")
	apiKey := c.Get("X-BUSCONNECT-APP-API-KEY")
	fmt.Println("AuthMiddleware invoked")
	fmt.Println("this is the origin: ", origin)
	fmt.Println("this is the apiKey: ", apiKey)

	if origin != "" {
		// Web client request
		fmt.Println("Processing web client request")
		if !allowedOrigins[origin] {
			fmt.Println("Origin not allowed: ", origin)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}

		// Perform authentication and authorization check for web client
		err := mid.WebRequireAuth(c)
		if err != nil {
			fmt.Println("WebRequireAuth failed: ", err)
			return err // Return the error if authentication fails
		}
	} else if apiKey != "" {
		// App client request
		fmt.Println("Processing app client request")
		expectedApiKey := os.Getenv("PAYUEE_APP_API_KEY")
		if apiKey != expectedApiKey {
			fmt.Println("API key mismatch: ", apiKey)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}

		// Perform authentication and authorization check for app client
		// err := mid.AppRequireAuth(c)
		// if err != nil {
		// 	fmt.Println("AppRequireAuth failed: ", err)
		// 	return err // Return the error if authentication fails
		// }
	} else {
		// No valid Origin or API Key found
		fmt.Println("No valid Origin or API Key found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	fmt.Println("AuthMiddleware completed successfully")
	return c.Next()
}

func Routers() *fiber.App {
	// Create a new Fiber application
	router := fiber.New(fiber.Config{
		// Adjust the maximum header size limit (default is 8 KB)
		// You can increase it as needed
		// BodyLimit: 4 * 1024 * 1024, // 16 KB, for example
		ReadBufferSize:  50 * 4096,
		WriteBufferSize: 2 * 4096,
		Prefork:         true, // Enable prefork mode for better performance
		AppName:         "Business Connect API",
	})

	// Configure the rate limiter
	limiterConfig := limiter.Config{
		Max:        60,              // Maximum number of requests
		Expiration: 1 * time.Minute, // Time duration before the rate limit is reset
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Limit based on client IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too Many Requests",
			})
		},
	}

	// Use the rate limiter middleware
	router.Use(limiter.New(limiterConfig))

	// Configure CORS.
	CORSconfig := cors.Config{
		// AllowOrigins:     "*", // Use a single string, not an array
		AllowOrigins:     `https://business-connect-eta.vercel.app/, https://payuee.shop`,
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE",
		// AllowHeaders:     "Content-Type, X-DORNG-APP-API-KEY",
	}
	router.Use(cors.New(CORSconfig))

	// Use the Fiber gzip middleware to compress responses.
	router.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression, // 2
	}))

	router.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Panic recovered:", r)
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error", "panic": r})
			}
		}()
		return c.Next()
	})

	// Create Redis storage
	    // Redis storage configured in code
    // redisStore := redis.New(redis.Config{
    //     Host:     "127.0.0.1",  // e.g., "127.0.0.1" or Render Redis host
    //     Port:     6379,               // default Redis port
    //     Password: "",  // leave empty if none
    //     Database: 0,                   // Redis DB index
    //     PoolSize: 10,                  // number of connections
    // })

	// securing all the web endpoint from being accessible to app cause of the origin is not included in the app requests

	// payuee web authentication using email and password
	router.Post("/sign-up", NotAuthMiddleware, authentication.SignUp)
	// CACHED ROUTE
	router.Get("/get-states-cities/:countryCode", cache.New(cache.Config{
		// Storage: redisStore,
		Expiration: 0, // 0 means it never expires
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.OriginalURL() // cache by URL
		},
		CacheControl: true,
	}), NotAuthMiddleware, authentication.GetStatesAndCitiesByCountryCode)
	router.Get("/get-states-cities/:countryCode", NotAuthMiddleware, authentication.GetStatesAndCitiesByCountryCode)
	router.Post("/email-verification", NotAuthMiddleware, authentication.EmailAuthentication)
	router.Post("/resend-otp", NotAuthMiddleware, authentication.ResendEmailVerification)
	router.Post("/sign-in", NotAuthMiddleware, authentication.SignIn)
	router.Post("/forgotten-password-email", NotAuthMiddleware, authentication.SendEmailPasswordChange)
	router.Post("/forgotten-password-verification", NotAuthMiddleware, authentication.VerifyForgotPassword)
	router.Get("/log-out", NotAuthMiddleware, authentication.Logout)

	// this is the magic login routes
	router.Post("/magic-link", NotAuthMiddleware, authentication.MagicLinkSignIn)
	router.Post("/verify/magic-link", NotAuthMiddleware, authentication.VerifySignInMagicLink)

	// web := router.Group("/web", AuthMiddleware)
	// get and update profile information
	// router.Get("/profile", mid.WebRequireAuth, profile.RetrievePersonalInformation)
	// router.Post("/profile/update", mid.WebRequireAuth, profile.UpdatePersonalInformation)
	// router.Post("/profile/update/password", mid.WebRequireAuth, profile.UpdatePassword)

	// all this are open routes to get products and make order
	router.Get("/next-product/:id", NotAuthMiddleware, profile.GetNextProductID)
	router.Get("/previous-product/:id", NotAuthMiddleware, profile.GetPreviousProductID)
	router.Get("/search-products", NotAuthMiddleware, profile.SearchProductsByTitleAndCategory)
	router.Post("/admin-product-search", NotAuthMiddleware, profile.SearchAdminProductsByTitle)
	router.Post("/admin-order-search", NotAuthMiddleware, profile.SearchAdminOrderByTitle)
	router.Post("/transaction/date", NotAuthMiddleware, profile.GetTransactionHistoryByDate)
	router.Post("/place-order", NotAuthMiddleware, order.AddOrder)
	router.Get("/get-order/:orderID", NotAuthMiddleware, order.GetOrder)

	// get dorng home products
	router.Get("/business-connect-product-home", NotAuthMiddleware, profile.GetBusinessConnectHomePageProducts)

	// get blog post by id
	router.Get("/get-blog/:blogID", NotAuthMiddleware, blog.GetBlogPost)
	router.Get("/get-blog-posts/:idLimit", NotAuthMiddleware, blog.GetBlogPosts)

	// make and group payments with paystack
	paystackGroup := router.Group("/paystack")

	// initialize transaction & and webhook
	paystackGroup.Get("/init-transaction/call-back", initTrans.PaystackCallbackHandler)
	paystackGroup.Post("/webhook/call-back", webHook.WebHookStatus)

	// get all subscriptions for auto renewal and update
	// router.Get("/subscription/:idLimit", mid.WebRequireAuth, profile.GetSubscriptionHistoryByLimit)           //get subscriptions
	// router.Get("/cancel-subscription/:subscriptionID", mid.WebRequireAuth, profile.UpdateSubscriptionStatus)  //cancel subscriptions
	// router.Get("/recharge-subscription/:subscriptionID", mid.WebRequireAuth, profile.RechargeSubscriptionNow) //recharge subscriptions now

	// analytics and add email subscribers
	router.Post("/email-subscriber", NotAuthMiddleware, profile.AddEmailSubscription)

	// router.Get("/dorng-analytics", NotAuthMiddleware, profile.AddSiteVisit)

	// ADMIN ROUTES

	// Get BusinessConnect Users Analytics
	router.Get("/get-dorng-analytics", mid.WebRequireAuth, home.GetBusinessConnectAnalytics)

	// post a product on BusinessConnect
	router.Post("/publish-product", NotAuthMiddleware, mid.WebRequireAuth, upload.CreatePost)
	router.Post("/upload-profile-photo", NotAuthMiddleware, mid.WebRequireAuth, upload.UpdateProfilePhoto)

	// set shipping fee
	router.Post("/set-shipping-fee", mid.WebRequireAuth, order.SetShippingPricePerKm)
	router.Get("/get-shipping-fee", NotAuthMiddleware, order.GetShippingPricePerKm)

	// AI GENERATION FOR PAYUEE VENDORS
	router.Post("/ai-description", NotAuthMiddleware, mid.WebRequireAuth, ai.GetVendorProductDescriptionAI)
	router.Post("/ai-tag", NotAuthMiddleware, mid.WebRequireAuth, ai.GetVendorProductTagAI)

	// update BusinessConnect product and status
	router.Post("/update-dorng-product", mid.WebRequireAuth, order.UpdateBusinessConnectProduct)
	router.Post("/update-dorng-status", mid.WebRequireAuth, order.UpdateBusinessConnectOrderStatus)

	// get all products and product by id
	router.Get("/product/:id", NotAuthMiddleware, profile.GetBusinessConnectProductByID)
	router.Get("/admin-product/:id", NotAuthMiddleware, profile.GetBusinessConnectAdminProductByID)
	router.Get("/products/:page", NotAuthMiddleware, profile.GetBusinessConnectProductsByLimit)
	router.Get("/admin-products/:idLimit", NotAuthMiddleware, profile.GetBusinessConnectAdminProductsByLimit)
	router.Post("/post-comment", NotAuthMiddleware, upload.AddBusinessConnectProductComment)
	router.Get("/get-comment/:idLimit/:proId", NotAuthMiddleware, upload.GetBusinessConnectProductCommentsByLimit)

	// Business Connect
	router.Get("/posts", NotAuthMiddleware, mid.WebRequireAuth, profile.GetPostsPaginated)
	router.Get("/status", NotAuthMiddleware, mid.WebRequireAuth, profile.GetStatusPaginated)
	router.Get("/get-friends", NotAuthMiddleware, mid.WebRequireAuth, profile.GetFriends)
	router.Post("/connect-friends", NotAuthMiddleware, mid.WebRequireAuth, profile.ConnectFriend)
	router.Get("/get-groups", NotAuthMiddleware, mid.WebRequireAuth, profile.GetGroups)
	router.Post("/join-groups", NotAuthMiddleware, mid.WebRequireAuth, profile.JoinGroupHandler)

	// blog post, retrieval and updating
	router.Post("/publish-blog", mid.WebRequireAuth, upload.BlogPost)
	router.Post("/update-dorng-blog", mid.WebRequireAuth, blog.UpdateBusinessConnectBlog)
	router.Get("/delete-dorng-blog/:blogID", mid.WebRequireAuth, blog.DeleteBusinessConnectBlog)
	router.Post("/blog-comment", NotAuthMiddleware, blog.AddBusinessConnectBlogComment)
	router.Get("/get-blog-comment/:idLimit/:proId", NotAuthMiddleware, blog.GetBusinessConnectBlogCommentsByLimit)

	// send emails and analytics
	router.Get("/dorng-analytics", NotAuthMiddleware, profile.AddSiteVisit)
	router.Get("/dorng-user-fingerprint/:fingerprint", NotAuthMiddleware, profile.GetBusinessConnectUserByFingerprint)
	router.Post("/dorng-user-analytics", NotAuthMiddleware, profile.AddClickHistory)
	router.Post("/send-dorng-email", mid.WebRequireAuth, email.SendEmails)

	// delete dorng product
	router.Get("/delete-dorng-product/:productID", mid.WebRequireAuth, order.DeleteBusinessConnectProduct)

	// get dorng order
	router.Get("/get-dorng-order/:orderID", NotAuthMiddleware, order.GetBusinessConnectOrder)
	router.Get("/get-orders/:orderLimit", NotAuthMiddleware, order.GetBusinessConnectOrdersByLimit)

	// router.Get("/send-sms/:phone", NotAuthMiddleware, order.SendSmsBusinessConnect)

	// check if user is authenticated
	router.Get("/auth-status", mid.WebRequireAuth, profile.CheckAuthStatus)

	// AI GENERATION FOR DORNG GLOBAL
	router.Post("/ai-description", NotAuthMiddleware, mid.WebRequireAuth, ai.GetVendorProductDescriptionAI)
	router.Post("/ai-tag", NotAuthMiddleware, mid.WebRequireAuth, ai.GetVendorProductTagAI)

	// Handle preflight requests (OPTIONS)
	router.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	return router
}
