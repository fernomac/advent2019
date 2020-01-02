package main

import (
	"fmt"
	"math"
	"testing"
)

func TestMul(t *testing.T) {
	res := mul(1<<63, 2)
	if res.hi != 1 {
		t.Errorf("expected 1, got %v", res.hi)
	}
	if res.lo != 0 {
		t.Errorf("expected 0, got %v", res.lo)
	}
}

func TestAdd(t *testing.T) {
	res := add(uint128{1, math.MaxUint64}, 1)
	if res.hi != 2 {
		t.Errorf("expected 2, got %v", res.hi)
	}
	if res.lo != 0 {
		t.Errorf("expected 0, got %v", res.lo)
	}
}

func TestModExp(t *testing.T) {
	if v := modexp(10, 3, 37); v != 1000%37 {
		t.Errorf("expected %v, got %v", 1000%37, v)
	}
	if v := modexp(2, 67, 3); v != 2 {
		t.Errorf("expected %v, got %v", 2, v)
	}
	if v := modexp(8067652995194, 101741582076661, 119315717514047); v != 71126595274540 {
		t.Errorf("expected %v, got %v", 71126595274540, v)
	}
}

func TestGCD(t *testing.T) {
	a := uint64(123456 * 119315717514047)
	b := uint64(98765 * 119315717514047)

	g, x, y := gcd(a, b)

	g2 := int64(a)*x + int64(b)*y
	if int64(g) != g2 {
		t.Errorf("expected %v, got %v", g, g2)
	}
}

func TestModInv(t *testing.T) {
	test := func(a, b, e uint64) {
		t.Run(fmt.Sprintf("(%v,%v)", a, b), func(t *testing.T) {
			v := modinv(a, b)
			if v != e {
				t.Errorf("expected %v, got %v", e, v)
			}
		})
	}

	test(1, 11, 1)   // 1*1 = 1 (mod 11)
	test(2, 11, 6)   // 2*6 = 12 = 1 (mod 11)
	test(3, 11, 4)   // 3*4 = 12 = 1 (mod 11)
	test(4, 11, 3)   // 4*3
	test(5, 11, 9)   // 5*9 = 45 = 1 (mod 11)
	test(6, 11, 2)   // 6*2
	test(7, 11, 8)   // 7*8 = 56 = 1 (mod 11)
	test(8, 11, 7)   // 8*7
	test(9, 11, 5)   // 9*5
	test(10, 11, 10) // 10*10 = 100 = 1 (mod 99)

	test(123, 119315717514047, 47532277708848)
}

func TestModDiv(t *testing.T) {
	test := func(a, b, c, e uint64) {
		t.Run(fmt.Sprintf("(%v,%v,%v)", a, b, c), func(t *testing.T) {
			if v := moddiv(a, b, c); v != e {
				t.Errorf("expected %v, got %v", e, v)
			}
		})
	}

	test(2, 2, 11, 1) // 2*1 = 2 (mod 11)
	test(3, 2, 11, 7) // 2*7 = 14 = 3 (mod 11)
	test(4, 2, 11, 2) // 2*2 = 4 (mod 11)
	test(5, 2, 11, 8) // 2*8 = 16 = 5 (mod 11)

	test(7, 3, 119315717514047, 39771905838018)
}

func TestLinearize(t *testing.T) {
	a, b := linearize([]string{
		"deal into new stack",
		"cut -2",
		"deal with increment 7",
		"cut 8",
		"cut -4",
		"deal with increment 7",
		"cut 3",
		"deal with increment 9",
		"deal with increment 3",
		"cut -1",
	}, 10)

	test := func(x, e uint64) {
		y := apply(a, b, x, 10)
		if y != e {
			t.Errorf("%v: expected %v, got %v", x, e, y)
		}
	}

	for i, x := range []uint64{9, 2, 5, 8, 1, 4, 7, 0, 3, 6} {
		test(x, uint64(i))
	}
}
