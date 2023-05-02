package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

func server() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Println("accept error:", err)
		return
	}
	//lecture du message envoyé par le client avec bufio
	status, _ := bufio.NewReader(conn).ReadString('\n')
	log.Println(status)

	defer conn.Close()
	log.Println("Le client s'est connecté")

	time.Sleep(10 * time.Second)

}
