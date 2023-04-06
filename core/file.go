package core

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Find the line where the specified text is located in a file
func Seach_line_text(path string, text string) int {
	f, _ := os.Open(path)
	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			return line
		}
		line++
	}
	return 0
}

// Add line in a specific position of the file
func Insert_str_file(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}
	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str + "\n"
		}
		fileContent += line
		fileContent += "\n"
	}
	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

// Add line at the end of the file
func Append_str_file(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

// Delete a line from a file
func Remove_line(path string, lineNumber int) {
	file, _ := ioutil.ReadFile(path)
	info, _ := os.Stat(path)
	mode := info.Mode()
	array := strings.Split(string(file), "\n")
	array = append(array[:lineNumber], array[lineNumber+1:]...)
	ioutil.WriteFile(path, []byte(strings.Join(array, "\n")), mode)
}
