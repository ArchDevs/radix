package validation

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyAddress    = errors.New("address cannot be empty")
	ErrInvalidAddress  = errors.New("invalid address format")
	ErrEmptyUsername   = errors.New("username cannot be empty")
	ErrInvalidUsername = errors.New("invalid username")
	ErrEmptyPublicKey  = errors.New("public key cannot be empty")
	ErrInvalidPubKey   = errors.New("invalid public key format")
	ErrEmptyNonce      = errors.New("nonce cannot be empty")
	ErrEmptySignature  = errors.New("signature cannot be empty")
	ErrInvalidBase64   = errors.New("invalid base64 encoding")
)

func ValidateAddress(address string) error {
	if strings.TrimSpace(address) == "" {
		return ErrEmptyAddress
	}

	if len(address) < 10 || len(address) > 256 {
		return ErrInvalidAddress
	}
	return nil
}

func ValidatePublicKey(pubKey string) error {
	if strings.TrimSpace(pubKey) == "" {
		return ErrEmptyPublicKey
	}

	if len(pubKey) < 20 {
		return ErrInvalidPubKey
	}
	return nil
}

func ValidateNonce(nonce string) error {
	if strings.TrimSpace(nonce) == "" {
		return ErrEmptyNonce
	}
	return nil
}

func ValidateSignature(signature string) error {
	if strings.TrimSpace(signature) == "" {
		return ErrEmptySignature
	}
	return nil
}

func ValidateUsername(username string) error {
	if strings.TrimSpace(username) == "" {
		return ErrEmptyUsername
	}

	if match := UsernameRegex.MatchString(username); match == false {
		return ErrInvalidUsername
	}

	return nil
}

func IsValidHex(s string) bool {
	s = strings.TrimPrefix(s, "0x")
	_, err := hex.DecodeString(s)
	return err == nil
}

var AddressRegex = regexp.MustCompile(`^(0x)?[a-fA-F0-9]{40,64}$`)
var UsernameRegex = regexp.MustCompile(`^[a-fA-F0-9_]{3,16}$`)
