package stringx

import (
	"strings"
	"sync"
)

var builderPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// Acquire 获取一个可复用的 Builder
func Acquire() *strings.Builder {
	b := builderPool.Get().(*strings.Builder)
	b.Reset()
	return b
}

// Release 释放 Builder 回池中
func Release(b *strings.Builder) {
	builderPool.Put(b)
}

func build(f func(*strings.Builder)) string {
	b := Acquire()
	defer Release(b)
	f(b)
	return b.String()
}

func Build(args ...string) string {
	return build(func(b *strings.Builder) {
		total := 0
		for _, arg := range args {
			total += len(arg)
		}

		b.Grow(total)
		for _, arg := range args {
			b.WriteString(arg)
		}
	})
}
