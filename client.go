package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn   net.Conn
	name   string
	runner Runner
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) sendMessage(message string) error {
	writer := bufio.NewWriter(c.conn)
	_, err := writer.WriteString(message)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
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
		fmt.Print(data)
	}
}
