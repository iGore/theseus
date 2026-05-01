package legacy

import (
	"errors"
	"testing"
)

func TestStackAggregatesResultsByAddOrder(t *testing.T) {
	t.Parallel()

	stack := NewStack()
	first := stack.Add(func(args ...any) any { return args[0] })
	second := stack.Add(func(args ...any) any { return args[0] })

	var gotErrors []error
	var gotResults []any
	var gotData any
	stack.Done(func(errs []error, results []any, data any) {
		gotErrors = errs
		gotResults = results
		gotData = data
	}, "context")

	// Complete in reverse order to verify positional aggregation by Add slot.
	second(nil, "two")
	first(nil, "one")

	if gotErrors != nil {
		t.Fatalf("expected nil errors, got %#v", gotErrors)
	}

	if len(gotResults) != 2 {
		t.Fatalf("expected 2 results, got %d", len(gotResults))
	}

	if gotResults[0] != "one" || gotResults[1] != "two" {
		t.Fatalf("unexpected results ordering: %#v", gotResults)
	}

	if gotData != "context" {
		t.Fatalf("unexpected callback data: %#v", gotData)
	}
}

func TestStackPreservesErrorPositions(t *testing.T) {
	t.Parallel()

	stack := NewStack()
	a := stack.Add(nil)
	b := stack.Add(nil)

	expectedErr := errors.New("boom")
	var gotErrors []error
	stack.Done(func(errs []error, _ []any, _ any) {
		gotErrors = errs
	}, nil)

	a(nil, "ok")
	b(expectedErr)

	if gotErrors == nil {
		t.Fatalf("expected non-nil errors slice")
	}

	if len(gotErrors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(gotErrors))
	}

	if gotErrors[0] != nil {
		t.Fatalf("expected first error nil, got %v", gotErrors[0])
	}

	if !errors.Is(gotErrors[1], expectedErr) {
		t.Fatalf("expected second error %v, got %v", expectedErr, gotErrors[1])
	}
}

func TestStackDoneWithNoAddsCallsImmediately(t *testing.T) {
	t.Parallel()

	stack := NewStack()
	called := 0

	stack.Done(func(errs []error, results []any, data any) {
		called++
		if errs != nil {
			t.Fatalf("expected nil errors, got %#v", errs)
		}
		if len(results) != 0 {
			t.Fatalf("expected empty results, got %#v", results)
		}
		if data != 123 {
			t.Fatalf("expected data 123, got %#v", data)
		}
	}, 123)

	if called != 1 {
		t.Fatalf("expected callback to be called once, got %d", called)
	}
}
