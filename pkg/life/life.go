package life

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
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
	cnt := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx, ny := x+dx, y+dy
			if nx < w.Width && nx >= 0 && ny < w.Height && ny >= 0 {
				if w.Cells[ny][nx] {
					cnt += 1
				}
			}
		}
	}

	return cnt
}

// Вычисление количества соседей, если поле сворачивается в тор
func (w *World) DonutNeighbors(x, y int) int {
	cnt := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx := (x + dx + w.Width) % w.Width
			ny := (y + dy + w.Height) % w.Height

			if w.Cells[ny][nx] {
				cnt++
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

func (w *World) SaveState(filename string) error {
	if filename == "" {
		return errors.New("empty filename")
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Failed to open file: %v", err)
		return err
	}
	defer func() {
		f.Close()
	}()

	value := ""
	for _, row := range w.Cells {
		for _, cell := range row {
			if cell {
				value = "1"
			} else {
				value = "0"
			}

			_, err = f.WriteString(value)
			if err != nil {
				fmt.Printf("Failed to write string in file: %v", err)
				return err
			}
		}

		f.WriteString("\n")
	}

	return nil
}

func (w *World) LoadState(filename string) error {
	if filename == "" {
		return errors.New("empty filename")
	}

	f, err := os.Open(filename)
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}()

	if err != nil {
		return errors.New("failed to open file")
	}

	var lines []string
	var maxWidth int

	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
		if len(fileScanner.Text()) > maxWidth {
			maxWidth = len(fileScanner.Text())
		}
	}

	height := len(lines)

	if maxWidth == 0 || height == 0 {
		return errors.New("empty file")
	}

	for _, line := range lines {
		if len(line) != maxWidth {
			return errors.New("different lenght of rows")
		}
	}

	w.Width = maxWidth
	w.Height = height
	w.Cells = make([][]bool, height)

	for y, line := range lines {
		w.Cells[y] = make([]bool, maxWidth)
		for x, cell := range line {
			if cell == '0' {
				w.Cells[y][x] = false
			} else if cell == '1' {
				w.Cells[y][x] = true
			} else {
				return errors.New("invalid symbols")
			}
		}
	}

	return nil
}
