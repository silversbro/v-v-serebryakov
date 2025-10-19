package MultivalContext

import (
	"context"
)

type multiCtx struct {
	context.Context
	kv map[any]any
}

func (m *multiCtx) Value(key any) any {
	for k, v := range m.kv {
		if k == key {
			return v
		}
	}
	return m.Context.Value(key)
}

func (m *multiCtx) AddValue(k any, v any) {
	m.kv[k] = v
}

func WithMultivalContext(baseCtx context.Context, kv map[any]any) context.Context {
	mc := multiCtx{}
	mc.Context = baseCtx
	mc.kv = make(map[any]any)
	for k, v := range kv {
		mc.kv[k] = v
	}
	return &mc
}
