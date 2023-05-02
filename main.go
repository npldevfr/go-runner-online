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

	var getTPS bool
	var isServer bool
	var isClient bool
	flag.BoolVar(&getTPS, "tps", false, "Afficher le nombre d'appel à Update par seconde")
	flag.BoolVar(&isServer, "server", false, "Lancer le jeu en mode serveur")
	flag.BoolVar(&isClient, "client", false, "Lancer le jeu en mode client avec l'adresse IP du serveur")
	flag.Parse()

	if isServer {
		InitServer()
		return
	}

	if isClient {
		InitClient()
		return
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("LP MiAR -- Programmation répartie (UE03EC2)")

	g := InitGame()
	g.getTPS = getTPS

	err := ebiten.RunGame(&g)
	log.Print(err)

}
