package bond

import (
	"fmt"
)

func (atc *ATCommander) Call(phoneNumber string) error {
	cmdCall := fmt.Sprintf("D%s;\r", phoneNumber)

	_, err := atc.AT.Command(cmdCall)
	if err != nil {
		return err
	}

	fmt.Printf("Calling %s...\n", phoneNumber)
	return nil
}

func (atc *ATCommander) Hangup() error {
	cmdCall := "+CHUP"

	_, err := atc.AT.Command(cmdCall)
	if err != nil {
		return err
	}

	fmt.Printf("Hung up the call.")
	return nil
}
