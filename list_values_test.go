package raphanus

import (
	"strings"
	"testing"
)

func Test_ListMethods01(t *testing.T) {
	raph := New("", 0)

	_ = raph.SetInt("key_int", 7, 0)
	if err := raph.SetList("key", []string{"value", "2"}, 1); err != nil {
		t.Errorf("SetList failed: %v", err)
	}

	if err := raph.SetList(" ", []string{"value", "2"}, 0); err == nil {
		t.Errorf("SetList validate key failed")
	}

	if _, err := raph.GetList("key_int"); err == nil {
		t.Errorf("GetList check type failed")
	}

	val, err := raph.GetList("key")
	if err != nil {
		t.Errorf("GetList failed: %v", err)
	}

	if strings.Join(val, "/") != strings.Join([]string{"value", "2"}, "/") {
		t.Errorf("List not equal, got: %v, expected: %v", val, []string{"value", "2"})
	}

	if _, err := raph.GetList("key_fake"); err == nil {
		t.Error("GetList check exists key failed")
	}
}

func Test_ListMethods02(t *testing.T) {
	raph := New("", 0)

	_ = raph.SetInt("key_int", 7, 0)
	if err := raph.SetList("key", []string{"value", "2"}, 0); err != nil {
		t.Errorf("SetList failed: %v", err)
	}

	if err := raph.SetListItem("key_fake", 1, "3"); err == nil {
		t.Errorf("SetListItem want error")
	}
	if err := raph.SetListItem("key_int", 1, "3"); err == nil {
		t.Errorf("SetListItem check type failed")
	}
	if err := raph.SetListItem("key", -1, "3"); err == nil {
		t.Errorf("SetListItem check index failed")
	}
	if err := raph.SetListItem("key", 10000, "3"); err == nil {
		t.Errorf("SetListItem check index failed")
	}

	if err := raph.SetListItem("key", 1, "3"); err != nil {
		t.Errorf("SetListItem error: %v", err)
	}

	val, err := raph.GetList("key")
	if err != nil {
		t.Errorf("GetList failed: %v", err)
	}

	if strings.Join(val, "/") != strings.Join([]string{"value", "3"}, "/") {
		t.Errorf("List not equal, got: %v, expected: %v", val, []string{"value", "3"})
	}

	if err = raph.UpdateList("key_fake", []string{"value", "5"}); err == nil {
		t.Errorf("UpdateList want error")
	}

	if err = raph.UpdateList("key", []string{"value", "5"}); err != nil {
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

	if _, err := raph.GetListItem("key", -1); err == nil {
		t.Errorf("GetListItem check index failed")
	}
	if _, err := raph.GetListItem("key", 100000); err == nil {
		t.Errorf("GetListItem check index failed")
	}
	if _, err := raph.GetListItem("key_int", 1); err == nil {
		t.Errorf("GetListItem check type failed")
	}
	if _, err := raph.GetListItem("key_fake", 1); err == nil {
		t.Errorf("GetListItem want error")
	}
}
