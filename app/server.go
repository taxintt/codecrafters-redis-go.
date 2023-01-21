package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func makeReponse(conn net.Conn) {
	buffer := make([]byte, 1024)
	defer conn.Close()

	n, err := conn.Read(buffer)
	if err != nil {
		log.Println(n, "read error", err)
	}
	fmt.Println(string(buffer))

	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		return
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go makeReponse(conn)
	}
}
