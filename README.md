# goirc
Teaching myself go. This is (will be) an attempt at building a super simple irc bot in golang.

This isn't even remotely close to working yet.

## Example

	c := goirc.Configuration{}
	c.Server = "irc.someserver.com"
	c.Port = 6667
	c.Nickname = "bot"
	c.RealName = "justabot"
	c.Channel = "#channel"
	c.Timeout = 3

	b := goirc.Bot{}
	b.Connect(c)