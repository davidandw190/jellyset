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

// Copy creates a copy of the set and returns it.
func (s set) copy() set {
	copy := newSet()
	for item := range s {
		copy.add(item)
	}
	return copy
}

// List returns all items in the set as a slice.
func (s set) list() []interface{} {
	list := make([]interface{}, 0, len(s))

	for item := range s {
		list = append(list, item)
	}

	return list
}

// Foreach iterates over the items in the set and calls the provided function for each set member.
// The iteration continues until all items in the set have been visited or the closure returns false.
func (s set) foreach(callback func(item interface{}) bool) {
	for item := range s {
		if callback(item) {
			break
		}
	}
}

// Merge merges the current set with another set.
func (s set) merge(t set) {
	t.foreach(func(item interface{}) bool {
		s.add(item)
		return true
	})
}
