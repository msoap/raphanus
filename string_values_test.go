package raphanus

import "testing"

func Test_StringMethods(t *testing.T) {
	raph := New()

	raph.SetStr("key", "value", 0)
	vStr, err := raph.GetStr("key")
	if err == ErrKeyNotExists {
		t.Error("Got ErrKeyNotExists error")
	}

	if vStr != "value" {
		t.Errorf("GetStr:\ngot:      %s\nexpected: %s", vStr, "value")
	}

	_, err = raph.GetStr("key_fake")
	if err != ErrKeyNotExists {
		t.Error("Not got ErrKeyNotExists error")
	}

	err = raph.UpdateStr("key", "new value")
	if err != nil {
		t.Errorf("UpdateStr got error: %v", err)
	}

	err = raph.UpdateStr("key_fake", "new value")
	if err != ErrKeyNotExists {
		t.Error("Not got ErrKeyNotExists error")
	}
}
