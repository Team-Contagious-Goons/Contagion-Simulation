package main

import "fmt"

type Agent struct {
	x, y     int
	infected bool
}

type Game struct {
	agents [][]Agent
	l, w   int
}

func (g Game) step() {}

func main() {

	fmt.Println("Hello World")
	game := Game{}

	for range 100 {
		game.step()
	}
}
