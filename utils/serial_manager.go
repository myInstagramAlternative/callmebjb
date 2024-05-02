package utils

import (
	"callmebjb/bond"
	"fmt"
	"sync"

	"github.com/warthog618/modem/at"
)

type Command struct {
	Action func() error
	Done   chan error
}

type SerialManager struct {
	Port       string
	BaudRate   int
	Connection *at.AT
	CommandCh  chan Command
	Mu         sync.Mutex
}

func NewSerialManager(port string, baudRate int) (*SerialManager, error) {
	manager := &SerialManager{
		Port:      port,
		BaudRate:  baudRate,
		CommandCh: make(chan Command, 100),
	}
	fmt.Printf("%v, %v", manager.Port, manager.BaudRate)
	conn, err := bond.InitAt(manager.Port, manager.BaudRate, false)
	if err != nil {
		return nil, err
	}

	if err = conn.Init(); err != nil {
		return nil, err
	}

	manager.Connection = conn

	go manager.processCommands()

	return manager, nil
}

func (sm *SerialManager) processCommands() {
	for cmd := range sm.CommandCh {
		err := cmd.Action()
		cmd.Done <- err
	}
}

func (sm *SerialManager) Execute(action func() error) error {
	done := make(chan error, 1)
	command := Command{
		Action: action,
		Done:   done,
	}

	sm.CommandCh <- command
	return <-done
}

func (sm *SerialManager) Close() {
	close(sm.CommandCh)
	if sm.Connection != nil {
		// Consider closing the serial connection if necessary
	}
}
