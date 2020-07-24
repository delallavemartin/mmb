package main

import (
	"fmt"
	"net"
)

func publish(msg string, c chan string) {
	c <- msg
	defer close(c)
}

func sendToConsumers(msg string) {
	Conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 9999, Zone: ""})
	defer Conn.Close()
	Conn.Write([]byte("hello"))
}

func dequeAndSendToConsumers(c chan string) {
	select {
	case msg := <-c:
		fmt.Println("Message: ", msg, " Enque")
		go sendToConsumers(msg)
	}
}

func main() {
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 48772, Zone: ""})
	defer ServerConn.Close()

	c := make(chan string)
	dequeAndSendToConsumers(c)

	buf := make([]byte, 1024)
	for {
		n, addr, _ := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		go publish(string(buf[0:n]), c)
	}

}
