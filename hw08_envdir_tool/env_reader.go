package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrUnsupportedSymbolsInFilename = errors.New("unsupported symbols in files name")
	ErrUnsupportedInputPath         = errors.New("unsupported input path")
)

var re = regexp.MustCompile(`\W`)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if len(dir) == 0 {
		return nil, ErrUnsupportedInputPath
	}
	if i, err := os.Stat(dir); err != nil {
		return nil, err
	} else if !i.IsDir() {
		return nil, ErrUnsupportedInputPath
	}

	envs := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range files {
		if dirEntry.IsDir() || !dirEntry.Type().IsRegular() {
			continue
		}
		entryName := dirEntry.Name()
		if re.MatchString(entryName) {
			return nil, ErrUnsupportedSymbolsInFilename
		}
		if i, err := dirEntry.Info(); err != nil {
			return nil, err
		} else if i.Size() == 0 {
			envs[entryName] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}
		file, err := os.Open(dir + "/" + entryName)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			file.Close()
			return nil, err
		}
		val := scanner.Text()
		s := strings.TrimRight(val, " \n\t\r\v")
		s = strings.ReplaceAll(s, "\x00", "\n")
		file.Close()
		envs[entryName] = EnvValue{
			Value: s,
		}
	}
	return envs, nil
}
