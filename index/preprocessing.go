package index

import (
    "github.com/jdkato/prose/v2"
    "github.com/kljensen/snowball"
	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
    "strings"
)

type StopWordsResponse map[string]bool

type PreProcessor struct {

}

func NewPreprocessor() *PreProcessor {
    return &PreProcessor{}
}

func (processor *PreProcessor) Lemmatize(text string) ([]string, error) {
    doc, err := prose.NewDocument(text)
    if err != nil {
        return nil, err
    }

    var lemmatizedWords []string
    for _, token := range doc.Tokens() {
        word := strings.ToLower(token.Text)
        lemmatizedWords = append(lemmatizedWords, processor.lemmatizeSimple(word))
    }

    return lemmatizedWords, nil
}

func (processor *PreProcessor) lemmatizeSimple(word string) string {
	lemmatizer, err := golem.New(en.New())
	if err != nil {
		panic(err)
	}
	res := lemmatizer.Lemma(word)
    return res
}

func (processor *PreProcessor) Stem(text string) ([]string, error) {
    words := strings.Fields(text)
    var stemmedWords []string

    for _, word := range words {
        stemmedWord, err := snowball.Stem(word, "english", true)
        if err != nil {
            return nil, err
        }
        stemmedWords = append(stemmedWords, stemmedWord)
    }

    return stemmedWords, nil
}

func (processor *PreProcessor) ClassifyStopWords(text string) (StopWordsResponse, error) {
    stopWords := map[string]bool{
        "the": true, "is": true, "at": true, "which": true, "on": true,
    }
    words := strings.Fields(text)
    response := make(StopWordsResponse)

    for _, word := range words {
        _, exists := stopWords[strings.ToLower(word)]
        response[word] = exists
    }

    return response, nil
}

func (processor *PreProcessor) LemmatizeAndRemoveStopWords(text string) ([]string, error) {
    lemmatizedWords, err := processor.Lemmatize(text)
    if err != nil {
        return nil, err
    }

    stopWords, err := processor.ClassifyStopWords(text)
    if err != nil {
        return nil, err
    }

    var result []string
    for _, word := range lemmatizedWords {
        if !stopWords[word] {
            result = append(result, word)
        }
    }

    return result, nil
}

func (processor *PreProcessor) StemAndRemoveStopWords(text string) ([]string, error) {
    stemmedWords, err := processor.Stem(text)
    if err != nil {
        return nil, err
    }

    stopWords, err := processor.ClassifyStopWords(text)
    if err != nil {
        return nil, err
    }

    var result []string
    for _, word := range stemmedWords {
        if !stopWords[word] {
            result = append(result, word)
        }
    }

    return result, nil
}
