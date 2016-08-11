package raphanus

import "testing"

func Test_UpdateInt(t *testing.T) {
	raph := New("", 0)

	raph.SetInt("key", 7, 0)
	if err := raph.UpdateInt("key", 8); err != nil {
		t.Errorf("UpdateInt got error: %v", err)
	}

	if v, err := raph.GetInt("key"); err != nil || v != 8 {
		t.Error("UpdateInt failed")
	}
}

func Test_IncrInt(t *testing.T) {
	raph := New("", 0)

	raph.SetInt("key", 7, 0)
	if err := raph.IncrInt("key_fake"); err == nil {
		t.Errorf("IncrInt got error failed")
	}
	if err := raph.IncrInt("key"); err != nil {
		t.Errorf("IncrInt got error: %v", err)
	}

	if v, err := raph.GetInt("key"); err != nil || v != 8 {
		t.Error("IncrInt failed")
	}
}

func Test_DecrInt(t *testing.T) {
	raph := New("", 0)

	raph.SetInt("key", 7, 0)
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
