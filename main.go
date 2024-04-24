package main

import (
	"callmebjb/api"
	"log"
	"net/http"
)

func main() {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer()

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8081",
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
