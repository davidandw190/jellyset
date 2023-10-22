// Package jellyset  provides a Redis-like Set data structure.
package jellyset

// keyExists is a placeholder to not write struct{}{} everywhere.
var keyExists = struct{}{}

// set is an alias for a map of elements to struct{} to mimic a set.
type set map[interface{}]struct{}

type Set struct {
	records map[string]set
}

// NewSet creates and returns a new empty set.
func newSet() set {
	return make(map[interface{}]struct{})
}

// Add adds one or more items to the set.
// If no items are provided, it has no effect.
func (s set) add(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		s[item] = keyExists
	}
}

// Remove removes one or more items from the set.
// If passed nothing, it has no effect.
func (s set) remove(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		delete(s, item)
	}
}

// Has looks for the existence of items passed.
// For multiple items, it returns true only if all of the items exist.
//
// It returns false if nothing is passed.
func (s set) has(items ...interface{}) bool {
	if len(items) == 0 {
		return false
	}

	exist := true
	for _, item := range items {
		if _, exist = s[item]; !exist {
			break
		}
	}

	return exist
}
