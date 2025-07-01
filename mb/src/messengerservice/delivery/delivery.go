package delivery

// Package delivery provides interfaces and implementations for delivering messages.
// It includes HttpDelivery for sending messages over HTTP.

import (
	"log"
	"net/http"
	"mllave.com/mllave/mmb/mb/src/messengerservice/mail"
)

// Delivery defines the interface for message delivery.
type Delivery interface {
	Delivers() error
}

// HttpDelivery is an implementation of the Delivery interface for HTTP communication.
// It contains the mail to be delivered.
type HttpDelivery struct {
	Mail mail.Mail
}

// Delivers sends the mail via an HTTP POST request.
// It returns an error if the request fails.
func (self *HttpDelivery) Delivers() error {
	resp, err := http.Post(self.Mail.Url, self.Mail.ContentType, self.Mail.Msg)
	if err != nil {
		log.Println("ERROR - send request to: ", self.Mail.Url, " FAILED")
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	log.Println("INFO - message SUCCESFULLY sent to: ", self.Mail.Url)
	return nil
}