package handler

import (
	"fmt"
	"net/http"
	"strings"
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
	defer func() {
		if rec := recover(); rec != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"success":false,"message":"Serverless Handler Panic: %v"}`, rec)
		}
	}()

	once.Do(func() {
		defer func() {
			if rec := recover(); rec != nil {
				initErr = fmt.Errorf("SetupApp Panic: %v", rec)
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

	// Fix Serverless environment properties for fasthttp / fiber adaptor
	if r.RemoteAddr == "" || !strings.Contains(r.RemoteAddr, ":") {
		r.RemoteAddr = "127.0.0.1:80"
	}

	// Restore original client URI on Vercel rewrites
	origURI := r.Header.Get("x-matched-path")
	if origURI == "" {
		origURI = r.Header.Get("x-forwarded-uri")
	}
	if origURI != "" && !strings.HasPrefix(origURI, "/api/index") {
		r.RequestURI = origURI
		if r.URL != nil {
			r.URL.Path = strings.Split(origURI, "?")[0]
		}
	} else if r.RequestURI == "" && r.URL != nil {
		r.RequestURI = r.URL.RequestURI()
	}

	httpHandler(w, r)
}
