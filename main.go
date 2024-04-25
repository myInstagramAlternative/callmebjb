package main

import (
	"callmebjb/api"
	// "callmebjb/bond"
	"callmebjb/utils"
	"fmt"
	// "time"

	"log"
	"net/http"
	"os"
	// "github.com/warthog618/sms/encoding/gsm7"
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

	// // Make call
	// timeout := 15 * time.Second
	// a, err := bond.InitAt(Config.Modem.Port, Config.Modem.BaudRate, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = a.Init(); err != nil {
	// 	log.Fatal(err)
	// }
	// atc := bond.NewATCommander(a, timeout)
	// if err = atc.Call(""); err != nil {
	// 	log.Fatal(err)
	// }
	// time.Sleep(5 * time.Second)

	// // Hangup
	// if err = atc.Hangup(); err != nil {
	// 	log.Fatal(err)
	// }

	// return

	// // SMS EXAMPLE
	// modem, err := bond.InitModem(Config.Modem.Port, Config.Modem.BaudRate, false)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// sender := bond.NewMessagingService(modem, 5*time.Second)

	// if err := sender.Send("", "Test alarm", false); err != nil {
	// 	log.Fatal(err)
	// }

	// // USSD EXAMPLE
	// msg := "*131*1#"
	// timeout := 15 * time.Second
	// a, err := bond.InitAt(Config.Modem.Port, Config.Modem.BaudRate, false)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = a.Init(); err != nil {
	// 	log.Fatal(err)
	// }
	// atc := bond.NewATCommander(a, timeout)
	// res, err := atc.SendUSSD(msg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", res)

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
