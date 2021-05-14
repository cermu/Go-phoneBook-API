package controllers

import (
	"encoding/json"
	"github.com/cermu/Go-phoneBook-API/models"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

// FetchContactById public handler variable for fetching a single contact by its id
var FetchContactById = func(w http.ResponseWriter, req *http.Request) {
	contact := &models.Contact{}

	// extract id from URI
	params := mux.Vars(req)
	contactId, err := strconv.Atoi(params["contactId"])
	if err != nil {
		response := utl.Message(101, "request failed, contact id missing in URI")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		utl.Respond(w, response)
		return
	}

	response := contact.FetchContactById(uint(contactId))
	utl.Respond(w, response)
	return

}

// UpdateContact public handler variable for updating an existing contact record
var UpdateContact = func(w http.ResponseWriter, req *http.Request) {
	contact := &models.Contact{}

	// decode the request body into a struct
	err := json.NewDecoder(req.Body).Decode(contact)
	if err != nil {
		response := utl.Message(102, "request failed, check your inputs")
		utl.Respond(w, response)
		return
	}

	// fetch contact id from URI
	params := mux.Vars(req)
	contactId, paramErr := strconv.Atoi(params["contactId"])
	if paramErr != nil {
		response := utl.Message(101, "request failed, contact id missing in URI")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		utl.Respond(w, response)
		return
	}

	// update the contact
	response := contact.UpdateContact(uint(contactId))
	utl.Respond(w, response)
	return
}

// DeleteContact public handler variable for deleting a contact record
var DeleteContact = func(w http.ResponseWriter, req *http.Request) {
	contact := &models.Contact{}

	// fetch contact id to be deleted from URI
	params := mux.Vars(req)
	contactId, err := strconv.Atoi(params["contactId"])
	if err != nil {
		response := utl.Message(101, "request failed, contact id missing in URI")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		utl.Respond(w, response)
		return
	}

	// delete the record
	response := contact.DeleteContact(uint(contactId))
	utl.Respond(w, response)
	return
}
