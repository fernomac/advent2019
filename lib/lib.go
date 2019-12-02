package lib

import (
	"bufio"
	"io/ioutil"
	"os"
)

// Load loads a file.
func Load(filename string) string {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

// ForLines calls f for each line in a file.
func ForLines(filename string, f func(str string)) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		f(s.Text())
	}

	if err := s.Err(); err != nil {
		panic(err)
	}
}
