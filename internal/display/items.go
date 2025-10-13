package display

import "time"


type items[T any] struct {
	values []*T
	getKey func (obj *T) string
	getRow func (obj *T, now time.Time) string
}

func (items *items[T]) ItemString(i int) string {
	return items.getKey(items.values[i])
}

func (items *items[T]) Len() int {
	return len(items.values)
}

