package hw04lrucache

type List interface {
	Len() int                          // Длина списка
	Front() *ListItem                  // Первый элемент списка
	Back() *ListItem                   // Последний элемент списка
	PushFront(v interface{}) *ListItem // Добавление значения в начало
	PushBack(v interface{}) *ListItem  // Добавление значения в конец
	Remove(i *ListItem)                // Удаление элемента
	MoveToFront(i *ListItem)           // Перемещение элемента в начало
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem // Указатель на голову списка
	tail *ListItem // Указатель на хвост списка
	len  int       // Длина списка
}

func NewList() List {
	return &list{
		head: nil,
		tail: nil,
		len:  0,
	}
}

func (list *list) Len() int {
	return list.len
}

func (list *list) Front() *ListItem {
	return list.head
}

func (list *list) Back() *ListItem {
	return list.tail
}

func (list *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  list.head,
		Prev:  nil,
	}

	if list.head != nil {
		list.head.Prev = newItem
	}

	list.head = newItem

	if list.tail == nil {
		list.tail = newItem
	}

	list.len++

	return newItem
}

func (list *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  list.tail,
	}

	if list.tail != nil {
		list.tail.Next = newItem
	}

	list.tail = newItem

	if list.head == nil {
		list.head = newItem
	}

	list.len++

	return newItem
}

func (list *list) Remove(item *ListItem) {
	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		list.head = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		list.tail = item.Prev
	}

	list.len--
}

func (list *list) MoveToFront(item *ListItem) {
	if item == list.head {
		return
	}

	list.Remove(item)
	list.PushFront(item.Value)
}
