package raphanus

import (
	"reflect"
	"testing"

	"github.com/msoap/raphanus/common"
)

func Test_DictMethods01(t *testing.T) {
	raph := New("", 0)

	if err := raph.SetDict("key", raphanuscommon.DictValue{"value": "v1"}, 0); err != nil {
		t.Errorf("SetDict failed: %v", err)
	}

	val, err := raph.GetDict("key")
	if err != nil {
		t.Errorf("GetDict failed: %v", err)
	}

	if !reflect.DeepEqual(val, raphanuscommon.DictValue{"value": "v1"}) {
		t.Errorf("List not equal, got: %v, expected: %v", val, raphanuscommon.DictValue{"value": "v1"})
	}

	err = raph.UpdateDict("key_fake", raphanuscommon.DictValue{"k1": "v1"})
	if err == nil {
		t.Errorf("UpdateDict fake key failed")
	}

	err = raph.UpdateDict("key", raphanuscommon.DictValue{"k1": "v1"})
	if err != nil {
		t.Errorf("UpdateDict failed: %v", err)
	}

	val, err = raph.GetDict("key")
	if err != nil {
		t.Errorf("GetDict failed: %v", err)
	}

	if !reflect.DeepEqual(val, raphanuscommon.DictValue{"k1": "v1"}) {
		t.Errorf("List not equal, got: %v, expected: %v", val, raphanuscommon.DictValue{"k1": "v1"})
	}

	_, err = raph.GetDict("key_fake")
	if err == nil {
		t.Errorf("GetDict not exists key failed")
	}

	_ = raph.SetInt("key_int", 33, 0)
	_, err = raph.GetDict("key_int")
	if err == nil {
		t.Errorf("GetDict check type failed")
	}
}

func Test_DictMethods02(t *testing.T) {
	raph := New("", 0)

	if err := raph.SetDict("key", raphanuscommon.DictValue{"k1": "v1", "k2": "v2"}, 0); err != nil {
		t.Errorf("SetDict failed: %v", err)
	}

	valStr, err := raph.GetDictItem("key", "k1")
	if err != nil {
		t.Errorf("GetDictItem failed: %v", err)
	}
	if valStr != "v1" {
		t.Errorf("GetDictItem, got %s, expected: %s", valStr, "v1")
	}

	if err = raph.SetDictItem("key", "k1", "new_val"); err != nil {
		t.Errorf("SetDictItem failed: %v", err)
	}

	if err = raph.SetDictItem("key_fake", "", "new_val"); err == nil {
		t.Errorf("SetDictItem want error")
	}

	valStr, err = raph.GetDictItem("key", "k1")
	if err != nil {
		t.Errorf("GetDictItem failed: %v", err)
	}
	if valStr != "new_val" {
		t.Errorf("GetDictItem, got %s, expected: %s", valStr, "new_val")
	}

	err = raph.RemoveDictItem("key_fake", "")
	if err == nil {
		t.Errorf("RemoveDictItem want error")
	}

	err = raph.RemoveDictItem("key", "k1")
	if err != nil {
		t.Errorf("RemoveDictItem failed: %v", err)
	}
	_, err = raph.GetDictItem("key", "k1")
	if err != raphanuscommon.ErrDictKeyNotExists {
		t.Errorf("Not error after RemoveDictItem: %v", err)
	}
}

func Test_validateDictParams(t *testing.T) {
	raph := New("", 0)
	if err := raph.SetDict("key", raphanuscommon.DictValue{"k1": "v1", "k2": "v2"}, 0); err != nil {
		t.Errorf("SetDict failed: %v", err)
	}
	_ = raph.SetStr("key_str", "value", 0)

	if err := raph.validateDictParams("key", "k1"); err != nil {
		t.Errorf("validateDictParams failed: %s", err)
	}

	if err := raph.validateDictParams("key", "k1_fake"); err == nil {
		t.Errorf("1. validateDictParams want error")
	}

	if err := raph.validateDictParams("key_fake", "k1"); err == nil {
		t.Errorf("2. validateDictParams want error")
	}

	if err := raph.validateDictParams("", "k1"); err == nil {
		t.Errorf("3. validateDictParams want error")
	}

	if err := raph.validateDictParams("key", ""); err == nil {
		t.Errorf("4. validateDictParams want error")
	}

	if err := raph.validateDictParams("key_str", "k1"); err == nil {
		t.Errorf("5. validateDictParams want error")
	}
}
