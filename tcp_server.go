package main

import (
	"crypto/tls"
	"log"
	"net"
	"fmt"
	"strings"
)

func ListenTcp() {
	// Load keys
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("Couldn't load cers: %s", err)
	}

	// Add them to the config
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	service := "0.0.0.0:9000"

	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("Can't listen: %s", err)
	}
	fmt.Println("Listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %s", err)
			break
		}
		defer conn.Close()
		fmt.Printf("Accepted bytes from %s\n", conn.RemoteAddr())
		_, ok := conn.(*tls.Conn)
		if ok {
			fmt.Println("Connection is secured")

		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		data := strings.Fields(string(buf[:n]))
		lenArgs := len(data)
		if lenArgs == 5 || lenArgs == 1 {
			action := data[0]

			switch action {
			case "send":
				SendHandler(action, data, conn)
			case "get-last":
				GetLastHandler(err, conn)
			default:
				conn.Write([]byte("Undefined action: " + action))
			}
		} else {
			conn.Write([]byte("Usage: send <from> <to> <amount> <password> or get-last"))
			continue
		}
	}
}
