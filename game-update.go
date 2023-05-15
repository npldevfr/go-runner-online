/*
//  Implementation of the Update method for the Game structure
//  This method is called once at every frame (60 frames per second)
//  by ebiten, juste before calling the Draw method (game-draw.go).
//  Provided with a few utilitary methods:
//    - CheckArrival
//    - ChooseRunners
//    - HandleLaunchRun
//    - HandleResults
//    - HandleWelcomeScreen
//    - Reset
//    - UpdateAnimation
//    - UpdateRunners
*/

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"strconv"
	"time"
)

// HandleWelcomeScreen waits for the player to push SPACE in order to
// start the game
func (g *Game) HandleWelcomeScreen() bool {
	return g.runners[0].client.globalState == GlobalChooseRunner && inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// ChooseRunners loops over all the runners to check which sprite each
// of them selected
func (g *Game) ChooseRunners() (done bool) {
	done = true
	for i := range g.runners {
		if i == 0 {
			done = g.runners[i].ManualChoose() && done
		} else {
			done = g.runners[i].RandomChoose() && done
		}
	}
	return done
}

// HandleLaunchRun countdowns to the start of a run
func (g *Game) HandleLaunchRun() bool {
	if time.Since(g.f.chrono).Milliseconds() > 1000 {
		g.launchStep++
		g.f.chrono = time.Now()
	}
	if g.launchStep >= 5 {
		g.launchStep = 0
		return true
	}
	return false
}

// UpdateRunners loops over all the runners to update each of them
func (g *Game) UpdateRunners() {
	for i := range g.runners {
		if i == 0 {
			g.runners[i].ManualUpdate()
		}

		if g.runners[i].client.isAI {
			g.runners[i].RandomUpdate()
		}
	}
}

// CheckArrival loops over all the runners to check which ones are arrived
func (g *Game) CheckArrival() (finished bool) {
	finished = true
	for i := range g.runners {
		g.runners[i].CheckArrival(&g.f)
		finished = finished && g.runners[i].arrived
	}
	return finished
}

// Reset resets all the runners and the field in order to start a new run
func (g *Game) Reset() {
	for i := range g.runners {
		g.runners[i].Reset(&g.f)
	}
	g.f.Reset()
}

// UpdateAnimation loops over all the runners to update their sprite
func (g *Game) UpdateAnimation() {
	for i := range g.runners {
		g.runners[i].UpdateAnimation(g.runnerImage)
	}
}

// HandleResults computes the resuls of a run and prepare them for
// being displayed
func (g *Game) HandleResults() bool {
	if time.Since(g.f.chrono).Milliseconds() > 1000 || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep++
		g.f.chrono = time.Now()
	}
	if g.resultStep >= 4 && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep = 0
		return true
	}
	return false
}

// Update is the main update function of the game. It is called by ebiten
// at each frame (60 times per second) just before calling Draw (game-draw.go)
// Depending of the current state of the game it calls the above utilitary
// function and then it may update the state of the game
func (g *Game) Update() error {
	var localRunner Runner = g.runners[0]
	switch g.state {
	case StateWelcomeScreen:
		done := g.HandleWelcomeScreen()

		if done || localRunner.client.globalState == GlobalChooseRunner {

			for i := 0; i < len(localRunner.client.otherClient); i++ {
				client := Client{name: localRunner.client.otherClient[i], isAI: false}
				g.runners[i+1].client = &client
			}

			if len(localRunner.client.otherClient) < 3 {
				for i := len(localRunner.client.otherClient); i < 3; i++ {
					client := Client{name: "IA " + strconv.Itoa(i), isAI: true}
					g.runners[i+1].client = &client
				}
			}

			g.state++
		}
	case StateChooseRunner:
		done := g.ChooseRunners()
		if done {
			g.UpdateAnimation()
			localRunner.client.send("updateSkin", localRunner.colorScheme)
			localRunner.client.send("readyToRun", true)
			g.state++
		}
	case StateLaunchRun:
		if localRunner.client.globalState == GlobalLaunchRun {
			done := g.HandleLaunchRun()
			if done {
				g.state++
			}
		}
	case StateRun:
		g.UpdateRunners()
		finished := g.CheckArrival()
		g.UpdateAnimation()
		if finished && g.noSend {
			localRunner.client.send("runnerLaneFinished", localRunner.runTime)
			g.noSend = false
		}
		if localRunner.client.globalState == GlobalStateResult {
			g.noSend = true
			g.state++
		}
	case StateResult:
		done := g.HandleResults()
		if done {
			g.Reset()
			g.state = StateLaunchRun
		}
	}
	return nil
}
