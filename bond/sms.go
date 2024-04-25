package bond

import (
	"log"
	"time"

	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/sms"
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
