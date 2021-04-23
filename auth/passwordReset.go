package auth

import (
	"crypto/rand"
	"encoding/base64"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"strconv"
	"time"
)

// ResetPasswordURLDetails struct to store details for generated reset password link
type ResetPasswordURLDetails struct {
	RandomString     string
	RandStringExpire int64
}

// generateRandomBytes private function, will generate bytes of a specified length
func generateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// generateRandomString private function, will generate a string of a specified length
func generateRandomString(length int) (string, error) {
	b, err := generateRandomBytes(length)
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateResetPasswordLink public function
func GenerateResetPasswordLink() (*ResetPasswordURLDetails, error) {
	rPD := &ResetPasswordURLDetails{}

	randString, err := generateRandomString(64)
	if err != nil {
		return nil, err
	}

	// link will be valid for 10 minutes only
	rPD.RandStringExpire = time.Now().Add(time.Minute * 10).Unix()
	rPD.RandomString = randString

	return rPD, nil
}

// SaveResetLinkMetadata public function that saves reset password link data in redis
func SaveResetLinkMetadata(accountId uint, rPD *ResetPasswordURLDetails) error {
	// convert unix to UTC
	lt := time.Unix(rPD.RandStringExpire, 0)
	now := time.Now()

	err := utl.RedisClient().Set(rPD.RandomString, strconv.Itoa(int(accountId)), lt.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}
