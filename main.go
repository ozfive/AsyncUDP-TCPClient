package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/inancgumus/screen"
)

const maxDatagramSize = 8192

func main() {
	// Set up UDP client
	ServerAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:10001")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	ClientAddr, err := net.ResolveUDPAddr("udp", ":10101")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	Conn, err := net.ListenUDP("udp", ClientAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer Conn.Close()

	DiscoveryMessage := "REQUEST"
	buf := make([]byte, maxDatagramSize)
	writeBuf := []byte(DiscoveryMessage)

	// Send discovery message
	_, err = Conn.WriteToUDP(writeBuf, ServerAddr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Discovery packet sent from", ClientAddr.String())

	// Wait for autodiscovery response
	var tcpAddr *net.TCPAddr
	for {
		n, _, err := Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		message := make([]byte, n)
		copy(message, buf[:n])
		fmt.Println("Received", string(message))

		if strings.HasPrefix(string(message), "THE TCP SERVER IS LOCATED AT ") {
			tcpAddrString := strings.TrimSpace(strings.TrimPrefix(string(message), "THE TCP SERVER IS LOCATED AT "))
			tcpAddr, err = net.ResolveTCPAddr("tcp", tcpAddrString)
			if err != nil {
				fmt.Println("Error resolving TCP address:", err.Error())
				return
			}

			fmt.Println("TCP address resolved:", tcpAddr.String())
			break
		}
	}

	// Connect to TCP server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error connecting to TCP server:", err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Connected to TCP server", tcpAddr.String())

	// Create channel for user input
	inputCh := make(chan string)

	// Start goroutine to read user input
	go readInput(inputCh)

	// Start goroutine to send data to server
	go sendDataToServer(conn, inputCh)

	// Start goroutine to receive data from server
	go receiveDataFromServer(conn)

	// Wait for goroutines to finish before exiting
	select {}
}

func readInput(inputCh chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to send: ")
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
			fmt.Println("Error sending data to server:", err.Error())
			return
		}
	}
}

func receiveDataFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadString('\n') // Receive data from server
		if err != nil {
			fmt.Println("Error receiving data from server:", err.Error())
			return
		}
		fmt.Print("\033[2K") // Clear current line
		fmt.Print("\033[G")  // Move cursor to beginning of line
		fmt.Print("Received from server:", data)
		fmt.Print("Enter message to send: ")
	}
}
