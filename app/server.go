package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type dataStore struct {
	data          map[string]string
	dataTimestamp map[string]time.Time
	mu            sync.RWMutex
	expireTime    map[string]int
}

var store = dataStore{
	mu:            sync.RWMutex{},
	data:          map[string]string{},
	dataTimestamp: map[string]time.Time{},
	expireTime:    map[string]int{},
}

func setData(data []string) {
	store.data[data[4]] = data[6]
	store.dataTimestamp[data[4]] = time.Now()
	if strings.ToUpper(strings.TrimSpace(data[8])) == "PX" {
		millisecond, err := strconv.Atoi(data[10])
		if err != nil {
			panic(err)
		}
		store.expireTime[data[4]] = millisecond
	} else {
		store.expireTime[data[4]] = 0
	}
}

func checkExpiration(data []string) bool {
	if store.expireTime[data[4]] == 0 {
		return false
	}

	retrieveDataTimeStamp := time.Now()
	setDataTimeStamp := store.dataTimestamp[data[4]]
	expireTime := store.expireTime[data[4]]
	return retrieveDataTimeStamp.Sub(setDataTimeStamp) > time.Duration(expireTime)*time.Millisecond
}

func retrieveData(data []string) string {
	var respData string
	var isExpired bool = checkExpiration(data)

	if value, ok := store.data[data[4]]; ok && !isExpired {
		respData = fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
	} else {
		respData = "$-1\r\n"
	}
	return respData
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
		respData := strings.Split(string(buffer), "\r\n")
		command := strings.ToUpper(strings.TrimSpace(respData[2]))

		var msg string
		switch command {
		case "PING":
			msg = "+PONG\r\n"
		case "ECHO":
			msg = fmt.Sprintf("$%d\r\n%s\r\n", len(respData[4]), respData[4])
		case "SET":
			setData(respData)
			msg = "+OK\r\n"
		case "GET":
			msg = retrieveData(respData)
		default:
			msg = "-ERR unknown command '" + command + "'\r\n"
		}

		if _, err := conn.Write([]byte(msg)); err != nil {
			if err == io.EOF {
				fmt.Println("Info: io.EOF")
				os.Exit(1)
			}
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
