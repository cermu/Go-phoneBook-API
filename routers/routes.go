package routers

import (
	"github.com/cermu/Go-phoneBook-API/controllers"
	"net/http"
)

type route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type routes []route

var routeList = routes{
	route{
		Name:        "HealthCheck",
		Method:      "GET",
		Pattern:     "/healthcheck",
		HandlerFunc: controllers.HealthCheck,
	},
}
