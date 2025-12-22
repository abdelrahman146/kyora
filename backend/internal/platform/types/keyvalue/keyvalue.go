package keyvalue

import "github.com/abdelrahman146/kyora/internal/platform/types/schema"

type KeyValue struct {
	Key   any
	Value any
}

func New(key any, value any) KeyValue {
	return KeyValue{
		Key:   key,
		Value: value,
	}
}

func KeysFromKeyValueSlice(kvs []KeyValue) []any {
	keys := make([]any, len(kvs))
	for i, kv := range kvs {
		keys[i] = kv.Key
	}
	return keys
}

func ValuesFromKeyValueSlice(kvs []KeyValue) []any {
	values := make([]any, len(kvs))
	for i, kv := range kvs {
		values[i] = kv.Value
	}
	return values
}

var Schema = struct {
	Key   schema.Field
	Value schema.Field
}{
	Key:   schema.NewField("key", "key"),
	Value: schema.NewField("value", "value"),
}
