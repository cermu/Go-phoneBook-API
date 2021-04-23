package controllers

import (
	"encoding/json"
	"github.com/cermu/Go-phoneBook-API/auth"
	"github.com/cermu/Go-phoneBook-API/models"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// CreateAccount public handler variable for creating new users
var CreateAccount = func(w http.ResponseWriter, req *http.Request) {
	account := &models.Account{}

	// decode the request into a struct
	err := json.NewDecoder(req.Body).Decode(account)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	// create an account
	response := account.CreateAccount()
	utl.Respond(w, response)
	return
}

// MyAccount public handler variable to fetch a specific account details
var MyAccount = func(w http.ResponseWriter, req *http.Request) {
	// fetch account id from URI
	params := mux.Vars(req)
	accountId, err := strconv.Atoi(params["accountId"])
	if err != nil {
		response := utl.Message(101, "request failed, account id missing in URI")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		utl.Respond(w, response)
		return
	}

	account := &models.Account{}
	response := account.FetchAccount(uint(accountId))
	utl.Respond(w, response)
	return
}

// Authenticate public handler variable to authenticate users
var Authenticate = func(w http.ResponseWriter, req *http.Request) {
	loginDetails := &models.LoginDetails{}

	err := json.NewDecoder(req.Body).Decode(loginDetails)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	response := models.Login(loginDetails.Email, loginDetails.Password)
	utl.Respond(w, response)
	return
}

// UserLogout public handler variable to log out a logged in user
var UserLogout = func(w http.ResponseWriter, req *http.Request) {
	response := models.Logout(req)
	utl.Respond(w, response)
	return
}

// RefreshToken public handler variable to refresh JWT token
var RefreshToken = func(w http.ResponseWriter, req *http.Request) {
	mapRefreshToken := &models.MapRefreshToken{}
	err := json.NewDecoder(req.Body).Decode(mapRefreshToken)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	newTokens, newTokensErr := auth.Refresh(mapRefreshToken.RefreshToken)
	if newTokensErr != nil {
		response := utl.Message(105, newTokensErr.Error())
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		utl.Respond(w, response)
		return
	}

	response := utl.Message(0, "access_token has been refreshed")
	response["tokens"] = newTokens
	utl.Respond(w, response)
	return
}

// Deactivate public handler variable to fetch a specific account details
var Deactivate = func(w http.ResponseWriter, req *http.Request) {
	account := &models.Account{}
	response := account.DeactivateAccount(req)
	utl.Respond(w, response)
	return
}

// UpdateAccount public handler variable to make updates on a n existing account
var UpdateAccount = func(w http.ResponseWriter, req *http.Request) {
	// fetch account id from request context
	accountId := req.Context().Value("account").(uint)

	updateDetails := &models.UpdateAccountDetails{}
	err := json.NewDecoder(req.Body).Decode(updateDetails)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	// update an account
	response := models.UpdateAccount(updateDetails, accountId)
	utl.Respond(w, response)
	return
}

// ChangePassword public handler variable to enable an account password change
var ChangePassword = func(w http.ResponseWriter, req *http.Request) {
	// fetch account id from request context
	accountId := req.Context().Value("account").(uint)

	// decode json body
	changePassword := &models.ChangePassword{}
	err := json.NewDecoder(req.Body).Decode(changePassword)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	// change password
	account := &models.Account{}
	response := account.ChangePassword(changePassword, accountId)
	utl.Respond(w, response)
	return
}

// SendResetPasswordLink public handler variable to send password reset link
var SendResetPasswordLink = func(w http.ResponseWriter, req *http.Request) {
	// decode json body
	resetPassword := &models.ResetPassword{}
	err := json.NewDecoder(req.Body).Decode(resetPassword)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	// send password reset link
	response := models.SendResetPasswordLink(resetPassword)
	utl.Respond(w, response)
	return
}

// ResetPassword public handler variable to enable a password reset
var ResetPassword = func(w http.ResponseWriter, req *http.Request) {
	// fetch link token from URI
	params := mux.Vars(req)
	linkToken, ok := params["linkToken"]
	if !ok {
		response := utl.Message(102, "request failed, try again")
		utl.Respond(w, response)
		return
	}
	// decode json body
	changePassword := &models.ChangePassword{}
	decodeErr := json.NewDecoder(req.Body).Decode(changePassword)
	if decodeErr != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	// reset password
	response := models.ResetAccountPassword(linkToken, changePassword)
	utl.Respond(w, response)
	return
}
