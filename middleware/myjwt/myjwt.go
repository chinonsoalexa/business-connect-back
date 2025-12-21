package myjwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	Data "business-connect/models"

	rand "business-connect/controllers/authentication/utils"
	dbFunc "business-connect/database/dbHelpFunc"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

const (
	RefreshTokenValidTime = time.Hour * 24 * 14 // 14 days
	AuthTokenValidTime    = time.Hour * 24 * 14 // 14 hours
	// use environmental variable to get private key
	privKeyPath = "keys/private_key.pem"
	pubKeyPath  = "keys/public_key.pem"
)

func GenerateCSRFSecrete() (string, error) {
	return rand.RandomAlphanumericString(32)
}

func InitJWT() error {
	privKeyFilePath, err := os.Open(privKeyPath)
	if err != nil {
		return fmt.Errorf("error opening private key file: %w", err)
	}
	defer privKeyFilePath.Close()

	signBytes, err := io.ReadAll(privKeyFilePath)
	if err != nil {
		return fmt.Errorf("error reading private key file: %w", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return fmt.Errorf("error parsing private key: %w", err)
	}

	pubKeyFilePath, err := os.Open(pubKeyPath)
	if err != nil {
		return fmt.Errorf("error opening public key file: %w", err)
	}
	defer pubKeyFilePath.Close()

	verifyBytes, err := io.ReadAll(pubKeyFilePath)
	if err != nil {
		return fmt.Errorf("error reading public key file: %w", err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return fmt.Errorf("error parsing public key: %w", err)
	}

	// fmt.Println("signKey:", signKey)

	return nil
}

func CreateNewTokens(ctx *fiber.Ctx, uuid, role string) (authTokenString, refreshTokenString, csrfSecrete string, err error) {
	// generate the csrf token
	csrfSecrete, err = GenerateCSRFSecrete()
	if err != nil {
		fmt.Println("Error generating the csrf secrete")
	}

	// generating the refresh token
	refreshTokenString, err = CreateRefreshTokenString(uuid, role, csrfSecrete)
	if err != nil {
		fmt.Println("Error generating the crate refresh token string")
	}

	// generating the auth token
	authTokenString, err = CreateAuthTokenString(uuid, role, csrfSecrete)
	if err != nil {
		fmt.Println("Error generating the create authentication token string")
	}

	return
}

func CheckAndRefreshTokens(oldAuthTokenString string, oldRefreshTokenString string, oldCsrfSecrete string) (newAuthTokenString, newRefreshTokenString, newCsrfSecret string, err error) {

	if oldCsrfSecrete == "" {
		log.Println("No CSRF Token!")
		// err = errors.New("Unauthorized")
	}

	// Assuming verifyKey is defined and is the correct key used to verify the tokens.

	authToken, err := jwt.ParseWithClaims(oldAuthTokenString, &Data.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	// fmt.Println("err1", err)

	// Check for parsing errors
	// log.Println("auth token:", authToken)
	if err != nil {
		log.Println("Error parsing old auth token:", err)
		err = errors.New("Unauthorized")
	}

	authTokenClaims, ok := authToken.Claims.(*Data.TokenClaims)
	if !ok || !authToken.Valid {
		log.Println("Invalid auth token")
		err = errors.New("Unauthorized")
	}

	// Verify CSRF token
	if oldCsrfSecrete != authTokenClaims.Csrf {
		log.Println("CSRF Token does not match")
		err = errors.New("Unauthorized")
	}

	if authToken.Valid {
		// log.Println("Auth token is valid")

		// Update the CSRF secret with the current token's CSRF value
		newCsrfSecret = authTokenClaims.Csrf

		// Update the refresh token's expiration
		newRefreshTokenString, err = UpdateRefreshTokenExp(oldRefreshTokenString)
		if err != nil {
			return
		}
		// fmt.Println("err2", err)

		newAuthTokenString = oldAuthTokenString
		return
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		// log.Println("Auth token is expired")
		if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			// log.Println("Auth token is expired")

			newAuthTokenString, newCsrfSecret, err = UpdateAuthTokenString(oldRefreshTokenString, oldAuthTokenString)
			if err != nil {
				return
			}

			newRefreshTokenString, err = UpdateRefreshTokenExp(oldRefreshTokenString)
			if err != nil {
				return
			}

			newRefreshTokenString, err = UpdateRefreshTokenCsrf(newRefreshTokenString, newCsrfSecret)
			return
		} else {
			log.Println("Error in auth token")
			err = errors.New("error in auth token")
			return
		}
	} else {
		log.Println("Error in auth token")
		// err = errors.New("error in auth token")
	}

	err = errors.New("Unauthorized")
	return
}

func CreateAuthTokenString(uuid string, role string, csrfSecrete string) (authTokenString string, err error) {
	authTokenExp := time.Now().Add(AuthTokenValidTime)
	authClaims := Data.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid,
			ExpiresAt: jwt.NewNumericDate(authTokenExp),
		},
		Role: role,
		Csrf: csrfSecrete,
	}
	authJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, authClaims)

	// Check for errors during signing
	if authTokenString, err = authJwt.SignedString(signKey); err != nil {
		// Add more detailed error logging here to see the actual error
		log.Println("Error signing auth token:", err)
		return "", err
	}

	return authTokenString, nil
}

func CreateRefreshTokenString(uuid string, role string, csrfString string) (refreshTokenString string, err error) {
	refreshTokenExp := time.Now().Add(RefreshTokenValidTime)
	refreshJti, err := dbFunc.DBHelper.StoreRefreshToken()
	if err != nil {
		fmt.Println("this is the Refresh Token Error: ", err)
		return
	}

	refreshClaims := Data.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshJti,
			Subject:   uuid,
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		},
		Role: role,
		Csrf: csrfString,
	}

	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func UpdateRefreshTokenExp(oldRefreshTokenString string) (newRefreshTokenString string, err error) {
	refreshToken, _ := jwt.ParseWithClaims(oldRefreshTokenString, &Data.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*Data.TokenClaims)
	if !ok {
		return
	}

	refreshTokenExp := time.Now().Add(RefreshTokenValidTime)

	refreshClaims := Data.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        oldRefreshTokenClaims.RegisteredClaims.ID,
			Subject:   oldRefreshTokenClaims.RegisteredClaims.Subject,
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		},
		Role: oldRefreshTokenClaims.Role,
		Csrf: oldRefreshTokenClaims.Csrf,
	}

	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func UpdateAuthTokenString(refreshTokenString string, oldAuthTokenString string) (newAuthTokenString, csrfSecrete string, err error) {

	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &Data.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*Data.TokenClaims)
	if !ok {
		err = errors.New("errors reading jwt claims")
	}

	// i need to call the check refresh token function to check for this users refresh token in the database i haven't written the function.
	if dbFunc.DBHelper.CheckRefreshToken(refreshTokenClaims.RegisteredClaims.ID) {
		if refreshToken.Valid {
			authToken, _ := jwt.ParseWithClaims(oldAuthTokenString, &Data.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})
			oldAuthTokenClaims, ok := authToken.Claims.(*Data.TokenClaims)
			if !ok {
				err = errors.New("errors reading jwt claims")
				return
			}

			csrfSecrete, err = GenerateCSRFSecrete()
			if err != nil {
				return
			}

			CreateAuthTokenString(oldAuthTokenClaims.RegisteredClaims.Subject, oldAuthTokenClaims.Role, csrfSecrete)
			return
		} else {
			log.Println("Refresh token has expired")
			// still need to write this database function to delete refresh token from the database
			dbFunc.DBHelper.DeleteRefreshToken(refreshTokenClaims.RegisteredClaims.ID)
			err = errors.New("Unauthorized")
			return
		}
	} else {
		log.Println("Refresh token has been revoked")
		err = errors.New("Unauthorized")
	}
	return
}

func RevokeRefreshToken(refreshTokenString string) error {
	// use the refresh token string that this function would receive to get your refresh token
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &Data.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return errors.New("could not parse refresh token with claims")
	}

	// we use the refresh token to get the refresh token claims
	refreshTokenClaims, ok := refreshToken.Claims.(*Data.TokenClaims)
	if !ok {
		return errors.New("could not read refresh token claims")
	}

	// deleting the refresh token using the method in the db package
	// would still need to write a function that delets refresh token from the database
	dbFunc.DBHelper.DeleteRefreshToken(refreshTokenClaims.RegisteredClaims.ID)

	return nil
}

func UpdateRefreshTokenCsrf(oldRefreshTokenString string, newCsrfString string) (newRefreshTokenString string, err error) {
	// get access to the old refresh token by using the parse token function
	refreshToken, err := jwt.ParseWithClaims(oldRefreshTokenString, &Data.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	// get access to the refresh token claims.
	oldRefreshTokenClaims, ok := refreshToken.Claims.(*Data.TokenClaims)
	if !ok {
		return
	}

	// refresh claims
	refreshClaims := Data.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        oldRefreshTokenClaims.RegisteredClaims.ID,
			Subject:   oldRefreshTokenClaims.RegisteredClaims.Subject,
			ExpiresAt: oldRefreshTokenClaims.RegisteredClaims.ExpiresAt,
		},
		Role: oldRefreshTokenClaims.Role,
		Csrf: newCsrfString,
	}
	// new refresh jwt
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	// new refresh token string
	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func GrabUUID(authTokenString string) (string, error) {
	authToken, err := jwt.ParseWithClaims(authTokenString, &Data.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return "", errors.New("error parsing auth token")
	}

	authTokenClaims, ok := authToken.Claims.(*Data.TokenClaims)
	if !ok {
		return "", errors.New("error fetching claims")
	}

	// fmt.Println("uuid123", authTokenClaims.Subject)
	return authTokenClaims.Subject, nil
}
