package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// MODEL DEFINITION

// Store all information needed to publish the request.
type Request struct {
	Url         string
	ContentType string
	Reader      io.Reader
}

// SubscribersList with Sync.
type SubscribersList struct {
	Addresses []string
	Mux       sync.Mutex
}

// SubscribersList methods
func (self *SubscribersList) Add(address string) {
	self.Mux.Lock()
	self.Addresses = append(self.Addresses, address)
	self.Mux.Unlock()
}

func (self *SubscribersList) NotifySubscribers(notifier func(address string)) {
	self.Mux.Lock()
	defer self.Mux.Unlock()
	//Iterate over subscribers list.
	for i := 0; i < len(self.Addresses); i++ {
		notifier(self.Addresses[i])
	}
}

// Delivery, delivers notifications to subscribers
type Delivery interface {
	Delivers()
}

// HttpDelivery, implementation of HTTP communication
type HttpDelivery struct {
	Request Request
}

// HttpDelivery methods
func (self *HttpDelivery) Delivers() {
	_, err := http.Post(self.Request.Url, self.Request.ContentType, self.Request.Reader)
	if err != nil {
		log.Println("ERROR - send request to: ", self.Request.Url, " FAILED")
		log.Println(err)
		return
	}
	log.Println("INFO - message SUCCESFULLY sent to: ", self.Request.Url)
}

// Notification Center, knows how to notify and received messages
type NotificationCenter struct {
	Channel chan Request
}

// Notification Center methods
func (self *NotificationCenter) OnMessageReceived(router func(request Request)) {
	// this loop receives values from the channel repeatedly until it is closed
	for request := range self.Channel {
		router(request)
	}
}

func (self *NotificationCenter) Notifier(msg string) func(address string) {
	return func(address string) {
		// Send request to the channel in order to proccess it.
		self.Channel <- Request{"http://localhost:" + address + "/notify", "text/plain", strings.NewReader(msg)}
	}
}

// CustomReader, used to convert a Reader to string
type CustomReader struct {
	Reader         io.Reader
	ReaderAsString string
}

// CustomReader Methods
func (self *CustomReader) ToString() string {
	if self.ReaderAsString == "" {
		stream, err := ioutil.ReadAll(self.Reader)
		if err != nil {
			log.Println("ERROR - reading reader FAILED")
			return ""
		}
		self.ReaderAsString = string(stream)
	}
	return self.ReaderAsString
}

// Global Variables
// SPEC: Only 10 subscribers are supported
var aListOfSubscribers = SubscribersList{Addresses: make([]string, 0, 10)}

// Notification Center, used to notify to each subscriber
var SubscribersNotificationCenter = NotificationCenter{Channel: make(chan Request)}

// Request Handlers
func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	msg := CustomReader{Reader: r.Body}

	go SubscribersNotificationCenter.OnMessageReceived(func(request Request) {
		// OBJECT INITIALIZATION & MESSAGE
		anHTTPDelivery := HttpDelivery{Request: request}
		// go routines added to improve request per second performance.
		go anHTTPDelivery.Delivers()
	})

	//CLOSURE
	aListOfSubscribers.NotifySubscribers(SubscribersNotificationCenter.Notifier(msg.ToString()))
}

func subscriberHandler(w http.ResponseWriter, r *http.Request) {
	// Read port number
	portNumber := CustomReader{Reader: r.Body}

	// Append new subscribers to subscribers list
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
