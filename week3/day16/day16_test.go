package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fernomac/advent2019/lib"
)

func TestSpanerator(t *testing.T) {
	spaneq := func(a, b []span) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	test := func(n int, eval ...span) {
		t.Run(fmt.Sprintf("%v", n), func(t *testing.T) {
			p := spanerate(n)
			val := make([]span, len(eval))

			for i := 0; i < len(eval); i++ {
				val[i] = p.next()
			}

			if !spaneq(val, eval) {
				t.Errorf("expected %v, got %v", eval, val)
			}
		})
	}

	test(1, span{0, 1, 1}, span{2, 1, -1}, span{4, 1, 1}, span{6, 1, -1})
	test(2, span{1, 2, 1}, span{5, 2, -1}, span{9, 2, 1}, span{13, 2, -1})
	test(3, span{2, 3, 1}, span{8, 3, -1}, span{14, 3, 1}, span{20, 3, -1})
	test(12345, span{12344, 12345, 1}, span{37034, 12345, -1})
}

func TestInput(t *testing.T) {
	test := func(in string, reps int, eval string) {
		t.Run(fmt.Sprintf("(%v,%v)", in, reps), func(t *testing.T) {
			src := parse(in, reps)
			bld := strings.Builder{}
			for i := 0; i < src.len(); i++ {
				bld.WriteByte(byte(src.at(i) + '0'))
			}

			val := bld.String()
			if val != eval {
				t.Fatalf("expected %v, got %v", eval, val)
			}
		})
	}

	test("123", 3, "123123123")
}

func TestSummer(t *testing.T) {
	src := parse("1", blocksize*3)
	summer := &summer{
		src:   src,
		cache: map[int]int{},
	}

	test := func(off, len, eval int) {
		t.Run(fmt.Sprintf("(%v,%v)", off, len), func(t *testing.T) {
			val := summer.sum(off, len)
			if val != eval {
				t.Fatalf("expected %v, got %v", eval, val)
			}
		})
	}

	test(blocksize-1, 1, 1)
	test(blocksize, 1, 1)
	test(blocksize+1, 1, 1)

	test(blocksize-1, 2, 2)

	test(blocksize, blocksize, blocksize)
	test(blocksize-1, blocksize, blocksize)
	test(blocksize+1, blocksize, blocksize)

	test(blocksize-1, blocksize+1, blocksize+1)
	test(blocksize, blocksize+1, blocksize+1)
	test(blocksize-1, blocksize+2, blocksize+2)

	test(blocksize*2, blocksize, blocksize)
	test(blocksize*2+1, blocksize, blocksize-1)
	test(blocksize*3-1, blocksize, 1)
}

func TestFfter(t *testing.T) {
	test := func(in, eval []int) {
		t.Run(fmt.Sprintf("%v", in), func(t *testing.T) {
			src := &input{in, 1}
			fft := newffter(src)

			val := make([]int, len(in))
			for i := 0; i < len(in); i++ {
				val[i] = fft.at(i)
			}

			if !lib.IntEq(val, eval) {
				t.Fatalf("expected %v, got %v", eval, val)
			}
		})
	}

	test([]int{1, 2, 3, 4, 5, 6, 7, 8}, []int{4, 8, 2, 2, 6, 1, 5, 8})
	test([]int{4, 8, 2, 2, 6, 1, 5, 8}, []int{3, 4, 0, 4, 0, 4, 3, 8})
	test([]int{3, 4, 0, 4, 0, 4, 3, 8}, []int{0, 3, 4, 1, 5, 5, 1, 8})
	test([]int{0, 3, 4, 1, 5, 5, 1, 8}, []int{0, 1, 0, 2, 9, 4, 9, 8})
}

func TestFFT(t *testing.T) {
	test := func(in string, iters int, eval string) {
		t.Run(in, func(t *testing.T) {
			src := parse(in, 1)
			out := fft(src, iters)
			val := collect(out, 0, 8)
			if val != eval {
				t.Fatalf("expected %v, got %v", eval, val)
			}
		})
	}

	test("80871224585914546619083218645595", 100, "24176176")
	test("19617804207202209144916044189917", 100, "73745418")
	test("69317163492948606335995924319873", 100, "52432133")
}
