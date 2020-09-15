package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func printMsgHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error parsing message."))
		return
	}
	defer r.Body.Close()
	fmt.Println(string(body))
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

func main() {
	fmt.Println("INFO - Server started.")

	port := os.Args[1]

	postMsg("http://localhost:8080/subscribe", "text/plain", strings.NewReader(port))

	http.HandleFunc("/notify", printMsgHandler)

	// listen to port
	http.ListenAndServe(":"+port, nil)

}
