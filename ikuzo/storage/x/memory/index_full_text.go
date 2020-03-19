package memory

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/delving/hub3/ikuzo/service/x/search"
)

var (
	ErrSearchNoMatch = errors.New("the search query does not match the index")
)

// TestIndex is a single document full-text index.
// This means that all data you append to it will have its position incremented
// and appends to the known state. It is not replaced. To reset the index to
// an empty state you have to call the reset method.
type TextIndex struct {
	terms map[string]*search.TermVector
	a     search.Analyzer
}

func NewTextIndex() *TextIndex {
	return &TextIndex{
		terms: make(map[string]*search.TermVector),
	}
}

// AppendBytes extract words from bytes and updates the TextIndex.
func (ti *TextIndex) AppendBytes(b []byte) error {
	tok := search.NewTokenizer()
	for _, token := range tok.ParseBytes(b).Tokens() {
		if !token.Ignored {
			err := ti.addTerm(token.Normal, token.WordPosition)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AppendString extract words from bytes and updates the TextIndex.
func (ti *TextIndex) AppendString(text string) error {
	tok := search.NewTokenizer()
	for _, token := range tok.ParseString(text).Tokens() {
		if !token.Ignored {
			err := ti.addTerm(token.Normal, token.WordPosition)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ti *TextIndex) reset() {
	ti.terms = make(map[string]*search.TermVector)
}

func (ti *TextIndex) size() int {
	return len(ti.terms)
}

func (ti *TextIndex) setTermVector(term string, pos int, split bool) {
	tv, ok := ti.terms[term]
	if !ok {
		tv = search.NewTermVector()
		ti.terms[term] = tv
		tv.Split = split
	}

	tv.Positions[pos] = true
}

func (ti *TextIndex) addTerm(word string, pos int) error {
	if word == "" {
		return fmt.Errorf("cannot index empty string")
	}

	analyzedTerm := ti.a.Transform(word)

	if analyzedTerm == "" {
		return nil
	}

	ti.setTermVector(analyzedTerm, pos, false)

	if strings.Contains(analyzedTerm, "-") {
		for _, p := range strings.Split(analyzedTerm, "-") {
			ti.setTermVector(p, pos, true)
		}
	}

	return nil
}

func (ti *TextIndex) match(qt *search.QueryTerm, hits *search.Matches) bool {
	switch qt.Type() {
	case search.WildCardQuery:
		return ti.matchWildcard(qt, hits)
	case search.PhraseQuery:
		return ti.matchPhrase(qt, hits)
	case search.FuzzyQuery:
		return ti.matchFuzzy(qt, hits)
	default:
		// search.TermQuery is the default
		return ti.matchTerm(qt, hits)
	}
}

func (ti *TextIndex) matchPhrase(qt *search.QueryTerm, hits *search.Matches) bool {
	var nextPositions []int

	phrasePositions := map[int]string{}

	words := strings.Fields(qt.Value)

	if len(words) == 1 {
		term, ok := ti.terms[qt.Value]
		if !ok {
			return false
		}

		hits.AppendTerm(qt.Value, term.Size(), term.Positions)

		return true
	}

	var previousTerm string

	for idx, word := range words {
		term, ok := ti.terms[word]
		if !ok {
			return false
		}

		if idx != 0 {
			var wordMatch bool

			for _, pos := range nextPositions {
				// posMatch determines if this position can be followed
				var posMatch bool

				for _, nextPos := range search.ValidPhrasePosition(pos, qt.Slop) {
					_, ok := term.Positions[nextPos]
					if ok {
						phrasePositions[nextPos] = word
						posMatch = true
						wordMatch = true
					}
				}

				if posMatch {
					phrasePositions[pos] = previousTerm
				}
			}

			if !wordMatch {
				return false
			}
		}

		nextPositions = []int{}
		for pos := range term.Positions {
			nextPositions = append(nextPositions, pos)
		}

		previousTerm = word
	}

	matches := len(phrasePositions)

	if matches != 0 {
		wordPositions := map[int]bool{}
		for pos := range phrasePositions {
			wordPositions[pos] = true
		}

		hits.ApppendPositions(wordPositions)

		phraseHits := sortAndCountPhrases(words, phrasePositions)

		for k, v := range phraseHits {
			hits.AppendTerm(k, v, map[int]bool{})
		}
	}

	return matches != 0
}

func sortAndCountPhrases(words []string, phrases map[int]string) map[string]int {
	positions := []int{}
	phraseSize := len(words)
	phraseHits := map[string]int{}

	for k := range phrases {
		positions = append(positions, k)
	}

	sort.Slice(positions, func(i, j int) bool {
		return positions[i] < positions[j]
	})

	phrase := []string{}

	for idx, p := range positions {
		phrase = append(phrase, phrases[p])
		if len(phrase) != phraseSize {
			if idx < len(positions) {
				continue
			}
		}

		currentPhrase := strings.Join(phrase, " ")
		phraseHits[currentPhrase]++

		phrase = []string{}
	}

	return phraseHits
}

func (ti *TextIndex) matchFuzzy(qt *search.QueryTerm, hits *search.Matches) bool {
	var hasMatch bool

	for k, v := range ti.terms {
		ok, _ := search.IsFuzzyMatch(k, qt.Value, float64(qt.Fuzzy), search.Levenshtein)
		if ok {
			hasMatch = true

			hits.AppendTerm(k, v.Size(), v.Positions)
		}
	}

	return hasMatch
}

func (ti *TextIndex) matchWildcard(qt *search.QueryTerm, hits *search.Matches) bool {
	var matcher func(s, prefix string) bool

	switch {
	case qt.SuffixWildcard:
		matcher = strings.HasSuffix
	default:
		matcher = strings.HasPrefix
	}

	var hasMatch bool

	for k, v := range ti.terms {
		if matcher(k, qt.Value) {
			hasMatch = true

			hits.AppendTerm(k, v.Size(), v.Positions)
		}
	}

	return hasMatch
}

func (ti *TextIndex) matchTerm(qt *search.QueryTerm, hits *search.Matches) bool {
	term, ok := ti.terms[qt.Value]
	if ok && qt.Prohibited {
		return false
	}

	if !ok && !qt.Prohibited {
		return false
	}

	var count int

	if term == nil {
		term = search.NewTermVector()
	}

	if ok {
		count = term.Size()
	}

	hits.AppendTerm(qt.Value, count, term.Positions)

	return true
}

func (ti *TextIndex) Search(query *search.QueryTerm) (*search.Matches, error) {
	hits := search.NewMatches()
	err := ti.search(query, hits)

	return hits, err
}

func (ti *TextIndex) searchMustNot(query *search.QueryTerm, hits *search.Matches) error {
	for _, qt := range query.MustNot() {
		switch {
		case qt.Type() == search.BoolQuery:
			// TODO(kiivihal): find out why this is a dead branch
			if err := ti.search(qt, hits); err != nil {
				return err
			}
		default:
			if !qt.Prohibited {
				return fmt.Errorf("mustNot clauses must be marked as prohibited")
			}

			if ok := ti.match(qt, hits); !ok {
				return ErrSearchNoMatch
			}
		}
	}

	return nil
}

func (ti *TextIndex) searchMust(query *search.QueryTerm, hits *search.Matches) error {
	for _, qt := range query.Must() {
		switch {
		case qt.Type() == search.BoolQuery:
			if err := ti.search(qt, hits); err != nil {
				return err
			}
		default:
			if ok := ti.match(qt, hits); !ok {
				return ErrSearchNoMatch
			}
		}
	}

	return nil
}

func (ti *TextIndex) searchShould(query *search.QueryTerm, hits *search.Matches) error {
	var matched bool

	for _, qt := range query.Should() {
		switch {
		case qt.Type() == search.BoolQuery:
			if err := ti.search(qt, hits); err != nil {
				return err
			}
		default:
			if ok := ti.match(qt, hits); ok {
				matched = true
			}
		}
	}

	if !matched {
		return ErrSearchNoMatch
	}

	return nil
}

// recursive search function
func (ti *TextIndex) search(query *search.QueryTerm, hits *search.Matches) error {
	if len(query.MustNot()) != 0 {
		if err := ti.searchMustNot(query, hits); err != nil {
			return err
		}
	}

	if len(query.Must()) != 0 {
		if err := ti.searchMust(query, hits); err != nil {
			return err
		}
	}

	if len(query.Should()) != 0 {
		if err := ti.searchShould(query, hits); err != nil {
			return err
		}
	}

	return nil
}

func (ti *TextIndex) writeTo(w io.Writer) error {
	e := gob.NewEncoder(w)

	err := e.Encode(ti.terms)
	if err != nil {
		return fmt.Errorf("unable to marshall TextIndex to GOB; %w", err)
	}

	return nil
}

func (ti *TextIndex) readFrom(r io.Reader) error {
	var terms map[string]*search.TermVector

	d := gob.NewDecoder(r)

	err := d.Decode(&terms)
	if err != nil {
		return err
	}

	ti.terms = terms

	return nil
}