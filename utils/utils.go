package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func URI(action string) string {
	return "/job"
}

func AESGCMEncrypt(plaintextStr string, keyStr string) (string, string) {
	// 将明文和密钥转换为字节切片
	plaintext := []byte(plaintextStr)
	key := []byte(keyStr)

	// 创建加密分组
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error()))
	}

	// 创建 GCM 模式的 AEAD
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	// 创建随机数,这里在实际应用中让它只生成一次.不然每次都需要进行修改
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	// 生成密文
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	// 返回密文及随机数的 base64 编码
	fmt.Println(ciphertext, nonce, key)
	return base64.RawURLEncoding.EncodeToString(ciphertext), base64.RawURLEncoding.EncodeToString(nonce)
}

func AESGCMDecrypt(ciphertextStr string, keyStr string, nonceStr string) string {
	// 将密文,密钥和生成的随机数转换为字节切片
	ciphertext, _ := base64.RawURLEncoding.DecodeString(ciphertextStr)
	nonce, _ := base64.RawURLEncoding.DecodeString(nonceStr)
	key := []byte(keyStr)

	// 创建加密分组
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	// 创建 GCM 模式的 AEAD
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	// 明文内容
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext)
}
