package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func SendEth(from string, to string, amount string) {
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
	log.Println("TLS connection established to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}

	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	message := fmt.Sprintf("%s %s %s %s", "send", from, to, amount)

	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("client: exiting")
}

func main() {
	args := os.Args
	lenArgs := len(args)

	if lenArgs != 5 {
		fmt.Print("Usage:\nsend <from> <to> <amount>\nget-last\n")
		os.Exit(1)
	}

	switch args[1] {
	case "send":
		var hexRegex = regexp.MustCompile("0[xX][0-9a-fA-F]+")
		var isWrongData bool
		from, to, amount := args[2], args[3], args[4]

		fmt.Println("Lets send some money!")
		if hexRegex.Find([]byte(from)) == nil {
			fmt.Println("Sender has wrong wallet")
			isWrongData = true
		}

		if hexRegex.Find([]byte(to)) == nil {
			fmt.Println("Recipient has wrong wallet")
			isWrongData = true
		}

		if amountConverted, err := strconv.ParseFloat(amount, 64); err == nil {
			if amountConverted <= 0 {
				fmt.Println("Amount is zero or lower than")
				isWrongData = true
			}
		} else {
			fmt.Println("Wrong amount")
			isWrongData = true
		}

		if isWrongData {
			os.Exit(1)
		}

		SendEth(from, to, amount)


	}
}
