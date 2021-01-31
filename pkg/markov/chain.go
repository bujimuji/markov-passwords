package markov

import (
	"sync"
)

// Pair represents a state transition between a set of one or more bytes and the
// next byte in the string.
type Pair struct {
	Current string
	Next    byte
}

// transition is a map that represents number of a byte occurrences
type transition map[byte]int

// Chain represents a Markov chain containing the ngram number and frequency map
type Chain struct {
	Ngram int
	Trans map[string]transition
}

// NewChain returns a new Markov Chain model with number of given ngram.
func NewChain(ngram int) *Chain {
	if ngram < 1 {
		ngram = 1
	}
	return &Chain{
		Ngram: ngram,
		Trans: make(map[string]transition),
	}
}

// Pairs adds the transitions extracted from input string to an output channel
func (c *Chain) Pairs(in <-chan string) <-chan Pair {
	out := make(chan Pair)
	go func() {
		for s := range in {
			l := len(s)
			n := c.Ngram
			// invalid string size for generating any ngrams
			if l < 2 {
				continue
			}

			// clamp ngram size

			for i := 0; i < l-n; i++ {
				out <- Pair{
					Current: s[i : i+n],
					Next:    s[i+n],
				}
			}
		}
		close(out)
	}()
	return out
}

// mergePair converts a list of channels to a single channel
func mergePair(cs ...<-chan Pair) <-chan Pair {
	var wg sync.WaitGroup
	out := make(chan Pair)

	output := func(c <-chan Pair) {
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

// Add count the number of pairs by inserting them into a frequency map
func (c *Chain) Add(cs ...<-chan Pair) {
	for p := range mergePair(cs...) {
		if c.Trans[p.Current] == nil {
			c.Trans[p.Current] = make(transition, 0)
		}

		c.Trans[p.Current][p.Next]++
	}
}
