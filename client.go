package main

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	conn   net.Conn
	name   string
	runner Runner
}

func NewClient(name string) *Client {
	return &Client{
		name: name,
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
	}
}
