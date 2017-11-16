package lib

import (
	"bytes"
	"errors"
	"fmt"
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
		err = removeLines(filename, lineNum+1, 1)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil
	}

	return true, nil
}

func removeLines(fn string, start, n int) (err error) {
	if start < 1 {
		return errors.New("invalid request.  line numbers start at 1.")
	}
	if n < 0 {
		return errors.New("invalid request.  negative number to remove.")
	}
	var f *os.File
	if f, err = os.OpenFile(fn, os.O_RDWR, 0); err != nil {
		return
	}
	defer func() {
		if cErr := f.Close(); err == nil {
			err = cErr
		}
	}()
	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}
	cut, ok := skip(b, start-1)
	if !ok {
		return fmt.Errorf("less than %d lines", start)
	}
	if n == 0 {
		return nil
	}
	tail, ok := skip(cut, n)
	if !ok {
		return fmt.Errorf("less than %d lines after line %d", n, start)
	}
	t := int64(len(b) - len(cut))
	if err = f.Truncate(t); err != nil {
		return
	}
	if len(tail) > 0 {
		_, err = f.WriteAt(tail, t)
	}
	return
}

func skip(b []byte, n int) ([]byte, bool) {
	for ; n > 0; n-- {
		if len(b) == 0 {
			return nil, false
		}
		x := bytes.IndexByte(b, '\n')
		if x < 0 {
			x = len(b)
		} else {
			x++
		}
		b = b[x:]
	}
	return b, true
}
