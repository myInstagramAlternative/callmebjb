package api

import (
	"callmebjb/bond"
	"callmebjb/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() Server {
	return Server{}
}

var serialMgr *utils.SerialManager

func init() {
	utils.InitConfig()

	var err error
	log.Printf("Initializing Serial Manager on port: %s", utils.Config.Modem.Port)
	serialMgr, err = utils.NewSerialManager(utils.Config.Modem.Port, utils.Config.Modem.BaudRate)
	if err != nil {
		log.Fatalf("Failed to initialize serial manager: %v", err)
	}
}

// (GET /ping)
func (Server) GetPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	resp := Pong{
		Ping: "pong",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// (GET /call)
func (Server) Call(w http.ResponseWriter, r *http.Request, params CallParams) {
	w.Header().Add("Content-Type", "application/json")

	timeout := 15 * time.Second
	atc := bond.NewATCommander(serialMgr.Connection, timeout)

	if err := serialMgr.Execute(func() error {
		return atc.Call(params.Number)
	}); err != nil {

		resp := Error{
			Code:    500,
			Message: fmt.Sprintf("%v", err),
		}

		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	status := "Calling now"
	resp := CallResponse{
		Status: &status,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// (GET /hangup)
func (Server) Hangup(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	timeout := 15 * time.Second
	atc := bond.NewATCommander(serialMgr.Connection, timeout)

	if err := serialMgr.Execute(func() error {
		return atc.Hangup()
	}); err != nil {

		resp := Error{
			Code:    500,
			Message: fmt.Sprintf("%v", err),
		}

		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	status := "Hung up"
	resp := CallResponse{
		Status: &status,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
	return
}

func (Server) Sms(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var sms SmsJSONRequestBody
	err := decoder.Decode(&sms)
	if err != nil {
		resp := Error{
			Code:    500,
			Message: fmt.Sprintf("%v", err),
		}

		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	timeout := 15 * time.Second
	atc := bond.NewATCommander(serialMgr.Connection, timeout)

	if err := serialMgr.Execute(func() error {
		return atc.SendSMS(sms.Number, sms.Message, false)
	}); err != nil {

		resp := Error{
			Code:    500,
			Message: fmt.Sprintf("%v", err),
		}

		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	status := "Sent"
	resp := SMSResponse{
		Status: &status,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (Server) Ussd(w http.ResponseWriter, r *http.Request, params UssdParams) {
	w.Header().Add("Content-Type", "application/json")

	timeout := 15 * time.Second
	atc := bond.NewATCommander(serialMgr.Connection, timeout)

	responseChannel := make(chan string, 1)
	errChannel := make(chan error, 1)

	err := serialMgr.Execute(func() error {
		result, err := atc.SendUSSD(params.Code)
		if err != nil {
			errChannel <- err
			return err
		}
		responseChannel <- result
		return nil
	})

	if err != nil {
		resp := Error{
			Code:    500,
			Message: fmt.Sprintf("%v", err),
		}

		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	status := "Success"
	responseValue := <-responseChannel
	resp := USSDResponse{
		Message: &responseValue,
		Status:  &status,
	}
	close(responseChannel)
	close(errChannel)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
