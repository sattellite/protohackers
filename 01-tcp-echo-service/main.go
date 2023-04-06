package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

// TCP echo service
func main() {
	debug := os.Getenv("DEBUG") != ""
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
		if debug {
			go debugHandler(conn)
		} else {
			go fastHandler(conn)
		}
	}
}

func fastHandler(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	defer fmt.Println(addr, "connection closed")
	fmt.Println(addr, "connected")

	recv, err := io.Copy(conn, conn)
	if err != nil {
		fmt.Printf("%s error reading data from connection: %v\n", addr, err)
	}
	fmt.Println(addr, "received", recv, "bytes")
}

func debugHandler(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	defer fmt.Println(addr, "connection closed")
	fmt.Println(addr, "connected")

	buf := make([]byte, 1024)
	recv := 0
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				fmt.Printf("%s error reading data from connection: %v\n", addr, err)
				break
			}
		}
		data := buf[:n]
		recv += n

		// write data to the connection
		_, err = conn.Write(data)
		if err != nil {
			fmt.Printf("%s error writing data to connection: %v\n", addr, err)
			break
		}
	}
	fmt.Println(addr, "received", recv, "bytes")
}
