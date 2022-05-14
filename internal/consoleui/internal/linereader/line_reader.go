package linereader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadLine() (string, error) { return r.readLine() }

var r = newLineReader(os.Stdin)

type lineReader struct {
	s *bufio.Scanner
}

func newLineReader(r io.Reader) *lineReader {
	return &lineReader{
		s: bufio.NewScanner(r),
	}
}

func (r *lineReader) readLine() (string, error) {
	if ok := r.s.Scan(); !ok {
		err := r.s.Err()
		if err == nil {
			return "", io.EOF
		}
		return "", err
	}

	return r.s.Text(), nil
}

func ErrMsgf(pattern string, args ...interface{}) error {
	fmt.Printf(pattern, args...)
	fmt.Printf("Press ENTER to continue\n")
	if _, err := ReadLine(); err != nil {
		return fmt.Errorf("readline error: %w", err)
	}
	return nil
}
