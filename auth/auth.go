package auth

import (
	"errors"
	"fmt"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"log"
	"strconv"
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
	TokenType          string
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

	authDetails.TokenType = "Bearer"

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

// SaveJWTMetadata public function that saves JWT metadata in redis
func SaveJWTMetadata(accountId uint, authenticationDetails *AuthenticationDetails) error {
	// convert unix to UTC
	at := time.Unix(authenticationDetails.AccessTokenExpire, 0)
	rt := time.Unix(authenticationDetails.RefreshTokenExpire, 0)

	now := time.Now()

	// save to redis
	atErr := utl.RedisClient().Set(authenticationDetails.AccessUuid, strconv.Itoa(int(accountId)), at.Sub(now)).Err()
	if atErr != nil {
		return atErr
	}

	rtErr := utl.RedisClient().Set(authenticationDetails.RefreshUuid, strconv.Itoa(int(accountId)), rt.Sub(now)).Err()
	if rtErr != nil {
		return rtErr
	}

	return nil
}

// DeleteAuthenticationDetails public function that is called
// when a user logs out to invalidate JWT token
func DeleteAuthenticationDetails(uuid string) (int64, error) {
	deleted, err := utl.RedisClient().Del(uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// Refresh public function that refreshes access_token using refresh_token
// when 15 minutes are over and user is still active
func Refresh(refreshToken string) (map[string]interface{}, error) {
	// verify the token
	if refreshToken == "" {
		return nil, errors.New("refresh_token missing")
	}
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(utl.ReadConfigs().GetString("JWT.REFRESH_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// check token validity
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errors.New("refresh_token is not valid")
	}

	// extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, refreshOk := claims["refresh_uuid"].(string)
		if !refreshOk {
			return nil, errors.New("malformed refresh_token, some parameters are missing")
		}
		accountId, accIdErr := strconv.ParseUint(fmt.Sprintf("%.f", claims["account_id"]), 10, 64)
		if accIdErr != nil {
			return nil, errors.New("malformed refresh_token, some parameters are missing")
		}

		// delete the old refresh_token
		deleted, delErr := DeleteAuthenticationDetails(refreshUuid)
		if delErr != nil || deleted == 0 {
			log.Printf("WARNING | The following error occurred while refreshing access token: %v\n", delErr)
			return nil, errors.New("failed to refresh access token, try again")
		}

		// create new refresh and access token
		authDetails, authDetailsErr := CreateToken(uint(accountId))
		if authDetailsErr != nil {
			return nil, authDetailsErr
		}

		// save metadata to redis
		saveErr := SaveJWTMetadata(uint(accountId), authDetails)
		if saveErr != nil {
			log.Printf("WARNING | The following error occurred while saving refresh metadata to redis: %v\n",
				saveErr)
			return nil, saveErr
		}

		// return the new tokens to the caller
		tokens := map[string]interface{}{
			"access_token":  authDetails.AccessToken,
			"refresh_token": authDetails.RefreshToken,
			"type":          authDetails.TokenType,
		}
		return tokens, nil
	} else {
		return nil, errors.New("refresh_token is not valid")
	}
}
