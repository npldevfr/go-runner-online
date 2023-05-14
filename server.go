package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	address string
	players []*Client
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("impossible de démarrer le serveur : %v", err)
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
		log.Printf("Nombre de joueurs connectés : %d", len(s.players))
		if len(s.players) == 2 {
			// Quand on a 2 joueurs, on envoie la liste des autres joueurs à chaque joueur
			for _, c := range s.players {
				s.broadcast("newPlayer", c.conn.RemoteAddr().String())
			}

			s.broadcast("gameStart", nil)
			log.Printf("Start...")
		}
	}
}

func (s *Server) addPlayer(conn net.Conn) {
	p := &Client{
		conn: conn,
		name: conn.RemoteAddr().String(),
	}
	s.players = append(s.players, p)
	p.send("newLocalRemote", p.name)

	go s.listen(p)
}

func (s *Server) listen(c *Client) {
	// read data sent by the client using a bufio reader
	reader := bufio.NewReader(c.conn)

	fmt.Printf("hello")

	for {
		// read data sent by the client
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading data from client: %v", err)
			break
		}

		// split the message into key and data
		parts := strings.Split(message, ":")
		if len(parts) < 2 {
			log.Printf("Invalid message format: %v", message)
			continue
		}
		key := parts[0]
		encodedData := parts[1]

		// decode the data using base64 and gob
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			log.Printf("Error decoding data: %v", err)
			continue
		}
		var data interface{}
		err = gob.NewDecoder(bytes.NewReader(decodedData)).Decode(&data)
		if err != nil {
			log.Printf("Error decoding data: %v", err)
			continue
		}

		// print the message received from the client
		log.Printf("Message reçu de %s avec la clé %s: %v (%T)", c.name, key, data, data)

		switch key {
		case "runnerLaneFinished":

		}

		s.broadcast(key, data)
	}
}

func (s *Server) broadcast(key string, data interface{}) {

	for _, c := range s.players {
		gob.Register(time.Duration(0))

		var buffer bytes.Buffer
		err := gob.NewEncoder(&buffer).Encode(&data)
		if err != nil {
			log.Printf("Error encoding data: %v", err)
			return
		}
		writer := bufio.NewWriter(c.conn)
		_, err = writer.WriteString(key + ":" + base64.StdEncoding.EncodeToString(buffer.Bytes()) + "\n")
		if err != nil {
			log.Printf("Error sending message to server: %v", err)
		} else {
			err := writer.Flush()
			if err != nil {
				return
			}
		}
	}
}

// sendMessage est déprécié, car remplacé par broadcast
func (s *Server) sendMessage(c *Client, message string) {
	// send message to server using bufio writer
	writer := bufio.NewWriter(c.conn)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		log.Printf("Error sending message to server: %v", err)
	} else {
		writer.Flush()
	}
}
