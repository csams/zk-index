package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"io/fs"

	"github.com/csams/zk-index/pkg/index"
)

func BuildCorpus(base string) (index.Corpus, error) {
    if tmp, e := filepath.Abs(base); e != nil {
        return nil, e
    } else {
        base = tmp
    }

	corpus := index.Corpus{}
    walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}

		if d.IsDir() {
            if strings.HasPrefix(d.Name(), ".") {
                return fs.SkipDir
            }
			return nil
		}

		if filepath.Ext(path) == ".md" {
			corpus = append(corpus, index.Document(path))
		}
		return nil
	}

	err := filepath.WalkDir(base, walk)
	if err != nil {
		return nil, err
	}
	return corpus, nil
}

func main() {
	flag.Parse()
	directory := flag.Arg(0)

	fmt.Println(directory)

	corpus, err := BuildCorpus(directory)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", corpus)

	index, err := index.BuildIndex(corpus)
	if err != nil {
		panic(err)
	}
	b, e := json.Marshal(index)
	if e != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(b))
}
