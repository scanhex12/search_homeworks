package patternindex

import (
	"testing"
    "reflect"

	"github.com/stretchr/testify/assert"
)

func TestSplitting(t *testing.T) {
	text := "Example text to split"
	splitted := SplitTextToListPrefixes(text)

	assert.Equal(t, splitted, []string{"E", "Ex", "Exa", "Exam", "Examp", "Exampl", "Example", "t", "te", "tex", "text", "t", "to", "s", "sp", "spl", "spli", "split"})
}

func TestSplitting2(t *testing.T) {
	text := "Example text to split"
	splitted := SplitTextToListNGrams(text, 3)
	assert.Equal(t, splitted, []string{"E", "Ex", "Exa", "x", "xa", "xam", "a", "am", "amp", "m", "mp", "mpl", "p", "pl", "ple", "l", "le", "e", "t", "te", "tex", "e", "ex", "ext", "x", "xt", "t", "t", "to", "o", "s", "sp", "spl", "p", "pl", "pli", "l", "li", "lit", "i", "it", "t"})
}

func TestGetImportantSubwordFromPattern(t *testing.T) {
    tests := []struct {
        pattern string
        n       int
        want    []string
    }{
        {
            pattern: "hello*world",
            n:       3,
            want:    []string{"hel", "ell", "llo", "wor", "orl", "rld"},
        },
        {
            pattern: "go*language",
            n:       4,
            want:    []string{"go", "lang", "angu", "ngua", "guag", "uage"},
        },
        {
            pattern: "short*test",
            n:       5,
            want:    []string{"short", "test"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.pattern, func(t *testing.T) {
            got := GetImportantSubwordFromPattern(tt.pattern, tt.n)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetImportantSubwordFromPattern(%q, %d) = %v, want %v", tt.pattern, tt.n, got, tt.want)
            }
        })
    }
}

func TestMatchPattern(t *testing.T) {
    tests := []struct {
        word    string
        pattern string
        want    bool
    }{
        {"hello", "he*o", true},
        {"hello", "h*llo", true},
        {"hello", "*ell*", true},
        {"hello", "h*l*p", false},
        {"language", "l*ng*ge", true},
        {"pattern", "patt*ern*", true},
        {"golang", "*lang", true},
        {"golang", "go*java", false},
        {"example", "ex*mpl", false},
    }

    for _, tt := range tests {
        t.Run(tt.pattern, func(t *testing.T) {
            got := MatchPattern(tt.word, tt.pattern)
            if got != tt.want {
                t.Errorf("MatchPattern(%q, %q) = %v, want %v", tt.word, tt.pattern, got, tt.want)
            }
        })
    }
}
