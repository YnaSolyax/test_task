package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"test_task/internal/server"
	"test_task/internal/storage"
	"time"
)

func main() {
	store := storage.NewInmemory()
	serv := server.NewServer(":8080", store)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := serv.Start(); err != nil {
			fmt.Printf("Error: start server: %v\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serv.Stop(shutdownCtx); err != nil {
		fmt.Printf("Error: stop server: %v\n", err)
	}
}
