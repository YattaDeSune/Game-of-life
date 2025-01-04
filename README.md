func (w *World) LoadState(filename string) error {
	if filename == "" {
		return errors.New("Empty filename")
	}

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var lines []string
	var maxWidth int

	// Чтение файла построчно
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Проверка на одинаковую длину строк
	for _, line := range lines {
		if len(line) != maxWidth {
			return errors.New("inconsistent line lengths in file")
		}
	}

	// Установка размеров мира
	w.Width = maxWidth
	w.Height = len(lines)
	w.Cells = make([][]bool, w.Height)

	// Заполнение ячеек
	for i, line := range lines {
		w.Cells[i] = make([]bool, w.Width)
		for j, char := range line {
			if char == '1' {
				w.Cells[i][j] = true // Живая клетка
			} else if char == '0' {
				w.Cells[i][j] = false // Мертвая клетка
			} else {
				return errors.New("invalid character in file, only '0' and '1' are allowed")
			}
		}
	}

	return nil
}

Найти еще