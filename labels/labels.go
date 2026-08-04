package labels

import (
	"fmt"
	"maps"
	"strings"
)

type Label struct {
	Key   string
	Value string
}

func ParseLabels(raw string) ([]Label, error) {
	if raw == "" {
		return make([]Label, 0), nil
	}
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

func MergeMaps[K, V comparable](l, r map[K]V) map[K]V {
	m := make(map[K]V)
	maps.Copy(m, l)
	maps.Copy(m, r)
	return m
}

func Filter[T Filterable](items []T, by string) ([]T, error) {
	if by == "" {
		return items, nil
	}
	selectors, err := ParseLabels(by)
	if err != nil {
		return nil, fmt.Errorf(`parse --by "%s": %w`, by, err)
	}
	filtered := make([]T, 0)
	for _, instance := range items {
		retain := true
		for _, selector := range selectors {
			if selector.Key == "name" {
				if instance.Name() != selector.Value {
					retain = false
					break
				}
			} else {
				if value, ok := instance.Labels()[selector.Key]; !ok {
					retain = false
					break
				} else if value != selector.Value {
					retain = false
					break
				}
			}
		}
		if retain {
			filtered = append(filtered, instance)
		}
	}
	return filtered, nil
}
