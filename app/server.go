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

func isSimplePing(array []string) bool {
	rex := rexLeadingDigits.Copy()
	argsLength, _ := strconv.Atoi(rex.FindString(array[0]))
	fmt.Println(argsLength)
	return argsLength == 1
}

func handleRequest(conn net.Conn) {
	buffer := make([]byte, 1500)
	if n, err := conn.Read(buffer); err != nil {
		log.Println(n, "read error", err)
	}

	args := strings.Split(string(buffer), "\n")

	// simple string case
	if isSimplePing(args) {
		if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
			fmt.Println("Error writing data: ", err.Error())
			os.Exit(1)
		}
		return
	}

	// bulk string case
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
		defer conn.Close()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
