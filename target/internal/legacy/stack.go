package legacy

import "sync"

// CompleteFn is the callback returned by Add.
//
// It records one completion slot in the stack.
type CompleteFn func(err error, args ...any)

// ResultTransform maps callback payload values into one result entry.
//
// This mirrors the JavaScript behavior where add(fn) can transform positional
// callback arguments before aggregation.
type ResultTransform func(args ...any) any

// DoneCallback receives aggregated completion state.
//
// errors is nil when no completion captured a non-nil error.
// results are ordered by Add call position.
type DoneCallback func(errors []error, results []any, data any)

// Stack is a compatibility utility that aggregates callback completions.
//
// It retains the observable lifecycle contract from the legacy stack utility:
// Add allocates completion slots, Test checks completion state, and Done
// registers the final callback/context.
type Stack struct {
	mu       sync.Mutex
	errors   []error
	results  []any
	finished int
	total    int
	done     DoneCallback
	data     any
	fired    bool
}

// NewStack creates an empty stack.
func NewStack() *Stack {
	return &Stack{}
}

// Add registers one completion slot and returns its completion function.
func (s *Stack) Add(transform ResultTransform) CompleteFn {
	s.mu.Lock()
	index := s.total
	s.total++
	s.errors = append(s.errors, nil)
	s.results = append(s.results, nil)
	s.mu.Unlock()

	return func(err error, args ...any) {
		s.mu.Lock()
		if s.fired {
			s.mu.Unlock()
			return
		}

		s.errors[index] = err
		s.results[index] = applyTransform(transform, args...)
		s.finished++
		s.mu.Unlock()

		s.Test()
	}
}

func applyTransform(transform ResultTransform, args ...any) any {
	if transform != nil {
		return transform(args...)
	}

	switch len(args) {
	case 0:
		return nil
	case 1:
		return args[0]
	default:
		copyArgs := make([]any, len(args))
		copy(copyArgs, args)
		return copyArgs
	}
}

// Done registers the final callback and context payload.
func (s *Stack) Done(callback DoneCallback, data any) {
	s.mu.Lock()
	s.done = callback
	s.data = data
	s.mu.Unlock()

	s.Test()
}

// Test checks completion state and fires the done callback once complete.
func (s *Stack) Test() {
	s.mu.Lock()
	if s.fired || s.done == nil || s.finished != s.total {
		s.mu.Unlock()
		return
	}

	s.fired = true
	errors := cloneErrorsOrNil(s.errors)
	results := cloneResults(s.results)
	data := s.data
	done := s.done
	s.mu.Unlock()

	done(errors, results, data)
}

func cloneErrorsOrNil(errors []error) []error {
	if len(errors) == 0 {
		return nil
	}

	anyError := false
	out := make([]error, len(errors))
	for i := range errors {
		out[i] = errors[i]
		if errors[i] != nil {
			anyError = true
		}
	}

	if !anyError {
		return nil
	}

	return out
}

func cloneResults(results []any) []any {
	out := make([]any, len(results))
	copy(out, results)
	return out
}
