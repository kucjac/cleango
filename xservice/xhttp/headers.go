package xhttp

import (
	"net/http"
	"sort"
	"strings"
)

var octetTypes [256]octetType

func init() {
	for c := 0; c < 256; c++ {
		var t octetType
		isCtl := c <= 31 || c == 127
		isChar := 0 <= c && c <= 127
		isSeparator := strings.ContainsRune(" \t\"(),/:;<=>?@[]\\{}", rune(c))
		if strings.ContainsRune(" \t\r\n", rune(c)) {
			t |= isSpace
		}
		if isChar && !isCtl && !isSeparator {
			t |= isToken
		}
		octetTypes[c] = t
	}
}

// ParseAcceptEncodingHeader parses 'Accept-Encoding' header values sorted by it's quality.
func ParseAcceptEncodingHeader(header http.Header) []QualityValue {
	return parseQVHeader(header, "Accept-Encoding")
}

// ParseAcceptHeader parses 'Accept' header values sorted by it's quality.
func ParseAcceptHeader(header http.Header) []QualityValue {
	return parseQVHeader(header, "Accept")
}

func parseQVHeader(header http.Header, headerName string) []QualityValue {
	var sorter acceptSorter
	s := header.Get(headerName)
	for {
		var spec QualityValue
		spec.Value, s = expectTokenSlash(s)
		if spec.Value == "" {
			break
		}
		spec.Quality = 1.0
		s = skipSpace(s)
		if strings.HasPrefix(s, ";") {
			s = skipSpace(s[1:])
			if !strings.HasPrefix(s, "q=") {
				break
			}
			spec.Quality, s = expectQuality(s[2:])
			if spec.Quality < 0.0 {
				break
			}
		}
		sorter = append(sorter, spec)
		s = skipSpace(s)
		if !strings.HasPrefix(s, ",") {
			break
		}
		s = skipSpace(s[1:])
	}
	sort.Sort(sorter)
	return sorter
}

func expectQuality(s string) (q float64, rest string) {
	switch {
	case len(s) == 0:
		return -1, ""
	case s[0] == '0':
		q = 0
	case s[0] == '1':
		q = 1
	default:
		return -1, ""
	}
	s = s[1:]
	if !strings.HasPrefix(s, ".") {
		return q, s
	}
	s = s[1:]
	i := 0
	n := 0
	d := 1
	for ; i < len(s); i++ {
		b := s[i]
		if b < '0' || b > '9' {
			break
		}
		n = n*10 + int(b) - '0'
		d *= 10
	}
	return q + float64(n)/float64(d), s[i:]
}

func expectTokenSlash(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		b := s[i]
		if (octetTypes[b]&isToken == 0) && b != '/' {
			break
		}
	}
	return s[:i], s[i:]
}

func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isSpace == 0 {
			break
		}
	}
	return s[i:]
}

type octetType byte

const (
	isToken octetType = 1 << iota
	isSpace
)

// QualityValue is the structure that contains quality - value pair.
type QualityValue struct {
	Value   string
	Quality float64
}

var _ sort.Interface = acceptSorter{}

type acceptSorter []QualityValue

// Len implements sort.Interface.
func (a acceptSorter) Len() int {
	return len(a)
}

// Swap implements sort.Interface.
func (a acceptSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less implements sort.Interface.
func (a acceptSorter) Less(i, j int) bool {
	return a.by(&a[i], &a[j])
}

func (a acceptSorter) by(a1, a2 *QualityValue) bool {
	return a1.Quality > a2.Quality
}
