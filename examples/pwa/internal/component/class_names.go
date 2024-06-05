package component

import "sort"

type C map[string]bool

func classNames(classes C) (classNames string) {
	keys := make([]string, 0, len(classes))
	for k := range classes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i > 0 {
			classNames += " "
		}
		if classes[key] {
			classNames += key
		}
	}
	return
}
