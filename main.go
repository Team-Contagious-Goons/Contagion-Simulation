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
	x, y           int
	infected       bool
	infectionTimer int
}

type Game struct {
	agents            [][]Agent
	agentsCount       int
	infectedCount     int
	h, w              int
	infectionChance   float64
	infectionDuration int
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

			if agent.infected {
				targetX, targetY, found := g.find(agent.x, agent.y)
				if found {
					if targetX > agent.x {
						xStep = RIGHT
					} else if targetX < agent.x {
						xStep = LEFT
					} else {
						xStep = 0
					}

					if targetY > agent.y {
						yStep = DOWN
					} else if targetY < agent.y {
						yStep = UP
					} else {
						yStep = 0
					}
				}
			}
			if !agent.infected {
				bestMoveX, bestMoveY, found := g.escape(agent.x, agent.y)
				if found {
					if bestMoveX > agent.x {
						xStep = RIGHT
					} else if bestMoveX < agent.x {
						xStep = LEFT
					} else {
						xStep = 0
					}

					if bestMoveY > agent.y {
						yStep = DOWN
					} else if bestMoveY < agent.y {
						yStep = UP
					} else {
						yStep = 0
					}
				}
			}

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

			moved := false
			// Only move agent
			if newAgents[newX][newY] == (Agent{}) {
				newAgents[newX][newY] = agent
				newAgents[x][y] = Agent{}
				moved = true
			}

			if agent.infected {

				if moved {
					newAgents[newX][newY].infectionTimer--
					if newAgents[newX][newY].infectionTimer == 0 {
						newAgents[newX][newY].infected = false
					}
				} else {
					newAgents[x][y].infectionTimer--
					if newAgents[x][y].infectionTimer == 0 {
						newAgents[x][y].infected = false
					}
				}
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

func (g *Game) find(x, y int) (int, int, bool) {
	minDist := 99
	targetX, targetY := -1, -1
	found := false

	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			nx, ny := x+dx, y+dy

			if nx >= 0 && nx < g.w && ny >= 0 && ny < g.h && (dx != 0 || dy != 0) {
				target := g.agents[nx][ny]

				if target != (Agent{}) && !target.infected {
					dist := dx*dx + dy*dy
					if dist < minDist {
						minDist = dist
						targetX, targetY = nx, ny
						found = true
					}
				}
			}
		}
	}
	return targetX, targetY, found

}
func (g *Game) escape(x, y int) (int, int, bool) {
	closestInfectedDist := 999
	infX, infY := -1, -1
	foundInfected := false

	scanRadius := 3
	for dx := -scanRadius; dx <= scanRadius; dx++ {
		for dy := -scanRadius; dy <= scanRadius; dy++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < g.w && ny >= 0 && ny < g.h && (dx != 0 || dy != 0) {
				target := g.agents[nx][ny]

				if target != (Agent{}) && target.infected {
					dist := dx*dx + dy*dy
					if dist < closestInfectedDist {
						closestInfectedDist = dist
						infX, infY = nx, ny
						foundInfected = true
					}
				}
			}
		}
	}
	if !foundInfected {
		return x, y, false
	}

	bestMoveX, bestMoveY := x, y
	maxDistFromInfected := -1

	for moveDX := -1; moveDX <= 1; moveDX++ {
		for moveDY := -1; moveDY <= 1; moveDY++ {
			nx, ny := x+moveDX, y+moveDY

			if nx >= 0 && nx < g.w && ny >= 0 && ny < g.h {
				distToInfected := (nx-infX)*(nx-infX) + (ny-infY)*(ny-infY)

				if distToInfected > maxDistFromInfected {
					maxDistFromInfected = distToInfected
					bestMoveX, bestMoveY = nx, ny
				}
			}
		}
	}

	return bestMoveX, bestMoveY, true
}
func main() {

	h := 50
	w := 50
	agentCount := 300
	infectedCount := 1
	infectionChance := 1.0
	infectionDuration := 3

	game := newGame(h, w, agentCount, infectedCount, infectionChance, infectionDuration)

	for i := range 50 {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			// ignore
		}
		game.step(i)
		fmt.Print(game.String()) // Clear screen and print field.
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println(game.stats())
}
