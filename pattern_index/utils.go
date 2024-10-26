package patternindex

import (
	"strings"
)

func getPrefixes(word string) []string {
    prefixes := make([]string, len(word))
    for i := 1; i <= len(word); i++ {
        prefixes[i-1] = word[:i]
    }
    return prefixes
}

func SplitTextToListPrefixes(text string) []string {
	words := strings.Fields(text)
    var allPrefixes []string
    for _, word := range words {
        allPrefixes = append(allPrefixes, getPrefixes(word)...)
    }
    return allPrefixes
}

func getNgrams(word string, ngram int) []string {
    prefixes := make([]string, 0)
    for i := 0; i < len(word); i++ {
		for j := i+1; j <= len(word) && j - i <= ngram; j++ {
        	prefixes = append(prefixes, word[i:j])
		}
    }
    return prefixes
}

func SplitTextToListNGrams(text string, ngram int) []string {
	words := strings.Fields(text)
    var allPrefixes []string
    for _, word := range words {
        allPrefixes = append(allPrefixes, getNgrams(word, ngram)...)
    }
    return allPrefixes
}

func GetImportantSubwordFromPattern(pattern string, n int) []string {
    var result []string
    parts := strings.Split(pattern, "*")

    for _, part := range parts {
        if len(part) <= n {
            result = append(result, part)
        } else {
            for i := 0; i <= len(part)-n; i++ {
                result = append(result, part[i:i+n])
            }
        }
    }

    return result
}

func MatchPattern(word, pattern string) bool {
    parts := strings.Split(pattern, "*")
    start := 0

    for _, part := range parts {
        if part == "" {
            continue 
        }
        index := strings.Index(word[start:], part)
        if index == -1 {
            return false 
        }

        start += index + len(part)
    }

    if !strings.HasSuffix(pattern, "*") && start != len(word) {
        return false
    }

    return true
}

func MatchPatternToText(text, pattern string) bool {
	words := strings.Fields(text)
	for _, word := range words {
		if MatchPattern(word, pattern) {
			return true
		}
	}
	return false
}
