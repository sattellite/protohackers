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

type request struct {
	Method *string  `json:"method,omitempty"`
	Number *float64 `json:"number,omitempty"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	defer fmt.Println(addr, "connection closed")
	fmt.Println(addr, "connected")

	data := bufio.NewScanner(conn)
	for data.Scan() {
		fmt.Println(addr, "received", len(data.Bytes()), "bytes")
		req, err := parseRequest(data.Bytes())
		if err != nil {
			fmt.Println(addr, err)
			err := conn.Close()
			if err != nil {
				fmt.Println(addr, "error closing connection", err)
			}
			break
		}

		resp := generateResponse(&req)
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Println(addr, "error writing response", err)
			err := conn.Close()
			if err != nil {
				fmt.Println(addr, "error closing connection", err)
			}
			break
		}
	}
}

var t = []byte("{\"method\":\"isPrime\",\"prime\":true}\n")
var f = []byte("{\"method\":\"isPrime\",\"prime\":false}\n")

func parseRequest(data []byte) (request, error) {
	var r request
	if err := json.Unmarshal(data, &r); err != nil {
		return r, err
	}

	if r.Method == nil {
		return r, fmt.Errorf("malformed request: method is missing")
	}
	if *r.Method != "isPrime" {
		return r, fmt.Errorf("malformed request: method invalid")
	}
	if r.Number == nil {
		return r, fmt.Errorf("malformed request: number is missing")
	}

	return r, nil
}

func generateResponse(req *request) []byte {
	// if number is float, return false
	if *req.Number != float64(int(*req.Number)) {
		return f
	}

	// check if number is prime
	if big.NewInt(int64(*req.Number)).ProbablyPrime(0) {
		return t
	}

	return f
}
