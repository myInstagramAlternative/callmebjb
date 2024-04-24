package bond

import (
	"log"
	"time"

	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
	"github.com/warthog618/sms"
	"io"
)

type MessagingService struct {
	Modem   *gsm.GSM
	Timeout time.Duration
}

func NewMessagingService(modem *gsm.GSM, timeout time.Duration) *MessagingService {
	return &MessagingService{
		Modem:   modem,
		Timeout: timeout,
	}
}

func (s *MessagingService) Send(number string, message string, pdu bool) error {
	var err error
	if pdu {
		err = s.sendPDU(number, message)
	} else {
		_, err = s.Modem.SendShortMessage(number, message)
	}
	return err
}

func (s *MessagingService) sendPDU(number string, message string) error {
	pdus, err := sms.Encode([]byte(message), sms.To(number), sms.WithAllCharsets)
	if err != nil {
		return err
	}
	for i, p := range pdus {
		tp, err := p.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = s.Modem.SendPDU(tp)
		if err != nil {
			return err
		}
		log.Printf("PDU %d sent successfully\n", i+1)
	}
	return nil
}

func InitModem(port string, baudRate int) (*gsm.GSM, error) {
	m, err := serial.New(serial.WithPort(port), serial.WithBaud(baudRate))
	if err != nil {
		return nil, err
	}
	var mio io.ReadWriter = m
	mio = trace.New(m)
	// gopts := []gsm.Option{} // if you want PDU mode for some reason you need to remove gsm.WithTextMode, not gonna bother with that
	gopts := []gsm.Option{gsm.WithTextMode}

	g := gsm.New(at.New(mio, at.WithTimeout(5*time.Second)), gopts...)
	if err := g.Init(); err != nil {
		return nil, err
	}
	return g, nil
}
