package achim

import (
	"fmt"
	"strings"
)

type Selector struct {
	Label string
	Value string
}

func ParseSelector(by string) ([]Selector, error) {
	pairs := strings.Split(by, ",")
	selectors := make([]Selector, len(pairs))
	for i, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf(`split selector "%s" by '='`, pair)
		}
		selectors[i] = Selector{
			Label: parts[0],
			Value: parts[1],
		}
	}
	return selectors, nil
}
