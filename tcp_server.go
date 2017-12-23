package main

import (
	"crypto/tls"
	"log"
	"net"
	"crypto/x509"
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
	log.Print("Listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("Accepted bytes from %s", conn.RemoteAddr())
		tlsConnection, ok := conn.(*tls.Conn)
		if ok {
			log.Print("Connection established")
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

		requestString := strings.Fields(string(buf[:n]))

		action, _, _, _ := requestString[0], requestString[1], requestString[2], requestString[3]

		switch action {
		case "send": {
			fmt.Println("Отправляю запрос на перевод в ноду")

			// TODO: Нормальный ответ
			res := []byte(fmt.Sprintf("Got your request: %s", action))

			GetBalance()
			conn.Write(res)
		}
		}
	}
}

