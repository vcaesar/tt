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
	// _ "net/http/pprof"
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

func typeMap(expect string) string {
	m := map[string]string{
		"str":   "string",
		"int":   "int",
		"i8":    "int8",
		"i16":   "int16",
		"i32":   "int32",
		"i64":   "int64",
		"u":     "uint",
		"u8":    "uint8",
		"u16":   "uint16",
		"ui32":  "uint32",
		"ui64":  "uint64",
		"f32":   "float32",
		"f64":   "float64",
		"b":     "bool",
		"m":     "map",
		"ch":    "chan",
		"stu":   "struct",
		"bytes": "[]byte",
		"u8s":   "[]uint8",
		"c64":   "complex64",
		"c128":  "complex128",
	}

	return m[expect]
}

// IsTypes return bool when actual is expect type
func IsTypes(expect string, actual interface{}) bool {
	mt := typeMap(expect)
	s := reflect.TypeOf(actual).String()
	if s == mt || s == expect {
		return true
	}

	return false
}

func typVal(expect string) string {
	ept := typeMap(expect)
	if ept == "" {
		ept = expect
	}

	return ept
}

// TypeOf equal two interface{} type
func TypeOf(expect, actual interface{}) bool {
	return reflect.TypeOf(expect) == reflect.TypeOf(actual)
}

// IsType asserts that two objects type are equal
func IsType(t TestingT, expect string, actual interface{}, args ...interface{}) bool {
	if actual == nil && expect == "nil" {
		return true
	}
	if actual == nil {
		return false
	}

	s := reflect.TypeOf(actual).String()
	expect = typVal(expect)

	info, call, cinfo := typeCall(expect, s, args...)
	return Equal(t, expect, s, info, call, cinfo)
}

// TypeF make Type false
func TypeF() {
	Type = false
}

func argsFn(args ...interface{}) (string, int, string) {
	info := ""
	if len(args) > 0 {
		info = args[0].(string)
	}

	call := 5
	callInfo := ""
	if len(args) > 1 {
		if TypeOf(args[1], callInfo) {
			callInfo = args[1].(string)
			return info, call, callInfo
		}

		call = args[1].(int)
	}

	return info, call, callInfo
}

func callSub(args ...interface{}) (string, int, string) {
	info, call, cinfo := argsFn(args...)
	if len(args) < 1 {
		call = call - 1
	}
	return info, call, cinfo
}

func callAdd(t bool, args ...interface{}) (string, int, string) {
	info, call, cinfo := argsFn(args...)
	if len(args) < 1 && !t {
		call = call + 1
	}

	return info, call, cinfo
}

func typeCall(expect, actual interface{}, args ...interface{}) (string, int, string) {
	b := Type && !TypeOf(expect, actual)
	return callAdd(b, args...)
}

// Fmt return error string
func Fmt(equal, expect string, call int, info ...string) (err string) {
	if len(info) > 1 && info[1] != "" {
		err = RedBold("\n Error Trace:		" + info[1] + ",")
	} else {
		err = RedBold("\n Error Trace:		" + CallerInfo()[call] + ",")
	}

	err += Yellow("\n Error:		" + equal + "; \n ")
	if len(info) > 0 && info[0] != "" {
		err += "Messages:	" + info[0] + "\n "
	}

	err += Blue(expect+":	'%s',\n ") + Red("but got:	'%s' \n\n")
	return
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
	info, call, cinfo := argsFn(args...)
	Type = true
	defer TypeF()

	if TypeOf(expect, actual) {
		call = call + 1
	}

	return Equal(t, expect, actual, info, call, cinfo)
}

// Equal asserts that two objects are equal.
//
//    tt.Equal(t *testing.T, 1, 1)
//
func Equal(t TestingT, expect, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := argsFn(args...)
	if len(args) > 2 {
		cinfo = args[2].(string)
	}

	if Type && !TypeOf(expect, actual) {
		if len(args) < 2 {
			call = call - 1
		}

		err := FmtErr(call, info, cinfo)
		t.Errorf(err, expect, actual)
		return false
	}

	expectStr := fmt.Sprint(expect)
	return Expect(t, expectStr, actual, info, call, cinfo)
}

// Expect asserts that string and objects are equal.
//
//    tt.Expect(t *testing.T, "1", 1)
//
func Expect(t TestingT, expect string, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := callSub(args...)
	if len(args) > 2 {
		cinfo = args[2].(string)
	}

	actualStr := fmt.Sprint(actual)
	if expect != actualStr {
		err := FmtErr(call, info, cinfo)

		t.Errorf(err, expect, actualStr)
		return false
	}

	return true
}

// Nil asserts that nil and objects are equal.
func Nil(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := typeCall(nil, actual, args...)
	return Equal(t, nil, actual, info, call, cinfo)
}

// NotNil asserts that not equal nil.
//
//    tt.NotNil(t *testing.T, 1)
//
func NotNil(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := callAdd(Type, args...)
	return NotEqual(t, nil, actual, info, call, cinfo)
}

// Error asserts that equal error.
func Error(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := callAdd(Type, args...)
	return IsType(t, "*errors.errorString", actual, info, call, cinfo)
}

// Empty asserts that empty and objects are equal.
func Empty(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := typeCall("", actual, args...)
	return Equal(t, "", actual, info, call, cinfo)
}

// NotEmpty asserts that empty and objects are not equal.
func NotEmpty(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := callAdd(Type, args...)
	return NotEqual(t, "", actual, info, call, cinfo)
}

// Zero asserts that zero and objects are equal.
func Zero(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := typeCall(0, actual, args...)
	return Equal(t, 0, actual, info, call, cinfo)
}

// NotZero asserts that zero and objects are not equal.
func NotZero(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := callAdd(Type, args...)
	return NotEqual(t, 0, actual, info, call, cinfo)
}

// Bool asserts that true and objects are equal.
func Bool(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := typeCall(true, actual, args...)

	return Equal(t, true, actual, info, call, cinfo)
}

// True asserts that true and objects are equal.
func True(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := typeCall(true, actual, args...)

	return Equal(t, true, actual, info, call, cinfo)
}

// False asserts that flase and objects are equal.
func False(t TestingT, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := typeCall(false, actual, args...)

	return Equal(t, false, actual, info, call, cinfo)
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
	info, call, cinfo := callAdd(Type, args...)

	return NotEqual(t, expect, actual, info, call, cinfo)
}

// NotEqual asserts that two objects are not equal.
//
//    tt.NotEqual(t *testing.T, 1, 1)
//
func NotEqual(t TestingT, expect, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := argsFn(args...)
	if len(args) > 2 {
		cinfo = args[2].(string)
	}

	if Type && TypeOf(expect, actual) {
		if len(args) < 2 {
			call = call - 1
		}

		err := NotErr(call, info, cinfo)
		t.Errorf(err, expect, actual)
		return false
	}

	if Type && !TypeOf(expect, actual) {
		return true
	}

	expectStr := fmt.Sprint(expect)
	return NotExpect(t, expectStr, actual, info, call, cinfo)
}

// NotExpect asserts that string and objects are not equal.
//
//    tt.NotExpect(t *testing.T, "1", 1)
//
func NotExpect(t TestingT, expect string, actual interface{}, args ...interface{}) bool {
	info, call, cinfo := callSub(args...)
	if len(args) > 2 {
		cinfo = args[2].(string)
	}

	actualStr := fmt.Sprint(actual)
	if expect == actualStr {
		err := NotErr(call, info, cinfo)

		t.Errorf(err, expect, actualStr)
		return false
	}

	return true
}
