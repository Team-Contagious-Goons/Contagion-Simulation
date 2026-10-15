package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	UP    = -1
	DOWN  = 1
	RIGHT = 1
	LEFT  = -1
)

type Agent struct {
	x, y     int
	infected bool
}

type Game struct {
	agents          [][]Agent
	agentsCount     int
	infectedCount   int
	h, w            int
	infectionChance float64
}

func (g *Game) step(step int) {
	fmt.Printf("Game %d steps\n", step)

	newAgents := g.agents

	for x := 0; x < g.w; x++ {
		for y := 0; y < g.h; y++ {
			agent := g.agents[x][y]

			if agent == (Agent{}) {
				continue
			}

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

			// Only move agent
			if newAgents[newX][newY] == (Agent{}) {
				newAgents[newX][newY] = agent
				newAgents[x][y] = Agent{}
			}

			if agent.infected {
				// check up down left right
				if agent.x != 0 && newAgents[newX+LEFT][newY] != (Agent{}) && !newAgents[newX+LEFT][newY].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX+LEFT][newY].infected = true
					}
				}

				if agent.x != g.w-1 && newAgents[newX+RIGHT][newY] != (Agent{}) && !newAgents[newX+RIGHT][newY].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX+RIGHT][newY].infected = true
					}
				}

				if agent.y != 0 && newAgents[newX][newY+UP] != (Agent{}) && !newAgents[newX][newY+UP].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX][newY+UP].infected = true
					}
				}

				if agent.y != g.h-1 && newAgents[newX][newY+DOWN] != (Agent{}) && !newAgents[newX][newY+DOWN].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX][newY+DOWN].infected = true
					}
				}
			}
		}
	}

	g.agents = newAgents
}

func newGame(width int, height int, agentCount int, infectedCount int, infectionChance float64) *Game {
	g := Game{}
	agents := make([][]Agent, width)
	for i := range width {
		agents[i] = make([]Agent, height)
	}

	g.infectedCount = infectedCount

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
	g.infectionChance = infectionChance
	g.agentsCount = agentCount
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

func (g *Game) stats() string {

	endAgents := 0
	endInfected := 0
	for i := 0; i < g.w; i++ {
		for j := 0; j < g.h; j++ {
			if g.agents[i][j] != (Agent{}) {
				endAgents++
				if g.agents[i][j].infected {
					endInfected++
				}
			}
		}
	}

	return fmt.Sprintf(
		"Game Stats: Size: {%d, %d}  Agents: %d  Starting-Infected: %d  Infection Rate: %f\n\tEnd-Agents: %d End-Infected: %d",
		g.w, g.h, g.agentsCount, g.infectedCount, g.infectionChance, endAgents, endInfected)

}
func main() {

	h := 50
	w := 90
	agentCount := 300
	infectedCount := 100
	infectionChance := 0.4

	game := newGame(h, w, agentCount, infectedCount, infectionChance)

	for i := range 50 {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			// ignore
		}
		game.step(i)
		fmt.Print(game.String()) // Clear screen and print field.
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println(game.stats())
}
