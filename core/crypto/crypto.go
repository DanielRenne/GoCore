// Package provides some basic DES_CBC or CBC_PKCS7 encryption and decryption functions.
package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"strings"
)

// DecryptDES_CBC will decrypt a string using DES_CBC encryption.
func DecryptDES_CBC(data []byte, key string, iv []byte) string {

	block, err := des.NewCipher([]byte(key))

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("%d bytes NewCipher key with block size of %d bytes\n", len(key), block.BlockSize)
	return DecryptDES_CBC_Byte(data, block, iv)
}

// DecryptDES_CBC_Byte will decrypt an []byte using DES_CBC encryption.
func DecryptDES_CBC_Byte(encrypted []byte, block cipher.Block, iv []byte) string {

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(encrypted, encrypted)
	return strings.Trim(string(encrypted), string(0x03))

}

// EncryptDES_CBC_PKCS7 will encrypt a string using DES_CBC_PKCS7 encryption.
func EncryptDES_CBC_PKCS7(plainText []byte, key string, iv []byte) string {

	block, err := des.NewCipher([]byte(key))

	if err != nil {
		fmt.Println(err)
	}

	paddedPlainText, paddingError := pkcs7Pad(plainText, 8)

	if paddingError != nil {
		fmt.Printf("Error at core.crypto.EncryptDES_CBC_PKCS7 Padding PKCS7:  " + paddingError.Error())
	}

	return EncryptDES_CBC_Byte(paddedPlainText, block, iv)
}

// EncryptDES_CBC_Byte will encrypt an []byte using DES_CBC encryption.
func EncryptDES_CBC_Byte(plainText []byte, block cipher.Block, iv []byte) string {

	encrypted := make([]byte, len(plainText))

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(encrypted, plainText)
	return string(encrypted)

}

// Appends padding.
func pkcs7Pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(data, pad...), nil
}
