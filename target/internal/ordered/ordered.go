package ordered

import (
	"bytes"
	"encoding/json"
	"sort"
)

// Record preserves JavaScript-style insertion order for observable fields.
type Record struct {
	keys []string
	vals map[string]any
}

func NewRecord() *Record { return &Record{vals: make(map[string]any)} }
func (r *Record) Set(k string, v any) {
	if _, ok := r.vals[k]; !ok {
		r.keys = append(r.keys, k)
	}
	r.vals[k] = v
}
func (r *Record) Get(k string) (any, bool) { v, ok := r.vals[k]; return v, ok }
func (r *Record) String(k string) string {
	if v, ok := r.vals[k]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
func (r *Record) Bool(k string) bool {
	if v, ok := r.vals[k]; ok {
		b, _ := v.(bool)
		return b
	}
	return false
}
func (r *Record) Keys() []string { out := append([]string(nil), r.keys...); return out }
func (r *Record) Delete(k string) {
	if _, ok := r.vals[k]; !ok {
		return
	}
	delete(r.vals, k)
	out := r.keys[:0]
	for _, x := range r.keys {
		if x != k {
			out = append(out, x)
		}
	}
	r.keys = out
}
func (r *Record) Len() int { return len(r.keys) }
func (r *Record) Clone() *Record {
	n := NewRecord()
	for _, k := range r.keys {
		n.Set(k, r.vals[k])
	}
	return n
}

func (r *Record) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range r.keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		vb, err := marshalNoEscape(r.vals[k])
		if err != nil {
			return nil, err
		}
		b.Write(kb)
		b.WriteByte(':')
		b.Write(vb)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func marshalNoEscape(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(b.Bytes(), []byte("\n")), nil
}

// Map preserves observable top-level package order.
type Map struct {
	keys []string
	vals map[string]*Record
}

func NewMap() *Map { return &Map{vals: make(map[string]*Record)} }
func (m *Map) Set(k string, v *Record) {
	if _, ok := m.vals[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.vals[k] = v
}
func (m *Map) Get(k string) (*Record, bool) { v, ok := m.vals[k]; return v, ok }
func (m *Map) Delete(k string) {
	if _, ok := m.vals[k]; !ok {
		return
	}
	delete(m.vals, k)
	out := m.keys[:0]
	for _, x := range m.keys {
		if x != k {
			out = append(out, x)
		}
	}
	m.keys = out
}
func (m *Map) Keys() []string { out := append([]string(nil), m.keys...); return out }
func (m *Map) Len() int       { return len(m.keys) }
func (m *Map) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range m.keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		vb, err := marshalNoEscape(m.vals[k])
		if err != nil {
			return nil, err
		}
		b.Write(kb)
		b.WriteByte(':')
		b.Write(vb)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}
func SortedKeys[V any](m map[string]V) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
