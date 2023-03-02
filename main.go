package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// comment - what does this const
const (
	NOTHING = 0
	WALL    = 1
	PLAYER  = 69

	MAX_SAMPLES = 100
)

// comment - what does this struct
type input struct {
	pressedKey byte
}

// comment - what does this func
func (i *input) update() {
	b := make([]byte, 1)
	os.Stdin.Read(b) // blocking until stdin has stuff in buffer
	i.pressedKey = b[0]
}

// comment - what does this struct
type position struct {
	x, y int
}

// comment - what does this struct
type player struct {
	pos   position
	level *level
	input *input

	reverse bool
}

// comment - what does this func
func (p *player) update() {
	if p.reverse {
		p.pos.y -= 1
		if p.pos.y == 2 {
			p.pos.y += 1
			p.reverse = false
		}
		return
	}

	p.pos.x += 1
	if p.pos.x == p.level.width-2 {
		p.pos.x -= 1
		p.reverse = true
	}
}

// comment - what does this struct
type stats struct {
	start  time.Time
	frames int
	fps    float64
}

// comment - what does this func
func newStats() *stats {
	return &stats{
		fps:   69,
		start: time.Now(),
	}
}

// comment - what does this func
func (s *stats) update() {
	s.frames++
	if s.frames == MAX_SAMPLES {
		s.fps = float64(s.frames) / time.Since(s.start).Seconds()
		s.frames = 0
		s.start = time.Now()
	}
}

// comment - what does this struct
type level struct {
	width, height int
	data          [][]int
}

// comment - what does this func
func newLevel(width, height int) *level {
	data := make([][]int, height)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]int, width)
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

// comment - what does this func
func (l *level) set(pos position, v int) {
	l.data[pos.y][pos.x] = v
}

// comment - what does this struct
type game struct {
	isRunning bool
	level     *level
	stats     *stats
	player    *player
	input     *input

	drawBuf *bytes.Buffer
}

// comment - what does this func
func newGame(width, height int) *game {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	var (
		lvl  = newLevel(width, height)
		inpu = &input{}
	)
	return &game{
		level:   lvl,
		drawBuf: new(bytes.Buffer),
		stats:   newStats(),
		input:   inpu,
		player: &player{
			input: inpu,
			level: lvl,
			pos:   position{x: 2, y: 5},
		},
	}
}

// comment - what does this func
func (g *game) start() {
	g.isRunning = true
	g.loop()
}

// comment - what does this func
func (g *game) loop() {
	for g.isRunning {
		g.input.update()
		g.update()
		g.render()
		g.stats.update()

		time.Sleep(time.Millisecond * 16) // limit FPS
	}
}

// comment - what does this func
func (g *game) update() {
	g.level.set(g.player.pos, NOTHING)
	g.player.update()
	g.level.set(g.player.pos, PLAYER)
}

// comment - what does this func
func (g *game) renderLevel() {
	for h := 0; h < g.level.height; h++ {
		for w := 0; w < g.level.width; w++ {
			if g.level.data[h][w] == NOTHING {
				g.drawBuf.WriteString(" ")
			}
			if g.level.data[h][w] == WALL {
				g.drawBuf.WriteString("▢")
			}
			if g.level.data[h][w] == PLAYER {
				g.drawBuf.WriteString("♿")
			}
		}
		g.drawBuf.WriteString("\n")
	}
}

// comment - what does this func
func (g *game) render() {
	g.drawBuf.Reset()
	fmt.Fprint(os.Stdout, "\033[2j\033[1;1H")

	g.renderLevel()
	g.renderStats()
	fmt.Fprint(os.Stdout, g.drawBuf.String())
}

// comment - what does this func
func (g *game) renderStats() {
	g.drawBuf.WriteString("--STATS\n")
	g.drawBuf.WriteString(fmt.Sprintf("FPS: %.2f\n", g.stats.fps))
}

// comment - what does this func
func main() {
	width := 80
	height := 18
	g := newGame(width, height)
	g.start()
}
