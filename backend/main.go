package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kasirpro/internal/app"
	"kasirpro/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	fiberApp := app.SetupApp()

	// Graceful Shutdown Channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Mematikan server KasirPro secara graceful...")
		_ = fiberApp.Shutdown()
	}()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("KasirPro Backend Server berjalan pada %s\n", addr)
	if err := fiberApp.Listen(addr); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}

