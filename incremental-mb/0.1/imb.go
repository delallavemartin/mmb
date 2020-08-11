package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

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

func publisherHandler(w http.ResponseWriter, r *http.Request) {
	// Reader converted to string to create one Reader per POST.
	// Since is requested by postMsg function firm
	body := readerToString(r.Body)
	postMsg("http://localhost:9995/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9996/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9997/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9998/notify", "text/plain", strings.NewReader(body))
	postMsg("http://localhost:9999/notify", "text/plain", strings.NewReader(body))
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

	// When server receives a notification, the msg will be published to his subscribers/consumers
	// TODO: handle HTTP error codes.
	http.HandleFunc("/notify", publisherHandler)

	// Each request its mapped to one lightweight thread trough go routines.
	http.ListenAndServe(":8080", nil)

}
