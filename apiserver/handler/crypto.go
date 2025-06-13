package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// 32 bytes - key
// openssl rand -hex 32
// openssl rand -base64 32
var secretKey = []byte("Y/Z5XzW63o9FCqukI2TcIvGZdcAk+T0=")

func Encrypt(plainText []byte) ([]byte, error) {
	key := secretKey
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate random nonce (12 bytes)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data
	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)
	return cipherText, nil
}

func Decrypt(cipherText []byte) ([]byte, error) {
	key := secretKey

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("cipherText too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
