package utils

import (
	"encoding/json"
	"net/http"
)

// Message public function builds json messages
func Message(code int32, description string) map[string]interface{} {
	return map[string]interface{}{"response_code": code, "response_description": description}
}

// Respond public function responds with json message
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
