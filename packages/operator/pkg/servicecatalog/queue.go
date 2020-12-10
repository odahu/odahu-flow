package servicecatalog

import (
	"container/list"
	"context"
	"sync"
	"time"
)

const BackoffFactor = 2

type Queue struct {
	linkedList *list.List
	mu *sync.RWMutex
}

func NewQueue() Queue {
	return Queue{
		linkedList: list.New(),
		mu: &sync.RWMutex{},
	}
}

type Item struct {
	backoff int
	Value   interface{}
}

func (q Queue) Pop() *Item {

	q.mu.RLock()
	el := q.linkedList.Front()
	q.mu.RUnlock()

	if el == nil {
		return nil
	}

	item, _ := el.Value.(Item)

	return &item

}

func (q Queue) push(item *Item) {
	q.mu.Lock()
	q.linkedList.PushBack(item)
	q.mu.Unlock()
}


// Item is added to queue immediately for the first time. Without delay
// For next pushes delay equals to (N-1)*BackoffFactor where
// N - attempt
func (q Queue) Push(ctx context.Context, item *Item) {

	if item.backoff == 0 {
		item.backoff = 2
		q.push(item)
	}

	t := time.NewTimer(time.Duration(item.backoff) * time.Second)

	go func() {
		select {
		case <-t.C:
			item.backoff *= BackoffFactor
			q.push(item)
		case <-ctx.Done():
			return
		}
	}()
}