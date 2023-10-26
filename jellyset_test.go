package jellyset_test

import (
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
