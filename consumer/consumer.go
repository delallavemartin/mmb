package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func printMsgHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error parsing message."))
	}
	defer r.Body.Close()
	fmt.Println(string(body))
	w.WriteHeader(200)
}

func main() {
	fmt.Println("INFO - Server started.")

	port := os.Args[1]

	http.HandleFunc("/notify", printMsgHandler)

	// listen to port
	http.ListenAndServe(":"+port, nil)

}
