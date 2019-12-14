package main

import (
	"strconv"
	"testing"
)

func TestTest(t *testing.T) {
	test := func(n int, ok bool) {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			num := newnum(n)
			if num.test1() != ok {
				t.Errorf("test(%v) should have returned %v", n, ok)
			}
		})
	}

	test(111111, true)
	test(223450, false)
	test(123789, false)
}
