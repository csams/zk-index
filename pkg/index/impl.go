package index

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"
)

// DocId is a unique identifier given to documents
type DocId int64

// Term is a single term from a document
type Term string

// Document is the file name of a document
type Document string

// Corpus is the set of documents to index
type Corpus []Document

// IdMap is a map from DocId → the document
type IdMap map[DocId]Document

// TermCountMap is a map DocId → number of terms in the doc
type TermCountMap map[DocId]int64

// TermIndex is a map from Term → DocId → Number of times the term occurred in DocId
type TermIndex map[Term]TermCountMap

type Index struct {
	IdMap      IdMap
	TermCounts TermCountMap
	TermIndex  TermIndex
}

func BuildIndex(corpus Corpus) (*Index, error) {
	docIdCounter := 0

	idMap := IdMap{}
	termCounts := TermCountMap{}
	termIndex := TermIndex{}

	for _, doc := range corpus {
		id := DocId(docIdCounter)
		termCounter := int64(0)

		err := WithContents(doc,
			func(contents io.Reader) error {
				terms := GetTerms(contents)
				for _, term := range terms {
					var freqMap map[DocId]int64
					if f, ok := termIndex[term]; ok {
						freqMap = f
					} else {
						freqMap = map[DocId]int64{}
						termIndex[term] = freqMap
					}

					if _, ok := freqMap[id]; ok {
						freqMap[id] += 1
					} else {
						freqMap[id] = 1
					}

					termCounter += 1
				}

				idMap[id] = doc
				termCounts[id] = termCounter
				docIdCounter += 1

				return nil
			})

		if err != nil {
			return nil, err
		}
	}

	return &Index{
		IdMap:      idMap,
		TermCounts: termCounts,
		TermIndex:  termIndex,
	}, nil
}

func WithContents(doc Document, cb func(r io.Reader) error) error {
	f, err := os.Open(string(doc))
	if err != nil {
		return err
	}
	defer f.Close()
	return cb(f)
}

func GetTerms(f io.Reader) []Term {
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	var terms []Term
	for scanner.Scan() {
		word := scanner.Text()

		word = strings.ToLower(word)
		word = strings.TrimFunc(word, func(r rune) bool {
			alphaNum := unicode.IsDigit(r) || unicode.IsLetter(r)
			return !alphaNum
		})

		if len(word) > 0 {
			term := Term(word)
			terms = append(terms, term)
		}
	}

	return terms
}
