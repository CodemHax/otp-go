package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {

	otp, err := rand.Int(rand.Reader, big.NewInt(1000000))

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", otp.Int64()), nil
}
