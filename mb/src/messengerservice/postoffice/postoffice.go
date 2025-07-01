package postoffice

// Package postoffice provides functionality for managing and sending messages.
// It defines the PostOffice type and its methods for handling message reception and notification.

import (
	"os"
	"strings"
	"mllave.com/mllave/mmb/mb/src/messengerservice/mail"
)

// PostOffice represents a message post office that handles sending and receiving mail.
// It contains a channel for incoming mail.
type PostOffice struct {
	Channel chan mail.Mail
}

// OnMessageReceived processes incoming messages from the PostOffice's channel.
// It takes a 'send' function as an argument, which is used to process each received mail.
func (self *PostOffice) OnMessageReceived(send func(mail mail.Mail)) {
	// this loop receives values from the channel repeatedly until it is closed
	for mail := range self.Channel {
		// go routines added to improve request per second performance.
		go send(mail)
	}
}

// NotificationAssistant creates a function that sends a notification message to a given address.
// The notification host can be configured via the MB_NOTIFICATION_HOST environment variable,
// defaulting to "http://localhost:" if not set.
func (self *PostOffice) NotificationAssistant(msg string) func(address string) {
	notificationHost := os.Getenv("MB_NOTIFICATION_HOST")
	if notificationHost == "" {
		notificationHost = "http://localhost:"
	}
	return func(address string) {
		// Send request to the channel in order to proccess it.
		self.Channel <- mail.Mail{Url: notificationHost + address + "/notify", ContentType: "text/plain", Msg: strings.NewReader(msg)}
	}
}