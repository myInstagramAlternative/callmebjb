package main

import (
	"callmebjb/api"
	// "callmebjb/bond"
	"callmebjb/utils"

	"fmt"
	"log"
	"net/http"
	"os"
	// "time"
)

var Config utils.Config

func initConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	utils.ReadConfig(&Config, configPath)
}

func main() {
	initConfig()

	// modem, err := bond.InitModem(Config.Modem.Port, Config.Modem.BaudRate)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// sender := bond.NewMessagingService(modem, 5*time.Second)

	// if err := sender.Send("", "Test alarm", false); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("SMS sent successfully!")

	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer()

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    fmt.Sprintf("%s:%s", Config.Server.Listen, Config.Server.Port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
