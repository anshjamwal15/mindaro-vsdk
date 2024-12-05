package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	server "github.com/aimbot1526/mindaro-vsdk/cmd"
)

func main() {
	// Start the server
	srv := server.StartServer()

	// Set up a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-stop

	// Gracefully shutdown the server
	log.Println("Signal received. Shutting down...")
	srv.ShutdownServer()
}
