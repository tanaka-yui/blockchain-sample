package test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"testing"
)

// RSAを使った公開鍵暗号化方式
func Test_RSA(t *testing.T) {
	testMessage := "Just some test message..."
	log.Println("Original message: ", testMessage)
	plainText := []byte(testMessage)

	size := 2048

	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}

	publicKey := &privateKey.PublicKey

	// 公開鍵を使って暗号化
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}
	log.Println(fmt.Sprintf("Cipher text: %x", cipherText))

	// 秘密鍵を使って複合化
	decryptedText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return
	}

	log.Println(fmt.Sprintf("Decrypted message: %s", decryptedText))
}

// AESを使った共通鍵暗号化方式
func Test_AES(t *testing.T) {
	testMessage := "Just some test message..."
	log.Println("Original message: ", testMessage)
	plainText := []byte(testMessage)

	key := []byte("abcdefghijklmnopqrstuvwxyz123456") // 256bit(32byte)

	// Create new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	// Create IV
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Printf("err: %s\n", err)
	}

	// Encrypt
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	fmt.Printf("Cipher text: %x \n", cipherText)

	// Decrpt
	decryptedText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptedText, cipherText[aes.BlockSize:])
	fmt.Printf("Decrypted text: %s\n", string(decryptedText))
}
