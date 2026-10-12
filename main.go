package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Agent struct {
	x, y     int
	infected bool
}

type Game struct {
	agents [][]Agent
	h, w   int
}

func (g *Game) step(step int) {
	fmt.Printf("Game %d steps\n", step)

	newAgents := make([][]Agent, g.w)
	for i := range newAgents {
		newAgents[i] = make([]Agent, g.h)
	}

	for x := 0; x < g.w; x++ {
		for y := 0; y < g.h; y++ {
			agent := g.agents[x][y]

			if agent != (Agent{}) {
				xStep := rand.Intn(3) - 1
				yStep := rand.Intn(3) - 1

				newX := agent.x + xStep
				newY := agent.y + yStep

				// Keep agent inside bounds
				if newX < 0 {
					newX = 0
				}
				if newY < 0 {
					newY = 0
				}
				if newX >= g.w {
					newX = g.w - 1
				}
				if newY >= g.h {
					newY = g.h - 1
				}

				agent.x = newX
				agent.y = newY

				if newAgents[newX][newY] == (Agent{}) {
					newAgents[newX][newY] = agent
				}
			}
		}
	}

	g.agents = newAgents
}

func newGame(width int, height int, agentCount int, infectedCount int) *Game {
	g := Game{}
	agents := make([][]Agent, width)
	for i := range width {
		agents[i] = make([]Agent, height)
	}

	for range agentCount {
		randX := rand.Intn(width)
		randY := rand.Intn(height)
		agents[randX][randY].x = randX
		agents[randX][randY].y = randY
		if infectedCount != 0 {
			agents[randX][randY].infected = true
			infectedCount--
		} else {
			agents[randX][randY].infected = false
		}

	}

	g.agents = agents
	g.w = width
	g.h = height
	return &g
}

func (g *Game) String() string {
	var str strings.Builder
	for i := range g.w {
		for j := range g.h {
			if g.agents[i][j] == (Agent{}) {
				str.WriteString("  ")
			} else {
				if g.agents[i][j].infected {
					str.WriteString(fmt.Sprintf("X "))
				} else {
					str.WriteString(fmt.Sprintf("O "))
				}
			}
		}
		str.WriteString("\n")
	}
	return str.String()
}

func main() {

	game := newGame(20, 20, 40, 10)

	for i := range 20 {
		game.step(i)
		fmt.Print("\x0c", game.String()) // Clear screen and print field.
		time.Sleep(500 * time.Millisecond)
	}
}
