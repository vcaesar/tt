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

package tt

import "testing"

// Assertions provides assertion methods around the
// TestingT interface.
type Assertions struct {
	t TestingT
}

// New makes a new Assertions object for the specified TestingT.
func New(t TestingT) *Assertions {
	return &Assertions{
		t: t,
	}
}

// BM is the encapsulation of the benchmark function
//	func Benchmark1(b *testing.B, fn func())
func (at *Assertions) BM(b *testing.B, fn func()) {
	for i := 0; i < b.N; i++ {
		fn()
	}
}

// Equal asserts that two objects are equal.
func (at *Assertions) Equal(expect, actual interface{}, args ...interface{}) bool {
	info, call := typeCall(expect, actual, args...)

	return Equal(at.t, expect, actual, info, call)
}

// Expect asserts that string and objects are equal.
func (at *Assertions) Expect(expect string, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)
	return Expect(at.t, expect, actual, info, call)
}

// Nil asserts that nil and objects are equal.
func (at *Assertions) Nil(actual interface{}, args ...interface{}) bool {
	info, call := typeCall(nil, actual, args...)
	return Equal(at.t, nil, actual, info, call)
}

// Empty asserts that empty and objects are equal.
func (at *Assertions) Empty(actual interface{}, args ...interface{}) bool {
	info, call := typeCall("", actual, args...)
	return Equal(at.t, "", actual, info, call)
}

// Bool asserts that true and objects are equal.
func (at *Assertions) Bool(actual interface{}, args ...interface{}) bool {
	info, call := typeCall(true, actual, args...)
	return Equal(at.t, true, actual, info, call)
}

// True asserts that true and objects are equal.
func (at *Assertions) True(actual interface{}, args ...interface{}) bool {
	info, call := typeCall(true, actual, args...)
	return Equal(at.t, true, actual, info, call)
}

// False asserts that flase and objects are equal.
func (at *Assertions) False(actual interface{}, args ...interface{}) bool {
	info, call := typeCall(true, actual, args...)
	return Equal(at.t, false, actual, info, call)
}

// Not asserts that two objects are not equal.
func (at *Assertions) Not(expect, actual interface{}, args ...interface{}) bool {
	info, call := callAdd(Type, args...)
	return NotEqual(at.t, expect, actual, info, call)
}

// NotEqual asserts that two objects are not equal.
func (at *Assertions) NotEqual(expect, actual interface{}, args ...interface{}) bool {
	info, call := callAdd(Type, args...)
	return NotEqual(at.t, expect, actual, info, call)
}

// NotExpect asserts that string and objects are not equal.
func (at *Assertions) NotExpect(expect string, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)
	return NotExpect(at.t, expect, actual, info, call)
}
