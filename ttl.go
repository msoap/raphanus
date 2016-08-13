package raphanus

import (
	"container/heap"
	"sync"
	"time"
)

type ttlQueueItem struct {
	unixtime int64
	key      *string // point to key in cache
}

type ttlQueue struct {
	queue         []ttlQueueItem // for container/heap
	forNextRemove []ttlQueueItem // elements with the minimum and the same unixtime (for removing)
	timer         *time.Timer
	*sync.Mutex
}

// heap.Interface
func (ttlQ *ttlQueue) Len() int           { return len(ttlQ.queue) }
func (ttlQ *ttlQueue) Swap(i, j int)      { ttlQ.queue[i], ttlQ.queue[j] = ttlQ.queue[j], ttlQ.queue[i] }
func (ttlQ *ttlQueue) Less(i, j int) bool { return ttlQ.queue[i].unixtime < ttlQ.queue[j].unixtime }

func (ttlQ *ttlQueue) Push(value interface{}) {
	ttlQ.queue = append(ttlQ.queue, value.(ttlQueueItem))
}

func (ttlQ *ttlQueue) Pop() interface{} {
	length := len(ttlQ.queue)
	item := ttlQ.queue[length-1]
	ttlQ.queue = ttlQ.queue[0 : length-1]
	return item
}

// newTTLQueue - get TTL queue
func newTTLQueue() ttlQueue {
	ttlQ := ttlQueue{
		Mutex: new(sync.Mutex),
	}
	heap.Init(&ttlQ)
	return ttlQ
}

func ttl2unixtime(ttl int) int64 {
	return time.Now().Add(time.Duration(ttl) * time.Second).Unix()
}

func (ttlQ *ttlQueue) add(item ttlQueueItem) {
	ttlQ.Lock()
	defer ttlQ.Unlock()

	heap.Push(ttlQ, item)
}

// run - handle ttl queue
func (ttlQ *ttlQueue) handle(fn func([]string)) {
	ttlQ.Lock()
	defer ttlQ.Unlock()

	if len(ttlQ.queue) == 0 {
		return
	}

	// if timer is running - stop it
	// and return items for removing back to the queue
	if ttlQ.timer != nil && ttlQ.timer.Stop() && len(ttlQ.forNextRemove) > 0 {
		for _, item := range ttlQ.forNextRemove {
			heap.Push(ttlQ, item)
		}
		ttlQ.forNextRemove = ttlQ.forNextRemove[:0]
	}

	minItem := heap.Pop(ttlQ).(ttlQueueItem)
	duration := minItem.unixtime - time.Now().Unix()
	ttlQ.forNextRemove = append(ttlQ.forNextRemove, minItem)
	for len(ttlQ.queue) > 0 && minItem.unixtime == ttlQ.queue[0].unixtime {
		nextSameItem := heap.Pop(ttlQ).(ttlQueueItem)
		ttlQ.forNextRemove = append(ttlQ.forNextRemove, nextSameItem)
	}

	if duration <= 0 {
		// ttl at this time is 0, handle next keys
		go func() {
			ttlQ.exec(fn)
			ttlQ.handle(fn)
		}()
		return
	}

	ttlQ.timer = time.AfterFunc(time.Duration(duration)*time.Second, func() {
		go func() {
			ttlQ.exec(fn)
			ttlQ.handle(fn)
		}()
	})
}

// exec - remove keys and clear ttlQ.forNextRemove
func (ttlQ *ttlQueue) exec(fn func([]string)) {
	ttlQ.Lock()
	defer ttlQ.Unlock()

	if len(ttlQ.forNextRemove) == 0 {
		return
	}
	keys := []string{}
	for _, item := range ttlQ.forNextRemove {
		keys = append(keys, *item.key)
	}
	ttlQ.forNextRemove = ttlQ.forNextRemove[:0]
	go fn(keys)
}
