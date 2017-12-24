package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args
	lenArgs := len(args)

	if lenArgs != 6 {
		fmt.Print("Usage:\nsend <from> <to> <amount> <password>\nget-last\n")
		os.Exit(1)
	}

	SendData(strings.Join(args[1:], " "))
}

func SendData(query string) {
	// Load keys
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("Couldn't load cers: %s", err)
	}

	// Add them to the config
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "127.0.0.1:9000", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	fmt.Println("TLS connection established to: ", conn.RemoteAddr())

	message := query

	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	fmt.Printf("%q (%d bytes)\n", string(reply[:n]), n)
	fmt.Print("Exit\n")
}
