package goirc

import (
	"fmt"
	"net"
)

// Bot - Main bot class
type Bot struct {
	config Configuration
}

// Connect - connect to a server
func (b Bot) Connect(config Configuration) {
	b.config = config

	go b.serverPump()

	var input string
	fmt.Scanln(&input)
}

func (b Bot) serverPump() {
	serverConnection := fmt.Sprintf("%s:%d", b.config.Server, b.config.Port)
	fmt.Printf("Connecting to %s...\n", serverConnection)

	conn, err := net.Dial("tcp", serverConnection)

	if err != nil {
		fmt.Printf("Failed to connect to " + serverConnection)
		return
	}

	shouldLoop := true
	for shouldLoop {
		var data []byte
		_, readErr := conn.Read(data)

		if readErr != nil {
			fmt.Printf("Failed reading: %s", readErr)
			shouldLoop = false
		}

		if len(data) != 0 {
			fmt.Printf("Server: %s\n", data)
		} else {
			shouldLoop = false
		}

	}

	fmt.Println("Closing connection...")
	conn.Close()
}
