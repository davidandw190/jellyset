package jellyset_test

import (
	"testing"

	"github.com/davidandw190/jellyset"
)

func TestSet_SAdd(t *testing.T) {
	set := jellyset.New()

	t.Run("Add to New Set", func(t *testing.T) {
		// Test adding elements to a new set.
		// It checks if the correct number of elements was added.
		count := set.SAdd("myset", "member1", "member2", "member3")
		if count != 3 {
			t.Errorf("Expected to add 3 elements, but got %d", count)
		}
	})

	t.Run("Add to Existing Set", func(t *testing.T) {
		// Test adding elements to an existing set.
		// It checks if the correct number of elements was added.
		count := set.SAdd("myset", "member3", "member4", "member5")
		if count != 2 {
			t.Errorf("Expected to add 2 elements, but got %d", count)
		}
	})

	t.Run("Add to Multiple Sets", func(t *testing.T) {
		// Test adding elements to multiple sets and counting them.
		// It ensures that the correct number of elements was added to each set.
		set.SAdd("myset1", "member1", "member2", "member3")
		set.SAdd("myset2", "member3", "member4", "member5")

		count1 := set.SAdd("myset1", "member4", "member5")
		count2 := set.SAdd("myset2", "member1", "member2")

		if count1 != 2 {
			t.Errorf("Expected to add 2 elements to myset1, but got %d", count1)
		}

		if count2 != 2 {
			t.Errorf("Expected to add 2 elements to myset2, but got %d", count2)
		}
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
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})

	t.Run("Pop 0 Elements", func(t *testing.T) {
		// Test popping 0 elements.
		// It ensures that no elements are popped when requesting 0 elements.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", 0)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})

	t.Run("Pop -1 Elements", func(t *testing.T) {
		// Test popping -1 elements.
		// It ensures that no elements are popped when requesting a negative number of elements.
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", -1)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
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
