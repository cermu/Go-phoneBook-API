package middlewares

import (
	"context"
	"errors"
	"fmt"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// AccessTokenDetails struct to unpack token details from a request
type AccessTokenDetails struct {
	AccessUuid string
	AccountId  uint
}

/*
Structure of a basic middleware

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff
		next.ServeHTTP(w, r)
	})
}
*/

// JWTAuthentication public function which is used to authenticate
// requests that are restricted.
func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// resources that do not require authentication
		params := mux.Vars(req)
		linkToken, _ := params["linkToken"]
		passwordReset := fmt.Sprintf("/phonebookapi/v1/reset/password/%s", linkToken)

		nonAuthResources := []string{"/phonebookapi/v1/account/create", "/phonebookapi/v1/healthcheck",
			"/phonebookapi/v1/authenticate", "/phonebookapi/v1/send/reset/password/link",
			passwordReset}

		requestedResource := req.URL.Path // requested resource
		for _, value := range nonAuthResources {
			if value == requestedResource {
				next.ServeHTTP(w, req)
				return
			}
		}

		// authorization for resources that are restricted
		response := make(map[string]interface{})
		accessTokenDetails, err := ExtractTokenFromRequest(req)
		if err != nil {
			response = utl.Message(106, err.Error())
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
		}

		// check if metadata is in redis
		accountIdFromRedis, redisErr := fetchAccessMetadata(accessTokenDetails)
		if redisErr != nil {
			if accountIdFromRedis == 0 {
				response = utl.Message(106, "authentication token is invalid, please make request for a new one")
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				utl.Respond(w, response)
				return
			}
			log.Printf("WARNING | No key to fetch from redis, account id is %d and message is %v\n", accountIdFromRedis, redisErr)
			response = utl.Message(105, "authentication token not recognized, please make request for a new one")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
		}

		// all went well with authentication
		// add a context variable(account) in the request
		ctx := context.WithValue(req.Context(), "account", accountIdFromRedis)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

// fetchAccessMetadata private function that fetches accountId from redis
func fetchAccessMetadata(accessTokenDetails *AccessTokenDetails) (uint, error) {
	accountId, err := utl.RedisClient().Get(accessTokenDetails.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	accId, _ := strconv.ParseUint(accountId, 10, 64)
	return uint(accId), nil
}

// ExtractTokenFromRequest public function which is used to authorize requests
// by extracting the token from request and validating/verifying it
func ExtractTokenFromRequest(req *http.Request) (*AccessTokenDetails, error) {
	tokenHeader := req.Header.Get("Authorization") // retrieve the token

	// check if token is missing
	if tokenHeader == "" {
		return nil, errors.New("authentication token missing")
	}

	// fetch the token
	splitHeader := strings.Split(tokenHeader, " ")
	if len(splitHeader) != 2 {
		return nil, errors.New("malformed authentication token, expected: Bearer <token>")
	}

	tokenString := splitHeader[1]
	accessTokenDetails := &AccessTokenDetails{}

	// verify token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(utl.ReadConfigs().GetString("JWT.ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// check token validity, whether it has expired
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("authentication token is not valid")
	}

	// extract token metadata
	accessUuid, accessOk := claims["access_uuid"].(string)
	if !accessOk {
		return nil, errors.New("malformed authentication token, some parameters are missing")
	}
	accountId, accIdErr := strconv.ParseUint(fmt.Sprintf("%.f", claims["account_id"]), 10, 64)
	if accIdErr != nil {
		return nil, errors.New("malformed authentication token, some parameters are missing")
	}

	accessTokenDetails.AccessUuid = accessUuid
	accessTokenDetails.AccountId = uint(accountId)
	return accessTokenDetails, nil
}
