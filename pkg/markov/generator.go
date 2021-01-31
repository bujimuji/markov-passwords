package markov

import (
	"os"
	"strings"
	"sync"

	"github.com/valyala/fastrand"
)

// nextByte is the struct version of transition map
type nextByte struct {
	b    byte
	freq int
}

// Generator is another representation of Markov chain to generate markov passwords
type Generator struct {
	Ngram int
	trans map[string][]nextByte
	seeds []string
}

// NewGenerator returns a Markov password generator from the given chain
func NewGenerator(c *Chain) *Generator {
	g := &Generator{
		Ngram: c.Ngram,
		trans: make(map[string][]nextByte),
		seeds: make([]string, len(c.Trans)),
	}
	// map to slice
	i := 0
	for seed, trans := range c.Trans {
		g.trans[seed] = make([]nextByte, len(trans))
		j := 0
		for b, freq := range trans {
			g.trans[seed][j] = nextByte{b, freq}
			j++
		}
		g.seeds[i] = seed
		i++
	}
	return g
}

// Sum returns the sum of frequencies inside a transition array
func Sum(s []nextByte) int {
	var sum int
	for i := 0; i < len(s); i++ {
		sum += s[i].freq
	}
	return sum
}

// genPassword generates a new password and puts the password in a channel
func (g *Generator) genPassword(n uint32, sc <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		for range sc {
			var b byte
			var password strings.Builder
			randi := fastrand.Uint32n(n)
			seed := g.seeds[randi]
			password.WriteString(seed)
			for b != '\n' {
				trans := g.trans[seed]
				total := Sum(trans)
				rand := int(fastrand.Uint32n(uint32(total)))

				// iterate all possibilities until the rand runs out
				for _, next := range trans {
					rand -= next.freq
					if rand < 0 {
						b = next.b
						break
					}
				}

				// add new byte to password and update seed for next byte
				password.WriteByte(b)
				l := password.Len()
				seed = password.String()
				seed = seed[l-g.Ngram : l]
			}
			out <- password.String()
		}
		close(out)
	}()
	return out
}

// genTask generates a task for each password to be generated
func (g *Generator) genTask(n uint32, maxAttempts uint64) <-chan int {
	out := make(chan int)
	go func() {
		var i uint64
		for i = 0; i < maxAttempts; i++ {
			out <- int(i)
		}
		close(out)
	}()
	return out
}

// mergeStr converts a list of channels to a single channel
func mergeStr(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for p := range c {
			out <- p
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Generate generates new passwords based on trained Markov model
// and prints passwords in stdout
func (g *Generator) Generate(chans int, maxAttempts uint64) {
	sc := g.genTask(uint32(len(g.seeds)), maxAttempts)
	pc := make([]<-chan string, chans)
	for i := 0; i < chans; i++ {
		pc[i] = g.genPassword(uint32(len(g.seeds)), sc)
	}

	var i int
	var l = 1 << 14
	var buf = make([]byte, l)
	for pass := range mergeStr(pc...) {
		if l-i < len(pass) {
			_, _ = os.Stdout.Write(buf[:i])
			i = 0
		}
		i += copy(buf[i:], pass)
	}
	if i != 0 {
		_, _ = os.Stdout.Write(buf[:i])
	}
}
