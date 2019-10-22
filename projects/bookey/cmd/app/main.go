package main

import (
	"fmt"
	"net/http"

	"github.com/Wallruzz9114/bookey/app/routes"
	"github.com/Wallruzz9114/bookey/app/server"
	"github.com/Wallruzz9114/bookey/config"
	lr "github.com/Wallruzz9114/bookey/util/logger"
)

func main() {
	appConfig := config.AppConfig()
	logger := lr.New(appConfig.Debug)
	application := server.New(logger)
	appRouter := routes.New(application)
	address := fmt.Sprintf(":%d", appConfig.Server.Port)

	logger.Info().Msgf("Starting server %v", address)

	server := &http.Server{
		Addr:         address,
		Handler:      appRouter,
		ReadTimeout:  appConfig.Server.TimeoutRead,
		WriteTimeout: appConfig.Server.TimeoutWrite,
		IdleTimeout:  appConfig.Server.TimeoutIdle,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server startup failed")
	}
}

// Greet ...
func Greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
