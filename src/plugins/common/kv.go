package common

type Kv struct {
	Key   string
	Value string
}

type Tag struct {
	KV     *Kv
	Mapper func(s string) string
}

func (t *Tag) GetTagValue(by string) string {
	if t.KV == nil {
		return by
	}
	if t.KV.Value != "" {
		return t.KV.Value
	}
	return t.Mapper(by)
}
