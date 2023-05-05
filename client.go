package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"net"
	"strings"
)

type Client struct {
	conn   net.Conn
	name   string
	runner Runner
	ready  bool
}

func NewClient() *Client {
	return &Client{
		ready: false,
	}
}

func (c *Client) connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	c.name = conn.RemoteAddr().String()
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
		case "gameStart":
			c.ready = true
		default:
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
