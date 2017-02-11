package raphanus

import "testing"

func Test_UpdateInt(t *testing.T) {
	raph := New("", 0)

	if err := raph.SetInt("key", 7, 0); err != nil {
		t.Errorf("SetInt got error: %v", err)
	}

	if err := raph.SetInt("", 7, 0); err == nil {
		t.Errorf("SetInt validate key failed")
	}

	if err := raph.UpdateInt("key_fake", 8); err == nil {
		t.Errorf("UpdateInt want error")
	}
	if err := raph.UpdateInt("key", 8); err != nil {
		t.Errorf("UpdateInt got error: %v", err)
	}

	if v, err := raph.GetInt("key"); err != nil || v != 8 {
		t.Error("UpdateInt failed")
	}

	if v, err := raph.GetInt("key_fake"); err == nil || v != 0 {
		t.Error("GetInt want error")
	}

	if err := raph.SetStr("key_str", "str", 0); err != nil {
		t.Errorf("SetStr got error: %v", err)
	}
	if _, err := raph.GetInt("key_str"); err == nil {
		t.Error("GetInt want error for type")
	}
}

func Test_IncrInt(t *testing.T) {
	raph := New("", 0)

	if err := raph.SetInt("key", 7, 0); err != nil {
		t.Errorf("SetInt got error: %v", err)
	}

	_ = raph.SetStr("key_str", "str", 0)

	if err := raph.IncrInt("key_fake"); err == nil {
		t.Errorf("IncrInt want error")
	}
	if err := raph.IncrInt("key_str"); err == nil {
		t.Errorf("IncrInt type check failed")
	}
	if err := raph.IncrInt("key"); err != nil {
		t.Errorf("IncrInt got error: %v", err)
	}

	if v, err := raph.GetInt("key"); err != nil || v != 8 {
		t.Error("IncrInt failed")
	}

	if v, err := raph.GetInt("key_fake"); err == nil || v != 0 {
		t.Error("IncrInt want error")
	}
}

func Test_DecrInt(t *testing.T) {
	raph := New("", 0)

	if err := raph.SetInt("key", 7, 0); err != nil {
		t.Errorf("SetInt got error: %v", err)
	}

	if err := raph.DecrInt("key_fake"); err == nil {
		t.Errorf("DecrInt got error failed")
	}
	if err := raph.DecrInt("key"); err != nil {
		t.Errorf("DecrInt got error: %v", err)
	}

	if v, err := raph.GetInt("key"); err != nil || v != 6 {
		t.Error("DecrInt failed")
	}
}
