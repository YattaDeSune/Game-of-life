package main

import (
	"fmt"
	"math/rand"
	"time"
)

type World struct {
	Width  int
	Height int
	Cells  [][]bool
}

func NewWorld(w, h int) *World {
	cells := make([][]bool, h)
	for i := range cells {
		cells[i] = make([]bool, w)
	}

	return &World{
		Width:  w,
		Height: h,
		Cells:  cells,
	}
}

// Вычисление количества соседей
func (w *World) Neighbors(x, y int) int {
	directions := [][2]int{
		{-1, 1}, {0, 1}, {1, 1},
		{-1, 0}, {1, 0},
		{-1, -1}, {0, -1}, {1, -1},
	}

	cnt := 0
	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx < w.Width && nx >= 0 && ny < w.Height && ny >= 0 {
			if w.Cells[ny][nx] {
				cnt += 1
			}
		}
	}

	return cnt
}

// Вычисление следующего состояния клетки
func (w *World) Next(x, y int) bool {
	n := w.Neighbors(x, y)
	alive := w.Cells[y][x]
	if n < 4 && n > 1 && alive { // Если соседей двое или трое, а клетка жива,
		return true // то следующее её состояние — жива
	}
	if n == 3 && !alive { // Если клетка мертва, но у неё трое соседей,
		return true // клетка оживает
	}

	return false // В любых других случаях — клетка мертва
}

func NextState(oldWorld, newWorld *World) {
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			newWorld.Cells[i][j] = oldWorld.Next(j, i)
		}
	}
}

func (w *World) Seed() {
	for _, row := range w.Cells {
		for i := range row {
			if rand.Intn(10) < 2 {
				row[i] = true
			}
		}
	}
}

func (w *World) String() string {
	deadCell := "❌"
	aliveCell := "✅"
	str := ""

	for _, row := range w.Cells {
		for _, cell := range row {
			if cell {
				str += aliveCell
				continue
			}

			str += deadCell
		}

		str += "\n"
	}

	return str
}

func main() {
	height := 20
	width := 15

	currentWorld := NewWorld(height, width)
	nextWorld := NewWorld(height, width)

	currentWorld.Seed()
	for {
		fmt.Println(currentWorld.String())
		NextState(currentWorld, nextWorld)
		currentWorld = nextWorld

		time.Sleep(1 * time.Second)
		// Специальная последовательность для очистки экрана после каждого шага
		fmt.Print("\033[H\033[2J")
	}
}
