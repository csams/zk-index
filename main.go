package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/csams/zk-index/pkg/index"
	"github.com/csams/zk-index/pkg/serve"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "host and port to listen on")
	flag.Parse()

	directory := flag.Arg(0)
	corpus, err := index.BuildCorpus(directory)
	if err != nil {
		panic(err)
	}

	fmt.Println("Indexing:")
	for _, path := range corpus {
		fmt.Println(path)
	}

	stats, err := index.BuildCorpusStats(corpus)
	if err != nil {
		panic(err)
	}

	tfidf := index.BuildIndex(stats)

	fmt.Printf("Listening on %s\n", *addr)
	r := serve.BuildQueryRouter(tfidf)
	if err := http.ListenAndServe(*addr, r); err != nil {
		panic(err)
	}
}
