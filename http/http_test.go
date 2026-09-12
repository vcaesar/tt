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

package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/vcaesar/tt"
	"github.com/vcaesar/tt/mock"
)

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "world"
		}
		w.Header().Set("X-Hello", name)
		w.Write([]byte("hello " + name))
	})
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Header().Set("X-Method", r.Method)
		if c, err := r.Cookie("session"); err == nil {
			http.SetCookie(w, &http.Cookie{Name: "seen", Value: c.Value})
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(b)
	})
	mux.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		w.Write([]byte(r.PostForm.Get("a") + r.PostForm.Get("b")))
	})
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/hello", http.StatusFound)
	})
	mux.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	return mux
}

func TestRequest(t *testing.T) {
	h := handler()

	r := New(h).Get("/hello").Param("name", "tt").Run()
	tt.Nil(t, r.Err)
	Status(t, 200, r)
	Success(t, r)
	Body(t, "hello tt", r)
	Contains(t, "tt", r)
	NotContains(t, "world", r)
	Header(t, "X-Hello", "tt", r)
	tt.Equal(t, "200 hello tt", r.String())
	tt.Equal(t, "hello tt", string(r.Bytes()))
	tt.Equal(t, "name=tt", r.Request.URL.RawQuery)

	r = New(h).Params(H{"name": "x"}).Do("GET", "/hello").Run()
	Body(t, "hello x", r)

	r = New(h).Post("/echo").JSON(map[string]int{"a": 1}).
		Set("X-A", "1").Headers(H{"X-B": "2"}).Cookie("session", "s1").Run()
	Status(t, 201, r)
	Header(t, "Content-Type", "application/json", r)
	Header(t, "X-Method", "POST", r)
	Cookie(t, "seen", "s1", r)
	JSON(t, `{"a":1}`, r)
	JSON(t, []byte(`{"a": 1}`), r)
	JSON(t, map[string]interface{}{"a": 1}, r)
	tt.Equal(t, "1", r.Request.Header.Get("X-A"))
	tt.Equal(t, "2", r.Request.Header.Get("X-B"))

	var got map[string]int
	tt.Nil(t, r.JSON(&got))
	tt.Equal(t, 1, got["a"])
	tt.True(t, r.Cookie("missing") == nil)
	tt.NotNil(t, New(h).Get("/hello").Run().Request.Body, "handlers need non-nil Body")

	Body(t, `{"b":2}`, New(h).Put("/echo").JSON(`{"b":2}`).Run())
	Body(t, "raw", New(h).Patch("/echo").JSON([]byte("raw")).Run())
	Body(t, "txt", New(h).Delete("/echo").Text("txt").Run())
	Body(t, "bs", New(h).Options("/echo").Bytes([]byte("bs")).Run())
	Body(t, "rd", New(h).Post("/echo").Reader(strings.NewReader("rd")).Run())
	Body(t, "", New(h).Head("/echo").Run())

	r = New(h).Post("/form").Form(H{"a": "1", "b": "2"}).Run()
	Body(t, "12", r)

	// A builder with a body can be Run more than once.
	req := New(h).Post("/echo").Text("again")
	Body(t, "again", req.Run())
	Body(t, "again", req.Run())

	Redirect(t, New(h).Get("/redirect").Run())
	Error(t, New(h).Get("/fail").Run())

	r = Fetch(h, "GET", "/hello", url.Values{"name": {"f"}})
	Body(t, "hello f", r)

	r = Func(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("fn"))
	}).Run()
	Body(t, "fn", r)
}

func TestBuildError(t *testing.T) {
	r := New(handler()).Post("/echo").JSON(make(chan int)).Run()
	tt.NotNil(t, r.Err)
	tt.Equal(t, -1, r.Code)

	r = New(handler()).Do("GET", "://bad").Run()
	tt.NotNil(t, r.Err)

	_, err := New(handler()).Get("/x").Build()
	tt.Nil(t, err)
}

func TestExpect(t *testing.T) {
	h := handler()
	ok := New(h).Post("/echo").JSON(`{"a":[1,2]}`).Cookie("session", "z").Run().
		Expect(t).
		Status(201).Success().
		Body(`{"a":[1,2]}`).Contains("[1,2]").NotContains("3").
		Header("Content-Type", "application/json").
		Cookie("seen", "z").
		JSON(map[string][]int{"a": {1, 2}}).
		Fn(func(t TestingT, r *Response) bool { return tt.Equal(t, 201, r.Code) }).
		OK()
	tt.True(t, ok)

	tt.True(t, New(h).Get("/redirect").Run().Expect(t).Redirect().OK())
	tt.True(t, New(h).Get("/fail").Run().Expect(t).Error().OK())
}

func TestFail(t *testing.T) {
	h := handler()
	r := New(h).Get("/hello").Run()
	mt := new(mock.T)

	tt.False(t, Status(mt, 404, r))
	tt.False(t, Redirect(mt, r))
	tt.False(t, Error(mt, r))
	tt.False(t, Body(mt, "nope", r))
	tt.False(t, Contains(mt, "nope", r))
	tt.False(t, NotContains(mt, "hello", r))
	tt.False(t, Header(mt, "X-Hello", "nope", r))
	tt.False(t, Cookie(mt, "seen", "x", r))
	tt.False(t, JSON(mt, `{}`, r))
	tt.False(t, JSON(mt, make(chan int), r))
	tt.False(t, Success(mt, New(h).Get("/fail").Run()))
	tt.False(t, Cookie(mt, "seen", "x",
		New(h).Post("/echo").Cookie("session", "y").Run()))
	tt.False(t, JSON(mt, `{"a":1}`, New(h).Post("/echo").JSON(`{"a":2}`).Run()))

	tt.Equal(t, 13, len(mt.Errors))
	tt.True(t, strings.Contains(mt.Errors[0], "Not Equal"))
	tt.True(t, strings.Contains(mt.Errors[0], "'404'"))
	tt.True(t, strings.Contains(mt.Errors[0], "'200'"))
	tt.True(t, strings.Contains(mt.Errors[1], "Not Redirect"))
	tt.True(t, strings.Contains(mt.Errors[2], "Not Error"))
	tt.True(t, strings.Contains(mt.Errors[4], "Not Contains"))
	tt.True(t, strings.Contains(mt.Errors[5], "not expected"))
	tt.True(t, strings.Contains(mt.Errors[6], "X-Hello: nope"))
	tt.True(t, strings.Contains(mt.Errors[7], "Cookie Not Set"))
	tt.True(t, strings.Contains(mt.Errors[8], "Invalid JSON"))
	tt.True(t, strings.Contains(mt.Errors[9], "Invalid JSON"))
	tt.True(t, strings.Contains(mt.Errors[10], "Not Success"))
	tt.True(t, strings.Contains(mt.Errors[11], "seen=y"))
	tt.True(t, strings.Contains(mt.Errors[12], `'{"a":2}'`))

	mt.Reset()
	ok := r.Expect(mt).Status(404).Body("x").Success().OK()
	tt.False(t, ok)
	tt.Equal(t, 2, len(mt.Errors))
	tt.True(t, r.Expect(mt).Status(200).OK())
}

func TestTrace(t *testing.T) {
	r := New(handler()).Get("/hello").Run()
	mt := new(mock.T)

	// Default depth points at the caller of the assertion.
	Status(mt, 404, r)
	tt.True(t, strings.Contains(mt.Errors[0], "http_test.go:202"))

	// Chained assertions point at the caller of the chain method.
	r.Expect(mt).Status(404)
	tt.True(t, strings.Contains(mt.Errors[1], "http_test.go:206"))

	// Explicit message, depth and trace.
	Status(mt, 404, r, "msg")
	tt.True(t, strings.Contains(mt.Errors[2], "Messages:	msg"))
	tt.True(t, strings.Contains(mt.Errors[2], "http_test.go:210"))

	Status(mt, 404, r, "", 3)
	tt.True(t, strings.Contains(mt.Errors[3], "http.go:"))

	Status(mt, 404, r, "", "custom:1")
	tt.True(t, strings.Contains(mt.Errors[4], "custom:1"))

	Status(mt, 404, r, "", 4, "custom:2")
	tt.True(t, strings.Contains(mt.Errors[5], "custom:2"))

	r.Expect(mt).Status(404, "m", "custom:3")
	tt.True(t, strings.Contains(mt.Errors[6], "custom:3"))
	tt.True(t, strings.Contains(mt.Errors[6], "Messages:	m"))
}

func TestJSONNorm(t *testing.T) {
	v, err := jsonNorm(struct{ A int }{1})
	tt.Nil(t, err)
	b, _ := json.Marshal(v)
	tt.Equal(t, `{"A":1}`, string(b))

	_, err = jsonNorm("{")
	tt.NotNil(t, err)
}
