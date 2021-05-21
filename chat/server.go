package chat

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) Run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name(cmd.client, cmd.args)

		case CMD_JOIN:
			s.join(cmd.client, cmd.args)

		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)

		case CMD_MSG:
			s.message(cmd.client, cmd.args)

		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)

		}
	}
}

func (s *server) NewClient(conn net.Conn) {
	log.Printf("New Client Has Connected: %s", conn.RemoteAddr())
	c := &client{
		conn:     conn,
		name:     "anon",
		commands: s.commands,
	}
	c.readInput()
}

func (s *server) name(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("Ok I will call you %s", c.name))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]

	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c
	s.quitRoom(c)

	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.name))
	c.msg(fmt.Sprintf("Welcome to the room %s %s", r.name, c.name))
}
func (s *server) listRooms(c *client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("Avaiable rooms are: %s", strings.Join(rooms, "\n")))
}
func (s *server) message(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("You must join a room first"))
		return
	}
	c.room.broadcast(c, c.name+": "+strings.Join(args[1:], " "))
}
func (s *server) quit(c *client, args []string) {
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr())
	s.quitRoom(c)
	c.msg("Sad, bye bye")
}

func (s *server) quitRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room :(", c.name))
	}
}
