package routers

import (
	"github.com/cermu/Go-phoneBook-API/middlewares"
	"github.com/gorilla/mux"
)

// NewRouter public function that returns a pointer to mux.Router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middlewares.EnableCORS)        // Attach the EnableCORS middleware
	router.Use(middlewares.JWTAuthentication) // Attach the JWTAuthentication middleware
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
