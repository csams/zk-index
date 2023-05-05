package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/csams/zk-index/pkg/index"
)

func main() {
	flag.Parse()
	directory := flag.Arg(0)

	corpus, err := index.BuildCorpus(directory)
	if err != nil {
		panic(err)
	}

	stats, err := index.BuildStats(corpus)
	if err != nil {
		panic(err)
	}

	tfidf := index.BuildIndex(stats)
	b, e := json.Marshal(tfidf)
	if e != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b))
}
