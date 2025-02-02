package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
)

const DefaultKey = "schedulx"
const AesKeyPepper = "bridgx"

var (
	ErrEncryptFailed = errors.New("encrypt failed")
	ErrDecryptFailed = errors.New("decrypt failed")
)

func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) (string, error) {
	if len(key) == 0 {
		return "", errors.New("invalid key")
	}
	key = getKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(cryptedBase64 string, key []byte) ([]byte, error) {
	crypted, err := base64.StdEncoding.DecodeString(cryptedBase64)
	if err != nil {
		return nil, err
	}

	key = getKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func getKey(key []byte) []byte {
	var result []byte
	for {
		if len(result) < 16 {
			result = append(result, key...)
		} else {
			return result[:16]
		}
	}
}

func AESEncrypt(key, plaintext string) (text string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrEncryptFailed
			return
		}
	}()
	keyB := ensureKeyLength(key)
	block, err := aes.NewCipher(keyB)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	iv := make([]byte, blockSize)
	origData := padding([]byte(plaintext), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryptText := make([]byte, len(origData))
	blockMode.CryptBlocks(cryptText, origData)
	return base64.StdEncoding.EncodeToString(cryptText), nil
}

func AESDecrypt(key string, ct16 string) (text string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrDecryptFailed
			return
		}
	}()
	ciphertext, err := base64.StdEncoding.DecodeString(ct16)
	if err != nil {
		return "", err
	}
	keyB := ensureKeyLength(key)
	block, err := aes.NewCipher(keyB)
	if err != nil {
		return "", err
	}
	iv := make([]byte, block.BlockSize())
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = unPadding(origData)
	return string(origData), nil
}

func ensureKeyLength(key string) []byte {
	keyB := md5.Sum([]byte(key))
	return keyB[:]
}

func padding(src []byte, blockSize int) []byte {
	pad := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(src, padText...)
}

func unPadding(src []byte) []byte {
	length := len(src)
	unPad := int(src[length-1])
	return src[:(length - unPad)]
}
