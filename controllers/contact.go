package controllers

import (
	"encoding/json"
	"github.com/cermu/Go-phoneBook-API/models"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"net/http"
)

// CreateContact public handler variable for creating/saving new contacts
var CreateContact = func(w http.ResponseWriter, req *http.Request) {
	contact := &models.Contact{}

	// decode the request body into a struct
	err := json.NewDecoder(req.Body).Decode(contact)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	// fetch account id from request context
	accountId := req.Context().Value("account").(uint)

	// save the contact passed
	response := contact.CreateContact(accountId)
	utl.Respond(w, response)
	return
}

// FetchContactsByAccountId public handler variable for fetching contacts for a specified account
var FetchContactsByAccountId = func(w http.ResponseWriter, req *http.Request) {
	contact := &models.Contact{}

	// fetch account id from request context
	accountId := req.Context().Value("account").(uint)

	// fetch contacts
	response := contact.FetchContactsByAccountId(accountId)
	utl.Respond(w, response)
	return
}
