package index

import (
	"math"
	"os"
	"sort"
)

// DocId is a unique identifier given to each document.
type DocId int64

// IdMap is a map from DocId → the document. It lets us look up document file paths.
type IdMap map[DocId]DocumentPath

// TermCounts is a map from Term → DocId → number of times the Term occurred in the doc.
type TermCounts map[Term]map[DocId]int64

// DocLengths stores the total number of terms in each document.
type DocLengths map[DocId]int64

// Stats are the statistics for the Corpus that are used to build the tfidf index.
type Stats struct {
	LookupDoc  IdMap
	DocLengths DocLengths
	TermCounts TermCounts
}

// Score is a particular Term's tfidf value for a given document.
type Score struct {
	DocId DocId
	Score float64
}

// Scores in a list of Scores across all documents for a given Term.
type Scores []Score

func (s Scores) Len() int {
	return len(s)
}

func (s Scores) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func (s Scores) Swap(i, j int) {
	s[j], s[i] = s[i], s[j]
}

// Index is the set of Scores for all Terms across all documents.
type Index map[Term]Scores

// BuildStats creates the document statistics for a Corpus that are used to build the tfidf Index.
func BuildStats(corpus Corpus) (Stats, error) {
	ids := 0

	idMap := IdMap{}
	docLengths := DocLengths{}
	termCounts := TermCounts{}

	for _, doc := range corpus {
		f, err := os.Open(string(doc))
		if err != nil {
			return Stats{}, err
		}
		defer f.Close()

		id := DocId(ids)
		totalTerms := int64(0)

		for _, term := range Terms(f) {
			var counts map[DocId]int64
			if f, ok := termCounts[term]; ok {
				counts = f
			} else {
				counts = map[DocId]int64{}
				termCounts[term] = counts
			}

			if _, ok := counts[id]; ok {
				counts[id] += 1
			} else {
				counts[id] = 1
			}

			totalTerms += 1
		}

		if err != nil {
			return Stats{}, err
		}

		docLengths[id] = totalTerms
		idMap[id] = doc
		ids += 1
	}

	return Stats{
		LookupDoc:  idMap,
		DocLengths: docLengths,
		TermCounts: termCounts,
	}, nil
}

// BuildIndex creates a tfidf based Index from the Stats of a Corpus.
func BuildIndex(stats Stats) Index {
	index := Index{}

	numDocs := float64(len(stats.LookupDoc))
	for term, entries := range stats.TermCounts {
		var scores Scores
		numDocsWithTerm := float64(len(entries))
		idf := math.Log(numDocs / numDocsWithTerm)
		for docId, freq := range entries {
			tf := float64(freq) / float64(stats.DocLengths[docId])
			scores = append(scores, Score{
				DocId: docId,
				Score: tf * idf,
			})
		}
		sort.Sort(sort.Reverse(scores))
		index[term] = scores
	}

	return index
}
