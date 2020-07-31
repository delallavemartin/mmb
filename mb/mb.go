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
	select {
	case msg := <-c:
		fmt.Println("Message: ", msg, " Deque")
		go sendToConsumers(msg)
	}
}

func main() {
	serverConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 48772, Zone: ""})
	defer serverConn.Close()

	fmt.Println("Server started.")

	c := make(chan string)
	defer close(c)

	buffer := make([]byte, 1024)
	
	for {
		n, addr, _ := serverConn.ReadFromUDP(buffer)
		fmt.Println("Received ", string(buffer[0:n]), " from ", addr)
		go publish(string(buffer[0:n]), c)
		dequeAndSendToConsumers(c)
	}

}
