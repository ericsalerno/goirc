package goirc

import (
	"fmt"
	"net"
	"strings"
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

	if b.server != nil {
		b.server.Close()
	}
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
		n, readErr := conn.Read(data[0:])

		if readErr != nil {
			fmt.Printf("Failed reading: %s", readErr)
			shouldLoop = false
		}

		if len(data) != 0 {
			b.processServerResponse(string(data[:n]))
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

func (b Bot) processServerResponse(response string) {

	if response == "" {
		return
	}

	response = strings.TrimRight(response, "\r\n")

	//Check to see if this is a compound set of lines seperated by \n and if so run them through individually
	if strings.Contains(response, "\n") {
		stringSlice := strings.Split(response, "\n")

		for i := 0; i < len(stringSlice); i++ {
			b.processServerResponse(stringSlice[i])
		}
		return
	}

	fmt.Println("<- " + response)

	segments := strings.Split(response, " ")

	if len(segments) == 0 {
		return
	}

	if segments[0] == "PING" {
		b.sendRawCommand("PONG", segments[1])
		return
	}

	if strings.HasPrefix(response, ":") {
		if strings.Contains(segments[0], "!") {
			//This is a username
			if segments[1] == "PRIVMSG" {
				message := IRCMessage{}
				message.FromServerResponse(response)

				b.onUserMessage(message)
			}
		} else {
			//This is a server message
			if len(segments) > 1 {
				b.respondToServerEvent(segments[1], segments)
			}
		}
	}
}

func (b Bot) respondToServerEvent(event string, parameters []string) {
	switch event {
	case "376":
		b.onMOTDEnd()
	case "396":
		b.onConnected()
	}
}

func (b Bot) onMOTDEnd() {
	//b.sendRawCommand("JOIN", b.config.Channel)
}

func (b Bot) onConnected() {
	b.sendRawCommand("JOIN", b.config.Channel)
}

func (b Bot) onUserMessage(message IRCMessage) {
	fmt.Printf("%s with %s\n", message.Command, message.Parameters)
	if message.Command == "!say" {
		b.sendRawCommand("PRIVMSG", message.Channel+" :"+message.Parameters)
	}

}
