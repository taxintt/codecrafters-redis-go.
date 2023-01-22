package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// multiple request
	for {
		buffer := make([]byte, 1500)
		if _, err := conn.Read(buffer); err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				fmt.Println("Error: timeout error")
				break
			} else if err == io.EOF {
				fmt.Println("Error: io.EOF error")
				break
			}
			panic(err)
		}

		if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
			fmt.Println("Error writing data: ", err.Error())
			os.Exit(1)
		}
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
		// multiple connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
