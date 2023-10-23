// Package jellyset  provides a Redis-like Set data structure.
package jellyset

// keyExists is a placeholder to not write struct{}{} everywhere.
var keyExists = struct{}{}

// set is an alias for a map of elements to struct{} to mimic a set.
type set map[interface{}]struct{}

type Set struct {
	records map[string]set
}

// newSet creates and returns a new empty set.
func newSet() set {
	return make(map[interface{}]struct{})
}

// add adds one or more items to the set.
// if no items are provided, it has no effect.
func (s set) add(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		s[item] = keyExists
	}
}

// remove removes one or more items from the set.
// if passed nothing, it has no effect.
func (s set) remove(items ...interface{}) {
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		delete(s, item)
	}
}

// has looks for the existence of items passed.
// for multiple items, it returns true only if all of the items exist.
//
// it returns false if nothing is passed.
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

// copy creates a copy of the set and returns it.
func (s set) copy() set {
	copy := newSet()
	for item := range s {
		copy.add(item)
	}
	return copy
}

// list returns all items in the set as a slice.
func (s set) list() []interface{} {
	list := make([]interface{}, 0, len(s))

	for item := range s {
		list = append(list, item)
	}

	return list
}

// foreach iterates over the items in the set and calls the provided function for each set member.
// the iteration continues until all items in the set have been visited or the closure returns false.
func (s set) foreach(callback func(item interface{}) bool) {
	for item := range s {
		if callback(item) {
			break
		}
	}
}

// merge merges the current set with another set.
// It is basically the implementation of the set union between 2 sets.
func (s set) merge(t set) {
	for item := range t {
		s[item] = keyExists
	}
}

// separate removes the set items containing in t from set s.
// Does not undo merge!!
func (s set) separate(t set) {
	s.remove(t.list()...)
}

// size just returns the size of the s set
func (s set) size() int {
	return len(s)
}

// union returns a new set that is the union of multiple sets. It combines all elements
// present in all the sets provided as arguments.
func union(sets ...set) set {
	if len(sets) == 0 {
		return newSet()
	}

	totalSize := 0
	for _, s := range sets {
		totalSize += len(s)
	}

	unionSet := make(set, totalSize)

	for _, s := range sets {
		for item := range s {
			unionSet[item] = keyExists
		}
	}

	return unionSet
}
