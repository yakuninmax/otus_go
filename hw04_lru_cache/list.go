package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

// Leingth of the list.
func (l *list) Len() int {
	return l.len
}

// First element.
func (l *list) Front() *ListItem {
	return l.front
}

// Last element.
func (l *list) Back() *ListItem {
	return l.back
}

// Insert item to first position of the list.
func (l *list) PushFront(v interface{}) *ListItem {
	// New item.
	newItem := new(ListItem)
	newItem.Value = v

	// Check if list is empty.
	if l.len == 0 {
		l.back = newItem
	} else {
		newItem.Next = l.front
		l.front.Prev = newItem
	}

	// newItem first.
	l.front = newItem

	l.len++

	return newItem
}

// Insert item to last position of the list.
func (l *list) PushBack(v interface{}) *ListItem {
	// New item.
	newItem := new(ListItem)
	newItem.Value = v

	// Check if list is empty.
	if l.len == 0 {
		l.front = newItem
	} else {
		newItem.Prev = l.back
		l.back.Next = newItem
	}

	// newItem last.
	l.back = newItem

	l.len++

	return newItem
}

// Remove item from list.
func (l *list) Remove(i *ListItem) {
	// Get prev and next items of removing item.
	prev := i.Prev
	next := i.Next

	switch {
	// Only one item.
	case prev == nil && next == nil:
		l.front = nil
		l.back = nil
		i = nil

	// First item.
	case prev == nil:
		next.Prev = nil
		l.front = next

	// Last item.
	case next == nil:
		prev.Next = nil
		l.back = prev

	default:
		next.Prev = prev
		prev.Next = next
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	// Check if last item.
	if i.Next == nil {
		l.back = i.Prev
	}

	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return new(list)
}
