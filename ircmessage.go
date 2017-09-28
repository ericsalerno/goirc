package goirc

import "strings"

// IRCMessage represents a private message string
type IRCMessage struct {
	User    IRCUser
	Channel string
	Message string

	Command    string
	Parameters string
}

// FromServerResponse creates an object from a response string
func (m *IRCMessage) FromServerResponse(response string) {

	responseSegments := strings.Split(response, " ")

	m.Channel = responseSegments[2]

	messageStart := strings.Index(response[1:], ":")
	m.Message = response[messageStart+2:]

	m.User.Create(responseSegments[0])

	if len(m.Message) != 0 {
		space := strings.Index(m.Message, " ")
		if space != -1 {
			m.Command = strings.ToLower(m.Message[0:space])
			m.Parameters = m.Message[space+1:]
		}
	}
}
