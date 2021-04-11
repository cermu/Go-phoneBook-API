package main

import (
	"context"
	"fmt"
	"github.com/cermu/Go-phoneBook-API/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// API server
	apiServer := &http.Server{
		Addr:    ":8081",
		Handler: routers.NewRouter(),
	}

	// start the server in a go routine
	go func() {
		fmt.Printf("INFO | Starting API server on port: 8081\n")
		err := apiServer.ListenAndServe()

		if err != nil {
			if err.Error() != "http: Server closed" {
				fmt.Printf("ERROR | Failed to start API server: %v\n", err)
				os.Exit(1)
			}
		}
	}()

	// shut down the server
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// block until a signal is received
	receivedSignal := <-ch

	fmt.Printf("WARNING | Shutting down API server %v signal received\n", receivedSignal)
	err := apiServer.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("ERROR | Failed to shut down API server: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("INFO | API server has shut down")
}
