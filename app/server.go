package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var memory = map[string]string{}

func returnPong(conn net.Conn, args []string) {
	if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
		fmt.Println("Error writing data: ", err.Error())
		os.Exit(1)
	}
}

func returnData(conn net.Conn, args []string) {
	var resultArray []string
	for i := 3; i < len(args)-1; i++ {
		responseItem := args[i] + "\n"
		resultArray = append(resultArray, responseItem)
	}

	allItems := strings.Join(resultArray, "")
	if _, err := conn.Write([]byte(allItems)); err != nil {
		fmt.Println("Error writing data: ", err.Error())
		os.Exit(1)
	}
}

func setData(conn net.Conn, args []string) {
	key := args[4]
	value := args[6]
	memory[key] = strings.TrimRight(value, "\r")
	if _, err := conn.Write([]byte("+OK\r\n")); err != nil {
		fmt.Println("Error writing data: ", err.Error())
		os.Exit(1)
	}
}

func retrieveData(conn net.Conn, args []string) {
	if value, ok := memory[args[4]]; ok {
		RESPstring := "$" + fmt.Sprint(len(value)) + "\r\n" + value + "\r\n"
		fmt.Println(RESPstring)
		if _, err := conn.Write([]byte(RESPstring)); err != nil {
			fmt.Println("Error writing data: ", err.Error())
			os.Exit(1)
		}
	} else {
		if _, err := conn.Write([]byte("+(nil)\r\n")); err != nil {
			fmt.Println("Error writing data: ", err.Error())
			os.Exit(1)
		}
	}
}

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
				fmt.Println("Info: io.EOF")
				break
			}
			panic(err)
		}

		// TODO
		// does not consider the difference of RESP data types and there is no parse section for it
		data := strings.Split(string(buffer), "\n")
		command := strings.ToUpper(data[2])

		switch {
		case strings.Contains(command, "PING"):
			returnPong(conn, data)
		case strings.Contains(command, "ECHO"):
			returnData(conn, data)
		case strings.Contains(command, "SET"):
			setData(conn, data)
		case strings.Contains(command, "GET"):
			retrieveData(conn, data)
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
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
