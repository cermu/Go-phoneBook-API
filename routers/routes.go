package routers

import (
	"github.com/cermu/Go-phoneBook-API/controllers"
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var routeSlice = routes{
	route{
		Name:        "HealthCheck",
		Method:      "GET",
		Pattern:     "/healthcheck",
		HandlerFunc: controllers.HealthCheck,
	},
	route{
		Name:        "CreateAccount",
		Method:      "POST",
		Pattern:     "/account/create",
		HandlerFunc: controllers.CreateAccount,
	},
	route{
		Name:        "MyAccount",
		Method:      "GET",
		Pattern:     "/account/{accountId}",
		HandlerFunc: controllers.MyAccount,
	},
	route{
		Name:        "Authenticate",
		Method:      "POST",
		Pattern:     "/authenticate",
		HandlerFunc: controllers.Authenticate,
	},
}
