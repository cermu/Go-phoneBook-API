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
		Pattern:     "/create/account",
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
	route{
		Name:        "Logout",
		Method:      "GET",
		Pattern:     "/logout",
		HandlerFunc: controllers.UserLogout,
	},
	route{
		Name:        "Refresh",
		Method:      "POST",
		Pattern:     "/token/refresh",
		HandlerFunc: controllers.RefreshToken,
	},
	route{
		Name:        "DeactivateAccount",
		Method:      "GET",
		Pattern:     "/deactivate/account",
		HandlerFunc: controllers.Deactivate,
	},
	route{
		Name:        "UpdateAccount",
		Method:      "POST",
		Pattern:     "/update/account",
		HandlerFunc: controllers.UpdateAccount,
	},
	route{
		Name:        "ChangePassword",
		Method:      "POST",
		Pattern:     "/change/password",
		HandlerFunc: controllers.ChangePassword,
	},
	route{
		Name:        "SendResetPasswordLink",
		Method:      "POST",
		Pattern:     "/send/reset/password/link",
		HandlerFunc: controllers.SendResetPasswordLink,
	},
	route{
		Name:        "ResetPassword",
		Method:      "POST",
		Pattern:     "/reset/password/{linkToken}",
		HandlerFunc: controllers.ResetPassword,
	},
	route{
		Name:        "CreateContact",
		Method:      "POST",
		Pattern:     "/contact/create",
		HandlerFunc: controllers.CreateContact,
	},
	route{
		Name:        "FetchContactsByAccountId",
		Method:      "GET",
		Pattern:     "/fetch/account/contacts",
		HandlerFunc: controllers.FetchContactsByAccountId,
	},
	route{
		Name:        "FetchContactById",
		Method:      "GET",
		Pattern:     "/contact/{contactId}",
		HandlerFunc: controllers.FetchContactById,
	},
	route{
		Name:        "UpdateContact",
		Method:      "POST",
		Pattern:     "/update/contact/{contactId}",
		HandlerFunc: controllers.UpdateContact,
	},
}
