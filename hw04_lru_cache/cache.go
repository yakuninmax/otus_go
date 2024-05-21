package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	// New cache item.
	newItem := cacheItem{key, value}

	// Check if item exists.
	_, exists := l.items[key]

	if exists {
		// If item exists, update and move to front.
		l.queue.MoveToFront(l.items[key])

		l.queue.Front().Value = newItem
	} else {
		// If item not exists, add it to front.
		l.queue.PushFront(newItem)

		l.items[key] = l.queue.Front()

		// If capacity reached, remove last item.
		if l.queue.Len() > l.capacity {
			delete(l.items, l.queue.Back().Value.(cacheItem).key)

			l.queue.Remove(l.queue.Back())
		}
	}

	l.items[key] = l.queue.Front()

	return exists
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	var itemValue interface{}

	// Check if item exists.
	_, exists := l.items[key]

	// If item exists, move to front.
	if exists {
		l.queue.MoveToFront(l.items[key])

		l.items[key] = l.queue.Front()

		itemValue = l.items[key].Value.(cacheItem).value
	}

	return itemValue, exists
}

func (l *lruCache) Clear() {
	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}
