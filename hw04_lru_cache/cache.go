package hw04lrucache

type Key string

// Интерфейс для кэша
type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу.
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу.
	Clear()                              // Очистить кэш.
}

// Реализация кэша
type cache struct {
	capacity int               // Емкость кэша
	queue    List              // Очередь последних использованных элементов
	items    map[Key]*ListItem // Словарь для быстрого поиска элементов
}

// Конструктор для создания нового кэша
func NewCache(capacity int) Cache {
	return &cache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Метод для добавления значения в кэш
func (c *cache) Set(key Key, value interface{}) bool {
	// Проверим, есть ли такой ключ в словаре
	item, ok := c.items[key]
	if !ok {
		// Если ключа нет, добавляем новое значение
		newItem := c.queue.PushFront(ListItem{Key: key, Value: value})
		c.items[key] = newItem
		// Если размер очереди превышает емкость, удаляем последний элемент
		if c.queue.Len() > c.capacity {
			lastItem := c.queue.Back()
			delete(c.items, lastItem.Key)
			c.queue.Remove(lastItem)
		}

		return false
	}
	// Обновляем значение существующего элемента и перемещаем его в начало очереди
	item.Value = value
	c.queue.MoveToFront(item)

	return true
}

// Метод для получения значения из кэша
func (c *cache) Get(key Key) (interface{}, bool) {
	// Проверим, есть ли такой ключ в словаре
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	// Переместим элемент в начало очереди
	c.queue.MoveToFront(item)
	return item.Value, true
}

// Метод для очистки кэша
func (c *cache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
	c.capacity = 0
}
