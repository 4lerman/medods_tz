package utils

import "net/http"

func WriteError(r http.ResponseWriter, status int, err error) {
	WriteJSON(r, status, map[string]string{"error": err.Error()})
}
