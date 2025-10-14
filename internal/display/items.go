package display

type items struct {
	values []interface{}
	getKey func(obj interface{}) string
}

func (items *items) ItemString(i int) string {
	return items.getKey(items.values[i])
}

func (items *items) Len() int {
	return len(items.values)
}
