# Telekasten Indexer

This is a silly little indexer for Telekasten files that uses [tf-idf][tfidf] and [cosine similarity][cosine-sim]

```
./bin/index --addr 0.0.0.0:8080 <path to markdown files>
```

```
curl -L "http://localhost:8080?q=your+query+here&limit=10" | jq .
```

[tfidf]: https://en.wikipedia.org/wiki/Tf%E2%80%93idf
[cosine-sim]: https://en.wikipedia.org/wiki/Cosine_similarity
