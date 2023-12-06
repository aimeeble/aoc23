package lib

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type TextIterator interface {
	NextLine() (string, error)
	Close()
}

type bufioIterator struct {
	ogr io.ReadCloser
	r   *bufio.Reader
}

func (bii *bufioIterator) NextLine() (string, error) {
	data, err := bii.r.ReadBytes('\n')
	if err != nil {
		if len(data) != 0 {
			return string(data), nil
		}
		return "", err
	}
	return string(data), nil
}

func (bii *bufioIterator) Close() {
	bii.ogr.Close()
}

func GetInputIterator(r io.ReadCloser) (TextIterator, error) {
	return &bufioIterator{r, bufio.NewReader(r)}, nil
}

func GetInputFileIterator(filename string) (TextIterator, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return GetInputIterator(fh)
}

func GetInputFileAll(filename string) ([]string, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return GetInputAll(fh)
}

func GetInputAll(r io.ReadCloser) ([]string, error) {
	buf, err := GetInputIterator(r)
	if err != nil {
		return nil, err
	}
	var lines []string
	for {
		line, err := buf.NextLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return lines, nil
			} else {
				return lines, err
			}
		}
		lines = append(lines, line)
	}
}
