package main

import (
	"net"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/log"
)

type connAdd struct {
	pid  *actor.PID
	conn net.Conn
}

type connRem struct {
	pid *actor.PID
}

type server struct {
	listenAddr string
	ln         net.Listener
	sessions   map[*actor.PID]net.Conn
}

func NewServer(listenAddr string) actor.Producer {
	return func() actor.Receiver {
		return &server{
			listenAddr: listenAddr,
			sessions:   make(map[*actor.PID]net.Conn),
		}
	}
}

func (s *server) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	// quando o servidor iniciar ele, começa recebendo a mensagem de inicialização
	// logo dps ele ja envia a mensagem de started
	case actor.Initialized:
		ln, err := net.Listen("tcp", s.listenAddr)
		if err != nil {
			panic(err)
		}
		s.ln = ln
	case actor.Started:
		log.Infow("server ", log.M{"addr": s.listenAddr})
		go s.acceptLoop(c)
	case actor.Stopped:
	case *connAdd:
		log.Tracew("added new connection to map", log.M{"addr": msg.conn.RemoteAddr(), "pid": msg.pid})
		s.sessions[msg.pid] = msg.conn
	case *connRem:
		log.Tracew("removing connection from the map", log.M{"pid": msg.pid})
		delete(s.sessions, msg.pid)
	}
}

func (s *server) acceptLoop(c *actor.Context) {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Errorw("accept error", log.M{"error": err})
		}
		// Spamando novo ator quando chegar uma nova conexão
		pid := c.SpawnChild(newSession(conn), "session", actor.WithTags(conn.RemoteAddr().String()))
		c.Send(c.PID(), &connAdd{
			pid:  pid,
			conn: conn,
		})
	}
}
