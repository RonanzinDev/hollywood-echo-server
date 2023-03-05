package main

import (
	"net"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/log"
)

// Toda nova conexão no server é uma nova session
// Toda nova sessão vai ser um child, e o server é o Parent
type session struct {
	conn net.Conn
}

// Gerando um novo ator talvez?
func newSession(conn net.Conn) actor.Producer {
	return func() actor.Receiver {
		return &session{
			conn: conn,
		}
	}
}

// A função receiver é a base de todo actor
func (s *session) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
	case actor.Started:
		log.Infow("new connection", log.M{"addr": s.conn.RemoteAddr()})
		go s.readLoop(c)
	case actor.Stopped:
	case []byte:
		s.conn.Write(msg)
	}
}

func (s *session) readLoop(c *actor.Context) {
	buffer := make([]byte, 2048)
	for {
		n, err := s.conn.Read(buffer)
		if err != nil {
			log.Errorw("conn read error", log.M{"err": err})
			break
		}
		msg := buffer[:n]
		c.Send(c.PID(), msg)
	}
	// quando o loop acabar, vamos remover o child do map de conexões
	// C.Parant() vai retonar o PID do Pai(que é o server) e o c.Pid() é o PID do ator atual(que no caso é um child/sessao)
	c.Send(c.Parent(), &connRem{pid: c.PID()})
}
