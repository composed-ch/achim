package achim

import (
	"fmt"
	"strings"
)

type Label struct {
	Key   string
	Value string
}

func ParseLabels(raw string) ([]Label, error) {
	pairs := strings.Split(raw, ",")
	selectors := make([]Label, len(pairs))
	for i, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf(`split selector "%s" by '='`, pair)
		}
		selectors[i] = Label{
			Key:   parts[0],
			Value: parts[1],
		}
	}
	return selectors, nil
}

func AsMap(labels []Label) map[string]string {
	m := make(map[string]string)
	for _, l := range labels {
		m[l.Key] = l.Value
	}
	return m
}
