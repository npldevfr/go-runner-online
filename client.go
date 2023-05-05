package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"net"
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

func (c *Client) listen() {
	defer c.conn.Close()

	reader := bufio.NewReader(c.conn)

	for {
		// Read data sent by the server
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading data from server: %v", err)
			break
		}
		// Print the data received from the server
		log.Printf(data)
		if data == "La partie peut commencer !\n" {
			c.ready = true
		}

	}

}
