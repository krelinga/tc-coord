package main

// spell-checker:ignore pbconnect connectrpc tccoord tccoordv1

import (
	"log"
	"net/http"

	pbconnect "buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/tccoord/v1/tccoordv1connect"
)

func main() {
	// Create a new service instance
	service := &tcCoord{}

	// Set up the ConnectRPC server
	mux := http.NewServeMux()
	path, handler := pbconnect.NewTCCoordServiceHandler(service)
	mux.Handle(path, handler)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
