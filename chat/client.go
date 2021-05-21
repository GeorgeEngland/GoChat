package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	name     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			fmt.Printf("ERROR %s\n", err.Error())
			c.conn.Close()
			return
		}
		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")

		cmd := strings.TrimSpace(args[0])

		com := command{client: c, args: args}

		switch cmd {
		case "/name":
			com.id = CMD_NAME
		case "/join":
			com.id = CMD_JOIN
		case "/rooms":
			com.id = CMD_ROOMS
		case "/quit":
			com.id = CMD_QUIT
		case "/msg":
			com.id = CMD_MSG
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
			c.msg("USAGE: /{COMMAND} {option}\nCommands\n/name {your name}\n/join {room name}\n/rooms\n/msg {your message}\n/quit")
			continue
		}
		c.commands <- com

	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
