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

// Created to store all information needed to publish the request.
type Request struct {
	Url         string
	ContentType string
	Reader      io.Reader
}

type SafeSubscribers struct {
	addresses []string
	mux       sync.Mutex
}

// SPEC: Only 10 subscribers are supported
var subscribers_addresses = SafeSubscribers{addresses: make([]string, 0, 10)}

func (s *SafeSubscribers) add(address string) {
	s.mux.Lock()
	s.addresses = append(s.addresses, address)
	s.mux.Unlock()
}

func (s *SafeSubscribers) getByIndex(index int) string {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.addresses[index]
}

func (s *SafeSubscribers) numberOfSubscribers() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.addresses)
}

func readerToString(reader io.Reader) string {
	stream, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("ERROR - reading reader FAILED")
		return ""
	}
	return string(stream)
}

func postMsg(url string, contentType string, reader io.Reader) {
	_, err := http.Post(url, contentType, reader)
	if err != nil {
		log.Println("ERROR - send request to port: ", url, " FAILED")
		log.Println(err)
		return
	}
	log.Println("INFO - message SUCCESFULLY sent to: ", url)
}

func publish(ch chan Request) {
	// this loop receives values from the channel repeatedly until it is closed
	for request := range ch {
		// go routines added to improve request per second performance.
		go postMsg(request.Url, request.ContentType, request.Reader)
	}

}

func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)

	ch := make(chan Request)

	go publish(ch)

	//Iterate over subscribers list.
	for i := 0; i < subscribers_addresses.numberOfSubscribers(); i++ {
		// Send request to the channel in order to proccess it.
		ch <- Request{"http://localhost:" + subscribers_addresses.getByIndex(i) + "/notify", "text/plain", strings.NewReader(body)}
	}
}

func subscriberHandler(w http.ResponseWriter, r *http.Request) {
	// Read port number
	port_number := readerToString(r.Body)

	// Append new subscribers to subscribers list
	subscribers_addresses.add(port_number)
	log.Println("INFO - port added: " + port_number)
}

func main() {
	log.Println("INFO - Server started.")

	//Set log file
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
