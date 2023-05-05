package index

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Term is a single term from a document.
type Term string

func Terms(f io.Reader) []Term {
	scanner := bufio.NewScanner(f)
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
