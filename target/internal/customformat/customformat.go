package customformat

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"theseus-target/license-checker/internal/cli"
)

func Load(path string) ([]cli.Field, error) {
	if path == "" {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}
func Parse(b []byte) ([]cli.Field, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	var order []string
	dec := json.NewDecoder(bytesReader(b))
	dec.UseNumber()
	var obj orderedDecode
	if err := dec.Decode(&obj); err == nil {
		order = obj.keys
	}
	if len(order) == 0 {
		for k := range raw {
			order = append(order, k)
		}
	}
	fields := make([]cli.Field, 0, len(order))
	for _, k := range order {
		var v any
		_ = json.Unmarshal(raw[k], &v)
		fields = append(fields, cli.Field{Key: k, Value: v})
	}
	return fields, nil
}
func ParseAny(path any) ([]cli.Field, error) {
	s, ok := path.(string)
	if !ok {
		return nil, errors.New("did not specify a path")
	}
	return Load(s)
}

type byteReader []byte

func bytesReader(b []byte) *byteReader { r := byteReader(b); return &r }
func (r *byteReader) Read(p []byte) (int, error) {
	if len(*r) == 0 {
		return 0, io.EOF
	}
	n := copy(p, *r)
	*r = (*r)[n:]
	return n, nil
}

type orderedDecode struct{ keys []string }

func (o *orderedDecode) UnmarshalJSON(b []byte) error {
	dec := json.NewDecoder(bytesReader(b))
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if d, ok := t.(json.Delim); !ok || d != '{' {
		return nil
	}
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return err
		}
		o.keys = append(o.keys, tk.(string))
		var skip any
		if err := dec.Decode(&skip); err != nil {
			return err
		}
	}
	return nil
}
