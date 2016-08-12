package raphanus

import (
	"sort"
	"testing"
)

type ttlList []ttlQueueItem

func (t ttlList) Len() int           { return len(t) }
func (t ttlList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ttlList) Less(i, j int) bool { return t[i].unixtime > t[j].unixtime }

func Test_sortSortedListWithNewLastItem(t *testing.T) {
	str := "str"
	cases := []struct {
		list []ttlQueueItem
	}{
		{list: []ttlQueueItem{
			{10, &str},
			{7, &str},
			{3, &str},
			{1, &str},
		}},
		{list: []ttlQueueItem{
			{10, &str},
			{7, &str},
			{3, &str},
			{10, &str},
		}},
		{list: []ttlQueueItem{
			{10, &str},
			{7, &str},
			{3, &str},
			{11, &str},
		}},
		{list: []ttlQueueItem{
			{10, &str},
			{7, &str},
			{3, &str},
			{5, &str},
		}},
		{list: []ttlQueueItem{
			{10, &str},
			{7, &str},
			{3, &str},
			{8, &str},
		}},
	}

	for i, item := range cases {
		ttlQueue := newTTLQueue()
		ttlQueue.queue = make(ttlList, len(item.list))
		copy(ttlQueue.queue, item.list)
		ttlQueue.sortSortedListWithNewLastItem()

		sort.Sort(ttlList(item.list))

		if !isEqualLists(ttlQueue.queue, item.list) {
			t.Errorf("%d. sortSortedListWithNewLastItem failed", i)
		}
	}
}

func isEqualLists(l1, l2 []ttlQueueItem) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i := range l1 {
		if l1[i].unixtime != l2[i].unixtime {
			return false
		}
	}

	return true
}
