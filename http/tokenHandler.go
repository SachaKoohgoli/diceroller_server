package http

import (
	accesstoken "diceroller_server/access_token"
	"encoding/json"
	"net/http"
	"time"
)

// +HttpToken+ is the web-only data class for the token details.
// This should not be used internally.
type HttpToken struct {
	Token string `json:"token"`
}

// Generate the web token, return 200 in a +HttpToken+
func HandleTokenGeneration(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(HttpToken{accesstoken.GenerateToken(time.Now().UTC()).EncodedVal})
}
