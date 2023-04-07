package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
)

// TCP prime time service
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)

	// listen to port 5001
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Println("Listening on", address)

	for {
		// accept a connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	defer fmt.Println(addr, "connection closed")
	fmt.Println(addr, "connected")

	data := bufio.NewScanner(conn)
	for data.Scan() {
		resp, err := generateResponse(data.Bytes())
		if err != nil {
			fmt.Println(err)
			conn.Close()
			break
		}
		fmt.Println(string(resp))
		conn.Write(resp)
	}
}

type in struct {
	Method *string  `json:"method,omitempty"`
	Number *float64 `json:"number,omitempty"`
}

var t = []byte("{\"method\":\"isPrime\",\"prime\":true}\n")
var f = []byte("{\"method\":\"isPrime\",\"prime\":false}\n")

var wrong = []byte("malformed request\n")

func generateResponse(data []byte) ([]byte, error) {
	var req in
	if err := json.Unmarshal(data, &req); err != nil {
		return wrong, err
	}

	// check that data is valid: all fields are present and method is "isPrime"
	if req.Method == nil || req.Number == nil || *req.Method != "isPrime" {
		return wrong, fmt.Errorf("malformed request")
	}

	// if number is zero or negative, return false
	if *req.Number <= 0 {
		return f, nil
	}

	// if number is float, return false
	if *req.Number != float64(int(*req.Number)) {
		return f, nil
	}

	// if number is 1, return false
	if *req.Number == 1 {
		return f, nil
	}
	// check if number is prime
	if !big.NewInt(int64(*req.Number)).ProbablyPrime(0) {
		return f, nil
	}

	fmt.Println("true")
	return t, nil
}
