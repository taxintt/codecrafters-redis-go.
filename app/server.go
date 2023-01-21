package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var rexLeadingDigits = regexp.MustCompile(`\d+`)
var isSinglePing = func(array []string) bool {
	rex := rexLeadingDigits.Copy()
	argsLength, _ := strconv.Atoi(rex.FindString(array[0]))
	return argsLength == 1
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	var buffer []byte
	if n, err := conn.Read(buffer); err != nil {
		log.Println(n, "read error", err)
	}

	args := strings.Split(string(buffer), "\n")
	if isSinglePing(args) {
		if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
			fmt.Println("Error writing data: ", err.Error())
			os.Exit(1)
		}
		return
	}

	for i := 4; i < len(args); i += 2 {
		fmt.Printf("1) %#v\n", args[i])
		if _, err := conn.Write([]byte("+" + args[i] + "\n")); err != nil {
			fmt.Println("Error writing data: ", err.Error())
			os.Exit(1)
		}
		continue
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
		go handleRequest(conn)
	}
}
