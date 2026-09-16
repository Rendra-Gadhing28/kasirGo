package handler

import (
	"fmt"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "BELUM_DISET"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "BELUM_DISET"
	}
	fmt.Fprintf(w, `{"status":"ok","message":"Vercel Go Runtime Aktif","env_db_host":"%s","env_db_user":"%s"}`, dbHost, dbUser)
}
