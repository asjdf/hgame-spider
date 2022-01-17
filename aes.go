package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func PKCS7Padding(org []byte, blockSize int) []byte {
	pad := blockSize - len(org)%blockSize
	padArr := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(org, padArr...)
}

func PKCS7UnPadding(org []byte) []byte {
	l := len(org)
	pad := org[l-1]
	return org[:l-int(pad)]
}

func AESDecrypt(cipherTxt []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockMode := cipher.NewCBCDecrypter(block, key)
	org := make([]byte, len(cipherTxt))
	blockMode.CryptBlocks(org, cipherTxt)
	org = PKCS7UnPadding(org)
	return org
}

func AESEncrypt(org []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	org = PKCS7Padding(org, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	encrypted := make([]byte, len(org))
	blockMode.CryptBlocks(encrypted, org)
	return encrypted

}
