package serve

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/csams/zk-index/pkg/index"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Result struct {
	Docs []index.QueryResult
}

func BuildQueryRouter(idx index.Index) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		limitParam := r.URL.Query().Get("limit")

		limit := 10

		if limitParam != "" {
			if l, err := strconv.Atoi(limitParam); err == nil {
				limit = l
			}
		}

		q := index.BuildTermFreqMap(strings.NewReader(query))
		v := index.BuildIndex(
			index.CorpusStats{
				LookupDoc:       nil,
				LookupTermFreqs: map[index.DocId]index.TermFreqs{0: q},
			})

		docs := idx.Query(v)
		if len(docs) > limit {
			docs = docs[:limit]
		}
		res := Result{
			Docs: docs,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	})
	return r
}
