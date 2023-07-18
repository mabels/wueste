package wueste

import (
	"encoding/json"
	"strconv"
)

type Literal[T any] struct {
	value T
	str   string
}

func (l Literal[T]) Value() T {
	return l.value
}

func (l Literal[T]) String() *string {
	return &l.str
}

func BoolLiteral(v bool) Literal[bool] {
	if v {
		return Literal[bool]{value: v, str: "true"}
	} else {
		return Literal[bool]{value: v, str: "false"}
	}
}

func NumberLiteral[T float32 | float64](v T) Literal[T] {
	return Literal[T]{value: v, str: strconv.FormatFloat(float64(v), 'e', 0, 64)}
}

func IntegerLiteral[T uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64](v T) Literal[T] {
	return Literal[T]{value: v, str: strconv.FormatInt(int64(v), 10)}
}

func StringLiteral(v string) Literal[string] {
	return Literal[string]{value: v, str: QuoteString(v)}
}

func QuoteString(s string) string {
	byteStr, _ := json.Marshal(s)
	return string(byteStr)
}
