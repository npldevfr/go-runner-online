package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn        net.Conn
	name        string
	globalState int
	otherClient []string
}

const (
	GlobalWelcomeScreen int = iota
	GlobalChooseRunner
	GlobalLaunchRun
)

func NewClient() *Client {
	return &Client{
		globalState: GlobalWelcomeScreen,
	}
}

func (c *Client) connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) sendMessage(message string) {
	// send message to server using bufio writer
	writer := bufio.NewWriter(c.conn)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		log.Printf("Error sending message to server: %v", err)
	} else {
		writer.Flush()
	}
}

func (c *Client) send(key string, data interface{}) {
	// encode data using gob
	gob.Register(time.Duration(0))

	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(&data)
	if err != nil {
		log.Printf("Error encoding data: %v", err)
		return
	}
	// send message to server using bufio writer
	writer := bufio.NewWriter(c.conn)
	_, err = writer.WriteString(key + ":" + base64.StdEncoding.EncodeToString(buffer.Bytes()) + "\n")
	if err != nil {
		log.Printf("Error sending message to server: %v", err)
	} else {
		writer.Flush()
	}
}

// Start listening for messages from the server
// Start listening for messages from the server
func (c *Client) listen() {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}(c.conn)

	reader := bufio.NewReader(c.conn)

	for {
		// Read data sent by the server
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading data from server: %v", err)
			break
		}

		// split the message into key and data
		parts := strings.Split(data, ":")
		if len(parts) < 2 {
			log.Printf("Invalid message format: %v", data)
			continue
		}
		key := parts[0]
		encodedData := parts[1]

		//Register "name" for gob
		//name not registered for interface: "[]interface {}"
		gob.Register([]interface{}{})
		gob.Register(map[string]interface{}{})
		gob.Register(time.Duration(0))

		// decode the data using base64 and gob
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			log.Printf("Error decoding data: %v", err)
			continue
		}

		var eventData interface{}
		err = gob.NewDecoder(bytes.NewReader(decodedData)).Decode(&eventData)
		if err != nil {
			log.Printf("Error decoding data: %v", err)
			continue
		}

		// Switch statement to handle different keys
		switch key {
		case "newLocalRemote":
			c.name = eventData.(string)
		case "newPlayer":
			if eventData.(string) != c.name {
				c.otherClient = append(c.otherClient, eventData.(string))
				log.Printf("New player: %v", eventData)
			}
		case "gameStart":
			c.globalState = GlobalChooseRunner
		case "gameEnd":
			for _, item := range eventData.([]interface{}) {
				if data, ok := item.(map[string]interface{}); ok {
					if data["name"] != c.name {
						log.Printf("%v a fait %v", data["name"], data["duration"])
					}
				}
			}
		}
	}
}

// Listen for a specific key and return the event data
func (c *Client) listenForKey(key string) interface{} {
	encodedData := ""

	// Wait for message with the specified key
	for {
		data, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Printf("Error reading data from server: %v", err)
			return nil
		}

		// split the message into key and data
		parts := strings.Split(data, ":")
		if len(parts) < 2 {
			log.Printf("Invalid message format: %v", data)
			continue
		}
		eventKey := parts[0]
		encodedData = parts[1]

		if eventKey == key {
			break
		}
	}

	// decode the data using base64 and gob
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		log.Printf("Error decoding data: %v", err)
		return nil
	}

	var eventData interface{}
	err = gob.NewDecoder(bytes.NewReader(decodedData)).Decode(&eventData)
	if err != nil {
		log.Printf("Error decoding data: %v", err)
		return nil
	}

	return eventData
}
