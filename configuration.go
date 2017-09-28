package goirc

// Configuration - goirc bot configuration class
type Configuration struct {
	Server   string
	Port     int
	Nickname string
	RealName string
	Channel  string

	Timeout int
	Debug   bool
	Identd  bool
}
