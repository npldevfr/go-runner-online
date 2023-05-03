package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	address string
	players []*Client
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		players: make([]*Client, 0),
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
	p := &Client{
		conn: conn,
		name: conn.RemoteAddr().String(),
	}
	s.players = append(s.players, p)
	go s.listen(p)
}

func (s *Server) listen(c *Client) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(c.conn)

	reader := bufio.NewReader(c.conn)

	for {
		// Read data sent by the client
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading data from client: %v", err)
			break
		}
		// Print the data received from the client
		log.Printf("%s: %s", c.name, data)
		// Broadcast the message to all the other clients
		s.broadcast(c, data)
	}
}

func (s *Server) broadcast(sender *Client, data string) {
	for _, p := range s.players {
		if p != sender {
			writer := bufio.NewWriter(p.conn)
			writer.WriteString(fmt.Sprintf("%s: %s", sender.name, data))
			writer.Flush()
		}
	}
}
