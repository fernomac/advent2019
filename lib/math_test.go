package lib

import "testing"

func TestGCD(t *testing.T) {
	if val := Gcd(15, 10); val != 5 {
		t.Errorf("gcd(15,10)==%v", val)
	}

	if val := Gcd(10, 15); val != 5 {
		t.Errorf("gcd(10,15)==%v", val)
	}

	if val := Gcd(-10, 15); val != 5 {
		t.Errorf("gcd(-10,15)==%v", val)
	}

	if val := Gcd(10, -15); val != 5 {
		t.Errorf("gcd(10,-15)==%v", val)
	}

	if val := Gcd(-10, -15); val != 5 {
		t.Errorf("gcd(-10,-15)==%v", val)
	}
}
