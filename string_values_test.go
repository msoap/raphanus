package raphanus

import (
	"testing"

	"github.com/msoap/raphanus/common"
)

func Test_StringMethods(t *testing.T) {
	raph := New("", 0)

	raph.SetInt("key_int", 10, 0)
	raph.SetStr("key", "value", 0)

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
