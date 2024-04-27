package bond

import (
	"log"

	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/sms"
)

func (atc *ATCommander) SendSMS(number string, message string, pdu bool) error {

	gopts := []gsm.Option{gsm.WithTextMode}
	g := gsm.New(atc.AT, gopts...)
	if err := g.Init(); err != nil {
		return err
	}

	var err error
	if pdu {
		err = atc.sendPDU(number, message)
	} else {
		_, err = g.SendShortMessage(number, message)
	}
	return err
}

func (atc *ATCommander) sendPDU(number string, message string) error {

	gopts := []gsm.Option{gsm.WithTextMode}
	g := gsm.New(atc.AT, gopts...)
	if err := g.Init(); err != nil {
		return err
	}

	pdus, err := sms.Encode([]byte(message), sms.To(number), sms.WithAllCharsets)
	if err != nil {
		return err
	}
	for i, p := range pdus {
		tp, err := p.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = g.SendPDU(tp)
		if err != nil {
			return err
		}
		log.Printf("PDU %d sent successfully\n", i+1)
	}
	return nil
}
