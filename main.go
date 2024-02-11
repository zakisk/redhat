package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/rs/zerolog"
	fileops "github.com/zakisk/redhat-server/file_ops"
	"github.com/zakisk/redhat-server/handlers"
	"github.com/zakisk/redhat-server/router"
	"golang.org/x/crypto/blake2b"
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	hasher, _ := blake2b.New256(nil)
	fileOps := fileops.NewFileOps(hasher)

	handler := handlers.NewHandlerInstance(log, fileOps)
	r := router.NewRouter(handler)
	ch := gohandlers.CORS(
		gohandlers.AllowedOrigins([]string{"*"}),
		gohandlers.AllowedHeaders([]string{"*"}),
		gohandlers.AllowedMethods([]string{"*"}),
	)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "9254"
	}
	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      ch(r.SM),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("Server is running on port: %s", port))
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal().Str("error", err.Error()).Msg("Unable to start server")
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Info().Any("Signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
