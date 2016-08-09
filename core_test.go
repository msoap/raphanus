package raphanus

import (
	"sort"
	"strings"
	"testing"

	"github.com/msoap/raphanus/common"
)

func Test_coreSimple01(t *testing.T) {
	raph := New()

	raph.SetStr("key", "value", 0)
	vStr, err := raph.GetStr("key")
	if err == raphanuscommon.ErrKeyNotExists {
		t.Error("Got ErrKeyNotExists error")
	}

	if vStr != "value" {
		t.Errorf("GetStr:\ngot:      %s\nexpected: %s", vStr, "value")
	}

	raph.SetInt("key01", 7, 0)
	vInt, err := raph.GetInt("key01")
	if err == raphanuscommon.ErrKeyNotExists {
		t.Error("Got ErrKeyNotExists error")
	}

	if vInt != 7 {
		t.Errorf("GetStr:\ngot:      %d\nexpected: %d", vInt, 7)
	}

	if len := raph.Len(); len != 2 {
		t.Errorf("Len() failed:\ngot:      %d\nexpected: %d", len, 2)
	}

	allKeys := raph.Keys()
	sort.Strings(allKeys)
	if keysSorted := strings.Join(allKeys, ","); keysSorted != "key,key01" {
		t.Errorf("Keys() failed:\ngot:      %s\nexpected: %s", keysSorted, "key,key01")
	}

	err = raph.Remove("fake_key")
	if err == nil {
		t.Error("Remove() fake key failed")
	}

	err = raph.Remove("key01")
	if err != nil {
		t.Error("Remove() failed")
	}
	if len := raph.Len(); len != 1 {
		t.Errorf("Len() after remove failed:\ngot:      %d\nexpected: %d", len, 1)
	}
}
