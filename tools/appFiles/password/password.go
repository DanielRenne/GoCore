package password

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"unicode"

	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	secret = "GoCorePasswordSecret"
)

var secretKeyBytes []byte

func init() {
	var err error
	secretKeyBytes, err = hex.DecodeString(secret)
	if err != nil {
		session_functions.Log("password.go", "Failed to decode secret key for password encryption and decryption.")
		return
	}
}

func Encrypt(value string) (string, error) {

	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		session_functions.Log("password.go", "Failed to read nonce on password.Encrypt:  "+err.Error())
		return "", err
	}

	// This encrypts value and appends the result to the nonce.
	encrypted := secretbox.Seal(nonce[:], []byte(value), &nonce, &secretKey)
	return string(encrypted[:]), nil

}

func Decrypt(value string) (string, error) {

	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)
	// When you decrypt, you must use the same nonce and key you used to
	// encrypt the message. One way to achieve this is to store the nonce
	// alongside the encrypted message. Above, we stored the nonce in the first
	// 24 bytes of the encrypted text.
	var decryptNonce [24]byte
	copy(decryptNonce[:], value[:24])
	decrypted, ok := secretbox.Open([]byte{}, []byte(value[24:]), &decryptNonce, &secretKey)

	if !ok {
		return "", errors.New("Failed to decrypt value")
	}

	return string(decrypted), nil

}

func EncryptBase64(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	encrypted, err := Encrypt(value)
	if err != nil {
		return "", err
	}
	encryptedBase64 := base64.StdEncoding.EncodeToString([]byte(encrypted))
	return encryptedBase64, nil
}

func DecryptBase64(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	decodedValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return Decrypt(string(decodedValue[:]))
}

// Verifies that the password has 8 characters, and at least 1 number, upper, lower, and special char within the password.
func VerifyPassword(s string) (number, upper, lower, special, eightOrMore bool) {
	letters := 0
	number, upper, lower, special, eightOrMore = false, false, false, false, false
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			number = true
			letters++
		case unicode.IsUpper(s):
			upper = true
			letters++
		case unicode.IsLower(s):
			lower = true
			letters++
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
			letters++
		case unicode.IsLetter(s) || s == ' ':
			letters++
		default:
			//return false, false, false, false
		}
	}
	eightOrMore = letters >= 8

	return number, upper, lower, special, eightOrMore
}
