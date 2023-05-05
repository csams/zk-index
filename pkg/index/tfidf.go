package index

import (
	"io"
	"math"
	"os"
	"sort"
)

// DocId is a unique identifier given to each document.
type DocId int64

// IdMap is a map from DocId → the document. It lets us look up document file paths.
type IdMap map[DocId]DocumentPath

type TermFreqs map[Term]float64

// CorpusStats are the statistics for the Corpus that are used to build the tfidf index.
type CorpusStats struct {
	LookupDoc       IdMap
	LookupTermFreqs map[DocId]TermFreqs
}

// Score is a particular Term's tfidf value for a given document.
type Score struct {
	DocId DocId
	Score float64
}

// Index is the set of Scores for all Terms across all documents.
type Index struct {
	LookupDoc IdMap
	Scores    map[Term][]Score
}

// BuildTermFreqMap creates a map of Term frequencies for the given input.
func BuildTermFreqMap(r io.Reader) TermFreqs {
	hist := TermFreqs{}
	numTerms := 0
	for _, term := range Terms(r) {
		hist[term] += 1
		numTerms += 1
	}
	freq := TermFreqs{}
	for t, c := range hist {
		freq[t] = float64(c) / float64(numTerms)
	}
	return freq
}

// BuildCorpusStats creates the document statistics for a Corpus that are used to build the tfidf Index.
func BuildCorpusStats(corpus Corpus) (CorpusStats, error) {
	ids := int64(0)
	idMap := IdMap{}

	docStats := map[DocId]TermFreqs{}
	for _, doc := range corpus {
		f, err := os.Open(string(doc))
		if err != nil {
			return CorpusStats{}, err
		}
		defer f.Close()

		id := DocId(ids)
		docStats[id] = BuildTermFreqMap(f)

		idMap[id] = doc
		ids += 1
	}

	return CorpusStats{
		LookupDoc:       idMap,
		LookupTermFreqs: docStats,
	}, nil
}

// BuildIndex creates a tfidf based Index from the Stats of a Corpus.
func BuildIndex(corpus CorpusStats) Index {
	lookupTermFreqs := corpus.LookupTermFreqs
	docsWithTerm := map[Term]float64{}
	for _, freqs := range lookupTermFreqs {
		for term := range freqs {
			docsWithTerm[term] += 1.0
		}
	}

	idx := Index{
		LookupDoc: corpus.LookupDoc,
		Scores:    map[Term][]Score{},
	}
	numDocs := float64(len(lookupTermFreqs))
	for docId, termFreqs := range lookupTermFreqs {
		for term, tf := range termFreqs {
			numDocsWithTerm := docsWithTerm[term]
			idf := math.Log(numDocs/(1+numDocsWithTerm)) + 1
			if _, ok := idx.Scores[term]; !ok {
				idx.Scores[term] = []Score{}
			}
			idx.Scores[term] = append(idx.Scores[term], Score{
				DocId: docId,
				Score: tf * idf,
			})
		}
	}

	return idx
}

type QueryResult struct {
	Document DocumentPath
	Score    float64
}

type QueryResults []QueryResult

func (r QueryResults) Len() int {
	return len(r)
}

func (r QueryResults) Less(i, j int) bool {
	return r[i].Score < r[j].Score
}

func (r QueryResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Query determines matches using cosine similarity. Its input is also an Index since user queries are
// tokenized and scored the same way as documents so they can be compared.
// See: https://en.wikipedia.org/wiki/Cosine_similarity
func (idx Index) Query(qry Index) []QueryResult {
	var res QueryResults

	numers := map[DocId]float64{}
	b_d := map[DocId]float64{}
	a_d := float64(0.0)
	for at, a_scores := range qry.Scores {
		// there's only one "document" in the query
		ascore := a_scores[0]
		a_d += ascore.Score * ascore.Score

		if b_scores, found := idx.Scores[at]; found {
			for _, bscore := range b_scores {
				numers[bscore.DocId] += ascore.Score * bscore.Score
				b_d[bscore.DocId] += bscore.Score * bscore.Score
			}
		}
	}

	ad := math.Sqrt(a_d)
	for doc, n := range numers {
		bd := math.Sqrt(b_d[doc])
		res = append(res, QueryResult{
			Document: idx.LookupDoc[doc],
			Score:    n / ad * bd,
		})
	}

	sort.Sort(sort.Reverse(res))
	return res
}
