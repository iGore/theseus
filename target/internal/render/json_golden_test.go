package render

import (
	"testing"

	"github.com/davglass/license-checker/internal/ordered"
)

func TestJSONGoldenIndentAndNewline(t *testing.T) {
	c := ordered.NewCollection()
	r := ordered.NewObject()
	r.Set("licenses", "UNKNOWN")
	r.Set("repository", "")
	c.Set("missing@1.0.0", r)
	got, err := JSON(c)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"missing@1.0.0\": {\n    \"licenses\": \"UNKNOWN\",\n    \"repository\": \"\"\n  }\n}\n"
	if got != want {
		t.Fatalf("JSON bytes mismatch\nwant %q\ngot  %q", want, got)
	}
}
