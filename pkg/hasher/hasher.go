package hasher

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"rest-api/internal/config"
)

var bytes = []byte{25, 34, 57, 74, 53, 35, 01, 07, 87, 90, 88, 11, 45, 77, 14, 05}
var secretCrypto = config.GetConfig().SecretCrypto

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secretCrypto))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secretCrypto))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		panic(err)
	}

	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
