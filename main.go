package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	w := 150
	h := 150
	agentCount := 400
	infectedCount := 200
	infectionChance := 0.9
	infectionDuration := 50

	game := newGame(w, h, agentCount, infectedCount, infectionChance, infectionDuration)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Contagion")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println(game.stats())
}
