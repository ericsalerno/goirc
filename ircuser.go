package goirc

import "strings"

// IRCUser represents a user sending a message
type IRCUser struct {
	Nickname string
	RealName string
	Host     string
	Identd   bool
}

// Create a user object from a string
func (u *IRCUser) Create(request string) {
	userString := request[1:]

	exclamationPoint := strings.Index(userString, "!")
	atSymbol := strings.Index(userString, "@")

	if exclamationPoint == -1 || atSymbol == -1 {
		return
	}

	u.Nickname = userString[0:exclamationPoint]
	u.RealName = userString[exclamationPoint+1 : atSymbol]
	u.Host = userString[atSymbol+1:]

	if strings.Contains(u.RealName, "~") {
		u.Identd = false
		u.RealName = u.RealName[1:]
	} else {
		u.Identd = true
	}
}

// AsString prints out the user as a string
func (u IRCUser) AsString() string {

	identified := ""
	if u.Identd == false {
		identified = "~"
	}

	return u.Nickname + "!" + identified + u.RealName + "@" + u.Host
}
