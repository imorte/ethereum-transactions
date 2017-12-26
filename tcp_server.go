package main

import (
	"crypto/tls"
	"log"
	"net"
	"crypto/x509"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"bytes"
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
				validatedResult, status := validateData(data[1], data[2], data[3])
				if !status {
					conn.Write([]byte("Please, check your data: " + strings.Join(validatedResult, "; ") ))
					conn.Close()
				}

				password := data[len(data)-1]

				hexAmount, err := strconv.Atoi(data[3])
				checkErr(err)

				result, status := SendEth(data[1], data[2], fmt.Sprintf("0x%X", hexAmount), password)

				if status {
					res := []byte(fmt.Sprintf("Success: %s", action))
					conn.Write(res)
					message, isStored := Store(data[1], data[2], result, hexAmount)

					go CatchDeliveryTime(result)

					fmt.Println(message)
					if !isStored {
						// TODO
						conn.Write([]byte(""))
						conn.Close()
					}
				} else {
					conn.Write([]byte("Transaction error: " + result))
					conn.Close()
				}
			case "get-last":
				var lastTransactions []LastTransactions
				var stringBuffer bytes.Buffer
				c := make(chan bool)

				go func() {
					lastTransactions = GetLast(c)
				}()
				<-c

				for i, v := range lastTransactions {
					marshaledTransaction, err := json.Marshal(v)
					checkErr(err)

					for _, j := range marshaledTransaction {
						stringBuffer.WriteByte(j)
					}

					if i != len(lastTransactions)-1 {
						// Comma
						stringBuffer.WriteByte(44)
					}
				}

				if len(stringBuffer.String()) != 0 {
					fmt.Println(fmt.Sprintf("[%s]", stringBuffer.String()))
				} else {
					conn.Write([]byte("I have no recent transactions"))
				}
				conn.Close()
			default:
				conn.Write([]byte("Undefined action: " + action))
			}
		} else {
			conn.Write([]byte("Usage: send <from> <to> <amount> <password> or get-last"))
			conn.Close()
			continue
		}
	}
}
