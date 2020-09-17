package delivery

import (
	"log"
	"net/http"
	"mllave.com/mllave/mmb/mb/src/messengerservice/mail"
)

// Delivery, delivers notifications to subscribers
type Delivery interface {
	Delivers()
}

// HttpDelivery, implementation of HTTP communication
type HttpDelivery struct {
	Mail mail.Mail
}

// HttpDelivery methods
func (self *HttpDelivery) Delivers() {
	_, err := http.Post(self.Mail.Url, self.Mail.ContentType, self.Mail.Msg)
	if err != nil {
		log.Println("ERROR - send request to: ", self.Mail.Url, " FAILED")
		log.Println(err)
		return
	}
	log.Println("INFO - message SUCCESFULLY sent to: ", self.Mail.Url)
}