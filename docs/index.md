Make a list of all the markdown documents beneath some root directory.

For a given document:

Associate an `Integer` identifier with the document. (`DocId → String`)

Tokenize the document into a list of words.
Normalize all of the words into `Terms`. (lowercase them, stem them, etc.)
Relate each `Term` to the number of times it appears in the document. (`Term → Integer`)

Associate the document's identifier with the number of `Terms` in it; e.g. the doc's length (`DocId → Integer`)


For each `Term` map (one for each document):

Relate each `Term` to a list of `Pairs`, each of which contains a document identifier and the number of times
the `Term` appeared in the document. (Build up a function `Term → DocId → Integer`)

Relate each `Term` to the number of documents in which it appeared.


`LookupDocName :: DocId → String`
`TermsInDoc :: DocId → Integer`
`TermFrequency :: Term → Integer`
`Term → DocId → Integer`
