package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"transaction_system/app/lib/db"
	"transaction_system/app/routes"
	"transaction_system/cmd"

	"github.com/julienschmidt/httprouter"
)

var ctx = context.Background()

func init() {
	// Call setupDBConnection during initialization
	cmd.SetupDBConnection()
}

func main() {

	defer db.Close()

	// Set up routes
	router := httprouter.New()

	// Initialize application routes
	routes.InitRoutes(router)

	// Set up HTTP server
	port := 8080
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server is running on :%d...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for an interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown:", err)
	}
	log.Println("Server exiting")
}
