package utils

import (
	"crypto/rand"
    "math/big"
)


func RandomAlphanumericString(length int) (string, error) {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

    max := big.NewInt(int64(len(chars)))
    b := make([]byte, length)
    for i := range b {
        n, err := rand.Int(rand.Reader, max)
        if err != nil {
            return "", err
        }
        b[i] = chars[n.Int64()]
    }
    return string(b), nil
}