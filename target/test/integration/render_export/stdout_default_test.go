package render_export

import (
	"bytes"
	"strings"
	"testing"

	"theseus/target/internal/core"
)

func TestStdoutDefaultWhenOutAbsent(t *testing.T) {
	svc := core.NewRenderService()
	var stdout bytes.Buffer

	_, _, err := svc.Execute(sampleModules(), core.RenderOptions{}, &stdout)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "alpha@1.0.0") {
		t.Fatalf("stdout did not receive rendered output: %q", stdout.String())
	}
}
