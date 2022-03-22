package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/IllicLanthresh/pack-and-go/internal/router"
	"github.com/labstack/echo/v4"
)

func main() {
	sv, err := router.New()
	if err != nil {
		panic(err)
	}
	defer func(sv *echo.Echo) {
		err := sv.Close()
		if err != nil {
			sv.Logger.Fatal(err)
		}
	}(sv)

	// Start router
	go func() {
		if err := sv.Start(":8989"); err != nil && err != http.ErrServerClosed {
			sv.Logger.Fatal("shutting down the router")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the router with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sv.Shutdown(ctx); err != nil {
		sv.Logger.Fatal(err)
	}
}
