package protobuf_ext

import (
	"github.com/samber/mo"
)

type WrappedValue[T any] interface {
	GetValue() T
}

func WrappedValueToOption[T any](wrappedValue WrappedValue[T]) mo.Option[T] {
	if wrappedValue == nil {
		return mo.None[T]()
	}
	return mo.Some(wrappedValue.GetValue())
}
