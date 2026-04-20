package main

import (
	"log"
	"net/http"
	"os"

	"tbankshariki/internal/demoapi"
)

func main() {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := demoapi.NewServer(demoapi.NewStore())

	log.Printf("tbankshariki api listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
