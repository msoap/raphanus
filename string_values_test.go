package raphanus

import (
	"testing"

	"github.com/msoap/raphanus/common"
)

func Test_StringMethods(t *testing.T) {
	raph := New()

	_ = raph.SetInt("key_int", 10, 0)
	if err := raph.SetStr("key", "value", 0); err != nil {
		t.Errorf("SetStr got error: %v", err)
	}

	if err := raph.SetStr("", "str", 0); err == nil {
		t.Errorf("SetStr validate key failed")
	}

	if _, err := raph.GetStr("key_int"); err == nil {
		t.Error("GetStr check type failed")
	}

	vStr, err := raph.GetStr("key")
	if err == raphanuscommon.ErrKeyNotExists {
		t.Error("Got ErrKeyNotExists error")
	}

	if vStr != "value" {
		t.Errorf("GetStr:\ngot:      %s\nexpected: %s", vStr, "value")
	}

	_, err = raph.GetStr("key_fake")
	if err != raphanuscommon.ErrKeyNotExists {
		t.Error("Not got ErrKeyNotExists error")
	}

	err = raph.UpdateStr("key", "new value")
	if err != nil {
		t.Errorf("UpdateStr got error: %v", err)
	}

	err = raph.UpdateStr("key_fake", "new value")
	if err != raphanuscommon.ErrKeyNotExists {
		t.Error("Not got ErrKeyNotExists error")
	}
}
