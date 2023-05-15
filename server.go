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
	address   string
	players   []*Runner
	nbPlayers int
}

func NewServer(address string, nbPlayers int) *Server {
	return &Server{
		nbPlayers: nbPlayers,
		address:   address,
	}
}

// start the server
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
		s.broadcast("nbPlayers", len(s.players))
		log.Printf("Nombre de joueurs connectés : %d", len(s.players))
		if len(s.players) == s.nbPlayers {

			// Quand on a x joueurs, on envoie la liste des autres joueurs à chaque joueur
			for _, c := range s.players {
				s.broadcast("newPlayer", c.client.conn.RemoteAddr().String())
			}

			// Début de la partie
			s.broadcast("gameCharacterSelection", nil)
		}
	}
}

// add new player to the server
func (s *Server) addPlayer(conn net.Conn) {
	p := &Client{
		conn: conn,
		name: conn.RemoteAddr().String(),
	}
	r := &Runner{
		client:  p,
		runTime: 0,
	}

	s.players = append(s.players, r)
	p.send("newLocalRemote", p.name)

	go s.listen(r)
}

// listen for incoming message from the client
func (s *Server) listen(r *Runner) {
	// read data sent by the client using a bufio reader
	reader := bufio.NewReader(r.client.conn)

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
		log.Printf("Message reçu de %s avec la clé %s: %v (%T)", r.client.name, key, data, data)

		// switch on the key to handle the message and select the correct state of game
		switch key {
		case "readyToRun":
			r.colorSelected = data.(bool)

			var allPlayersReady = true
			for _, p := range s.players {
				if !p.colorSelected {
					allPlayersReady = false
				}
			}

			if allPlayersReady {
				fmt.Printf("allPlayersReady")
				s.broadcast("gameStart", nil)
			}

		case "updateSkin":
			s.broadcast("updateSkin", map[string]interface{}{
				"name": r.client.name,
				"skin": data,
			})
		case "updatePos":
			s.broadcast("updatePos", map[string]interface{}{
				"name": r.client.name,
				"pos":  data,
			})
		case "runnerLaneFinished":
			r.runTime = data.(time.Duration)

			var allPlayersFinished = true
			for _, p := range s.players {
				if p.runTime == 0 {
					allPlayersFinished = false
				}
			}

			if allPlayersFinished {
				var runDurations []interface{}
				for _, p := range s.players {
					runDurations = append(runDurations, map[string]interface{}{
						"name":     p.client.name,
						"duration": p.runTime,
					})
				}

				s.broadcast("gameEnd", runDurations)
			}
		case "readyToReRun":
			if len(s.players) == s.nbPlayers {
				s.players = []*Runner{}
			}
			s.addPlayer(r.client.conn)
			if len(s.players) == s.nbPlayers {
				time.Sleep(1 * time.Second)
				s.broadcast("gameStart", nil)
			}
		}
	}
}

// broadcast a message to all clients and encode the data using gob and base64
func (s *Server) broadcast(key string, data interface{}) {
	gob.Register(time.Duration(0))
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	for _, c := range s.players {

		var buffer bytes.Buffer
		err := gob.NewEncoder(&buffer).Encode(&data)
		if err != nil {
			log.Printf("Error encoding data broadcast: %v", err)
			return
		}
		writer := bufio.NewWriter(c.client.conn)
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
