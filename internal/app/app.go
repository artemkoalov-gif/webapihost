package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"webServerGo/internal/pkg/http_server"
)

func Run(port string) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := http_server.StartHttpServer(ctx, port); err != nil {
			fmt.Printf("❌ HealthServer error: %v\n", err)
		}
	}()

	go func() {
		defer wg.Done()
		<-quit
		cancel()
		close(quit)
	}()

	wg.Wait()
}
