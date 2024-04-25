package bond

import (
	"fmt"
	"strings"
	"time"

	"github.com/warthog618/modem/info"
)

func (atc *ATCommander) SendUSSD(msg string) (string, error) {

	rspChan := make(chan string)
	handler := func(info []string) {
		rspChan <- info[0]
	}
	atc.AT.AddIndication("+CUSD:", handler)

	cmd := fmt.Sprintf("+CUSD=1,\"%s\"", msg)

	_, err := atc.AT.Command(cmd)
	if err != nil {
		return "", err
	}
	select {
	case <-time.After(atc.Timeout):
		return "", fmt.Errorf("Timeout")
	case rsp := <-rspChan:
		fields := strings.Split(info.TrimPrefix(rsp, "+CUSD"), ",")
		rspb := string(fields[1])
		return rspb, nil
	}
}
