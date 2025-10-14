package analytics

type KeyValue struct {
	Key   string
	Value float64
}

func NewKeyValue(key string, value float64) *KeyValue {
	return &KeyValue{
		Key:   key,
		Value: value,
	}
}
