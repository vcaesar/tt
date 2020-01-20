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
	"reflect"
	"testing"
	"time"

	"net/http"
	_ "net/http/pprof"
)

const (
	// Version get the tt version
	Version = "v0.20.0.77, Sierra Nevada!"
)

var (
	// Type type must
	Type bool
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

func typeOf(expect, actual interface{}) bool {
	if reflect.TypeOf(expect) == reflect.TypeOf(actual) {
		return true
	}
	return false
}

// TypeF make Type false
func TypeF() {
	Type = false
}

func argsFn(args ...interface{}) (string, int) {
	call := 5
	if len(args) > 1 {
		call = args[1].(int)
	}

	info := ""
	if len(args) > 0 {
		info = args[0].(string)
	}

	return info, call
}

// Fmt return error string
func Fmt(equal, expect string, call int, info ...string) string {
	err := RedBold("\n Error Trace:		" + CallerInfo()[call] + ",")
	err += Yellow("\n Error:		" + equal + "; \n ")
	if len(info) > 0 && info[0] != "" {
		err += "Messages:	" + info[0] + "\n "
	}

	err += Blue(expect+":	'%s',\n ") + Red("but got:	'%s' \n\n")
	return err
}

// FmtErr return error string
func FmtErr(call int, info ...string) string {
	return Fmt("Not Equal", "expected", call, info...)
}

// DEqual asserts that two objects are deep equal.
//
//    tt.DEqual(t *testing.T, 1, 1)
//
func DEqual(t TestingT, expect, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)
	Type = true
	defer TypeF()

	if typeOf(expect, actual) {
		call = call + 1
	}

	return Equal(t, expect, actual, info, call)
}

// Equal asserts that two objects are equal.
//
//    tt.Equal(t *testing.T, 1, 1)
//
func Equal(t TestingT, expect, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	if Type && !typeOf(expect, actual) {
		if len(args) < 1 {
			call = call - 1
		}

		err := FmtErr(call, info)
		t.Errorf(err, expect, actual)
		return false
	}

	expectStr := fmt.Sprint(expect)
	return Expect(t, expectStr, actual, info, call)
}

// Expect asserts that string and objects are equal.
//
//    tt.Expect(t *testing.T, "1", 1)
//
func Expect(t TestingT, expect string, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)
	if len(args) < 1 {
		call = call - 1
	}

	actualStr := fmt.Sprint(actual)
	if expect != actualStr {
		err := FmtErr(call, info)

		t.Errorf(err, expect, actualStr)
		return false
	}

	return true
}

// Nil asserts that nil and objects are equal.
func Nil(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	return Equal(t, nil, actual, info, call)
}

// Empty asserts that empty and objects are equal.
func Empty(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	return Equal(t, "", actual, info, call)
}

// Bool asserts that true and objects are equal.
func Bool(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	return Equal(t, true, actual, info, call)
}

// True asserts that true and objects are equal.
func True(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	return Equal(t, true, actual, info, call)
}

// False asserts that flase and objects are equal.
func False(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	return Equal(t, false, actual, info, call)
}

// NotErr return not equal error string
func NotErr(call int, info ...string) string {
	return Fmt("Equal", "not expected", call, info...)
}

// Not asserts that two objects are not equal.
//
//    tt.NotEqual(t *testing.T, 1, 1)
//
func Not(t TestingT, expect, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	return NotEqual(t, expect, actual, info, call)
}

// NotEqual asserts that two objects are not equal.
//
//    tt.NotEqual(t *testing.T, 1, 1)
//
func NotEqual(t TestingT, expect, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)

	if Type && typeOf(expect, actual) {
		if len(args) < 1 {
			call = call - 1
		}

		err := FmtErr(call, info)
		t.Errorf(err, expect, actual)
		return false
	}

	expectStr := fmt.Sprint(expect)
	return NotExpect(t, expectStr, actual, info, call)
}

// NotExpect asserts that string and objects are not equal.
//
//    tt.NotExpect(t *testing.T, "1", 1)
//
func NotExpect(t TestingT, expect string, actual interface{}, args ...interface{}) bool {
	info, call := argsFn(args...)
	if len(args) < 1 {
		call = call - 1
	}

	actualStr := fmt.Sprint(actual)
	if expect == actualStr {
		err := NotErr(call, info)

		t.Errorf(err, expect, actualStr)
		return false
	}

	return true
}
