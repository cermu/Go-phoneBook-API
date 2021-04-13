package auth

import (
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// CreateToken public function that returns a JWT auth token
func CreateToken(accountId uint) (string, error) {
	jwtSecret := utl.ReadConfigs().GetString("JWT.ACCESS_SECRET") // fetch secret key

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["account_id"] = accountId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() // token valid for 15 mins only

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}
