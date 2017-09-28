# goirc
Teaching myself go. This is (will be) an attempt at building a super simple irc bot in golang.

## Example

	c := goirc.Configuration{}
	c.Server = "irc.someserver.com"
	c.Port = 6667
	c.Nickname = "bot"
	c.RealName = "justabot"
	c.Channel = "#channel"
	c.Timeout = 3
    c.Debug = true
    c.Identd = false

	b := goirc.Bot{}
	b.Connect(c)