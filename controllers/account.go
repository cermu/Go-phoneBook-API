package controllers

import (
	"encoding/json"
	"github.com/cermu/Go-phoneBook-API/models"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"net/http"
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
var MyAccount = func(w http.ResponseWriter, req * http.Request) {
	response := utl.Message(0, "coming soon")
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
