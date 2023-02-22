package main

import (
	"bytes"
	"fmt"
	"os"
)

const (
	NOTHING = 0
	WALL    = 1
	PLAYER  = 69
)

type level struct {
	width, height int
	data          [][]byte
}

func newLevel(width, height int) *level {
	data := make([][]byte, height)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]byte, width)
		}
	}
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			if w == 0 {
				data[h][w] = WALL
			}
			if h == 0 {
				data[h][w] = WALL
			}
			if w == width-1 {
				data[h][w] = WALL
			}
			if h == height-1 {
				data[h][w] = WALL
			}
		}
	}
	return &level{
		width:  width,
		height: height,
		data:   data,
	}
}

func (l *level) x() {}

type game struct {
	isRunning bool
	level     *level

	drawBuf *bytes.Buffer
}

func newGame(width, height int) *game {
	lvl := newLevel(width, height)
	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
	}
}

func (g *game) start() {
	g.isRunning = true
	g.loop()
}

func (g *game) loop() {
	for g.isRunning {
		g.update()
		g.render()

	}
}

func (g *game) update() {}

func (g *game) renderLevel() {
	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			if g.level.data[h][w] == NOTHING {
				g.drawBuf.WriteString(" ")
			}
			if g.level.data[h][w] == WALL {
				g.drawBuf.WriteString("▢")
			}
		}
		g.drawBuf.WriteString("\n")
	}
}

func (g *game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[2j\033[1;1H")
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}

func main() {
	width := 80
	height := 18
	g := newGame(width, height)
	g.start()
}
