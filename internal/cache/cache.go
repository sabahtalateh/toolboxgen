package cache

type Key interface {
	Key() string
}

type Cache[T Key] struct {
	tt map[string]T
}

func New[T Key]() *Cache[T] {
	return &Cache[T]{tt: map[string]T{}}
}

func (c *Cache[T]) Put(t T) {
	c.tt[t.Key()] = t
}

func (c *Cache[T]) Get(key string) (T, bool) {
	t, ok := c.tt[key]
	return t, ok
}
