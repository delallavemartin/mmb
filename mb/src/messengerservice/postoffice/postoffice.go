package postoffice

import (
	"strings"
	"mllave.com/mllave/mmb/mb/src/messengerservice/mail"
)

// Post Office, knows how to notify and received messages
type PostOffice struct {
	Channel chan mail.Mail
}

// Post Office methods
func (self *PostOffice) OnMessageReceived(send func(mail mail.Mail)) {
	// this loop receives values from the channel repeatedly until it is closed
	for mail := range self.Channel {
		// go routines added to improve request per second performance.
		go send(mail)
	}
}

func (self *PostOffice) NotificationAssistant(msg string) func(address string) {
	return func(address string) {
		// Send request to the channel in order to proccess it.
		self.Channel <- mail.Mail{"http://localhost:" + address + "/notify", "text/plain", strings.NewReader(msg)}
	}
}