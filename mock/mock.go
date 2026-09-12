// Copyright 2018 The go-ego Project Developers.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package mock is a simple mock object system for tt.
//
// Embed Mock in a struct that implements the interface you want to fake,
// record expectations with On(...).Return(...), and forward calls with
// Called(...):
//
//	type Store struct{ mock.Mock }
//
//	func (s *Store) Get(id int) (string, error) {
//		args := s.Called(id)
//		return args.String(0), args.Error(1)
//	}
//
//	s := new(Store)
//	s.On("Get", 1).Return("one", nil).Once()
//	v, err := s.Get(1)
//	s.AssertExpectations(t)
//
// The package also provides T, a TestingT that records failures instead of
// reporting them, for testing tt assertions themselves.
package mock

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

const (
	// Anything matches any argument value in On / AssertCalled.
	Anything = "mock.Anything"
)

// TestingT is an interface wrapper around *testing.T,
// compatible with tt.TestingT.
type TestingT interface {
	Errorf(format string, args ...interface{})
}

// T is a TestingT that records failures instead of reporting them.
// It is useful for testing assertions that are expected to fail.
type T struct {
	mu     sync.Mutex
	Errors []string
}

// Errorf records the formatted failure message.
func (t *T) Errorf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Errors = append(t.Errors, fmt.Sprintf(format, args...))
}

// Failed reports whether Errorf was called.
func (t *T) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.Errors) > 0
}

// Reset clears all recorded failures.
func (t *T) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Errors = nil
}

// AnythingOfTypeArgument matches any argument whose type string equals it,
// e.g. AnythingOfType("int") or AnythingOfType("*errors.errorString").
type AnythingOfTypeArgument string

// AnythingOfType returns a matcher for any value of the given type.
func AnythingOfType(t string) AnythingOfTypeArgument {
	return AnythingOfTypeArgument(t)
}

// Matches reports whether arg's type string equals the matcher.
func (a AnythingOfTypeArgument) Matches(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).String() == string(a)
}

// ArgumentMatcher matches an argument with a custom function.
type ArgumentMatcher struct {
	fn reflect.Value
}

// MatchedBy returns a matcher that calls fn with the actual argument.
// fn must be a func with one argument returning a bool.
func MatchedBy(fn interface{}) ArgumentMatcher {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func || v.Type().NumIn() != 1 ||
		v.Type().NumOut() != 1 || v.Type().Out(0).Kind() != reflect.Bool {
		panic(fmt.Sprintf("mock: MatchedBy expects func(T) bool, got %T", fn))
	}
	return ArgumentMatcher{fn: v}
}

// Matches reports whether the matcher accepts the argument.
func (m ArgumentMatcher) Matches(arg interface{}) bool {
	in := m.fn.Type().In(0)
	var v reflect.Value
	if arg == nil {
		switch in.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice,
			reflect.Chan, reflect.Func:
			v = reflect.Zero(in)
		default:
			return false
		}
	} else {
		v = reflect.ValueOf(arg)
		if !v.Type().AssignableTo(in) {
			return false
		}
	}
	return m.fn.Call([]reflect.Value{v})[0].Bool()
}

// Arguments is a list of call or return values.
type Arguments []interface{}

// Get returns the i-th argument, panicking if out of range.
func (a Arguments) Get(i int) interface{} {
	if i < 0 || i >= len(a) {
		panic(fmt.Sprintf("mock: arguments: Get(%d) out of range, have %d", i, len(a)))
	}
	return a[i]
}

// String returns the i-th argument as a string.
func (a Arguments) String(i int) string {
	v, ok := a.Get(i).(string)
	if !ok {
		panic(fmt.Sprintf("mock: arguments: String(%d) is %T, not string", i, a.Get(i)))
	}
	return v
}

// Int returns the i-th argument as an int.
func (a Arguments) Int(i int) int {
	v, ok := a.Get(i).(int)
	if !ok {
		panic(fmt.Sprintf("mock: arguments: Int(%d) is %T, not int", i, a.Get(i)))
	}
	return v
}

// Bool returns the i-th argument as a bool.
func (a Arguments) Bool(i int) bool {
	v, ok := a.Get(i).(bool)
	if !ok {
		panic(fmt.Sprintf("mock: arguments: Bool(%d) is %T, not bool", i, a.Get(i)))
	}
	return v
}

// Error returns the i-th argument as an error, or nil.
func (a Arguments) Error(i int) error {
	v := a.Get(i)
	if v == nil {
		return nil
	}
	err, ok := v.(error)
	if !ok {
		panic(fmt.Sprintf("mock: arguments: Error(%d) is %T, not error", i, v))
	}
	return err
}

// Is reports whether the arguments match the expected objects,
// honouring Anything, AnythingOfType and MatchedBy.
func (a Arguments) Is(objects ...interface{}) bool {
	if len(a) != len(objects) {
		return false
	}
	for i, expect := range a {
		if !match(expect, objects[i]) {
			return false
		}
	}
	return true
}

func match(expect, actual interface{}) bool {
	switch e := expect.(type) {
	case string:
		if e == Anything {
			return true
		}
	case AnythingOfTypeArgument:
		return e.Matches(actual)
	case ArgumentMatcher:
		return e.Matches(actual)
	}
	return objectsAreEqual(expect, actual)
}

func objectsAreEqual(expect, actual interface{}) bool {
	if expect == nil || actual == nil {
		return expect == actual
	}
	eb, eok := expect.([]byte)
	ab, aok := actual.([]byte)
	if eok && aok {
		return string(eb) == string(ab)
	}
	return reflect.DeepEqual(expect, actual)
}

// Call represents an expected or recorded method call.
type Call struct {
	parent *Mock

	// Method is the name of the method.
	Method string
	// Arguments the method is expected to be called with.
	Arguments Arguments
	// ReturnArguments returned by Called for this expectation.
	ReturnArguments Arguments
	// Repeatability is how many times this expectation may be consumed;
	// 0 means unlimited, -1 means exhausted.
	Repeatability int
	// RunFn, if set, is run with the actual arguments on each call.
	RunFn func(Arguments)

	totalCalls int
}

// Return sets the values returned by Called for this expectation.
func (c *Call) Return(returnArguments ...interface{}) *Call {
	c.lock()
	defer c.unlock()
	c.ReturnArguments = returnArguments
	return c
}

// Once marks the expectation as consumable exactly once.
func (c *Call) Once() *Call {
	return c.Times(1)
}

// Twice marks the expectation as consumable exactly twice.
func (c *Call) Twice() *Call {
	return c.Times(2)
}

// Times marks the expectation as consumable exactly i times.
func (c *Call) Times(i int) *Call {
	c.lock()
	defer c.unlock()
	c.Repeatability = i
	return c
}

// Run sets a function that is invoked with the actual arguments on each call.
func (c *Call) Run(fn func(args Arguments)) *Call {
	c.lock()
	defer c.unlock()
	c.RunFn = fn
	return c
}

// On adds another expectation to the parent mock, for chaining.
func (c *Call) On(methodName string, arguments ...interface{}) *Call {
	return c.parent.On(methodName, arguments...)
}

func (c *Call) lock() {
	if c.parent != nil {
		c.parent.mu.Lock()
	}
}

func (c *Call) unlock() {
	if c.parent != nil {
		c.parent.mu.Unlock()
	}
}

// Mock records expectations and calls. Embed it in your fake type.
type Mock struct {
	mu sync.Mutex

	// ExpectedCalls registered by On.
	ExpectedCalls []*Call
	// Calls actually made through Called / MethodCalled.
	Calls []Call
}

// On registers an expectation that methodName is called with arguments.
func (m *Mock) On(methodName string, arguments ...interface{}) *Call {
	m.mu.Lock()
	defer m.mu.Unlock()

	c := &Call{
		parent:    m,
		Method:    methodName,
		Arguments: arguments,
	}
	m.ExpectedCalls = append(m.ExpectedCalls, c)
	return c
}

// Called records the call to the calling method with the given arguments
// and returns the configured return values. It panics if no matching
// expectation was registered.
func (m *Mock) Called(arguments ...interface{}) Arguments {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("mock: Called: could not get the caller information")
	}
	name := runtime.FuncForPC(pc).Name()
	name = strings.TrimSuffix(name, "-fm")
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return m.MethodCalled(name, arguments...)
}

// MethodCalled records the call to methodName with the given arguments
// and returns the configured return values. It panics if no matching
// expectation was registered.
func (m *Mock) MethodCalled(methodName string, arguments ...interface{}) Arguments {
	m.mu.Lock()

	call := m.find(methodName, arguments)
	if call == nil {
		m.mu.Unlock()
		panic(m.unexpected(methodName, arguments))
	}

	call.totalCalls++
	if call.Repeatability == 1 {
		call.Repeatability = -1
	} else if call.Repeatability > 1 {
		call.Repeatability--
	}

	m.Calls = append(m.Calls, Call{
		Method:          methodName,
		Arguments:       arguments,
		ReturnArguments: call.ReturnArguments,
	})
	runFn := call.RunFn
	ret := call.ReturnArguments
	m.mu.Unlock()

	if runFn != nil {
		runFn(arguments)
	}
	return ret
}

// find returns the first non-exhausted expectation matching the call.
// Callers must hold m.mu.
func (m *Mock) find(methodName string, arguments Arguments) *Call {
	for _, c := range m.ExpectedCalls {
		if c.Method == methodName && c.Repeatability > -1 && c.Arguments.Is(arguments...) {
			return c
		}
	}
	return nil
}

// unexpected builds the panic message for an unmatched call.
// Callers must hold m.mu.
func (m *Mock) unexpected(methodName string, arguments Arguments) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "mock: unexpected call %s(%s)", methodName, formatArgs(arguments))

	closest := ""
	for _, c := range m.ExpectedCalls {
		if c.Method != methodName {
			continue
		}
		if c.Repeatability == -1 {
			closest = fmt.Sprintf("\n  expectation %s(%s) has already been consumed",
				c.Method, formatArgs(c.Arguments))
			continue
		}
		closest = fmt.Sprintf("\n  closest expectation %s(%s)",
			c.Method, formatArgs(c.Arguments))
		break
	}
	if closest == "" {
		closest = "\n  no expectations registered for " + methodName +
			"; add mock.On(\"" + methodName + "\", ...)"
	}
	sb.WriteString(closest)
	return sb.String()
}

func formatArgs(args Arguments) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = fmt.Sprintf("%#v", a)
	}
	return strings.Join(parts, ", ")
}

// AssertExpectations asserts that every expectation registered with On
// was met (called at least once, or exactly Times(n) times).
func (m *Mock) AssertExpectations(t TestingT) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	ok := true
	for _, c := range m.ExpectedCalls {
		switch {
		case c.totalCalls == 0:
			t.Errorf("mock: expected call not made: %s(%s)", c.Method, formatArgs(c.Arguments))
			ok = false
		case c.Repeatability > 0:
			t.Errorf("mock: %s(%s) expected %d more call(s), was called %d time(s)",
				c.Method, formatArgs(c.Arguments), c.Repeatability, c.totalCalls)
			ok = false
		}
	}
	return ok
}

// AssertCalled asserts that methodName was called with the given arguments.
func (m *Mock) AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.wasCalled(methodName, arguments) {
		return true
	}
	t.Errorf("mock: expected %s(%s) to have been called, calls: %s",
		methodName, formatArgs(arguments), m.callsString())
	return false
}

// AssertNotCalled asserts that methodName was not called with the given arguments.
func (m *Mock) AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.wasCalled(methodName, arguments) {
		return true
	}
	t.Errorf("mock: expected %s(%s) not to have been called",
		methodName, formatArgs(arguments))
	return false
}

// AssertNumberOfCalls asserts that methodName was called expectedCalls times,
// regardless of arguments.
func (m *Mock) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	n := 0
	for _, c := range m.Calls {
		if c.Method == methodName {
			n++
		}
	}
	if n == expectedCalls {
		return true
	}
	t.Errorf("mock: expected %s to be called %d time(s), but was called %d time(s)",
		methodName, expectedCalls, n)
	return false
}

// wasCalled reports whether a recorded call matches. Callers must hold m.mu.
func (m *Mock) wasCalled(methodName string, arguments Arguments) bool {
	for _, c := range m.Calls {
		if c.Method == methodName && Arguments(arguments).Is(c.Arguments...) {
			return true
		}
	}
	return false
}

func (m *Mock) callsString() string {
	if len(m.Calls) == 0 {
		return "none"
	}
	parts := make([]string, len(m.Calls))
	for i, c := range m.Calls {
		parts[i] = fmt.Sprintf("%s(%s)", c.Method, formatArgs(c.Arguments))
	}
	return strings.Join(parts, "; ")
}

// Reset clears all expectations and recorded calls.
func (m *Mock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ExpectedCalls = nil
	m.Calls = nil
}
