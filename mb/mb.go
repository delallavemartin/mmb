package main

import (
	"log"
	"net/http"
	"os"
	"mllave.com/mllave/mmb/mb/src/messengerservice/mail"
	"mllave.com/mllave/mmb/mb/src/messengerservice/postoffice"
	"mllave.com/mllave/mmb/mb/src/messengerservice/delivery"
	"mllave.com/mllave/mmb/mb/src/reader"
	"mllave.com/mllave/mmb/mb/src/model/subscribers"
)

//  Global Variables - SPEC: Only 10 subscribers are supported
var aListOfSubscribers = subscribers.SubscribersList{Addresses: make([]string, 0, 10)}

func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	msg := reader.CustomReader{Reader: r.Body}

	// Post Office, used to notify to each subscriber
	aSubscribersPostOffice := postoffice.PostOffice{Channel: make(chan mail.Mail)}

	go aSubscribersPostOffice.OnMessageReceived(func(aSubscriberMail mail.Mail) {
		// OBJECT INITIALIZATION & MESSAGE
		anHTTPDelivery := delivery.HttpDelivery{Mail: aSubscriberMail}
		anHTTPDelivery.Delivers()
	})

	//CLOSURE
	aListOfSubscribers.NotifySubscribers(aSubscribersPostOffice.NotificationAssistant(msg.ToString()))
}

func subscriberHandler(w http.ResponseWriter, r *http.Request) {
	portNumber := reader.CustomReader{Reader: r.Body}

	aListOfSubscribers.Add(portNumber.ToString())
	log.Println("INFO - port added: " + portNumber.ToString())
}

func main() {
	log.Println("INFO - Server started.")

	// Set log file
	// TODO refactor logging.
	f, err := os.OpenFile("log/imb.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// When server receives a notification, the msg will be published to his subscribers
	// TODO: handle HTTP error codes.
	http.HandleFunc("/notify", publisherHandler)

	// When server receives a subscription, port will be added to subscribers list
	http.HandleFunc("/subscribe", subscriberHandler)

	// Each request its mapped to one lightweight thread trough go routines.
	http.ListenAndServe(":8080", nil)

}
