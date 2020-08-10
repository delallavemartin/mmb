package main

import (
	"fmt"
	"io"
	"net/http"
)

func postMsg(url string, contentType string, reader io.Reader) int {
	resp, err := http.Post(url, contentType, reader)
	if err != nil {
		fmt.Println("ERROR - send request to port: ", url, " FAILED")
		return 500
	}
	fmt.Println("INFO - message SUCCESFULLY sent to: ", url)
	defer resp.Body.Close()
	return 200
}

func publisherHandler(w http.ResponseWriter, r *http.Request) {
	statusCode := postMsg("http://localhost:9999/notify", "text/plain", r.Body)
	w.WriteHeader(statusCode)
}

func main() {
	fmt.Println("INFO - Server started.")

	// When server receives a notification, the msg will be published to his subscribers/consumers
	http.HandleFunc("/notify", publisherHandler)

	// Each request its mapped to one lightweight thread trough go routines.
	http.ListenAndServe(":8080", nil)

}
