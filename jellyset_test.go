package jellyset_test

import (
	"reflect"
	"testing"

	"github.com/davidandw190/jellyset"
)

func TestJellyset_SAdd(t *testing.T) {

	set := jellyset.New()

	t.Run("Adding elements to a new set", func(t *testing.T) {
		count := set.SAdd("myset", "member1", "member2", "member3")
		if count != 3 {
			t.Errorf("Expected to add 3 elements, but got %d", count)
		}
	})

	t.Run("Adding elements to an existing set", func(t *testing.T) {
		count := set.SAdd("myset", "member3", "member4", "member5")
		if count != 2 {
			t.Errorf("Expected to add 2 elements, but got %d", count)
		}
	})

	t.Run("Adding elements to multiple sets", func(t *testing.T) {
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

	t.Run("Popping elements from the existing set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", 3)
		expected := []interface{}{"member1", "member2", "member3"}

		if !reflect.DeepEqual(popped, expected) {
			t.Errorf("Expected to pop %#v, but got %#v", expected, popped)
		}
	})

	t.Run("Popping elements from a non-existing set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("nonexistent", 2)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})

	t.Run("Popping 0 elements", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", 0)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})

	t.Run("Popping -1 elements", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

		popped := set.SPop("myset", -1)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})
}

func TestSet_SRem(t *testing.T) {
	set := jellyset.New()

	t.Run("Removing a member from a set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3")
		removed := set.SRem("myset", "member2")
		if !removed {
			t.Errorf("Expected to remove 'member2' from the set, but it was not removed")
		}
	})

	t.Run("Removing a non-existent member", func(t *testing.T) {
		set.SAdd("myset", "member1", "member3")
		removed := set.SRem("myset", "nonexistent")
		if removed {
			t.Errorf("Expected not to remove 'nonexistent' from the set, but it was removed")
		}
	})

	t.Run("Removing from a non-existent set", func(t *testing.T) {
		removed := set.SRem("nonexistent", "member1")
		if removed {
			t.Errorf("Expected not to remove from a non-existent set, but removal occurred")
		}
	})
}

func TestSet_SRandMember(t *testing.T) {
	set := jellyset.New()

	t.Run("Retrieving from an existing set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
		randomMembers := set.SRandMember("myset", 3)
		if len(randomMembers) != 3 {
			t.Errorf("Expected to retrieve 3 random members, but got %d", len(randomMembers))
		}
	})

	t.Run("Retrieving from a non-existing set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
		randomMembers := set.SRandMember("nonexistent", 3)
		if len(randomMembers) != 0 {
			t.Errorf("Expected to retrieve 0 random members, but got %d", len(randomMembers))
		}
	})

	t.Run("Retrieving 0 random members", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")
		randomMembers := set.SRandMember("myset", 0)
		if len(randomMembers) != 0 {
			t.Errorf("Expected to retrieve 0 random members, but got %d", len(randomMembers))
		}
	})
}

func TestSet_SIsMember(t *testing.T) {
	set := jellyset.New()

	t.Run("Checking for an existing member in a set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3")
		exists := set.SIsMember("myset", "member2")

		if !exists {
			t.Errorf("Expected member2 to exist in the set, but it doesn't")
		}
	})

	t.Run("Checking for a non-existing member in a set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3")
		exists := set.SIsMember("myset", "nonexistent")

		if exists {
			t.Errorf("Expected nonexistent to not exist in the set, but it does")
		}
	})

	t.Run("Checking for a member in a non-existing set", func(t *testing.T) {
		set.SAdd("myset", "member1", "member2", "member3")
		exists := set.SIsMember("nonexistent", "member1")

		if exists {
			t.Errorf("Expected member1 to not exist in the nonexistent set, but it does")
		}
	})
}

func Test_SMove(t *testing.T) {
	set := jellyset.New()

	t.Run("Moving a member from a source set to a destination set", func(t *testing.T) {
		set.SAdd("sourceSet", "member1", "member2")
		moved := set.SMove("sourceSet", "destSet", "member2")

		if !moved {
			t.Errorf("Expected to move member2 from sourceSet to destSet, but the operation was not successful")
		}
	})

	t.Run("Moving a non-existing member from a source set to a destination set", func(t *testing.T) {
		set.SAdd("sourceSet", "member1", "member2")
		moved := set.SMove("sourceSet", "destSet", "nonexistent")

		if moved {
			t.Errorf("Expected not to move a nonexistent member, but the operation was successful")
		}
	})

	t.Run("Moving a member from a non-existing source set to a destination set", func(t *testing.T) {
		set.SAdd("sourceSet", "member1", "member2")
		moved := set.SMove("nonexistentSource", "destSet", "member1")

		if moved {
			t.Errorf("Expected not to move from a nonexistent source set, but the operation was successful")
		}
	})

	t.Run("Moving a member from a source set to a non-existing destination set", func(t *testing.T) {
		set.SAdd("sourceSet", "member1", "member2")
		moved := set.SMove("sourceSet", "nonexistentDest", "member1")

		if !moved {
			t.Errorf("Expected to move from sourceSet to a new destination set, but the operation was not successful")
		}
	})
}
