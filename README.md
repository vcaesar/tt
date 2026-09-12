# tt
Simple and colorful test tools

[![Build Status](https://github.com/vcaesar/tt/actions/workflows/go.yml/badge.svg)](https://github.com/vcaesar/tt/actions/workflows/go.yml)
[![CircleCI Status](https://circleci.com/gh/vcaesar/tt.svg?style=shield)](https://circleci.com/gh/vcaesar/tt)
[![codecov](https://codecov.io/gh/vcaesar/tt/branch/master/graph/badge.svg)](https://codecov.io/gh/vcaesar/tt)
[![Go Report Card](https://goreportcard.com/badge/github.com/vcaesar/tt)](https://goreportcard.com/report/github.com/vcaesar/tt)
[![Go Reference](https://pkg.go.dev/badge/github.com/vcaesar/tt.svg)](https://pkg.go.dev/github.com/vcaesar/tt)
[![Release](https://github-release-version.herokuapp.com/github/vcaesar/tt/release.svg?style=flat)](https://github.com/vcaesar/tt/releases/latest)
[![Join the chat at https://gitter.im/go-ego/ego](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-ego/ego?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Installation/Update

```
go get -u github.com/vcaesar/tt
```

## Usage:

#### [Look at an example](/example/)

```go
package tt

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt"
	"github.com/vcaesar/tt/example"
)

func TestAdd(t *testing.T) {
	fmt.Println(add.Add(1, 1))

	tt.Expect(t, "1", add.Add(1, 1))
	tt.Expect(t, "2", add.Add(1, 1))

	tt.Equal(t, 1, add.Add(1, 1))
	tt.Equal(t, 2, add.Add(1, 1))

	at := tt.New(t)
	at.Expect("2", add.Add(1, 1))
	at.Equal(2, add.Add(1, 1))
}

func Benchmark1(b *testing.B) {
	at := tt.New(b)
	fn := func() {
		at.Equal(2, add.Add(1, 1))
	}

	tt.BM(b, fn)
	// at.BM(b, fn)
}

func Benchmark2(b *testing.B) {
	at := tt.New(b)
	for i := 0; i < b.N; i++ {
		at.Equal(2, Add(1, 1))
	}
}

```
## Mock

`github.com/vcaesar/tt/mock` is a small mock object system. Embed `mock.Mock`
in a fake, register expectations with `On(...).Return(...)`, and forward calls
with `Called(...)`:

```go
package store

import (
	"errors"
	"testing"

	"github.com/vcaesar/tt"
	"github.com/vcaesar/tt/mock"
)

type Store struct{ mock.Mock }

func (s *Store) Get(id int) (string, error) {
	args := s.Called(id)
	return args.String(0), args.Error(1)
}

func TestStore(t *testing.T) {
	s := new(Store)
	s.On("Get", 1).Return("one", nil).Once()
	s.On("Get", mock.Anything).Return("", errors.New("not found"))

	v, err := s.Get(1)
	tt.Equal(t, "one", v)
	tt.Nil(t, err)

	s.AssertExpectations(t)
	s.AssertCalled(t, "Get", 1)
	s.AssertNumberOfCalls(t, "Get", 1)
}
```

Argument matchers: `mock.Anything`, `mock.AnythingOfType("int")`,
`mock.MatchedBy(func(v T) bool)`. Call options: `Once()`, `Twice()`,
`Times(n)`, `Run(func(mock.Arguments))`.

`mock.T` is a `tt.TestingT` that records failures instead of reporting them,
for testing assertions that are expected to fail:

```go
mt := new(mock.T)
tt.Equal(mt, 1, 2)
tt.True(t, mt.Failed())
```

## HTTP

`github.com/vcaesar/tt/http` tests `http.Handler`s in-process with
`net/http/httptest`. Build a request, `Run()` it, and assert on the response:

```go
package api

import (
	"testing"

	"github.com/vcaesar/tt"
	tthttp "github.com/vcaesar/tt/http"
)

func TestUsers(t *testing.T) {
	r := tthttp.New(handler).Get("/users").Param("id", "1").Run()
	tthttp.Status(t, 200, r)
	tthttp.Header(t, "Content-Type", "application/json", r)
	tthttp.JSON(t, `{"id":1,"name":"alice"}`, r)

	tthttp.New(handler).Post("/users").JSON(map[string]string{"name": "bob"}).
		Set("Authorization", "Bearer x").Run().
		Expect(t).Status(201).Contains("bob").Cookie("session", "s1")

	var u struct{ ID int }
	tt.Nil(t, r.JSON(&u))
}
```

Request builders: `Get/Post/Put/Patch/Delete/Head/Options/Do`, `Set/Headers`,
`Param/Params`, `Cookie`, `Text/Bytes/Reader/JSON/Form`. Assertions:
`Status`, `Success`, `Redirect`, `Error`, `Body`, `Contains`, `NotContains`,
`Header`, `Cookie`, `JSON` (structural), each also available on
`r.Expect(t)` as a chain ending in `OK()`. `Fetch(h, method, path, values)`
mirrors testify's helper signature.

## Thanks

[Testify](https://github.com/stretchr/testify), the code has some inspiration.