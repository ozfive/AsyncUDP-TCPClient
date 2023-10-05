# AsyncUDP-TCPClient

The AsyncUDP-TCPClient is a versatile client application that serves the following purposes:

1. **Discovery:** It initiates the discovery process by sending a UDP discovery packet to the server.

2. **TCP Connection:** Upon receiving the server's response, it upgrades to a more reliable TCP connection with the AsyncTCPServer.

## Usage

1. **Start the Client:**

Run the AsyncUDP-TCPClient to initiate the discovery process and establish a TCP connection with the server.

```bash
./async-udp-tcp-client
```

Discovery Process:
The client sends a UDP discovery packet to the server, requesting the TCP address.
Upon receiving the response, it extracts the TCP address and prepares to establish a TCP connection.

TCP Connection:
The client connects to the AsyncTCPServer using the obtained TCP address.
Once connected, it enables reliable communication with the server over TCP.

## Features

Structured Protocol: The client uses structured message formats (DiscoveryRequest and DiscoveryResponse) to communicate with the server, enhancing protocol clarity and robustness.

Graceful Resource Cleanup: The client handles termination signals gracefully. When a termination signal (e.g., SIGINT or SIGTERM) is received, it closes network connections and exits cleanly, ensuring resource cleanup.

User Interaction: The client allows users to input messages to send to the server, facilitating interactive communication.

Clear Screen: Typing /clear clears the terminal screen, providing a cleaner user interface.

## Dependencies

Go (Golang): The client is implemented in Go and requires a Go environment to build and run.

## Configuration

Default UDP Discovery Port: 10101
Default UDP Server Broadcast Address: 255.255.255.255
Default TCP Connection Port: 12345

## Notes

Ensure that the AsyncTCPServer is running and reachable when starting the AsyncUDP-TCPClient for successful TCP connection.

Customize the UDP discovery port, server broadcast address, and TCP connection port as needed by modifying the source code.

## License

This project is licensed under the MIT License.
