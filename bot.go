package goirc

import (
	"fmt"
	"net"
)

// Bot - Main bot class
type Bot struct {
	config Configuration
	server net.Conn
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
	b.server = conn

	if err != nil {
		fmt.Printf("Failed to connect to " + serverConnection)
		return
	}

	b.sendRawCommand("USER", b.config.Nickname+" ericsalerno/goirc_1.0 "+b.config.Nickname+" :"+b.config.RealName)
	b.sendRawCommand("NICK", b.config.Nickname)

	shouldLoop := true
	for shouldLoop {
		var data [512]byte
		_, readErr := conn.Read(data[0:])

		if readErr != nil {
			fmt.Printf("Failed reading: %s", readErr)
			shouldLoop = false
		}

		if len(data) != 0 {
			fmt.Printf("Server: %s\n", data)
		} else {
			//shouldLoop = false
		}

	}

	fmt.Println("Closing connection...")
	conn.Close()
}

func (b Bot) sendRawCommand(command string, message string) {

	var commandString = command
	if message != "" {
		commandString = commandString + " " + message
	}

	b.server.Write([]byte(commandString + "\r\n"))
	fmt.Println("-> " + commandString)
}
