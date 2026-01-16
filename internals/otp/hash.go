package otp

import "golang.org/x/crypto/bcrypt"

func HashOTP(code string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}

func VerifyOTP(hash string, code string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(code))
	return err == nil
}
