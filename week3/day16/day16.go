package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

// The phases of the pattern. Conveniently, one and negone have the right
// integer values.
type phase int

const (
	firstzero  phase = 0
	one        phase = 1
	secondzero phase = 2
	negone     phase = -1
)

func (p phase) next() phase {
	if p == secondzero {
		return negone
	}
	return p + 1
}

// A span of inputs from the level below to be summed together
// and multiplied by the given pattern value.
type span struct {
	off int
	len int
	mul int
}

// A spanerator generates the spans from the underlying input needed to
// produce a given value at index n. Spans with a mul of zero are skipped.
type spanerator struct {
	n int
	i int
	p phase
}

func spanerate(n int) *spanerator {
	return &spanerator{
		n: n,
		i: 0,
		p: firstzero,
	}
}

func (p *spanerator) next() span {
	p.i += p.n
	p.p = p.p.next()

	if p.p == firstzero || p.p == secondzero {
		p.i += p.n
		p.p = p.p.next()
	}

	return span{p.i - 1, p.n, int(p.p)}
}

// A source of input values, either the raw input or a lower level of fft.
type source interface {
	at(n int) int
	len() int
}

// An input is a source that directly returns input values.
type input struct {
	ns   []int
	reps int
}

var _ source = &input{}

func parse(in string, reps int) *input {
	ns := make([]int, len(in))
	for i, c := range in {
		ns[i] = int(c - '0')
	}
	return &input{ns, reps}
}

func (i *input) at(n int) int {
	x := i.ns[n%len(i.ns)]
	// fmt.Printf("in[%v] = %v\n", n, x)
	return x
}

func (i *input) len() int {
	return len(i.ns) * i.reps
}

const blocksize = 500

// A summer sums up spans of input values, caching in blocks of 500 to save
// a bit of time.
type summer struct {
	src   source
	cache map[int]int
}

func (s *summer) sum(off, len int) int {
	avail := s.src.len() - off
	if avail < len {
		len = avail
	}

	sum := 0

	mod := off % blocksize
	if mod != 0 {
		next := (off - mod) + blocksize
		span := next - off
		if span > len {
			span = len
		}

		sum += s.rawsum(off, span)
		off += span
		len -= span
	}

	for len >= blocksize {
		sum += s.blocksum(off)
		off += blocksize
		len -= blocksize
	}

	if len > 0 {
		sum += s.rawsum(off, len)
	}

	return sum
}

func (s *summer) rawsum(off, len int) int {
	sum := 0
	for i := 0; i < len; i++ {
		sum += s.src.at(off + i)
	}
	return sum
}

func (s *summer) blocksum(off int) int {
	val, ok := s.cache[off]
	if ok {
		return val
	}

	for i := 0; i < blocksize; i++ {
		val += s.src.at(off + i)
	}

	s.cache[off] = val
	return val
}

// An ffter lazily performs a round of fft.
type ffter struct {
	summer summer
	cache  map[int]int
}

var _ source = &ffter{}

func newffter(src source) *ffter {
	return &ffter{
		summer: summer{
			src:   src,
			cache: map[int]int{},
		},
		cache: map[int]int{},
	}
}

func (f *ffter) at(n int) int {
	val, ok := f.cache[n]
	if ok {
		return val
	}

	sum := 0
	spans := spanerate(n + 1)
	end := f.summer.src.len()

	for {
		span := spans.next()
		if span.off >= end {
			break
		}
		spansum := f.summer.sum(span.off, span.len)
		sum += spansum * span.mul
	}

	res := lib.Abs(sum % 10)

	f.cache[n] = res
	return res
}

func (f *ffter) len() int {
	return f.summer.src.len()
}

func fft(src source, iters int) source {
	for i := 1; i <= iters; i++ {
		src = newffter(src)
	}
	return src
}

func collect(src source, off, len int) string {
	bld := strings.Builder{}
	for i := 0; i < len; i++ {
		bld.WriteByte(byte(src.at(off+i) + '0'))
	}
	return bld.String()
}

// The offset we're interested in (5978261) is more than half of the input size (6500000),
// so fft.at(n) simplifies to just the sum from src.at(n) to src.at(src.len()). We can
// calculate that a lottttt faster than the general fft case...
func fft2(in []int) []int {
	len := len(in)

	for i := 0; i < 100; i++ {
		out := make([]int, len)
		val := 0

		for j := 1; j <= len; j++ {
			val += in[len-j]
			out[len-j] = val % 10
		}

		in = out
	}

	return in
}

func main() {
	in := strings.TrimSpace(lib.Load("input.txt"))
	{
		src := parse(in, 1)
		out := fft(src, 100)
		sum := collect(out, 0, 8)
		fmt.Println(sum)
	}

	offset, err := strconv.Atoi(in[:7])
	if err != nil {
		panic(err)
	}

	// This way is still quite slow even with optimizations, but it works.
	//
	// {
	// 	src := parse(in, 10000)
	// 	out := fft(src, 100)
	// 	msg := collect(out, offset, 8)
	// 	fmt.Println(msg)
	// }

	{
		ns := make([]int, len(in)*10000)
		for i := 0; i < 10000; i++ {
			base := i * len(in)
			for j := 0; j < len(in); j++ {
				ns[base+j] = int(in[j] - '0')
			}
		}

		out := fft2(ns[offset:])
		msg := collect(&input{out, 1}, 0, 8)

		fmt.Println(msg)
	}
}
