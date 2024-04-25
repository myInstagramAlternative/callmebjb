package bond

import (
	"time"

	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
	"io"
)

type ATCommander struct {
	AT      *at.AT
	Timeout time.Duration
}

func NewATCommander(at *at.AT, timeout time.Duration) *ATCommander {
	return &ATCommander{
		AT:      at,
		Timeout: timeout,
	}
}

func InitModem(port string, baudRate int, verbose bool) (*gsm.GSM, error) {

	a, err := InitAt(port, baudRate, verbose)
	if err != nil {
		return nil, err
	}

	// gopts := []gsm.Option{} // if you want PDU mode for some reason you need to remove gsm.WithTextMode, not gonna bother with that
	gopts := []gsm.Option{gsm.WithTextMode}
	g := gsm.New(a, gopts...)
	if err := g.Init(); err != nil {
		return nil, err
	}
	return g, nil
}

func InitAt(port string, baudRate int, verbose bool) (*at.AT, error) {

	m, err := serial.New(serial.WithPort(port), serial.WithBaud(baudRate))
	if err != nil {
		return nil, err
	}
	var mio io.ReadWriter = m
	if verbose == true {
		mio = trace.New(m)
	}
	a := at.New(mio, at.WithTimeout(5*time.Second))

	return a, nil
}
