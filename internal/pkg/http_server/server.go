package http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Payload struct {
	Name   string   `json:"name"`
	Links  []string `json:"links"`
	Status int      `json:"status"`
}

func StartHttpServer(ctx context.Context, port string) error {
	mux := http.NewServeMux()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	links := Links{
		Data: make(map[int]ResponseLinks),
	}

	mux.HandleFunc("POST /status/v1/", links.HandleAdd(logger))
	mux.HandleFunc("GET /get/v1/", links.HandleGet(logger))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      HttpMiddleware(logger, mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Println("✅ Http Server Started on " + port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("❌ ListenAndServe error: %v\n", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("🛑 Http Server Stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
