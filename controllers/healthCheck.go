package controllers

import (
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"net/http"
)

var HealthCheck = func(w http.ResponseWriter, r *http.Request) {
	response := utl.Message(0, "Phone book API is up and running")
	utl.Respond(w, response)
	return
}
