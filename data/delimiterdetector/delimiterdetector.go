package delimiterdetector

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// DetectDelimiters parses src (CSV file usually) and gives one possible delimiter of values
// if it can be exactly identified.
// It uses default Detector, parsed delimiters are: "," ";" "\t" "|"
// If your file may have another possible delimiter use New() to specify them.
// First nLines of src will be parsed (-1 for parse until EOF)
func DetectDelimiters(src io.Reader, nLines int) (string, error) {
	return Default.Parse(src, nLines)
}

// Detector is main structure of package and it's entry point
// Contains possible delimiters which are parsed while Detector.Parse()
type Detector struct {
	PossibleDelimiters []string
}

var (
	ErrEmptySource = errors.New("empty source given")
	// ErrDefiningDelimiter occures if there are no delimiters from PossibleDelimiters
	// or input is ragged (each row must have an equal number of delimiters)
	// See Examples in ./testdata/CSVshouldfail
	ErrDefiningDelimiter = errors.New("unable to define delimiter")
	// ErrMultipleDelimiterOptions occures if there are two or more options of possible
	// delimiters (each row has an equal number of two or more delimiters)
	// See Examples in ./testdata/CSVshouldfail
	ErrMultipleDelimiterOptions = errors.New("multiple delimiter options")
)

var Default = Detector{
	PossibleDelimiters: []string{",", ";", "\t", "|"},
}

// New() returns Detector with given possible delimiters
// Run Detector.Parse(src, nLines) to detect delimiter
// of src from the list of possible delimiters
func New(possibleDelimiters []string) *Detector {
	return &Detector{
		PossibleDelimiters: possibleDelimiters,
	}
}

type delimiterCounter map[string]int

// Parse() parses src (CSV file usually) and gives one possible delimiter of values
// if it can be exactly identified.
// Only nLines will be parsed (-1 for parse until EOF)
func (d Detector) Parse(src io.Reader, nLines int) (string, error) {
	scanner := bufio.NewScanner(src)
	var counter delimiterCounter
	if ok := scanner.Scan(); ok {
		counter = d.initDelimiterCounterByFirstLine(d.PossibleDelimiters, scanner.Text())
		nLines-- // first line read
	} else {
		return "", ErrEmptySource
	}

	for ; scanner.Scan() && nLines != 0; nLines-- { // read until EOF or nLines limit
		counter.filterDelimitersByString(scanner.Text())
		if counter.empty() {
			return "", ErrDefiningDelimiter
		}
	}
	return counter.result()
}

func (d Detector) initDelimiterCounterByFirstLine(delimiters []string, fline string) delimiterCounter {
	delimiterCounter := make(map[string]int)
	for _, key := range d.PossibleDelimiters {
		entriesCount := strings.Count(fline, key)
		if entriesCount > 0 {
			delimiterCounter[key] = entriesCount
		}
	}
	return delimiterCounter
}

func (dc delimiterCounter) filterDelimitersByString(line string) {
	for key := range dc {
		entriesCount := strings.Count(line, key)
		if entriesCount != dc[key] {
			delete(dc, key)
		}
	}
}

func (dc delimiterCounter) result() (key string, err error) {
	switch len(dc) {
	case 1:
		return dc.singleKey(), nil
	case 0:
		return "", ErrDefiningDelimiter
	default:
		return "", ErrMultipleDelimiterOptions
	}
}

func (dc delimiterCounter) empty() bool {
	return len(dc) == 0
}

func (dc delimiterCounter) containsSingleKey() bool {
	return len(dc) == 1
}

// returns one random key from delimiterCounter's map
func (dc delimiterCounter) singleKey() string {
	for key := range dc {
		return key
	}
	return ""
}
