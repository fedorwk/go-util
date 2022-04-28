package delimiterdetector_test

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	detector "github.com/fedorwk/go-util/data/delimiterdetector"
)

type TestCaseShouldPass struct {
	name  string
	input string
	want  string
}

type TestCaseShouldFail struct {
	name  string
	input string
}

func ExampleMain() {
	yourCSV :=
		`product;price;amount
apple;5;100
orange;7;30
something expensive;2,000;5 
`
	// Note: First line of csv shouldn't be blank line
	csvAsReader := strings.NewReader(yourCSV)
	delimiter, err := detector.Parse(csvAsReader, -1) // second parameter: lines of csv to analyse, -1 for all lines
	if err != nil {
		panic(err)
	}

	raggedCSV := strings.Replace(yourCSV, delimiter, ",", -1)
	// Input is ragged now, last line contains more delimiters than expected
	_, errRaggedCSV := detector.Parse(strings.NewReader(raggedCSV), -1)
	if errRaggedCSV == nil {
		panic("ragged input should throw error")
	}

	fmt.Printf("\nDelimiter of correct CSV: %s\nRagged CSV Error: %s\n", delimiter, errRaggedCSV.Error())
	// Output:
	// Delimiter of correct CSV: ;
	// Ragged CSV Error: unable to define delimiter
}

// How to use with custom delimiters
func ExampleDetector() {
	yourCSV :=
		`h1-h2-h3
11-12-13
21-22-23`
	detector := detector.New([]string{"-", "!", "&"})
	delimiter, err := detector.Parse(strings.NewReader(yourCSV), -1) // parse yourCSV as reader until EOF
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delimiter: %s\n", delimiter)
	// Output:
	// Delimiter: -
}

func TestDetermineShouldPass(t *testing.T) {
	testCases := parseTestCasesShouldPass("testdata/CSVshouldpass")
	for i, test := range testCases {
		t.Logf("#%d TEST: %s", i, test.name)
		got, err := detector.Default.Parse(strings.NewReader(test.input), -1)
		if err != nil {
			t.Error(err)
		}
		if got != test.want {
			t.Errorf("Want:\n%q\nGot:\n%q", test.want, got)
		}
	}
}

func TestDetermineShouldFail(t *testing.T) {
	testCases := parseTestCasesShouldFail("testdata/CSVshouldfail")
	for i, test := range testCases {
		t.Logf("#%d TEST: %s", i, test.name)
		_, err := detector.Default.Parse(strings.NewReader(test.input), -1)
		if err == nil {
			t.Errorf("Should Fail")
		}
	}
}

func parseTestCasesShouldPass(path string) []TestCaseShouldPass {
	testCases := make([]TestCaseShouldPass, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var testcase TestCaseShouldPass
		testcase.name = scanner.Text()
		scanner.Scan()
		testcase.want = scanner.Text()
		input := strings.Builder{}

		for scanner.Scan() && scanner.Text() != "" {
			_, err := input.WriteString(scanner.Text() + "\n")
			if err != nil {
				panic(err)
			}
		}
		testcase.input = input.String()
		testCases = append(testCases, testcase)
	}
	return testCases
}

func parseTestCasesShouldFail(path string) []TestCaseShouldFail {
	testCases := make([]TestCaseShouldFail, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var testcase TestCaseShouldFail
		testcase.name = scanner.Text()

		input := strings.Builder{}
		for scanner.Scan() && scanner.Text() != "" {
			_, err := input.WriteString(scanner.Text() + "\n")
			if err != nil {
				panic(err)
			}
		}
		testcase.input = input.String()
		testCases = append(testCases, testcase)
	}
	return testCases
}
