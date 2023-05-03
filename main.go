/*
// Implementation of a main function setting a few characteristics of
// the game window, creating a game, and launching it
*/

package main

import (
	"flag"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800 // Width of the game window (in pixels)
	screenHeight = 160 // Height of the game window (in pixels)
)

func main() {

	var (
		getTPS   bool
		isServer bool
		clientIP string
	)
	flag.BoolVar(&getTPS, "tps", false, "Afficher le nombre d'appel à Update par seconde")
	flag.BoolVar(&isServer, "server", false, "Lancer le jeu en mode serveur")
	flag.StringVar(&clientIP, "client", "", "Connecter le jeu en tant que client à l'adresse IP spécifiée")
	flag.Parse()

	// If the game is launched in server mode, start the server
	if isServer {
		s := NewServer(":8080")
		err := s.Start()
		if err != nil {
			return
		}
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("LP MiAR -- Programmation répartie (UE03EC2)")

	g := InitGame()
	g.getTPS = getTPS

	// If the clientIP flag is set, create a client
	if clientIP != "" {
		g.createClient(clientIP)
	}

	err := ebiten.RunGame(&g)
	log.Print(err)

}
