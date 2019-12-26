package lib

// Abs returns the absolute value of a.
func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// Gcd computes the greatest common divisor of a and b.
func Gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	return gcd(a, b)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// Lcm computes the lowest common multiple of a set of integers.
func Lcm(as ...int) int {
	if len(as) == 0 {
		return 0
	}
	a := as[0]

	for i := 1; i < len(as); i++ {
		gcd := Gcd(a, as[i])
		a = a * (as[i] / gcd)
	}

	return a
}
