package achim

import (
	"slices"
	"testing"
)

var tests = []struct {
	By       string
	Expected []Label
}{
	{
		"foo=bar",
		[]Label{
			{"foo", "bar"},
		},
	},
	{
		"foo=bar,qux=baz",
		[]Label{
			{"foo", "bar"},
			{"qux", "baz"},
		},
	},
}

func TestParseLabels(t *testing.T) {
	for _, test := range tests {
		actual, err := ParseLabels(test.By)
		if err != nil {
			t.Errorf(`parse by "%s": %v\n`, test.By, err)
		}
		if !slices.Equal(actual, test.Expected) {
			t.Errorf(`expected "%s" parse to %v, was %v\n`,
				test.By, test.Expected, actual)
		}
	}
}
