package main

import (
	"crypto/tls"
	"log"
	"net"
	"crypto/x509"
	"fmt"
	"strings"
	"strconv"
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
		tlsConnection, ok := conn.(*tls.Conn)
		if ok {
			fmt.Println("Connection established")
			state := tlsConnection.ConnectionState()

			// TODO: Del this
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		fmt.Println("Waiting")
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		data := strings.Fields(string(buf[:n]))

		validatedResult, status := validateData(data[1], data[2], data[3])

		if !status {
			conn.Write([]byte("Please, check your data: " + strings.Join(validatedResult, "; ") ))
			conn.Close()
		}

		action := data[0]
		password := data[len(data)-1]

		switch action {
		case "send":
			hexAmount, err := strconv.Atoi(data[3])
			checkErr(err)

			result, status := SendEth(data[1], data[2], fmt.Sprintf("0x%X", hexAmount), password)

			if status {
				res := []byte(fmt.Sprintf("Success: %s", action))
				conn.Write(res)

				message, isStored := Store(data[1], data[2], result, hexAmount)

				fmt.Println(message)
				if !isStored {
					conn.Write([]byte(""))
					conn.Close()
				}
			} else {
				conn.Write([]byte("Transaction error: " + result))
				conn.Close()
			}
		default:
			conn.Write([]byte("Undefined action: " + action))
		}
	}
}
