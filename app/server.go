package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var rexLeadingDigits = regexp.MustCompile(`\d+`)

func isSimplePing(array []string) bool {
	rex := rexLeadingDigits.Copy()
	argsLength, _ := strconv.Atoi(rex.FindString(array[0]))
	return argsLength == 1
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// multiple request
	for {
		buffer := make([]byte, 1500)
		_, err := conn.Read(buffer)

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
		fmt.Println("bulk string case")
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

		// process before checking EOF
		if err != nil {
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
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
