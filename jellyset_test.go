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
	set.SAdd("myset", "member1", "member2", "member3", "member4", "member5")

	t.Run("Popping elements from the existing set", func(t *testing.T) {
		popped := set.SPop("myset", 3)
		expected := []interface{}{"member1", "member2", "member3"}

		if !reflect.DeepEqual(popped, expected) {
			t.Errorf("Expected to pop %#v, but got %#v", expected, popped)
		}
	})

	t.Run("Popping elements from a non-existing set", func(t *testing.T) {
		popped := set.SPop("nonexistent", 2)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})

	t.Run("Popping 0 elements", func(t *testing.T) {
		popped := set.SPop("myset", 0)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})

	t.Run("Popping -1 elements", func(t *testing.T) {
		popped := set.SPop("myset", -1)
		if len(popped) != 0 {
			t.Errorf("Expected to pop 0 elements, but got %d", len(popped))
		}
	})

}
