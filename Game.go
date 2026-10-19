package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	UP    = -1
	DOWN  = 1
	RIGHT = 1
	LEFT  = -1
)

const (
	screenWidth  = 320
	screenHeight = 340
)

type Game struct {
	agents            [][]Agent
	agentsCount       int
	infectedCount     int
	h, w              int
	infectionChance   float64
	infectionDuration int
	iteration         int
	lastStepTime      time.Time
	stepDelay         int
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Clear()
	gridOffset := 50

	// Draw the agents to the screen
	for _, agent := range g.agents {
		for _, agentState := range agent {
			if agentState != (Agent{}) {
				if agentState.infected {
					vector.FillCircle(screen, float32(agentState.x+gridOffset), float32(agentState.y+gridOffset), 2.5, colornames.Red, true)
				} else {
					vector.FillCircle(screen, float32(agentState.x+gridOffset), float32(agentState.y+gridOffset), 2.5, colornames.Green, true)
				}
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nStep: %d", g.iteration))
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	// step every 10 millis
	//if time.Duration(g.lastStepTime)
	g.stepDelay--
	if g.stepDelay == 0 {
		g.step()
		g.lastStepTime = time.Now()
		g.iteration++
		g.stepDelay = 5
	}
	return nil
}

func (g *Game) step() {
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

			moved := false
			// Only move agent
			if newAgents[newX][newY] == (Agent{}) {
				agent.x = newX
				agent.y = newY
				newAgents[newX][newY] = agent
				newAgents[x][y] = Agent{}
				moved = true
			} else {
				newX = x
				newY = y
			}

			if agent.infected {

				if moved {
					newAgents[newX][newY].infectionTimer--
					if newAgents[newX][newY].infectionTimer <= 0 {
						newAgents[newX][newY].infected = false
					}
				} else {
					newAgents[x][y].infectionTimer--
					if newAgents[x][y].infectionTimer <= 0 {
						newAgents[x][y].infected = false
					}
				}
				// check up down left right
				if agent.x != 0 && newAgents[newX+LEFT][newY] != (Agent{}) && !newAgents[newX+LEFT][newY].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX+LEFT][newY].infected = true
						newAgents[newX+LEFT][newY].infectionTimer = g.infectionDuration

					}
				}

				if agent.x != g.w-1 && newAgents[newX+RIGHT][newY] != (Agent{}) && !newAgents[newX+RIGHT][newY].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX+RIGHT][newY].infected = true
						newAgents[newX+RIGHT][newY].infectionTimer = g.infectionDuration
					}
				}

				if agent.y != 0 && newAgents[newX][newY+UP] != (Agent{}) && !newAgents[newX][newY+UP].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX][newY+UP].infected = true
						newAgents[newX][newY+UP].infectionTimer = g.infectionDuration

					}
				}

				if agent.y != g.h-1 && newAgents[newX][newY+DOWN] != (Agent{}) && !newAgents[newX][newY+DOWN].infected {
					shouldInfect := rand.Float64() < g.infectionChance
					if shouldInfect {
						newAgents[newX][newY+DOWN].infected = true
						newAgents[newX][newY+DOWN].infectionTimer = g.infectionDuration
					}
				}
			}
		}
	}

	g.agents = newAgents
}

func newGame(width int, height int, agentCount int, infectedCount int, infectionChance float64, infectionDuration int) *Game {
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
			agents[randX][randY].infectionTimer = infectionDuration
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
	g.infectionDuration = infectionDuration
	g.lastStepTime = time.Now()
	g.stepDelay = 20
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
		"Game Stats: Size: {%d, %d}  Agents: %d  Starting-Infected: %d  Infection Rate: %f  Infection Duration: %d \nEnd-Agents: %d End-Infected: %d",
		g.w, g.h, g.agentsCount, g.infectedCount, g.infectionChance, g.infectionDuration, endAgents, endInfected)

}
