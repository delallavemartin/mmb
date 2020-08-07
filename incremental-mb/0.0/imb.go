package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func postMsg(reader io.Reader) int {
	resp, err := http.Post("http://localhost:9999", "text/plain", reader)
	if err != nil {
		fmt.Println("ERROR - send request to port: 9999 FAILED")
		return 500
	}
	fmt.Println("INFO - message SUCCESFULLY sent to port: 9999")
	defer resp.Body.Close()
	return 200
}

func sendToHandler(w http.ResponseWriter, r *http.Request) {
	statusCode := postMsg(r.Body)
	io.WriteString(w, "Status: "+strconv.Itoa(statusCode))
}

func main() {
	fmt.Println("INFO - Server started.")

	http.HandleFunc("/notify", sendToHandler)

	// listen to port
	http.ListenAndServe(":8080", nil)

}
