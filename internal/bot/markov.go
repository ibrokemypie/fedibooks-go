// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Generating random text: a Markov chain algorithm

Based on the program presented in the "Design and Implementation" chapter
of The Practice of Programming (Kernighan and Pike, Addison-Wesley 1999).
See also Computer Recreations, Scientific American 260, 122 - 125 (1989).

A Markov chain algorithm generates text by creating a statistical model of
potential textual suffixes for a given prefix. Consider this text:

	I am not a number! I am a free man!

Our Markov chain algorithm would arrange this text into this set of prefixes
and suffixes, or "chain": (This table assumes a prefix length of two words.)

	Prefix       Suffix

	"" ""        I
	"" I         am
	I am         a
	I am         not
	a free       man!
	am a         free
	am not       a
	a number!    I
	number! I    am
	not a        number!

To generate text using this table we select an initial prefix ("I am", for
example), choose one of the suffixes associated with that prefix at random
with probability determined by the input statistics ("a"),
and then create a new prefix by removing the first word from the prefix
and appending the suffix (making the new prefix is "am a"). Repeat this process
until we can't find any suffixes for the current prefix or we exceed the word
limit. (The word limit is necessary as the chain table may contain cycles.)

Our version of this program reads text from standard input, parsing it into a
Markov chain, and writes generated text to standard output.
The prefix and output lengths can be specified using the -prefix and -words
flags on the command-line.
*/
package bot

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strings"

	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
	sentences [][]string
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen, [][]string{}}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		c.sentences = append(c.sentences, strings.Fields(s))
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(maxWords, minWords int, maxOverlapRatio float64, maxOverLapTotal, tries int) string {
	if tries == 0 {
		tries = 10
	}
	if maxOverlapRatio == 0 {
		maxOverlapRatio = 0.7
	}
	if maxOverLapTotal == 0 {
		maxOverLapTotal = 15
	}
	if maxOverLapTotal == 0 {
		maxOverLapTotal = 15
	}

	for i := 0; i < tries; i++ {
		p := make(Prefix, c.prefixLen)
		var words []string
		for i := 0; i < maxWords; i++ {
			choices := c.chain[p.String()]
			if len(choices) == 0 {
				break
			}
			next := choices[rand.Intn(len(choices))]
			words = append(words, next)
			p.Shift(next)
		}

		if c.TestOutput(words, maxOverlapRatio, float64(maxOverLapTotal)) && len(words) > minWords {
			return strings.Join(words, " ")
		}
	}
	return ""
}

// Taken from https://github.com/jsvine/markovify/blob/master/markovify/text.py#L155
func (c *Chain) TestOutput(words []string, maxOverlapRatio, maxOverlapTotal float64) bool {
	overlapRatio := math.Round(maxOverlapRatio * float64(len(words)))
	overlapMax := math.Min(maxOverlapTotal, overlapRatio)
	overlapOver := overlapMax + 1

	gramCount := math.Max(float64(len(words))-overlapMax, 1)
	var grams [][]string
	for i := 0; i < int(gramCount); i++ {
		num := len(words)
		if overlapOver < float64(num) {
			num = int(overlapOver)
		}
		grams = append(grams, words[i:i+num])
	}
	for _, g := range grams {
		gramJoined := strings.Join(g, "")
		if strings.Contains(c.RejoinedText(), gramJoined) {
			return false
		}
	}
	return true
}

func (c *Chain) RejoinedText() string {
	var sentencesJoined []string
	for _, sentence := range c.sentences {
		sentencesJoined = append(sentencesJoined, strings.Join(sentence, " "))
	}

	return strings.Join(sentencesJoined, " ")
}

func GenQuote(history *History, followedUsers []fedi.Account, maxWords int) string {
	c := NewChain(2)

	for _, s := range history.Statuses {
		for _, u := range followedUsers {
			if s.AuthorID == u.ID {
				c.Build(strings.NewReader(s.Text))
			}
		}

	}
	text := c.Generate(maxWords, 1, 0, 0, 10000)
	// break generated mentions
	text = strings.ReplaceAll(text, "@", "@\u200B")
	return text
}
