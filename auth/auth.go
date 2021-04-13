package auth

import (
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"time"
)

// AuthenticationDetails struct to store token claims
type AuthenticationDetails struct {
	AccessToken        string
	RefreshToken       string
	AccessUuid         string
	RefreshUuid        string
	AccessTokenExpire  int64
	RefreshTokenExpire int64
}

// CreateToken public function that returns a JWT auth token
func CreateToken(accountId uint) (*AuthenticationDetails, error) {
	var err error
	authDetails := &AuthenticationDetails{}
	jwtAccessSecret := utl.ReadConfigs().GetString("JWT.ACCESS_SECRET")
	jwtRefreshSecret := utl.ReadConfigs().GetString("JWT.REFRESH_SECRET")

	// access token valid for 15 minutes only
	authDetails.AccessTokenExpire = time.Now().Add(time.Minute * 15).Unix()
	authDetails.AccessUuid = uuid.NewV4().String()

	// refresh token valid for 7 days only
	authDetails.RefreshTokenExpire = time.Now().Add(time.Hour * 24 * 7).Unix()
	authDetails.RefreshUuid = uuid.NewV4().String()

	// creating access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["account_id"] = accountId
	atClaims["access_uuid"] = authDetails.AccessUuid
	atClaims["exp"] = authDetails.AccessTokenExpire
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	authDetails.AccessToken, err = at.SignedString([]byte(jwtAccessSecret))
	if err != nil {
		return nil, err
	}

	// creating refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["account_id"] = accountId
	rtClaims["refresh_uuid"] = authDetails.RefreshUuid
	rtClaims["exp"] = authDetails.RefreshTokenExpire
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	authDetails.RefreshToken, err = rt.SignedString([]byte(jwtRefreshSecret))
	if err != nil {
		return nil, err
	}

	return authDetails, nil
}
