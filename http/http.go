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

// Package http is a small http.Handler test helper for tt.
//
// Build a request against a handler, run it in-process with
// net/http/httptest and assert on the recorded response:
//
//	r := http.New(handler).Get("/users").Query("id", "1").Run()
//	http.Status(t, 200, r)
//	http.Contains(t, "alice", r)
//
// or fluently:
//
//	http.New(handler).Post("/users").JSON(user).Run().
//		Expect(t).Status(201).Header("Content-Type", "application/json").
//		JSON(map[string]interface{}{"id": 1}).OK()
//
// Every assertion accepts tt's trailing args: [0] message, [1] caller depth
// int or trace string, [2] trace string.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"

	"github.com/vcaesar/tt"
)

// TestingT is an interface wrapper around *testing.T,
// compatible with tt.TestingT.
type TestingT interface {
	Errorf(format string, args ...interface{})
}

// H is a shorthand for a map of header, query or form values.
type H map[string]string

// Request is a builder for an in-process test request.
type Request struct {
	Method  string
	Path    string
	Header  http.Header
	Query   url.Values
	Cookies []*http.Cookie
	Body    io.Reader

	handler http.Handler
	err     error
}

// New returns a Request builder bound to handler.
func New(handler http.Handler) *Request {
	return &Request{
		Method:  http.MethodGet,
		Path:    "/",
		Header:  http.Header{},
		Query:   url.Values{},
		handler: handler,
	}
}

// Func returns a Request builder bound to an http.HandlerFunc.
func Func(fn func(http.ResponseWriter, *http.Request)) *Request {
	return New(http.HandlerFunc(fn))
}

// Do sets the method and path.
func (r *Request) Do(method, path string) *Request {
	r.Method = method
	r.Path = path
	return r
}

// Get sets a GET request to path.
func (r *Request) Get(path string) *Request { return r.Do(http.MethodGet, path) }

// Post sets a POST request to path.
func (r *Request) Post(path string) *Request { return r.Do(http.MethodPost, path) }

// Put sets a PUT request to path.
func (r *Request) Put(path string) *Request { return r.Do(http.MethodPut, path) }

// Patch sets a PATCH request to path.
func (r *Request) Patch(path string) *Request { return r.Do(http.MethodPatch, path) }

// Delete sets a DELETE request to path.
func (r *Request) Delete(path string) *Request { return r.Do(http.MethodDelete, path) }

// Head sets a HEAD request to path.
func (r *Request) Head(path string) *Request { return r.Do(http.MethodHead, path) }

// Options sets an OPTIONS request to path.
func (r *Request) Options(path string) *Request { return r.Do(http.MethodOptions, path) }

// Set adds a request header.
func (r *Request) Set(key, value string) *Request {
	r.Header.Set(key, value)
	return r
}

// Headers adds request headers.
func (r *Request) Headers(h H) *Request {
	for k, v := range h {
		r.Header.Set(k, v)
	}
	return r
}

// Param adds a query parameter.
func (r *Request) Param(key, value string) *Request {
	r.Query.Add(key, value)
	return r
}

// Params adds query parameters.
func (r *Request) Params(h H) *Request {
	for k, v := range h {
		r.Query.Add(k, v)
	}
	return r
}

// Cookie adds a request cookie.
func (r *Request) Cookie(name, value string) *Request {
	r.Cookies = append(r.Cookies, &http.Cookie{Name: name, Value: value})
	return r
}

// Text sets a raw string body.
func (r *Request) Text(body string) *Request {
	r.Body = strings.NewReader(body)
	return r
}

// Bytes sets a raw []byte body.
func (r *Request) Bytes(body []byte) *Request {
	r.Body = bytes.NewReader(body)
	return r
}

// Reader sets the body from an io.Reader.
func (r *Request) Reader(body io.Reader) *Request {
	r.Body = body
	return r
}

// JSON marshals v as the body and sets Content-Type: application/json.
// string and []byte values are sent as-is.
func (r *Request) JSON(v interface{}) *Request {
	switch b := v.(type) {
	case string:
		r.Body = strings.NewReader(b)
	case []byte:
		r.Body = bytes.NewReader(b)
	default:
		buf, err := json.Marshal(v)
		if err != nil {
			r.err = err
			return r
		}
		r.Body = bytes.NewReader(buf)
	}

	r.Header.Set("Content-Type", "application/json")
	return r
}

// Form sets an url-encoded form body and
// Content-Type: application/x-www-form-urlencoded.
func (r *Request) Form(h H) *Request {
	v := url.Values{}
	for k, val := range h {
		v.Set(k, val)
	}
	r.Body = strings.NewReader(v.Encode())
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// Build returns the *http.Request described by the builder.
func (r *Request) Build() (*http.Request, error) {
	if r.err != nil {
		return nil, r.err
	}

	body := r.Body
	if body == nil {
		// Server-side handlers always see a non-nil Body.
		body = http.NoBody
	} else if s, ok := body.(io.Seeker); ok {
		// Rewind so the builder can be Run more than once.
		if _, err := s.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(r.Method, r.Path, body)
	if err != nil {
		return nil, err
	}

	if len(r.Query) > 0 {
		q := req.URL.Query()
		for k, vs := range r.Query {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	for k, vs := range r.Header {
		req.Header[k] = append([]string(nil), vs...)
	}
	for _, c := range r.Cookies {
		req.AddCookie(c)
	}

	return req, nil
}

// Run executes the request against the handler and records the response.
// A request that cannot be built yields a Response with Code -1 and Err set.
func (r *Request) Run() *Response {
	req, err := r.Build()
	if err != nil {
		return &Response{Code: -1, Header: http.Header{}, Err: err}
	}

	w := httptest.NewRecorder()
	r.handler.ServeHTTP(w, req)

	res := w.Result()
	return &Response{
		Code:     res.StatusCode,
		Header:   res.Header,
		Body:     w.Body.String(),
		Cookies:  res.Cookies(),
		Request:  req,
		Recorder: w,
	}
}

// Response is the recorded result of Request.Run.
type Response struct {
	Code    int
	Header  http.Header
	Body    string
	Cookies []*http.Cookie
	// Err is set when the request could not be built or run.
	Err error

	Request  *http.Request
	Recorder *httptest.ResponseRecorder
}

// Bytes returns the body as []byte.
func (r *Response) Bytes() []byte {
	return []byte(r.Body)
}

// JSON unmarshals the body into v.
func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal([]byte(r.Body), v)
}

// Cookie returns the response cookie with name, or nil.
func (r *Response) Cookie(name string) *http.Cookie {
	for _, c := range r.Cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// Expect returns an assertion chain for the response.
func (r *Response) Expect(t TestingT) *Expect {
	return &Expect{t: t, r: r, ok: true}
}

// Fetch runs the handler with method, path and query values and returns the
// response. It mirrors testify's http assertions helper signature.
func Fetch(handler http.Handler, method, path string, values url.Values) *Response {
	r := New(handler).Do(method, path)
	r.Query = values
	return r.Run()
}

// argsFn parses tt trailing args: [0] message, [1] depth int or trace string,
// [2] trace string.
func argsFn(depth int, args ...interface{}) (string, int, string) {
	info, cinfo := "", ""
	if len(args) > 0 {
		info, _ = args[0].(string)
	}

	if len(args) > 1 {
		switch v := args[1].(type) {
		case int:
			depth = v
		case string:
			cinfo = v
		}
	}
	if len(args) > 2 {
		cinfo, _ = args[2].(string)
	}

	return info, depth, cinfo
}

// depth is the CallerInfo index of the test that called a package-level
// assertion: 0 CallerInfo, 1 tt.Fmt, 2 fail, 3 assertion, 4 test.
const depth = 4

func fail(t TestingT, equal, expect string, call int, info, cinfo string,
	ept, got interface{}) bool {
	err := tt.Fmt(equal, expect, call, info, cinfo)
	t.Errorf(err, fmt.Sprint(ept), fmt.Sprint(got))
	return false
}

// Status asserts the response status code equals expect.
func Status(t TestingT, expect int, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if r.Code != expect {
		return fail(t, "Not Equal", "expected", call, info, cinfo, expect, r.Code)
	}
	return true
}

// Success asserts the status code is 2xx.
func Success(t TestingT, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if r.Code < http.StatusOK || r.Code > 299 {
		return fail(t, "Not Success", "expected", call, info, cinfo, "2xx", r.Code)
	}
	return true
}

// Redirect asserts the status code is 3xx.
func Redirect(t TestingT, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if r.Code < http.StatusMultipleChoices || r.Code > 399 {
		return fail(t, "Not Redirect", "expected", call, info, cinfo, "3xx", r.Code)
	}
	return true
}

// Error asserts the status code is 4xx or 5xx.
func Error(t TestingT, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if r.Code < http.StatusBadRequest {
		return fail(t, "Not Error", "expected", call, info, cinfo, ">= 400", r.Code)
	}
	return true
}

// Body asserts the response body equals expect.
func Body(t TestingT, expect string, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if r.Body != expect {
		return fail(t, "Not Equal", "expected", call, info, cinfo, expect, r.Body)
	}
	return true
}

// Contains asserts the response body contains str.
func Contains(t TestingT, str string, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if !strings.Contains(r.Body, str) {
		return fail(t, "Not Contains", "expected", call, info, cinfo, str, r.Body)
	}
	return true
}

// NotContains asserts the response body does not contain str.
func NotContains(t TestingT, str string, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if strings.Contains(r.Body, str) {
		return fail(t, "Contains", "not expected", call, info, cinfo, str, r.Body)
	}
	return true
}

// Header asserts the response header key equals expect.
func Header(t TestingT, key, expect string, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	if got := r.Header.Get(key); got != expect {
		return fail(t, "Not Equal", "expected", call, info, cinfo, key+": "+expect, key+": "+got)
	}
	return true
}

// Cookie asserts the response sets cookie name with value expect.
func Cookie(t TestingT, name, expect string, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)
	c := r.Cookie(name)
	if c == nil {
		return fail(t, "Cookie Not Set", "expected", call, info, cinfo, name+"="+expect, "")
	}
	if c.Value != expect {
		return fail(t, "Not Equal", "expected", call, info, cinfo, name+"="+expect, name+"="+c.Value)
	}
	return true
}

// jsonNorm decodes v (or marshals then decodes non-string/[]byte values)
// into a generic interface{} so that JSON documents compare structurally.
func jsonNorm(v interface{}) (interface{}, error) {
	var raw []byte
	switch b := v.(type) {
	case string:
		raw = []byte(b)
	case []byte:
		raw = b
	default:
		var err error
		raw, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}

	var out interface{}
	err := json.Unmarshal(raw, &out)
	return out, err
}

// JSON asserts the response body is JSON structurally equal to expect.
// expect may be a JSON string, []byte or any value that marshals to JSON.
func JSON(t TestingT, expect interface{}, r *Response, args ...interface{}) bool {
	info, call, cinfo := argsFn(depth, args...)

	ept, err := jsonNorm(expect)
	if err != nil {
		return fail(t, "Invalid JSON", "expected", call, info, cinfo, expect, err)
	}
	got, err := jsonNorm(r.Body)
	if err != nil {
		return fail(t, "Invalid JSON", "expected", call, info, cinfo, r.Body, err)
	}

	if !reflect.DeepEqual(ept, got) {
		e, _ := json.Marshal(ept)
		g, _ := json.Marshal(got)
		return fail(t, "Not Equal", "expected", call, info, cinfo, string(e), string(g))
	}
	return true
}

// Expect is a chain of assertions on a Response.
type Expect struct {
	t  TestingT
	r  *Response
	ok bool
}

// chain shifts the default depth by one for the extra Expect frame.
func chain(args []interface{}) []interface{} {
	if len(args) > 1 {
		return args
	}
	if len(args) == 0 {
		args = append(args, "")
	}
	return append(args, depth+1)
}

func (e *Expect) set(ok bool) *Expect {
	e.ok = e.ok && ok
	return e
}

// Status asserts the status code equals code.
func (e *Expect) Status(code int, args ...interface{}) *Expect {
	return e.set(Status(e.t, code, e.r, chain(args)...))
}

// Success asserts the status code is 2xx.
func (e *Expect) Success(args ...interface{}) *Expect {
	return e.set(Success(e.t, e.r, chain(args)...))
}

// Redirect asserts the status code is 3xx.
func (e *Expect) Redirect(args ...interface{}) *Expect {
	return e.set(Redirect(e.t, e.r, chain(args)...))
}

// Error asserts the status code is 4xx or 5xx.
func (e *Expect) Error(args ...interface{}) *Expect {
	return e.set(Error(e.t, e.r, chain(args)...))
}

// Body asserts the body equals expect.
func (e *Expect) Body(expect string, args ...interface{}) *Expect {
	return e.set(Body(e.t, expect, e.r, chain(args)...))
}

// Contains asserts the body contains str.
func (e *Expect) Contains(str string, args ...interface{}) *Expect {
	return e.set(Contains(e.t, str, e.r, chain(args)...))
}

// NotContains asserts the body does not contain str.
func (e *Expect) NotContains(str string, args ...interface{}) *Expect {
	return e.set(NotContains(e.t, str, e.r, chain(args)...))
}

// Header asserts the header key equals expect.
func (e *Expect) Header(key, expect string, args ...interface{}) *Expect {
	return e.set(Header(e.t, key, expect, e.r, chain(args)...))
}

// Cookie asserts the cookie name equals expect.
func (e *Expect) Cookie(name, expect string, args ...interface{}) *Expect {
	return e.set(Cookie(e.t, name, expect, e.r, chain(args)...))
}

// JSON asserts the body is JSON structurally equal to expect.
func (e *Expect) JSON(expect interface{}, args ...interface{}) *Expect {
	return e.set(JSON(e.t, expect, e.r, chain(args)...))
}

// Fn runs a custom assertion against the response.
func (e *Expect) Fn(fn func(t TestingT, r *Response) bool) *Expect {
	return e.set(fn(e.t, e.r))
}

// OK reports whether every assertion in the chain passed.
func (e *Expect) OK() bool {
	return e.ok
}

// String returns a short summary of the response for debugging.
func (r *Response) String() string {
	return fmt.Sprintf("%d %s", r.Code, r.Body)
}
