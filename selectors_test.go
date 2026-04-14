package achim

import (
	"slices"
	"testing"
)

var tests = []struct {
	By       string
	Expected []Selector
}{
	{
		"foo=bar",
		[]Selector{
			{"foo", "bar"},
		},
	},
	{
		"foo=bar,qux=baz",
		[]Selector{
			{"foo", "bar"},
			{"qux", "baz"},
		},
	},
}

func TestParseSelectors(t *testing.T) {
	for _, test := range tests {
		actual, err := ParseSelector(test.By)
		if err != nil {
			t.Errorf(`parse by "%s": %v\n`, test.By, err)
		}
		if !slices.Equal(actual, test.Expected) {
			t.Errorf(`expected "%s" parse to %v, was %v\n`,
				test.By, test.Expected, actual)
		}
	}
}
