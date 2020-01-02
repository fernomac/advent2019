package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"

	"github.com/fernomac/advent2019/lib"
)

type uint128 struct {
	hi, lo uint64
}

func mul(a, b uint64) uint128 {
	hi, lo := bits.Mul64(a, b)
	return uint128{hi, lo}
}

func add(a uint128, b uint64) uint128 {
	lo, carry := bits.Add64(a.lo, b, 0)
	hi := a.hi + carry
	return uint128{hi, lo}
}

func mod(a uint128, b uint64) uint64 {
	_, rem := bits.Div64(a.hi, a.lo, b)
	return rem
}

// r := a^b (mod c)
func modexp(a uint64, b uint64, c uint64) uint64 {
	bit := uint64(1)
	pow := a
	out := uint64(1)

	for bit <= b {
		if b&bit != 0 {
			out = mod(mul(out, pow), c)
		}
		pow = mod(mul(pow, pow), c)
		bit <<= 1
	}

	return out
}

// r1*a + r2*b = r0
func gcd(a, b uint64) (uint64, int64, int64) {
	if a == 0 {
		return b, 0, 1
	}

	g, x, y := gcd(b%a, a)
	return g, y - (int64(b/a) * x), x
}

// a*r % b == 1
func modinv(a, b uint64) uint64 {
	g, x, _ := gcd(a, b)
	if g != 1 {
		panic("a and b are not coprime")
	}
	if x < 0 {
		return b - uint64(-x)
	}
	return uint64(x)
}

// b*r % c == a
func moddiv(a, b, c uint64) uint64 {
	inv := modinv(b, c)
	return mod(mul(a, inv), c)
}

func linearize(ops []string, len uint64) (uint64, uint64) {
	aa := uint64(1)
	bb := uint64(0)

	for _, op := range ops {
		var a, b uint64
		switch {
		case op == "deal into new stack":
			// y = -1*x + (len-1) (mod len)
			a = len - 1
			b = len - 1

		case strings.HasPrefix(op, "cut "):
			n, err := strconv.ParseInt(op[4:], 10, 64)
			if err != nil {
				panic(err)
			}

			if n < 0 {
				n += int64(len)
			}

			// y = x + (len-n) (mod len)
			a = 1
			b = len - uint64(n)

		case strings.HasPrefix(op, "deal with increment "):
			n, err := strconv.ParseUint(op[20:], 10, 64)
			if err != nil {
				panic(err)
			}

			// y = n*x + 0 (mod len)
			a = n
			b = 0

		default:
			panic(fmt.Sprintf("unrecognized operation '%v'", op))
		}

		// y = ax+b
		// y = a(aa*x+bb) + b
		// y = (a*aa)x + (a*bb + b)
		aa = mod(mul(a, aa), len)
		bb = mod(add(mul(a, bb), b), len)
	}

	return aa, bb
}

func apply(a, b, x, n uint64) uint64 {
	// y = ax + b (mod n)
	return mod(add(mul(a, x), b), n)
}

func iterate(a, b, t, n uint64) (uint64, uint64) {
	// t=1: y=ax + b
	// t=2: y=a*(ax+b)+b=a^2x + ab + b
	// t=3: y=a*(a^2x+ab+b)+b = a^3x + a^2b + ab + b
	// ...
	// t=n: a=a^nx + b*((a^n - 1) / (a - 1)))

	aa := modexp(a, t, n)
	bb := mod(mul(moddiv(aa-1, a-1, n), b), n)

	return aa, bb
}

func unapply(a, b, y, n uint64) uint64 {
	// y = ax + b
	// ax = y - b
	// x = (y-b) / a

	return moddiv((n+y)-b, a, n)
}

func main() {
	in := strings.TrimSpace(lib.Load("input.txt"))
	ops := strings.Split(in, "\n")

	{
		cards := uint64(10007)

		a, b := linearize(ops, cards)
		y := apply(a, b, 2019, cards)

		fmt.Println(y)
	}

	{
		cards := uint64(119315717514047)
		shuffles := uint64(101741582076661)

		a, b := linearize(ops, cards)
		a, b = iterate(a, b, shuffles, cards)

		x := unapply(a, b, 2020, cards)
		fmt.Println(x)
	}
}
