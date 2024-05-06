package main

import (
	"callmebjb/api"
	"callmebjb/utils"
	"fmt"

	"log"
	"net/http"
	// "github.com/warthog618/sms/encoding/gsm7"
)

func init() {
	utils.InitConfig()
}

func main() {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer()

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(server, r)

	addr := fmt.Sprintf("%s:%s", utils.Config.Server.Listen, utils.Config.Server.Port)
	log.Printf("Listening on %s", addr)
	s := &http.Server{
		Handler: h,
		Addr:    addr,
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())

}
