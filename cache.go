package tmpls

type Cache[T any] interface {
	Get(key string) (T, bool)
	Set(key string, val T)
}

type MapCache[T any] map[string]T

func (m MapCache[T]) Get(key string) (T, bool) {
	v, ok := m[key]
	return v, ok
}

func (m MapCache[T]) Set(key string, val T) {
	m[key] = val
}
