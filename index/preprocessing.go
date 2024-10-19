package index

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type PreProcessor struct {
}

func NewPreprocessor() *PreProcessor {
	return &PreProcessor{};
}

type LemmatizeRequest struct {
    Text string `json:"text"`
}

func (procesor *PreProcessor) Lemmatize(words string) ([]string, error) {
	requestBody, _ := json.Marshal(LemmatizeRequest{Text: words})
    resp, err := http.Post("http://localhost:5000/lemmatize", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result []string
    json.NewDecoder(resp.Body).Decode(&result)
    return result, nil
}

type StemRequest struct {
    Text string `json:"text"`
}

func (procesor *PreProcessor) Stem(text string) ([]string, error) {
    requestBody, _ := json.Marshal(StemRequest{Text: text})
    resp, err := http.Post("http://localhost:5000/stem", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result []string
    json.NewDecoder(resp.Body).Decode(&result)
    return result, nil
}

type StopWordsRequest struct {
    Text string `json:"text"`
}

type StopWordsResponse map[string]bool

func (procesor *PreProcessor) ClassifyStopWords(text string) (StopWordsResponse, error) {
    requestBody, _ := json.Marshal(StopWordsRequest{Text: text})
    resp, err := http.Post("http://localhost:5000/classify_stopwords", "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result StopWordsResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return result, nil
}

func (procesor *PreProcessor) LemmatizeAndRemoveStopWords(text string) ([]string, error) {
    lemmatizedWords, err := procesor.Lemmatize(text)
    if err != nil {
        return nil, err
    }

    stopWords, err := procesor.ClassifyStopWords(text)
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

// Стемминг с удалением стоп-слов
func (procesor *PreProcessor) StemAndRemoveStopWords(text string) ([]string, error) {
    stemmedWords, err := procesor.Stem(text)
    if err != nil {
        return nil, err
    }

    stopWords, err := procesor.ClassifyStopWords(text)
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
