package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type layer struct {
	bits []uint8
}

func newlayer(str string) *layer {
	bits := make([]uint8, 25*6)
	for i := 0; i < 25*6; i++ {
		bits[i] = str[i] - '0'
	}
	return &layer{bits}
}

func (l *layer) get(x, y int) uint8 {
	return l.bits[(y*25)+x]
}

func (l *layer) checksum() (int, int, int) {
	zeroes := 0
	ones := 0
	twos := 0

	for _, b := range l.bits {
		switch b {
		case 0:
			zeroes++
		case 1:
			ones++
		case 2:
			twos++
		}
	}

	return zeroes, ones, twos
}

type image struct {
	layers []*layer
}

func newimage(str string) *image {
	layers := []*layer{}
	for len(str) > 0 {
		layers = append(layers, newlayer(str))
		str = str[25*6:]
	}
	return &image{layers}
}

func (i *image) get(x, y int) uint8 {
	for _, l := range i.layers {
		val := l.get(x, y)
		if val == 0 || val == 1 {
			return val
		}
	}
	return 2
}

func (i *image) checksum() int {
	min := math.MaxInt32
	score := -1

	for _, l := range i.layers {
		zeroes, ones, twos := l.checksum()
		if zeroes < min {
			min = zeroes
			score = ones * twos
		}
	}

	return score
}

func main() {
	str := strings.TrimSpace(lib.Load("input.txt"))
	img := newimage(str)

	fmt.Println("checksum:", img.checksum())

	for y := 0; y < 6; y++ {
		for x := 0; x < 25; x++ {
			switch img.get(x, y) {
			case 0:
				fmt.Print(" ")
			case 1:
				fmt.Print("X")
			case 2:
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
