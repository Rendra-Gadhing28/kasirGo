package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"kasirpro/internal/app"
)

var (
	httpHandler http.HandlerFunc
	initErr     error
	once        sync.Once
)

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		defer func() {
			if rec := recover(); rec != nil {
				initErr = fmt.Errorf("%v", rec)
			}
		}()
		fiberApp := app.SetupApp()
		httpHandler = adaptor.FiberApp(fiberApp)
	})

	if initErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"success":false,"message":"Inisialisasi Backend Gagal: %s"}`, initErr.Error())
		return
	}

	httpHandler(w, r)
}
