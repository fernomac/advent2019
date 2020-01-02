package lib

import "strings"

// An ASCIIPuter makes it easy to do ASCII input to and output from a Puter.
type ASCIIPuter struct {
	p      *Puter
	in     chan int
	result int
}

// ASCIIfy wraps a Puter in an ASCIIPuter.
func ASCIIfy(p *Puter) *ASCIIPuter {
	in := make(chan int, 1024)
	p.StdinCh(in)
	return &ASCIIPuter{p, in, 0}
}

// Run runs the puter until it wants more input.
func (a *ASCIIPuter) Run() bool {
	return a.p.RunNB()
}

// Input feeds some ASCII input to the Puter.
func (a *ASCIIPuter) Input(in string) {
	for _, i := range in {
		a.in <- int(i)
	}
}

// Output returns the current ASCII-interpreted output of the Puter.
func (a *ASCIIPuter) Output() string {
	out := a.p.Stdout()
	a.p.DropStdout()

	b := strings.Builder{}

	for _, o := range out {
		if o >= 128 {
			a.result = o
			break
		}
		b.WriteRune(rune(o))
	}

	return b.String()
}

// Result returns the non-ascii result.
func (a *ASCIIPuter) Result() int {
	return a.result
}
