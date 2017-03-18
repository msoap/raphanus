package raphanus

import (
	"testing"

	"github.com/msoap/raphanus/common"
)

func Test_BytesMethods(t *testing.T) {
	raph := New("", 0)

	_ = raph.SetInt("key_int", 10, 0)
	if err := raph.SetBytes("key", []byte("value"), 0); err != nil {
		t.Errorf("SetBytes got error: %v", err)
	}

	if err := raph.SetBytes("", []byte("str"), 0); err == nil {
		t.Errorf("SetBytes validate key failed")
	}

	if _, err := raph.GetBytes("key_int"); err == nil {
		t.Error("GetBytes check type failed")
	}

	vBytes, err := raph.GetBytes("key")
	if err == raphanuscommon.ErrKeyNotExists {
		t.Error("Got ErrKeyNotExists error")
	}

	if string(vBytes) != "value" {
		t.Errorf("GetBytes:\ngot:      %s\nexpected: %s", string(vBytes), "value")
	}

	_, err = raph.GetBytes("key_fake")
	if err != raphanuscommon.ErrKeyNotExists {
		t.Error("Not got ErrKeyNotExists error")
	}

	if err = raph.UpdateBytes("key", []byte("new value")); err != nil {
		t.Errorf("UpdateBytes got error: %v", err)
	}

	err = raph.UpdateBytes("key_fake", []byte("new value"))
	if err != raphanuscommon.ErrKeyNotExists {
		t.Error("Not got ErrKeyNotExists error")
	}
}
