package main

import (
	// "callmebjb/api"
	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
	"github.com/warthog618/sms"
	"io"
	"log"
	// "net/http"
	"time"
)

// func main() {
// 	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
// 	server := api.NewServer()

// 	r := http.NewServeMux()

// 	// get an `http.Handler` that we can use
// 	h := api.HandlerFromMux(server, r)

// 	s := &http.Server{
// 		Handler: h,
// 		Addr:    "0.0.0.0:8081",
// 	}

// 	// And we serve HTTP until the world ends.
// 	log.Fatal(s.ListenAndServe())
// }

func main() {
	// testing modem
	m, err := serial.New(serial.WithPort("/dev/tty.usbserial-0001"), serial.WithBaud(115200))
	if err != nil {
		log.Fatal(err)
	}
	var mio io.ReadWriter = m
	mio = trace.New(m)
	gopts := []gsm.Option{}
	pduMode := false
	num := "0616693192"
	msg := "Test alarm"
	if !pduMode {
		gopts = append(gopts, gsm.WithTextMode)
	}
	g := gsm.New(at.New(mio, at.WithTimeout(5*time.Second)), gopts...)
	if err = g.Init(); err != nil {
		log.Fatal(err)
	}
	if pduMode {
		sendPDU(g, num, msg)
		return
	}
	mr, err := g.SendShortMessage(num, msg)
	// !!! check CPIN?? on failure to determine root cause??  If ERROR 302
	log.Printf("%v %v\n", mr, err)
}

func sendPDU(g *gsm.GSM, number string, msg string) {
	pdus, err := sms.Encode([]byte(msg), sms.To(number), sms.WithAllCharsets)
	if err != nil {
		log.Fatal(err)
	}
	for i, p := range pdus {
		tp, err := p.MarshalBinary()
		if err != nil {
			log.Fatal(err)
		}
		mr, err := g.SendPDU(tp)
		if err != nil {
			// !!! check CPIN?? on failure to determine root cause??  If ERROR 302
			log.Fatal(err)
		}
		log.Printf("PDU %d: %v\n", i+1, mr)
	}
}
