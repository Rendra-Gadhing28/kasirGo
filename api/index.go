package handler

import (
	"net/http"

	backendHandler "kasirpro/api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	backendHandler.Handler(w, r)
}
