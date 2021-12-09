package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

const (
	start  = 4
	revive = 3
	starve = 2
)

var size = flag.Int("size", 11, "board size")

func main() {
	flag.Parse()

	g := newGame(*size)
	g.board.draw()
	fmt.Printf("%s's turn ", g.curPlayer)
	g.curPlayer.Draw()
	fmt.Print(" > ")

	for {
		var r, c int
		_, err := fmt.Scanf("%d %d\n", &r, &c)

		for err != nil {
			fmt.Println(err)
			fmt.Print("Try again: ")
			_, err = fmt.Scanf("%d %d\n", &r, &c)
		}

		w, err := g.next(r, c)
		if err != nil {
			fmt.Println(err)
		}
		g.board.draw()

		if w != none && g.round >= start {
			fmt.Printf("%s wins!\n", w)
			break
		} else {
			fmt.Printf("%s's turn ", g.curPlayer)
			g.curPlayer.Draw()
			fmt.Print(" > ")
		}
	}
}

type game struct {
	round     int
	curPlayer player
	board     board
}

func newGame(size int) *game {
	return &game{
		round:     0,
		curPlayer: p1,
		board:     newBoard(size),
	}
}

func (g *game) next(r, c int) (winner player, err error) {
	err = g.board.set(r, c, g.curPlayer)
	if err != nil {
		return
	}

	if g.round >= start {
		g.board = g.board.next()
	}

	if g.curPlayer == p1 {
		g.curPlayer = p2
	} else {
		g.curPlayer = p1
	}

	g.round++

	return g.board.winner(), nil
}

type player int

const (
	none player = iota
	p1
	p2
)

func (p player) String() string {
	if p == p1 {
		return "Player 1"
	} else if p == p2 {
		return "Player 2"
	}
	return "None"
}

func (p player) Draw() {
	if p == p1 {
		fmt.Print("O")
	} else if p == p2 {
		fmt.Print("x")
	} else {
		fmt.Print(" ")
	}
}

type board [][]player

func newBoard(size int) board {
	b := make([][]player, 0, size)
	for i := 0; i < size; i++ {
		b = append(b, make([]player, size))
	}
	return b
}

func (b board) set(r, c int, player player) error {
	if r < 0 || c < 0 || r >= len(b) || c >= len(b[0]) {
		return errors.New("Out of bounds!")
	}

	if b[r][c] != 0 {
		return errors.New("Occupied!")
	}

	b[r][c] = player

	return nil
}

func (b board) get(r, c int) player {
	if r < 0 || r >= len(b) || c < 0 || c >= len(b[0]) {
		return none
	}

	return b[r][c]
}

func (b board) next() board {
	nb := newBoard(len(b))

	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[0]); j++ {
			neighbors := make(map[player]int)
			neighbors[b.get(i-1, j)]++
			neighbors[b.get(i+1, j)]++
			neighbors[b.get(i-1, j-1)]++
			neighbors[b.get(i+1, j-1)]++
			neighbors[b.get(i-1, j+1)]++
			neighbors[b.get(i+1, j+1)]++
			neighbors[b.get(i, j-1)]++
			neighbors[b.get(i, j+1)]++
			nb[i][j] = nextState(b[i][j], neighbors)
		}
	}

	return nb
}

func (b board) draw() {
	fmt.Println()
	fmt.Print("    ")
	for i := 0; i < len(b[0]); i++ {
		fmt.Printf("%2d", i%10)
	}
	fmt.Printf("\n   +-%s+\n", strings.Repeat("--", len(b[0])))

	for i := 0; i < len(b); i++ {
		fmt.Printf("%3d| ", i)
		for j := 0; j < len(b[0]); j++ {
			b[i][j].Draw()
			fmt.Print(" ")
		}
		fmt.Println("|")
	}

	fmt.Printf("   +-%s+\n", strings.Repeat("--", len(b[0])))
}

func nextState(self player, neighbors map[player]int) player {
	p1Nbs, p2Nbs := neighbors[p1], neighbors[p2]
	if self == none {
		if p1Nbs+p2Nbs == revive {
			if p1Nbs > p2Nbs {
				return p1
			} else {
				return p2
			}
		}
		return none
	}

	if p1Nbs+p2Nbs < starve {
		return none
	}

	if p1Nbs+p2Nbs > revive {
		return none
	}

	return self
}

func (b board) winner() player {
	ps := make(map[player]int)
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[0]); j++ {
			if p := b[i][j]; p != 0 {
				ps[p]++
			}
		}
	}
	// if all pieces are eliminated, then p2 wins since p1 moves first
	if ps[p1] == 0 && ps[p2] == 0 {
		return p2
	}
	if ps[p1] == 0 && ps[p2] != 0 {
		return p2
	}
	if ps[p2] == 0 && ps[p1] != 0 {
		return p1
	}
	return none
}
