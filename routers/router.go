package routers

import (
	"github.com/cermu/Go-phoneBook-API/middlewares"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middlewares.EnableCORS) // Attach the EnableCORS middleware
	api := router.PathPrefix("/phonebookapi/v1").Subrouter()

	for _, route := range routeSlice {
		api.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}


