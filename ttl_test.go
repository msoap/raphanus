package raphanus

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/msoap/raphanus/common"
)

func Test_TTL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	raph := New("", 0)

	_ = raph.SetInt("key01", 42, 5)
	_ = raph.SetInt("key02", 43, 2)
	_ = raph.SetInt("key03", 44, 1)

	time.Sleep(time.Second + 100*time.Millisecond)

	if _, err := raph.GetInt("key01"); err != nil {
		t.Error("TTL dont work")
	}
	if _, err := raph.GetInt("key02"); err != nil {
		t.Error("TTL dont work")
	}
	if _, err := raph.GetInt("key03"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL dont work")
	}

	time.Sleep(time.Second + 100*time.Millisecond)
	if _, err := raph.GetInt("key01"); err != nil {
		t.Error("TTL dont work")
	}
	if _, err := raph.GetInt("key02"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL dont work")
	}
	if _, err := raph.GetInt("key03"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL dont work")
	}

	time.Sleep(3*time.Second + 100*time.Millisecond)
	if _, err := raph.GetInt("key01"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL dont work")
	}
	if _, err := raph.GetInt("key02"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL dont work")
	}
	if _, err := raph.GetInt("key03"); err != raphanuscommon.ErrKeyNotExists {
		t.Error("TTL dont work")
	}
}

func Test_TTLqueue(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	queue := newTTLQueue()
	k0, k1, k2, k3, k4 := "0", "1", "2", "3", "4"

	queue.add(ttlQueueItem{key: &k1, unixtime: ttl2unixtime(1)})
	queue.add(ttlQueueItem{key: &k4, unixtime: ttl2unixtime(4)})
	queue.add(ttlQueueItem{key: &k2, unixtime: ttl2unixtime(2)})
	queue.add(ttlQueueItem{key: &k3, unixtime: ttl2unixtime(3)})
	queue.add(ttlQueueItem{key: &k0, unixtime: ttl2unixtime(0)})
	queue.add(ttlQueueItem{key: &k2, unixtime: ttl2unixtime(2)})

	mutex := new(sync.Mutex)
	result := []string{}
	queue.handle(func(keys []string) {
		go func() {
			mutex.Lock()
			result = append(result, strings.Join(keys, "/"))
			mutex.Unlock()
		}()
	})

	time.Sleep(4*time.Second + 100*time.Millisecond)

	mutex.Lock()
	if strings.Join(result, ",") != "0,1,2/2,3,4" {
		t.Errorf("ttlQueue failed, got: %s, expected: %s", strings.Join(result, ","), "0,1,2/2,3,4")
	}
	mutex.Unlock()
}
