package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/inancgumus/screen"
)

type DiscoveryRequest struct {
	Command string
}

type DiscoveryResponse struct {
	Address string
}

const maxDatagramSize = 8192

func main() {
	// Handle termination signals gracefully
	setupSignalHandler()

	// Set up UDP client
	serverAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:10001")
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}
	clientAddr, err := net.ResolveUDPAddr("udp", ":10101")
	if err != nil {
		log.Fatalf("Error resolving UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", clientAddr)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer conn.Close()

	// Create a DiscoveryRequest struct
	discoveryRequest := DiscoveryRequest{
		Command: "REQUEST",
	}

	// Initialize encoder and decoder for structured communication
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// Send discovery message as a structured request
	if err := encoder.Encode(discoveryRequest); err != nil {
		log.Fatalf("Error encoding discovery request: %v", err)
	}
	log.Printf("Discovery packet sent from %s", clientAddr.String())

	// Wait for autodiscovery response
	var tcpAddr *net.TCPAddr
	for {
		var response DiscoveryResponse

		// Receive the structured response
		if err := decoder.Decode(&response); err != nil {
			log.Fatalf("Error decoding discovery response: %v", err)
		}
		log.Printf("Received TCP server address: %s", response.Address)

		tcpAddr, err = net.ResolveTCPAddr("tcp", response.Address)
		if err != nil {
			log.Fatalf("Error resolving TCP address: %v", err)
		}

		log.Printf("TCP address resolved: %s", tcpAddr.String())
		break
	}

	// Connect to TCP server
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Error connecting to TCP server: %v", err)
	}
	defer tcpConn.Close()
	log.Printf("Connected to TCP server %s", tcpAddr.String())

	// Create channel for user input
	inputCh := make(chan string)

	// Start goroutine to read user input
	go readInput(inputCh)

	// Start goroutine to send data to server
	go sendDataToServer(tcpConn, inputCh)

	// Start goroutine to receive data from server
	go receiveDataFromServer(tcpConn)
}

func readInput(inputCh chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Print("Enter message to send: ")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) == "/clear" {
			screen.Clear()
			screen.MoveTopLeft()
			continue
		}
		inputCh <- input
	}
}

func sendDataToServer(conn net.Conn, inputCh <-chan string) {
	for input := range inputCh { // Wait for user input from channel
		_, err := conn.Write([]byte(input)) // Send input to server
		if err != nil {
			log.Fatalf("Error sending data to server: %v", err)
		}
	}
}

func receiveDataFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadString('\n') // Receive data from server
		if err != nil {
			log.Fatalf("Error receiving data from server: %v", err)
		}
		log.Printf("Received from server: %s", data)
		log.Print("Enter message to send: ")
	}
}

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\nReceived termination signal. Cleaning up resources...")
		os.Exit(0)
	}()
}
