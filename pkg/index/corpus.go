package index

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// DocumentPath is the absolute name of the file on the filesystem.
type DocumentPath string

// Corpus is the set of documents to index.
type Corpus []DocumentPath

func BuildCorpus(basePath string) (Corpus, error) {
	if tmp, e := filepath.Abs(basePath); e != nil {
		return nil, e
	} else {
		basePath = tmp
	}

	corpus := Corpus{}
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
			corpus = append(corpus, DocumentPath(path))
		}
		return nil
	}

	if err := filepath.WalkDir(basePath, walk); err != nil {
		return nil, err
	}
	return corpus, nil
}
