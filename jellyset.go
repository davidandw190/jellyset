// Package jellyset  provides a Redis-like Set data structure.
package jellyset

import "math"

// keyExists is a placeholder to not write struct{}{} everywhere.
var keyExists = struct{}{}

// set represents individual sets within the Set type.
// Each set is implemented as a map, with keys representing the elements in the set.
type set map[interface{}]struct{}

// Set represents the high-level interface for interacting with sets.
// It encapsulates multiple sets, each associated with a unique key.
type Set struct {
	records map[string]set
}

func New() *Set {
	return &Set{
		records: make(map[string]set),
	}
}

// newSet creates and returns a new empty set.
func newSet() set {
	return make(map[interface{}]struct{})
}

// SAdd adds one or more members to the set associated with the provided key. If the key does not exist,
// it creates a new set and adds the specified members to it. This function returns the number of elements
// that were successfully added to the set.
//
// Parameters:
//   - key: 	The key associated with the set.
//   - members: One or more members to be added to the set.
//
// Returns:
//   - The number of elements added to the set.
//
// Example:
//   set := New()
//   count := set.SAdd("myset", "member1", "member2", "member3")
//
// In this example, three members are added to the set "myset," and the function returns the count of elements added.

func (s *Set) SAdd(key string, members ...interface{}) int {
	if !s.exists(key) {
		s.records[key] = newSet()
	}

	added := 0
	set := s.records[key]

	for _, member := range members {
		if _, exists := set[member]; !exists {
			set[member] = keyExists
			added++
		}
	}

	return added
}

// SPop removes and returns one or more random members from the set associated with the given key.
// If the key does not exist or the count is less than or equal to 0, it returns an empty slice.
//
// Parameters:
//   - key: 	The key associated with the set.
//   - count: 	The number of random members to pop from the set. If count is 0 or negative, no members are popped.
//
// Returns:
//   - A slice containing the popped members. If the set is empty or the count is zero, an empty slice is returned.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
//	popped := set.SPop("myset", 3)
//
// In this example, three random members are removed and returned from the set "myset," and they are stored in the 'popped' slice.
func (s *Set) SPop(key string, count int) []interface{} {
	if !s.exists(key) || count <= 0 {
		return []interface{}{}
	}

	set := s.records[key]
	members := make([]interface{}, count)

	i := 0
	for k := range set {
		members[i] = k
		delete(set, k)
		i++

		if i == count {
			break
		}
	}

	return members
}

// SRandMember returns one or more random members from the set associated with the given key.
// If the key does not exist or the count is less than 1, it returns an empty slice.
//
// Parameters:
//   - key: 	The key associated with the set.
//   - count: 	The number of random members to retrieve from the set. If count is less than 1, no members are retrieved.
//
// Returns:
//   - A slice containing the random members. If the set is empty or the count is less than 1, an empty slice is returned.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
//	randomMembers := set.SRandMember("myset", 3)
//
// In this example, three random members are retrieved from the set "myset," and they are stored in the 'randomMembers' slice.
func (s *Set) SRandMember(key string, count int) []interface{} {
	if !s.exists(key) || count < 1 {
		return []interface{}{}
	}

	set := s.records[key]
	members := make([]interface{}, count)

	if count > 0 {
		i := 0
		for k := range set {
			members[i] = k
			i++

			if i == count {
				break
			}
		}
	} else {
		count = -count
		for i := 0; i < count; i++ {
			randomVal := randomElement(set)
			if randomVal == nil {
				break
			}
			members[i] = randomVal
		}
	}
	return members
}

// SIsMember checks if the specified member exists in the set associated with the given key.
// If the key does not exist, it returns false.
//
// Parameters:
//   - key: 	The key associated with the set.
//   - member: 	The member to check for existence in the set.
//
// Returns:
//   - true if the member exists in the set, false otherwise.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3")
//	exists := set.SIsMember("myset", "member2")
//
// In this example, it checks if "member2" exists in the set "myset," and 'exists' will be true.
func (s *Set) SIsMember(key string, member interface{}) bool {
	if !s.exists(key) {
		return false
	}

	set := s.records[key]
	_, exists := set[member]

	return exists
}

// SRem removes the specified member from the set associated with the given key.
// If the key does not exist or the member is not in the set, it returns false.
//
// Parameters:
//   - key: 	The key associated with the set.
//   - member: 	The member to remove from the set.
//
// Returns:
//   - true if the member was successfully removed, false otherwise.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3")
//	removed := set.SRem("myset", "member2")
//
// In this example, it removes "member2" from the set "myset," and 'removed' will be true.
func (s *Set) SRem(key string, member interface{}) bool {
	if !s.exists(key) {
		return false
	}

	set := s.records[key]

	if _, exists := set[member]; exists {
		delete(set, member)
		return true
	}

	return false
}

// SMove moves a member from the source set to the destination set.
// If the source set does not exist or the member is not in the source set, it returns false.
// If the destination set does not exist, it creates a new set.
//
// Parameters:
//   - src: 	The key associated with the source set.
//   - dest: 	The key associated with the destination set.
//   - member: 	The member to move from the source set to the destination set.
//
// Returns:
//   - true if the member was successfully moved, false otherwise.
//
// Example:
//
//	set := New()
//	set.SAdd("sourceSet", "member1", "member2")
//	set.SMove("sourceSet", "destSet", "member2")
//
// In this example, it moves "member2" from the "sourceSet" to the "destSet," and it returns true.
func (s *Set) SMove(src, dest string, member interface{}) bool {
	if !s.fieldExists(src, member) {
		return false
	}

	if !s.exists(dest) {
		s.records[dest] = make(set)
	}

	srcSet := s.records[src]
	destSet := s.records[dest]

	srcSet.remove(member)
	destSet.add(member)

	return true
}

// SCard returns the number of elements in the set associated with the given key.
// If the key does not exist, it returns 0, indicating an empty set.
//
// Parameters:
//   - key:		The key associated with the set.
//
// Returns:
//   - The number of elements in the set.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3")
//	size := set.SCard("myset")
//
// In this example, it retrieves the size of the set "myset," which contains three members, and 'size' will be 3.
func (s *Set) SCard(key string) int {
	if !s.exists(key) {
		return 0
	}

	set := s.records[key]
	return set.size()
}

// SMembers returns a slice containing all the members of the set associated with the given key.
// If the key does not exist, it returns an empty slice.
//
// Parameters:
//   - key: 	The key associated with the set.
//
// Returns:
//   - A slice containing all the members of the set. If the set is empty or the key does not exist, an empty slice is returned.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3")
//	members := set.SMembers("myset")
//
// In this example, it retrieves all members from the set "myset," and 'members' will be a slice containing ["member1", "member2", "member3"].
func (s *Set) SMembers(key string) []interface{} {
	if !s.exists(key) {
		return []interface{}{}
	}

	set := s.records[key]
	members := make([]interface{}, 0, len(set))
	for item := range set {
		members = append(members, item)
	}

	return members
}

// SUnion returns a new set that is the union of multiple sets. It combines all elements
// present in all the sets provided as arguments.
//
// Parameters:
//   - keys: 	The keys associated with the sets to be combined in the union.
//
// Returns:
//   - A slice containing the union of elements from all the specified sets.
//
// Example:
//
//	set := New()
//	set.SAdd("set1", "member1", "member2", "member3")
//	set.SAdd("set2", "member2", "member3", "member4")
//	result := set.SUnion("set1", "set2")
//
// In this example, the union of "set1" and "set2" is computed, and 'result' contains all unique elements from both sets.
func (s *Set) SUnion(keys ...string) []interface{} {
	if len(keys) == 0 {
		return []interface{}{}
	}

	unionSet := newSet()

	for _, key := range keys {
		if s.exists(key) {
			set := s.records[key]
			for member := range set {
				unionSet[member] = keyExists
			}
		}
	}

	return unionSet.list()
}

// SUnionStore computes the union of multiple sets and stores the result in a new set.
//
// Parameters:
//   - storeKey: 	The key associated with the destination set where the result will be stored.
//   - keys: 		The keys associated with the sets to be combined in the union.
//
// Returns:
//   - The number of elements in the resulting union set.
//
// Example:
//
//	set := New()
//	set.SAdd("set1", "member1", "member2", "member3")
//	set.SAdd("set2", "member2", "member3", "member4")
//	count := set.SUnionStore("unionSet", "set1", "set2")
//
// In this example, the union of "set1" and "set2" is computed and stored in "unionSet," and 'count' contains the number of elements in the resulting union set.
func (s *Set) SUnionStore(storeKey string, keys ...string) int {
	union := s.SUnion(keys...)
	for _, unionKey := range union {
		s.SAdd(storeKey, unionKey)
	}

	return len(union)
}

// SKeyExists checks if the specified key exists in the Set.
//
// Parameters:
//   - key: 	The key to check for existence.
//
// Returns:
//   - true if the key exists in the Set, false otherwise.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3")
//	exists := set.SKeyExists("myset")
//
// In this example, it checks if the key "myset" exists in the Set, and 'exists' will be true.
func (s *Set) SKeyExists(key string) bool {
	return s.exists(key)
}

// SClear deletes the specified key and its associated set from the records.
//
// Parameters:
//   - key: 	The key associated with the set to be cleared.
//
// Example:
//
//	set := New()
//	set.SAdd("myset", "member1", "member2", "member3")
//	set.SClear("myset")
//
// In this example, the set associated with the key "myset" is deleted from the records.
func (s *Set) SClear(key string) {
	if s.exists(key) {
		delete(s.records, key)
	}
}

// SDiff returns a new set that contains items which are in the first set but not in the others.
//
// Parameters:
//   - keys: 	The keys associated with the sets to be used in the difference operation.
//
// Returns:
//   - A slice containing the elements that are present in the first set but not in the other specified sets.
//
// Example:
//
//	set := New()
//	set.SAdd("set1", "member1", "member2", "member3")
//	set.SAdd("set2", "member2", "member3", "member4")
//	result := set.SDiff("set1", "set2")
//
// In this example, the difference between "set1" and "set2" is computed, and 'result' contains elements unique to "set1."
func (s *Set) SDiff(keys ...string) []interface{} {
	if len(keys) == 0 {
		return []interface{}{}
	}

	if len(keys) == 1 {
		if s.exists(keys[0]) {
			return s.records[keys[0]].list()
		}

		return []interface{}{}
	}

	excludeMap := make(map[interface{}]bool)

	for _, key := range keys {
		if key != keys[0] {
			nextSet, ok := s.records[key]
			if !ok {
				return []interface{}{}
			}

			for item := range nextSet {
				excludeMap[item] = true
			}
		}

	}

	firstSet := s.records[keys[0]]
	result := make([]interface{}, 0, len(firstSet))

	for item := range firstSet {
		if !excludeMap[item] {
			result = append(result, item)
		}
	}

	return result
}

// SDiffStore computes the set difference between the first key provided and all the other keys.
// It stores the result in a new set identified by storeKey.
//
// Parameters:
//   - storeKey: 	The key where the resulting set difference will be stored.
//   - keys: 		One or more keys associated with the sets to calculate the difference.
//
// Returns:
//   - The number of elements in the resulting difference set.
//
// Example:
//
//	set := New()
//	set.SAdd("set1", "member1", "member2", "member3")
//	set.SAdd("set2", "member2", "member3", "member4")
//	count := set.SDiffStore("resultSet", "set1", "set2")
//
// In this example, it calculates the difference between "set1" and "set2" and stores the result in "resultSet."
// The resulting difference set contains "member1," and 'count' will be 1.
func (s *Set) SDiffStore(storeKey string, keys ...string) int {
	difference := s.SDiff(keys...)

	for _, diffKey := range difference {
		s.SAdd(storeKey, diffKey)
	}

	return len(difference)
}

// SInter returns a new set that contains items present in all the specified sets.
//
// Parameters:
//   - keys: 	The keys associated with the sets to be intersected.
//
// Returns:
//   - A slice containing the intersection of elements from all the specified sets.
//
// Example:
//
//	set := New()
//	set.SAdd("set1", "member1", "member2", "member3")
//	set.SAdd("set2", "member2", "member3", "member4")
//	result := set.SInter("set1", "set2")
//
// In this example, the intersection of "set1" and "set2" is computed, and 'result'.
func (s *Set) SInter(keys ...string) []interface{} {
	if len(keys) == 0 {
		return []interface{}{}
	}

	if len(keys) == 1 {
		if s.exists(keys[0]) {
			return s.records[keys[0]].list()
		}
		return []interface{}{}
	}

	var smallestSet set
	var smallestKey string
	var smallestSize = math.MaxInt

	for _, key := range keys {
		currentSet, ok := s.records[key]
		if !ok {
			return []interface{}{}
		}

		if len(currentSet) < smallestSize {
			smallestSize = len(currentSet)
			smallestSet = currentSet
			smallestKey = key
		}
	}

	inAllSets := make(map[interface{}]bool)

	for item := range smallestSet {
		inAllSets[item] = true
	}

	for _, key := range keys {
		if key != smallestKey {
			nextSet, ok := s.records[key]
			if !ok {
				return []interface{}{}
			}

			for item := range inAllSets {
				if !nextSet.has(item) {
					delete(inAllSets, item)
				}
			}
		}
	}

	result := make([]interface{}, 0, len(inAllSets))
	for item := range inAllSets {
		result = append(result, item)
	}

	return result
}

// SInterStore computes the intersection of sets specified by the provided keys
// and stores the result in a new set identified by storeKey.
//
// Parameters:
//   - storeKey: 	The key where the resulting intersection will be stored.
//   - keys: 		One or more keys associated with the sets to be intersected.
//
// Returns:
//   - The number of elements in the resulting intersection set.
//
// Example:
//
//	set := New()
//	set.SAdd("set1", "member1", "member2", "member3")
//	set.SAdd("set2", "member2", "member3", "member4")
//	count := set.SInterStore("resultSet", "set1", "set2")
//
// In this example, it calculates the intersection of "set1" and "set2" and stores the result in "resultSet."
// The resulting intersection set contains "member2" and "member3," and 'count' will be 2.
func (s *Set) SInterStore(storeKey string, keys ...string) int {
	intersection := s.SInter(keys...)

	for _, interKey := range intersection {
		s.SAdd(storeKey, interKey)
	}

	return len(intersection)
}

// existsInAll checks if an item exists in all given sets.
func existsInAll(item interface{}, currentKey string, keys []string, s *Set) bool {
	for _, key := range keys {
		if key != currentKey {
			nextSet, ok := s.records[key]
			if !ok || !nextSet.has(item) {
				return false
			}
		}
	}
	return true
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

	i := 0
	for item := range s {
		list[i] = item
		i++
	}

	return list
}

// // foreach iterates over the items in the set and calls the provided function for each set member.
// // the iteration continues until all items in the set have been visited or the closure returns false.
// func (s set) foreach(callback func(item interface{}) bool) {
// 	for item := range s {
// 		if callback(item) {
// 			break
// 		}
// 	}
// }

// // merge merges the current set with another set.
// // It is basically the implementation of the set union between 2 sets.
// func (s set) merge(t set) {
// 	for item := range t {
// 		s[item] = keyExists
// 	}
// }

// // separate removes the set items containing in t from set s.
// // Does not undo merge!!
// func (s set) separate(t set) {
// 	s.remove(t.list()...)
// }

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

// difference returns a new set that contains items which are in the first set but not in the others.
// It precomputes the size of the resulting set based on the number of elements in the input sets.
func difference(sets ...set) set {
	if len(sets) == 0 {
		return newSet()
	}

	totalSize := len(sets[0])

	for i := 1; i < len(sets); i++ {
		totalSize -= len(sets[i])
	}

	if totalSize < 0 {
		totalSize = 0
	}

	resultSet := make(set, totalSize)

	for item := range sets[0] {
		resultSet[item] = keyExists
	}

	for i := 1; i < len(sets); i++ {
		for item := range sets[i] {
			delete(resultSet, item)
		}
	}

	return resultSet
}

// exists checks if a key exists in the Set's records.
func (s *Set) exists(key string) bool {
	_, exist := s.records[key]
	return exist
}

// fieldExists checks if the specified member exists in the set associated with the given key.
// If the key does not exist, it returns false.
func (s *Set) fieldExists(key string, member interface{}) bool {
	if !s.exists(key) {
		return false
	}

	set := s.records[key]
	_, exists := set[member]

	return exists
}

// randomElement returns a random element from the set.
func randomElement(set set) interface{} {
	for k := range set {
		return k
	}
	return nil
}
