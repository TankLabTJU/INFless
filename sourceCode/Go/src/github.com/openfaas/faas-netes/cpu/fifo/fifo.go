package fifo

import (
	"sync"
	"errors"
	// "strconv"
	// "fmt"
	// "time"
	// "math/rand"
)

// ErrFIFOClosed used when FIFO is closed
var ErrFIFOClosed = errors.New("DeltaFIFO: manipulating with closed queue")

type KeyFunc func(obj interface{}) (string,error)

type KeyError struct {
	Obj interface{}
	Err error
}

func (e KeyError) Error() string {
	return e.Err.Error()
}


type FIFO struct {
	lock sync.RWMutex
	cond sync.Cond

	items map[string]interface{}
	queue []string

	// populated is true if the first batch of items inserted by Replace() has been populated
	// or Delete/Add/Update was called first.
	populated bool
	// initialPopulationCount is the number of items inserted by the first call of Replace()
	initialPopulationCount int

	// keyFunc is used to make the key used for queued item insertion and retrieval, and
	// should be deterministic.
	keyFunc KeyFunc

	closed     bool
	closedLock sync.Mutex
}


// Close the queue.
func (f *FIFO) Close() {
	f.closedLock.Lock()
	defer f.closedLock.Unlock()
	f.closed = true
	f.cond.Broadcast()
}

// HasSynced returns true if an Add/Update/Delete/AddIfNotPresent are called first,
// or an Update called first but the first batch of items inserted by Replace() has been popped
func (f *FIFO) HasSynced() bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.populated && f.initialPopulationCount == 0
}

// Add inserts an item, and puts it in the queue. The item is only enqueued
// if it doesn't already exist in the set.
func (f *FIFO) Add(obj interface{}) error {
	id, err := f.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.populated = true
	if _, exists := f.items[id]; !exists {
		f.queue = append(f.queue, id)
	}
	f.items[id] = obj
	f.cond.Broadcast()
	return nil
}

// AddIfNotPresent inserts an item, and puts it in the queue. If the item is already
// present in the set, it is neither enqueued nor added to the set.
//
// This is useful in a single producer/consumer scenario so that the consumer can
// safely retry items without contending with the producer and potentially enqueueing
// stale items.
func (f *FIFO) AddIfNotPresent(obj interface{}) error {
	id, err := f.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.addIfNotPresent(id, obj)
	return nil
}

// addIfNotPresent assumes the fifo lock is already held and adds the provided
// item to the queue under id if it does not already exist.
func (f *FIFO) addIfNotPresent(id string, obj interface{}) {
	f.populated = true
	if _, exists := f.items[id]; exists {
		return
	}

	f.queue = append(f.queue, id)
	f.items[id] = obj
	f.cond.Broadcast()
}

// Update is the same as Add in this implementation.
func (f *FIFO) Update(obj interface{}) error {
	return f.Add(obj)
}

// Delete removes an item. It doesn't add it to the queue, because
// this implementation assumes the consumer only cares about the objects,
// not the order in which they were created/added.
func (f *FIFO) Delete(obj interface{}) error {
	id, err := f.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.populated = true
	delete(f.items, id)
	return err
}

// List returns a list of all the items.
func (f *FIFO) List() []interface{} {
	f.lock.RLock()
	defer f.lock.RUnlock()
	list := make([]interface{}, 0, len(f.items))
	for _, item := range f.items {
		list = append(list, item)
	}
	return list
}

// ListKeys returns a list of all the keys of the objects currently
// in the FIFO.
func (f *FIFO) ListKeys() []string {
	f.lock.RLock()
	defer f.lock.RUnlock()
	list := make([]string, 0, len(f.items))
	for key := range f.items {
		list = append(list, key)
	}
	return list
}

// Get returns the requested item, or sets exists=false.
func (f *FIFO) Get(obj interface{}) (item interface{}, exists bool, err error) {
	key, err := f.keyFunc(obj)
	if err != nil {
		return nil, false, KeyError{obj, err}
	}
	return f.GetByKey(key)
}

// GetByKey returns the requested item, or sets exists=false.
func (f *FIFO) GetByKey(key string) (item interface{}, exists bool, err error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	item, exists = f.items[key]
	return item, exists, nil
}

// IsClosed checks if the queue is closed
func (f *FIFO) IsClosed() bool {
	f.closedLock.Lock()
	defer f.closedLock.Unlock()
	if f.closed {
		return true
	}
	return false
}

// Pop waits until an item is ready and processes it. If multiple items are
// ready, they are returned in the order in which they were added/updated.
// The item is removed from the queue (and the store) before it is processed,
// so if you don't successfully process it, it should be added back with
// AddIfNotPresent(). process function is called under lock, so it is safe
// update data structures in it that need to be in sync with the queue.
func (f *FIFO) Pop() (interface{}, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	for {
		for len(f.queue) == 0 {
			// When the queue is empty, invocation of Pop() is blocked until new item is enqueued.
			// When Close() is called, the f.closed is set and the condition is broadcasted.
			// Which causes this loop to continue and return from the Pop().
			if f.IsClosed() {
				return nil, ErrFIFOClosed
			}

			f.cond.Wait()
		}
		id := f.queue[0]
		f.queue = f.queue[1:]
		if f.initialPopulationCount > 0 {
			f.initialPopulationCount--
		}
		item, ok := f.items[id]
		if !ok {
			// Item may have been deleted subsequently.
			continue
		}
		delete(f.items, id)
		return item, nil
	}
}

//Pop() without lock
func (f *FIFO) pop() (interface{}, error) {
	for {
		id := f.queue[0]
		f.queue = f.queue[1:]
		if f.initialPopulationCount > 0 {
			f.initialPopulationCount--
		}
		item, ok := f.items[id]
		if !ok {
			// Item may have been deleted subsequently.
			continue
		}
		delete(f.items, id)
		return item, nil
	}
}

func (f *FIFO) PopAll() ([]interface{}, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	for len(f.queue) == 0 {
		// When the queue is empty, invocation of Pop() is blocked until new item is enqueued.
		// When Close() is called, the f.closed is set and the condition is broadcasted.
		// Which causes this loop to continue and return from the Pop().
		if f.IsClosed() {
			return nil, ErrFIFOClosed
		}

		f.cond.Wait()
	}
	var items []interface{}
	for len(f.queue) > 0 {
		item,err := f.pop()
		if err != nil {
			return []interface{}{},err 
		}
		items = append(items,item)
	}
	return items,nil
}

// Replace will delete the contents of 'f', using instead the given map.
// 'f' takes ownership of the map, you should not reference the map again
// after calling this function. f's queue is reset, too; upon return, it
// will contain the items in the map, in no particular order.
func (f *FIFO) Replace(list []interface{}, resourceVersion string) error {
	items := make(map[string]interface{}, len(list))
	for _, item := range list {
		key, err := f.keyFunc(item)
		if err != nil {
			return KeyError{item, err}
		}
		items[key] = item
	}

	f.lock.Lock()
	defer f.lock.Unlock()

	if !f.populated {
		f.populated = true
		f.initialPopulationCount = len(items)
	}

	f.items = items
	f.queue = f.queue[:0]
	for id := range items {
		f.queue = append(f.queue, id)
	}
	if len(f.queue) > 0 {
		f.cond.Broadcast()
	}
	return nil
}


// NewFIFO returns a Store which can be used to queue up items to
// process.
func NewFIFO(keyFunc KeyFunc) *FIFO {
	f := &FIFO{
		items:   map[string]interface{}{},
		queue:   []string{},
		keyFunc: keyFunc,
	}
	f.cond.L = &f.lock
	return f
}