package jellyset_test

import (
	"testing"

	"github.com/davidandw190/jellyset"
)

// Helper function to assert that a slice is empty.
func assertEmptySlice(t *testing.T, slice []interface{}) {
	t.Helper()
	if len(slice) != 0 {
		t.Errorf("Expected an empty slice, but got %v", slice)
	}
}

// Helper function to assert that two slices are equal.
func assertSlicesEqual(t *testing.T, slice1, slice2 []interface{}) {
	t.Helper()
	if len(slice1) != len(slice2) {
		t.Errorf("Slices are of different lengths. Expected %v, but got %v", slice2, slice1)
		return
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			t.Errorf("Slices differ at index %d. Expected %v, but got %v", i, slice2[i], slice1[i])
		}
	}
}

// Helper function to assert that a set has the expected size.
func assertSetSize(t *testing.T, s *jellyset.Set, key string, expectedSize int) {
	t.Helper()
	if size := s.SCard(key); size != expectedSize {
		t.Errorf("Expected the size of set %s to be %d, but got %d", key, expectedSize, size)
	}
}

// Helper function to assert that a count is equal to the expected value.
func assertCountEqual(t *testing.T, count, expected int) {
	t.Helper()
	if count != expected {
		t.Errorf("Expected count to be %d, but got %d", expected, count)
	}
}

func assertSlicesEqualIgnoreOrder(t *testing.T, expected, actual []interface{}, message string) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Errorf("%s: Slices have different lengths. Expected %d, got %d.", message, len(expected), len(actual))
		return
	}

	expectedMap := make(map[interface{}]int)
	actualMap := make(map[interface{}]int)

	for _, item := range expected {
		expectedMap[item]++
	}

	for _, item := range actual {
		actualMap[item]++
	}

	for key, expectedCount := range expectedMap {
		if actualCount, exists := actualMap[key]; !exists || actualCount != expectedCount {
			t.Errorf("%s: Expected slice does not match the actual slice.", message)
			return
		}
	}
}

func TestSet_SAdd(t *testing.T) {
	set := jellyset.New()

	t.Run("Add to New Set", func(t *testing.T) {
		// Test adding elements to a new set.
		// It checks if the correct number of elements was added.
		count := set.SAdd("myset", "member1", "member2", "member3")
		assertCountEqual(t, count, 3)
	})

	t.Run("Add to Existing Set", func(t *testing.T) {
		// Test adding elements to an existing set.
		// It checks if the correct number of elements was added.
		count := set.SAdd("myset", "member3", "member4", "member5")
		assertCountEqual(t, count, 2)
	})

	t.Run("Add to Multiple Sets", func(t *testing.T) {
		// Test adding elements to multiple sets and counting them.
		// It ensures that the correct number of elements was added to each set.
		set.SAdd("myset1", "member1", "member2", "member3")
		set.SAdd("myset2", "member3", "member4", "member5")

		count1 := set.SAdd("myset1", "member4", "member5")
		count2 := set.SAdd("myset2", "member1", "member2")

		assertCountEqual(t, count1, 2)
		assertCountEqual(t, count2, 2)
	})
}

func TestSet_SPop(t *testing.T) {
	set := jellyset.New()

	t.Run("Pop from Existing Set", func(t *testing.T) {
		// Test popping elements from an existing set.
		// It ensures that the correct number of elements is popped and returns the expected elements.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", 3)
		expected := []interface{}{"member1", "member2", "member3"}

		if len(popped) != len(expected) {
			t.Errorf("Expected to pop %d, but got %d", len(expected), len(popped))
		}

	})

	t.Run("Pop from Non-Existing Set", func(t *testing.T) {
		// Test popping elements from a non-existing set.
		// It checks if an empty slice is returned when popping from a set that doesn't exist.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("nonexistent", 2)
		assertEmptySlice(t, popped)
	})

	t.Run("Pop 0 Elements", func(t *testing.T) {
		// Test popping 0 elements.
		// It ensures that no elements are popped when requesting 0 elements.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", 0)
		assertEmptySlice(t, popped)
	})

	t.Run("Pop -1 Elements", func(t *testing.T) {
		// Test popping -1 elements.
		// It ensures that no elements are popped when requesting a negative number of elements.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", -1)
		assertEmptySlice(t, popped)
	})
}

func TestSet_SRem(t *testing.T) {
	set := jellyset.New()

	t.Run("Remove Member", func(t *testing.T) {
		// Test removing a member from a set.
		// It checks if the specified member is removed.
		set.SAdd("myset", "member1", "member2", "member3")
		removed := set.SRem("myset", "member2")
		if !removed {
			t.Errorf("Expected to remove 'member2' from the set, but it was not removed")
		}
	})

	t.Run("Remove Non-Existent Member", func(t *testing.T) {
		// Test removing a non-existent member from a set.
		// It ensures that the removal operation doesn't affect the set.
		set.SAdd("myset", "member1", "member3")
		removed := set.SRem("myset", "nonexistent")
		if removed {
			t.Errorf("Expected not to remove 'nonexistent' from the set, but it was removed")
		}
	})

	t.Run("Remove from Non-Existent Set", func(t *testing.T) {
		// Test removing a member from a non-existent set.
		// It checks that removal doesn't occur when the set doesn't exist.
		removed := set.SRem("nonexistent", "member1")
		if removed {
			t.Errorf("Expected not to remove from a non-existent set, but removal occurred")
		}
	})
}

func TestSet_SRandMember(t *testing.T) {
	set := jellyset.New()

	t.Run("Retrieve from Existing Set", func(t *testing.T) {
		// Test retrieving random members from an existing set.
		// It ensures the correct number of random members is retrieved.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
		randomMembers := set.SRandMember("myset", 3)
		if len(randomMembers) != 3 {
			t.Errorf("Expected to retrieve 3 random members, but got %d", len(randomMembers))
		}
	})

	t.Run("Retrieve from Non-Existent Set", func(t *testing.T) {
		// Test retrieving random members from a non-existent set.
		// It ensures that an empty slice is returned when retrieving from a set that doesn't exist.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
		randomMembers := set.SRandMember("nonexistent", 3)
		if len(randomMembers) != 0 {
			t.Errorf("Expected to retrieve 0 random members, but got %d", len(randomMembers))
		}
	})

	t.Run("Retrieve 0 Random Members", func(t *testing.T) {
		// Test retrieving 0 random members.
		// It ensures that no elements are retrieved when requesting 0 random members.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
		randomMembers := set.SRandMember("myset", 0)
		if len(randomMembers) != 0 {
			t.Errorf("Expected to retrieve 0 random members, but got %d", len(randomMembers))
		}
	})
}

func Test_SMove(t *testing.T) {
	set := jellyset.New()

	t.Run("Move Member", func(t *testing.T) {
		// Test moving a member from a source set to a destination set.
		// It checks if the member is successfully moved from the source set to the destination set.
		set.SAdd("sourceSet", "member1", "member2")
		moved := set.SMove("sourceSet", "destSet", "member2")

		if !moved {
			t.Errorf("Expected to move member2 from sourceSet to destSet, but the operation was not successful")
		}
	})

	t.Run("Move Non-Existent Member", func(t *testing.T) {
		// Test moving a non-existent member from a source set to a destination set.
		// It ensures that no movement occurs when the member doesn't exist.
		set.SAdd("sourceSet", "member1", "member2")
		moved := set.SMove("sourceSet", "destSet", "nonexistent")

		if moved {
			t.Errorf("Expected not to move a nonexistent member, but the operation was successful")
		}
	})

	t.Run("Move from Non-Existent Source Set", func(t *testing.T) {
		// Test moving a member from a non-existent source set to a destination set.
		// It ensures that no movement occurs when the source set doesn't exist.
		moved := set.SMove("nonexistentSource", "destSet", "member1")

		if moved {
			t.Errorf("Expected not to move from a nonexistent source set, but the operation was successful")
		}
	})

	t.Run("Move to Non-Existent Dest Set", func(t *testing.T) {
		// Test moving a member from a source set to a non-existent destination set.
		// It checks that the destination set is created, and the member is moved.
		set.SAdd("sourceSet", "member1", "member2")
		moved := set.SMove("sourceSet", "nonexistentDest", "member1")

		if !moved {
			t.Errorf("Expected to move from sourceSet to a new destination set, but the operation was not successful")
		}
	})
}
func TestSet_SCard(t *testing.T) {
	set := jellyset.New()

	t.Run("Count Non-Existent Set", func(t *testing.T) {
		// Test counting elements in a non-existent set.
		// It ensures that the count is 0 for a set that doesn't exist.
		count := set.SCard("nonexistent_set")
		if count != 0 {
			t.Errorf("Expected count to be 0, but got %d", count)
		}
	})

	t.Run("Count Empty Set", func(t *testing.T) {
		// Test counting elements in an empty set.
		// It verifies that the count is 0 for an empty set.
		set.SAdd("empty_set")
		count := set.SCard("empty_set")
		if count != 0 {
			t.Errorf("Expected count to be 0, but got %d", count)
		}
	})

	t.Run("Count Set with Elements", func(t *testing.T) {
		// Test counting elements in a set with multiple members.
		// It ensures the correct count for a set with elements.
		set.SAdd("set_with_elements", "member1", "member2", "member3")
		count := set.SCard("set_with_elements")
		if count != 3 {
			t.Errorf("Expected count to be 3, but got %d", count)
		}
	})

	t.Run("Count Multiple Sets", func(t *testing.T) {
		// Test counting elements in multiple sets.
		// It checks the count of multiple sets with different numbers of members.
		set.SAdd("set1", "a", "b", "c")
		set.SAdd("set2", "c", "d")
		set.SAdd("set3", "d", "e", "f", "g")

		count1 := set.SCard("set1")
		count2 := set.SCard("set2")
		count3 := set.SCard("set3")

		if count1 != 3 {
			t.Errorf("Expected count1 to be 3 but got %d", count1)
		}

		if count2 != 2 {
			t.Errorf("Expected count2 to be 2, but got %d", count2)
		}

		if count3 != 4 {
			t.Errorf("Expected count3 to be 4, but got %d", count3)
		}
	})
}

func TestSet_SMembers(t *testing.T) {
	set := jellyset.New()

	t.Run("Retrieve Members of Non-Existent Set", func(t *testing.T) {
		// Test retrieving members from a non-existent set.
		// It verifies that an empty slice is returned for a set that doesn't exist.
		members := set.SMembers("nonexistent_set")
		if len(members) != 0 {
			t.Errorf("Expected an empty slice, but got %v", members)
		}
	})

	t.Run("Retrieve Members of Empty Set", func(t *testing.T) {
		// Test retrieving members from an empty set.
		// It ensures that an empty slice is returned for an empty set.
		set.SAdd("empty_set")
		members := set.SMembers("empty_set")
		if len(members) != 0 {
			t.Errorf("Expected an empty slice, but got %v", members)
		}
	})

	t.Run("Retrieve Members of Set with Elements", func(t *testing.T) {
		// Test retrieving members from a set with multiple members.
		// It ensures that the correct members are retrieved.
		set.SAdd("set_with_elements", "member1", "member2", "member3")
		members := set.SMembers("set_with_elements")
		expectedMembers := []interface{}{"member1", "member2", "member3"}

		if len(members) != len(expectedMembers) {
			t.Errorf("Expected members in number of %d, but got %d", len(expectedMembers), len(members))
		}
	})

	t.Run("Retrieve Members of Multiple Sets", func(t *testing.T) {
		// Test retrieving members from multiple sets.
		// It checks the members of multiple sets with different members.
		set.SAdd("set1", "a", "b", "c")
		set.SAdd("set2", "c", "d")
		set.SAdd("set3", "d", "e", "f", "g")

		members1 := set.SMembers("set1")
		members2 := set.SMembers("set2")
		members3 := set.SMembers("set3")

		expectedMembers1 := []interface{}{"a", "b", "c"}
		expectedMembers2 := []interface{}{"c", "d"}
		expectedMembers3 := []interface{}{"d", "e", "f", "g"}

		if len(members1) != len(expectedMembers1) {
			t.Errorf("Expected members1 of size %d, but got %d", len(expectedMembers1), len(members1))
		}

		if len(members2) != len(expectedMembers2) {
			t.Errorf("Expected members2 of size %d, but got %d", len(expectedMembers2), len(members2))
		}

		if len(members3) != len(expectedMembers3) {
			t.Errorf("Expected members3 of size %d, but got %d", len(expectedMembers3), len(members3))
		}
	})
}

func TestSet_SUnion(t *testing.T) {
	set := jellyset.New()

	t.Run("Union of Two Non-Existent Sets", func(t *testing.T) {
		result := set.SUnion("nonexistent_set1", "nonexistent_set2")
		assertEmptySlice(t, result)
	})

	t.Run("Union of Non-Existent Set with Empty Set", func(t *testing.T) {
		set.SAdd("empty_set")
		result := set.SUnion("nonexistent_set", "empty_set")
		assertEmptySlice(t, result)
	})

	t.Run("Union of Empty Sets", func(t *testing.T) {
		set.SAdd("empty_set1")
		set.SAdd("empty_set2")
		result := set.SUnion("empty_set1", "empty_set2")
		assertEmptySlice(t, result)
	})

	t.Run("Union of Non-Empty Sets", func(t *testing.T) {
		set.SAdd("set1", "a", "b", "c")
		set.SAdd("set2", "c", "d", "e")
		set.SAdd("set3", "e", "f", "g")

		result := set.SUnion("set1", "set2", "set3")
		expectedResult := []interface{}{"a", "b", "c", "d", "e", "f", "g"}

		assertSlicesEqualIgnoreOrder(t, expectedResult, result, "Union of Non-Empty Sets")
	})

	t.Run("Union of Sets with Duplicates", func(t *testing.T) {
		set.SAdd("set1", "a", "b", "c")
		set.SAdd("set2", "c", "d", "e", "a")
		set.SAdd("set3", "e", "f", "g", "a")

		result := set.SUnion("set1", "set2", "set3")
		expectedResult := []interface{}{"a", "b", "c", "d", "e", "f", "g"}
		assertSlicesEqualIgnoreOrder(t, expectedResult, result, "Union of Sets with Duplicates")
	})

}
