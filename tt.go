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
	"runtime"
	"strings"
	"testing"
	"time"
	"unicode"

	"net/http"
	_ "net/http/pprof"
	"unicode/utf8"
)

// Pprof use:
// pprof -http=:8090 http://127.0.0.1:6060/debug/pprof/heap
//
// go tool pprof http://localhost:6060/debug/pprof/heap
//
// debug/pprof/profile
// set time:
// debug/pprof/profile\?seconds\=10
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

// Blue returns a blue string
func Blue(message string) string {
	return fmt.Sprintf("\x1b[34m%s\x1b[0m", message)
}

// Red returns a red string
func Red(message string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", message)
}

// Yellow returns a yellow string
func Yellow(message string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", message)
}

//Bold returns a blod string
func Bold(message string) string {
	return fmt.Sprintf("\x1b[1m%s\x1b[21m", message)
}

// RedBold returns a red bold string
func RedBold(message string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", Bold(message))
}

func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}

// CallerInfo returns an array of strings containing the file and line number
// of each stack frame leading from the current test to the assert call that
// failed.
func CallerInfo() (callers []string) {
	var (
		pc         uintptr
		file, name string
		line       int
		ok         bool
	)

	for i := 0; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// The breaks below failed to terminate the loop, and we ran off the
			// end of the call stack.
			break
		}

		// This is a huge edge case, but it will panic if this is the case, see #180
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name = f.Name()

		// testing.tRunner is the standard library function that calls
		// tests. Subtests are called directly by tRunner, without going through
		// the Test/Benchmark/Example function that contains the t.Run calls, so
		// with subtests we should break when we hit tRunner, without adding it
		// to the list of callers.
		if name == "testing.tRunner" {
			break
		}

		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
		if len(parts) > 1 {
			dir := parts[len(parts)-2]
			if (dir != "assert" && dir != "mock" && dir != "require") ||
				file == "mock_test.go" {
				callers = append(callers, fmt.Sprintf("%s:%d", file, line))
			}
		}

		// Drop the package
		segments := strings.Split(name, ".")
		name = segments[len(segments)-1]
		if isTest(name, "Test") ||
			isTest(name, "Benchmark") ||
			isTest(name, "Example") {
			break
		}
	}

	return
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
