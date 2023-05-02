package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	address string
	players []*player
}

type player struct {
	conn net.Conn
	name string
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		players: make([]*player, 0),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("Impossible de démarrer le serveur : %v", err)
	}
	log.Printf("Serveur en attente de connexions sur %s", s.address)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Nouvelle connexion : %s", conn.RemoteAddr().String())
		s.addPlayer(conn)
	}
}

func (s *Server) addPlayer(conn net.Conn) {
	p := &player{
		conn: conn,
		name: conn.RemoteAddr().String(),
	}
	s.players = append(s.players, p)
	go p.listen(s)
}

func (p *player) listen(s *Server) {
	defer p.conn.Close()

	reader := bufio.NewReader(p.conn)

	for {
		// Read data sent by the player
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading data from player %s: %v", p.name, err)
			break
		}
		// Process the data and send it to the other players
		s.broadcast(p, data)
	}
}

func (s *Server) broadcast(sender *player, data string) {
	for _, p := range s.players {
		if p != sender {
			writer := bufio.NewWriter(p.conn)
			writer.WriteString(fmt.Sprintf("%s: %s", sender.name, data))
			writer.Flush()
		}
	}
}
