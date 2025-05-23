package main

import (
	"log"
	"net/http"
)

const (
	webPort = "80"
)

type Config struct {
}

func main() {
	app := Config{}
	log.Printf("Starting mail service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
