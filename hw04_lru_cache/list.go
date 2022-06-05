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
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFront := &ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.front,
	}
	if l.len != 0 {
		l.front.Prev = newFront
	}
	l.front = newFront
	l.len++
	if l.len == 1 {
		l.back = newFront
	}
	return newFront
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{
		Value: v,
		Prev:  l.back,
		Next:  nil,
	}
	if l.len != 0 {
		l.back.Next = newBack
	}
	l.back = newBack
	l.len++
	if l.len == 1 {
		l.front = newBack
	}
	return newBack
}

func (l *list) Remove(i *ListItem) {
	if l.len == 0 {
		return
	}
	curNext := i.Next
	curPrev := i.Prev
	if curNext != nil {
		curNext.Prev = i.Prev
	}
	if curPrev != nil {
		curPrev.Next = i.Next
	}
	if l.front == i {
		l.front = curNext
	}
	if l.back == i {
		l.back = curPrev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	curNext := i.Next
	curPrev := i.Prev
	if curNext != nil {
		curNext.Prev = curPrev
	}
	if curPrev != nil {
		curPrev.Next = curNext
	}
	if l.back == i && curPrev != nil {
		l.back = curPrev
	}
	if l.front != nil && l.front != i {
		l.front.Prev = i
		i.Next = l.front
	}
	i.Prev = nil
	l.front = i
}

func NewList() List {
	a := new(list)
	return a
}
