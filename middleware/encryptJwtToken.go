package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
)

// Pad data to make its length a multiple of blockSize
func PadJwtToken(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Unpad removes the padding from the data
func UnpadJwtToken(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty input")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-padding], nil
}

func EncryptJwtToken(data string) (string, error) {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Printf("Failed to load .env file: %v\n", envErr)
		return "", errors.New("error loading .env file")
	}

	key := os.Getenv("JWT_ENCRYPTION_KEY")

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Pad the data to a multiple of block size
	paddedData := PadJwtToken([]byte(data), aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(paddedData))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedData)

	// Encode the ciphertext as a base64 string
	encryptedString := base64.StdEncoding.EncodeToString(ciphertext)

	return encryptedString, nil
}

func DecryptJwtToken(binaryText string) ([]byte, error) {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Printf("Failed to load .env file: %v\n", envErr)
		return nil, errors.New("error loading .env file")
	}

	key := os.Getenv("JWT_ENCRYPTION_KEY")

	// Decode the Base64 string to get the ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(binaryText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad the decrypted data
	unpaddedData, err := UnpadJwtToken(ciphertext)
	if err != nil {
		return nil, err
	}

	return unpaddedData, nil
}
