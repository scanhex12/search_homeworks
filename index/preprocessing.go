package index

import (
    "github.com/jdkato/prose/v2"
    "github.com/kljensen/snowball"
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

// Пример простого подхода для английской лемматизации
func (processor *PreProcessor) lemmatizeSimple(word string) string {
    // Примеры замены окончаний
    if strings.HasSuffix(word, "ing") {
        return strings.TrimSuffix(word, "ing")
    } else if strings.HasSuffix(word, "ed") {
        return strings.TrimSuffix(word, "ed")
    }
    return word
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
        // Дополните список нужными стоп-словами
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
