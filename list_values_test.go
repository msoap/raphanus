package raphanus

import (
	"strings"
	"testing"
	"time"
)

func Test_ListMethods01(t *testing.T) {
	raph := New()

	raph.SetList("key", []string{"value", "2"}, 1)
	val, err := raph.GetList("key")
	if err != nil {
		t.Errorf("GetList failed: %v", err)
	}

	if strings.Join(val, "/") != strings.Join([]string{"value", "2"}, "/") {
		t.Errorf("List not equal, got: %v, expected: %v", val, []string{"value", "2"})
	}

	time.Sleep(time.Second + 100*time.Millisecond)
	_, err = raph.GetList("key")
	if err != ErrKeyNotExists {
		t.Error("TTL dont work")
	}
}

func Test_ListMethods02(t *testing.T) {
	raph := New()

	raph.SetList("key", []string{"value", "2"}, 0)
	err := raph.SetListItem("key", 1, "3")
	if err != nil {
		t.Errorf("SetListItem error: %v", err)
	}

	val, err := raph.GetList("key")
	if err != nil {
		t.Errorf("GetList failed: %v", err)
	}

	if strings.Join(val, "/") != strings.Join([]string{"value", "3"}, "/") {
		t.Errorf("List not equal, got: %v, expected: %v", val, []string{"value", "3"})
	}

	err = raph.UpdateList("key", []string{"value", "5"})
	if err != nil {
		t.Errorf("UpdateList failed: %v", err)
	}
	val, err = raph.GetList("key")
	if err != nil {
		t.Errorf("GetList failed: %v", err)
	}
	if strings.Join(val, "/") != strings.Join([]string{"value", "5"}, "/") {
		t.Errorf("List not equal, got: %v, expected: %v", val, []string{"value", "5"})
	}

	valStr, err := raph.GetListItem("key", 1)
	if err != nil {
		t.Errorf("GetListItem failed: %v", err)
	}
	if valStr != "5" {
		t.Errorf("GetListItem failed, got: %s, expected: %s", valStr, "5")
	}
}
