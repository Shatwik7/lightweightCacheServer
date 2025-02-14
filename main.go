package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/gofiber/fiber/v2/log"
	"github.com/tidwall/resp"
)

const defaultListenAddress = ":2345"

type Message struct {
	cmd  Command
	peer *Peer
}

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message
	delPeerCh chan *Peer
	KeyVal    *KeyVal
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		delPeerCh: make(chan *Peer),
		KeyVal:    NewKeyVal(10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()
	log.Infof("Running on PORT : ", s.Config.ListenAddress)
	return s.acceptLoop()
}
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message eror", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Error("Accept Error : Can not accept :", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case ClientCommand:
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil {
			return err
		}
	case SetCommand:
		if err := s.KeyVal.Set(v.key, v.val); err != nil {
			return err
		}
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil {
			return err
		}
	case GetCommand:
		val, ok := s.KeyVal.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString(string(val)); err != nil {
			return err
		}
	case HelloCommand:
		spec := map[string]string{
			"server": "redis",
		}
		_, err := msg.peer.Send(respWriteMap(spec))
		if err != nil {
			return fmt.Errorf("peer send error: %s", err)
		}
	}
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- peer
	slog.Info("new peer connected", "remote Addr :", conn.RemoteAddr().String())
	if err := peer.readLoop(); err != nil {
		log.Error("Peer read error", err)
	}
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start().Error())
	fmt.Println("hi!")
}
