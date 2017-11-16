package lib

import (
	"io/ioutil"
	"os"
	"strings"
)

func ReadLines(filename string) (map[int]string, error) {
	lines := make(map[int]string)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var n = 0
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			lines[n] = line
			n++
		}
	}

	return lines, nil
}

func AddLine(filename string, lineToAdd string) (bool, error) {
	lines, err := ReadLines(filename)
	if err != nil {
		return false, err
	}

	// Check to make sure we don't add the same file twice
	// TODO: This should be expanded to full path
	for _, line := range lines {
		if line == lineToAdd {
			return false, nil
		}
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return false, err
	}

	defer f.Close()

	if _, err = f.WriteString(lineToAdd + "\n"); err != nil {
		return false, err
	}

	return true, nil
}

func RemoveLine(filename string, lineToRemove string) (bool, error) {
	lines, err := ReadLines(filename)
	if err != nil {
		return false, err
	}

	var lineNum = -1

	for num, line := range lines {
		if line == lineToRemove {
			lineNum = num
		}
	}

	if lineNum > -1 {
		err = RemoveLines(filename, lineNum+1, 1)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil
	}

	return true, nil
}
