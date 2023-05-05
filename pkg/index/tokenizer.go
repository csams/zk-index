package index

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Term is a single term from a document.
type Term string

// Terms parses the input into Terms using simple text splitting and a couple of rules specific to Telekasten.
func Terms(input io.Reader) []Term {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)

	var terms []Term
	for scanner.Scan() {
		word := scanner.Text()

		word = strings.ToLower(word)
		word = strings.TrimFunc(word, func(r rune) bool {
			keep := unicode.IsDigit(r) || unicode.IsLetter(r) || r == '#'
			return !keep
		})

		if len(word) > 0 {
			var ts []string
			for _, w := range strings.Split(word, "|") {
				ts = append(ts, w)
			}

			for _, t := range ts {
				for _, s := range strings.Split(t, "][") {
					terms = append(terms, Term(s))
				}
			}
		}
	}

	return terms
}
