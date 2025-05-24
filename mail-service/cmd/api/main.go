package main

import (
	"log"
	"net/http"
)

const (
	webPort = "80"
)

type Config struct {
	Mailler Mail
}

func main() {
	app := Config{
		Mailler: createMail(),
	}
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
