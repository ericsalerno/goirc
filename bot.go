package goirc

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Bot - Main bot class
type Bot struct {
	config Configuration
	server net.Conn
	start  time.Time
}

// Connect - connect to a server
func (b Bot) Connect(config Configuration) {
	b.config = config
	b.start = time.Now()

	if b.config.Identd {
		go b.runIdentdServer()
	}

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
		}
	}

	fmt.Println("Closing connection...")
	conn.Close()
}

func (b Bot) runIdentdServer() {
	fmt.Println("Running identd server...")

	ln, err := net.Listen("tcp", ":113")
	if err != nil || ln == nil {
		fmt.Println("Identd failed to bind on port 113...")
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Identd failed to accept connection: %s\n", err)
			return
		}

		request := ""

		for {
			var data [512]byte
			n, readErr := conn.Read(data[0:])

			if readErr != nil {
				fmt.Printf("Identd failed to read incoming connection: %s\n", readErr)
			}

			if n != 0 {
				request = string(data[:n])

				if strings.Contains(request, ", ") {
					break
				}
			}
		}

		fmt.Println(request)
		requestPorts := strings.Split(request, ", ")

		response := requestPorts[0] + ", " + requestPorts[1] + " : USERID : UNIX : " + b.config.RealName

		fmt.Println("Identd " + request + " -> " + response)
		fmt.Fprintf(conn, response)

		break
	}

	fmt.Println("Closing down identd server...")
}

func (b Bot) sendRawCommand(command string, message string) {

	var commandString = command
	if message != "" {
		commandString = commandString + " " + message
	}

	b.server.Write([]byte(commandString + "\r\n"))

	if b.config.Debug {
		fmt.Println("-> " + commandString)
	}
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

	if b.config.Debug {
		fmt.Println("<- " + response)
	}

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
	fmt.Printf("Connected to %s:%d!\n", b.config.Server, b.config.Port)
	b.sendRawCommand("JOIN", b.config.Channel)
}

func (b Bot) onUserMessage(message IRCMessage) {
	if message.Command == "!say" {
		b.sendRawCommand("PRIVMSG", message.Channel+" :"+message.Parameters)
	}

	if message.Command == "!join" {
		b.sendRawCommand("JOIN", message.Parameters)
	}

	if message.Command == "!part" {
		b.sendRawCommand("PART", message.Parameters)
	}

	if message.Command == "!sv" {
		b.sendRawCommand("PRIVMSG", message.Channel+" :ericsalerno/go-irc 1.0, fork me on github... or don't, I'm not the police.")
	}

	if message.Command == "!uptime" {
		now := time.Now()
		now.Sub(b.start)

		uptime := fmt.Sprintf("I've been running for %s seconds!", now.String())
		b.sendRawCommand("PRIVMSG", message.Channel+" :I've been running for "+uptime)
	}

	/*if message.Command == "!quit" {
		b.sendRawCommand("QUIT", message.Parameters)
	}*/

}
