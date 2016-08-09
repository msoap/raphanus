package raphanus

import (
	"reflect"
	"testing"
)

func Test_DictMethods01(t *testing.T) {
	raph := New()

	raph.SetDict("key", DictValue{"value": "v1"}, 0)
	val, err := raph.GetDict("key")
	if err != nil {
		t.Errorf("GetDict failed: %v", err)
	}

	if !reflect.DeepEqual(val, DictValue{"value": "v1"}) {
		t.Errorf("List not equal, got: %v, expected: %v", val, DictValue{"value": "v1"})
	}

	err = raph.UpdateDict("key", DictValue{"k1": "v1"})
	if err != nil {
		t.Errorf("UpdateDict failed: %v", err)
	}

	val, err = raph.GetDict("key")
	if err != nil {
		t.Errorf("GetDict failed: %v", err)
	}

	if !reflect.DeepEqual(val, DictValue{"k1": "v1"}) {
		t.Errorf("List not equal, got: %v, expected: %v", val, DictValue{"k1": "v1"})
	}
}

func Test_DictMethods02(t *testing.T) {
	raph := New()

	raph.SetDict("key", DictValue{"k1": "v1", "k2": "v2"}, 0)

	valStr, err := raph.GetDictItem("key", "k1")
	if err != nil {
		t.Errorf("GetDictItem failed: %v", err)
	}
	if valStr != "v1" {
		t.Errorf("GetDictItem, got %s, expected: %s", valStr, "v1")
	}

	err = raph.SetDictItem("key", "k1", "new_val")
	if err != nil {
		t.Errorf("SetDictItem failed: %v", err)
	}
	valStr, err = raph.GetDictItem("key", "k1")
	if err != nil {
		t.Errorf("GetDictItem failed: %v", err)
	}
	if valStr != "new_val" {
		t.Errorf("GetDictItem, got %s, expected: %s", valStr, "new_val")
	}

	err = raph.RemoveDictItem("key", "k1")
	if err != nil {
		t.Errorf("RemoveDictItem failed: %v", err)
	}
	_, err = raph.GetDictItem("key", "k1")
	if err != ErrDictKeyNotExists {
		t.Errorf("Not error after RemoveDictItem: %v", err)
	}
}
