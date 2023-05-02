package main

import (
	"log"
	"net"
)

func client() {

	conn, err := net.Dial("tcp", "172.21.65.39:8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	//envoi du message au serveur
	_, err = conn.Write([]byte("Je suis le client"))

	if err != nil {
		log.Println("Flush error:", err)
		return
	}
	defer conn.Close()

	log.Println("Je suis connecté")

}
