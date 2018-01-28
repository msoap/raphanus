package raphanus

import (
	"testing"
	"time"

	"github.com/msoap/raphanus/common"
)

func Test_TTL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	raph := New()

	_ = raph.SetInt("key01", 42, 5)
	_ = raph.SetInt("key02", 43, 2)
	_ = raph.SetInt("key03", 44, 1)

	time.Sleep(time.Second + 100*time.Millisecond)

	if _, err := raph.GetInt("key01"); err != nil {
		t.Error("TTL don't work")
	}
	if _, err := raph.GetInt("key02"); err != nil {
		t.Error("TTL don't work")
	}
	if _, err := raph.GetInt("key03"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL don't work")
	}

	time.Sleep(time.Second + 100*time.Millisecond)
	if _, err := raph.GetInt("key01"); err != nil {
		t.Error("TTL don't work")
	}
	if _, err := raph.GetInt("key02"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL don't work")
	}
	if _, err := raph.GetInt("key03"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL don't work")
	}

	time.Sleep(3*time.Second + 100*time.Millisecond)
	if _, err := raph.GetInt("key01"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL don't work")
	}
	if _, err := raph.GetInt("key02"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL don't work")
	}
	if _, err := raph.GetInt("key03"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL don't work")
	}
}
