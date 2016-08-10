package raphanus

import (
	"sort"
	"sync"
	"time"
)

type ttlQueueItem struct {
	unixtime int64
	key      *string // point to key in cache
}

type ttlQueue struct {
	queue []ttlQueueItem
	timer *time.Timer
	*sync.Mutex
}

func newTTLQueue() ttlQueue {
	return ttlQueue{
		Mutex: new(sync.Mutex),
	}
}

func (ttlQ ttlQueue) Len() int           { return len(ttlQ.queue) }
func (ttlQ ttlQueue) Swap(i, j int)      { ttlQ.queue[i], ttlQ.queue[j] = ttlQ.queue[j], ttlQ.queue[i] }
func (ttlQ ttlQueue) Less(i, j int) bool { return ttlQ.queue[i].unixtime > ttlQ.queue[j].unixtime }

func ttl2unixtime(ttl int) int64 {
	return time.Now().Add(time.Duration(ttl) * time.Second).Unix()
}

func (ttlQ *ttlQueue) add(item ttlQueueItem) {
	ttlQ.Lock()
	defer ttlQ.Unlock()

	ttlQ.queue = append(ttlQ.queue, item)
	// TODO: use custom sort
	sort.Sort(ttlQ)
}

func (ttlQ *ttlQueue) removeLast(n int) {
	ttlQ.queue = ttlQ.queue[:len(ttlQ.queue)-n]
}

// run - handle ttl queue
func (ttlQ *ttlQueue) run(fn func([]string)) {
	ttlQ.Lock()
	defer ttlQ.Unlock()

	if len(ttlQ.queue) == 0 {
		return
	}
	if ttlQ.timer != nil {
		ttlQ.timer.Stop()
	}

	lastItem := ttlQ.queue[len(ttlQ.queue)-1]
	duration := lastItem.unixtime - time.Now().Unix()
	theSameLastCnt := 1
	keysForDelete := []string{*lastItem.key}
	prevIndex := len(ttlQ.queue) - theSameLastCnt - 1
	for prevIndex >= 0 && lastItem.unixtime == ttlQ.queue[prevIndex].unixtime {
		keysForDelete = append(keysForDelete, *(ttlQ.queue[prevIndex].key))
		theSameLastCnt++
		prevIndex--
	}

	if duration <= 0 {
		// ttl at this time is 0, remove from queue and handle next keys
		ttlQ.removeLast(theSameLastCnt)
		go ttlQ.run(fn)
		return
	}

	ttlQ.timer = time.AfterFunc(time.Duration(duration)*time.Second, func() {
		ttlQ.Lock()
		defer ttlQ.Unlock()

		ttlQ.removeLast(theSameLastCnt)
		go fn(keysForDelete)
		go ttlQ.run(fn) // handle next keys with ttl
	})
}
