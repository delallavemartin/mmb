package main

import (
	"fmt"
	"net"
)

func publish(msg string, c chan string) {
	c <- msg
	fmt.Println("Msg: ",msg," Published to ", c, "channel")
}

func sendToConsumers(msg string) {
	conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 9999, Zone: ""})
	defer conn.Close()
	fmt.Fprintf(conn,msg)
}

func dequeAndSendToConsumers(c chan string) {
	// - Select: blocks until one of its cases can run, then it executes that case. 
	// 	It chooses one at random if multiple are ready.
	// select {
	// case msg := <-c:
	msg := <-c
	fmt.Println("Message: ", msg, " Deque")
	go sendToConsumers(msg)
	// }
}

func main() {
	serverConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 48772, Zone: ""})
	defer serverConn.Close()

	fmt.Println("Server started.")


	// - Regular Channels : Sends and receives block until the other side is ready. 
	// 	This allows goroutines to synchronize without explicit locks or condition variables.
	// - Buffered Channel : Sends to a buffered channel block only when the buffer is full. Receives block when the buffer is empty.
	c := make(chan string)
	//  Only the sender should close a channel, never the receiver. Sending on a closed channel will cause a panic.
	//	Closing is only necessary when the receiver must be told there are no more values coming, such as to terminate a range loop.
	defer close(c)

	buffer := make([]byte, 1024)
	
	for {
		n, addr, _ := serverConn.ReadFromUDP(buffer)
		fmt.Println("Received ", string(buffer[0:n]), " from ", addr)
		// - Goroutine:  is a lightweight thread managed by the Go runtime.
		// 	Run in the same address space, so access to shared memory must be synchronized. 
		go publish(string(buffer[0:n]), c)
		dequeAndSendToConsumers(c)
	}

}
