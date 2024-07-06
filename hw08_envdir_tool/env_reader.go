package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// Errors.
var errInvalidFileName = errors.New("invalid file name")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Read files list from dir.
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Create environment variables map.
	envVarsMap := make(Environment)

	// Create Environment map.
	for _, file := range files {
		// Get file name.
		fileName := file.Name()

		// Check file name for "=" symbol.
		if strings.Contains(fileName, "=") {
			return nil, errInvalidFileName
		}

		// Get file path.
		filePath := filepath.Join(dir, fileName)

		// Check if file is empty.
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		// Create EnvValue structur
		var envVarValue EnvValue

		// If file empty, add Environment structure with NeedRemove to map.
		if fileInfo.Size() == 0 {
			envVarValue = EnvValue{NeedRemove: true}
		} else {
			// If file not empty, read first line from file.
			// Open file for reading.
			readFile, err := os.Open(filePath)
			if err != nil {
				return nil, err
			}
			defer readFile.Close()

			// Create file scanner.
			fileScanner := bufio.NewScanner(readFile)

			// Split file into lines.
			fileScanner.Split(bufio.ScanLines)

			// Get first line.
			fileScanner.Scan()
			firstLine := fileScanner.Text()

			// Trim right.
			firstLine = strings.TrimRightFunc(firstLine, unicode.IsSpace)

			// Replace terminal nulls (0x00) to end of line (\n).
			firstLine = string(bytes.ReplaceAll([]byte(firstLine), []byte("\x00"), []byte("\n")))

			// Create Environment structure.
			envVarValue = EnvValue{Value: firstLine}
		}

		// Add Environment structure to map.
		envVarsMap[fileName] = envVarValue
	}

	return envVarsMap, nil
}
