package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// Đây là hàm mã hóa dữ liệu cho thông tin trong claims
func EncryptAES(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block) // block và aesGCM đều bắt buộc phải có mã hóa và giải mã.
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize()) // nonce viết tắt number used once, có nghĩa đây là những chuỗi ngẫu nhiên.
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}
	// cipherText = nonce + plaintext + nonce
	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil) // Seal dùng để mã hóa , còn Open dùng để giải mã
	return base64.URLEncoding.EncodeToString(cipherText), nil

}

func DecryptAES(cipherBase64 string, key []byte) ([]byte, error) {
	cipherText, err := base64.URLEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()

	nonce := cipherText[:nonceSize]
	cipherText = cipherText[nonceSize:]
	return aesGCM.Open(nil, nonce, cipherText, nil)

}
