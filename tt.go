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

// package tt is simple and colorful test tools

package tt

import (
	"fmt"
	"log"
	"testing"
	"time"

	"net/http"
	_ "net/http/pprof"
)

const (
	// Version get the tt version
	Version = "v0.10.0.54, Sierra Nevada!"
)

// Pprof use:
// Mem:
// pprof -http=:8090 http://127.0.0.1:6060/debug/pprof/heap
//
// go tool pprof http://localhost:6060/debug/pprof/heap
//
// CPU:
//		debug/pprof/profile
// set time:
//		debug/pprof/profile\?seconds\=10
// pprof -http=:8090 http://127.0.0.1:6060/debug/pprof/profile\?seconds\=10
//
// debug/pprof/block
// debug/pprof/mutex
func Pprof(tm ...int) bool {
	go func() {
		log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
	}()

	if len(tm) > 0 {
		time.Sleep(time.Duration(tm[0]) * time.Second)
	}

	return true
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
}

// BM func Benchmark1(b *testing.B, fn func())
func BM(b *testing.B, fn func()) {
	for i := 0; i < b.N; i++ {
		fn()
	}
}

// FmtErr return error string
func FmtErr(call int) string {
	err := RedBold("\n Error Trace:		" + CallerInfo()[call] + ",")
	err += Yellow("\n Error:		Not equal; \n ")
	err += Blue("expected:	'%s',\n ") + Red("but got:	'%s' \n\n")
	return err
}

// Equal asserts that two objects are equal.
//
//    tt.Equal(t *testing.T, 1, 1)
//
func Equal(t TestingT, expect, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	expectStr := fmt.Sprint(expect)
	return Expect(t, expectStr, actual, call)
}

// Expect asserts that string and objects are equal.
//
//    tt.Expect(t *testing.T, "1", 1)
//
func Expect(t TestingT, expect string, actual interface{}, args ...int) bool {
	call := 3
	if len(args) > 0 {
		call = args[0]
	}

	actualStr := fmt.Sprint(actual)
	if expect != actualStr {
		err := FmtErr(call)
		t.Errorf(err, expect, actualStr)

		return false
	}

	return true
}

// Nil asserts that nil and objects are equal.
func Nil(t TestingT, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(t, "<nil>", actual, call)
}

// Empty asserts that empty and objects are equal.
func Empty(t TestingT, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(t, "", actual, call)
}

// Bool asserts that true and objects are equal.
func Bool(t TestingT, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(t, "true", actual, call)
}

// True asserts that true and objects are equal.
func True(t TestingT, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(t, "true", actual, call)
}

// False asserts that flase and objects are equal.
func False(t TestingT, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(t, "false", actual, call)
}

// NotErr return not equal error string
func NotErr(call int) string {
	err := RedBold("\n Error Trace:		" + CallerInfo()[call] + ",")
	err += Yellow("\n Error:		Equal; \n ")
	err += Blue("not expected:	'%s',\n ") + Red("but got:	'%s' \n\n")
	return err
}

// Not asserts that two objects are not equal.
//
//    tt.NotEqual(t *testing.T, 1, 1)
//
func Not(t TestingT, expect, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	expectStr := fmt.Sprint(expect)
	return NotExpect(t, expectStr, actual, call)
}

// NotEqual asserts that two objects are not equal.
//
//    tt.NotEqual(t *testing.T, 1, 1)
//
func NotEqual(t TestingT, expect, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	expectStr := fmt.Sprint(expect)
	return NotExpect(t, expectStr, actual, call)
}

// NotExpect asserts that string and objects are not equal.
//
//    tt.NotExpect(t *testing.T, "1", 1)
//
func NotExpect(t TestingT, expect string, actual interface{}, args ...int) bool {
	call := 3
	if len(args) > 0 {
		call = args[0]
	}

	actualStr := fmt.Sprint(actual)
	if expect == actualStr {
		err := NotErr(call)
		t.Errorf(err, expect, actualStr)

		return false
	}

	return true
}
