package middlewares

import (
	"context"
	"fmt"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/dgrijalva/jwt-go"
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
		nonAuthResources := []string{"/phonebookapi/v1/account/create", "/phonebookapi/v1/healthcheck",
			"/phonebookapi/v1/authenticate"}

		requestedResource := req.URL.Path // requested resource
		for _, value := range nonAuthResources {
			if value == requestedResource {
				next.ServeHTTP(w, req)
				return
			}
		}

		// authentication for resources that are restricted
		response := make(map[string]interface{})
		tokenHeader := req.Header.Get("Authorization") // retrieve the token

		// check if token is missing
		if tokenHeader == "" {
			response = utl.Message(106, "authentication token missing")
			w.Header().Add("Content-Type", "application/json") // headers should come before status
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
		}

		// fetch the token
		splitHeader := strings.Split(tokenHeader, " ")
		if len(splitHeader) != 2 {
			response = utl.Message(106, "malformed authentication token, expected: Bearer <token>")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
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
			response = utl.Message(106, err.Error())
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
		}

		// check token validity, whether it has expired
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			response = utl.Message(106, "authentication token is not valid")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
		}

		// extract token metadata
		accessUuid, accessOk := claims["access_uuid"].(string)
		if !accessOk {
			response = utl.Message(106, "malformed authentication token, some parameters are missing")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
		}
		accountId, accIdErr := strconv.ParseUint(fmt.Sprintf("%.f", claims["account_id"]), 10, 64)
		if accIdErr != nil {
			response = utl.Message(106, "malformed authentication token, some parameters are missing")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			utl.Respond(w, response)
			return
		}

		accessTokenDetails.AccessUuid = accessUuid
		accessTokenDetails.AccountId = uint(accountId)

		// check if metadata is in redis
		accountIdFromRedis, redisErr := fetchAccessMetadata(accessTokenDetails)
		if redisErr != nil {
			log.Printf("WARNING | An error occurred while fetching from redis: %v\n", redisErr)
			response = utl.Message(105, "authentication token not recognized")
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
