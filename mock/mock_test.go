package mock

import (
	"errors"
	"strings"
	"testing"

	"github.com/vcaesar/tt"
)

// store is an example fake built on Mock.
type store struct {
	Mock
}

func (s *store) Get(id int) (string, error) {
	args := s.Called(id)
	return args.String(0), args.Error(1)
}

func (s *store) Put(key string, val []byte) bool {
	return s.Called(key, val).Bool(0)
}

func (s *store) Count() int {
	return s.Called().Int(0)
}

func TestT(t *testing.T) {
	mt := new(T)
	tt.False(t, mt.Failed())

	// tt assertions must accept mock.T as a TestingT.
	tt.False(t, tt.Equal(mt, 1, 2))
	tt.True(t, mt.Failed())
	tt.Equal(t, 1, len(mt.Errors))
	tt.True(t, strings.Contains(mt.Errors[0], "Not Equal"))

	mt.Reset()
	tt.False(t, mt.Failed())
	tt.Equal(t, 0, len(mt.Errors))

	tt.True(t, tt.Equal(mt, 1, 1))
	tt.False(t, mt.Failed())
}

func TestCalledReturn(t *testing.T) {
	s := new(store)
	s.On("Get", 1).Return("one", nil)
	s.On("Get", 2).Return("", errors.New("not found"))
	s.On("Count").Return(2)

	v, err := s.Get(1)
	tt.Equal(t, "one", v)
	tt.Nil(t, err)

	v, err = s.Get(2)
	tt.Empty(t, v)
	tt.Expect(t, "not found", err)

	tt.Equal(t, 2, s.Count())
	tt.Equal(t, 3, len(s.Calls))
	tt.True(t, s.AssertExpectations(t))
	tt.True(t, s.AssertCalled(t, "Get", 1))
	tt.True(t, s.AssertNotCalled(t, "Get", 3))
	tt.True(t, s.AssertNumberOfCalls(t, "Get", 2))
	tt.True(t, s.AssertNumberOfCalls(t, "Count", 1))
}

func TestUnexpectedCallPanics(t *testing.T) {
	s := new(store)
	s.On("Get", 1).Return("one", nil)

	defer func() {
		r := recover()
		tt.NotNil(t, r)
		msg, _ := r.(string)
		tt.True(t, strings.Contains(msg, "unexpected call Get(2)"))
		tt.True(t, strings.Contains(msg, "closest expectation Get(1)"))
	}()
	s.Get(2)
}

func TestNoExpectationPanics(t *testing.T) {
	s := new(store)
	defer func() {
		msg, _ := recover().(string)
		tt.True(t, strings.Contains(msg, `no expectations registered for Count`))
	}()
	s.Count()
}

func TestOnceTimes(t *testing.T) {
	s := new(store)
	s.On("Count").Return(1).Once().
		On("Count").Return(2).Times(2)

	tt.Equal(t, 1, s.Count())
	mt := new(T)
	tt.False(t, s.AssertExpectations(mt))
	tt.True(t, mt.Failed())

	tt.Equal(t, 2, s.Count())
	tt.Equal(t, 2, s.Count())
	tt.True(t, s.AssertExpectations(t))

	// All expectations are exhausted now.
	defer func() {
		msg, _ := recover().(string)
		tt.True(t, strings.Contains(msg, "already been consumed"))
	}()
	s.Count()
}

func TestMatchers(t *testing.T) {
	s := new(store)
	s.On("Put", Anything, []byte("v")).Return(true)
	s.On("Put", AnythingOfType("string"), MatchedBy(func(b []byte) bool {
		return len(b) > 1
	})).Return(false)

	tt.True(t, s.Put("any", []byte("v")))
	tt.False(t, s.Put("k", []byte("long")))
	tt.True(t, s.AssertCalled(t, "Put", "any", Anything))
	tt.True(t, s.AssertCalled(t, "Put", Anything, AnythingOfType("[]uint8")))

	// nil against a nilable MatchedBy type is passed as the zero value;
	// nil against a non-nilable type never matches and must not panic.
	tt.True(t, MatchedBy(func(b []byte) bool { return b == nil }).Matches(nil))
	tt.False(t, MatchedBy(func(i int) bool { return true }).Matches(nil))
	tt.False(t, MatchedBy(func(i int) bool { return true }).Matches("s"))
	tt.False(t, AnythingOfType("int").Matches(nil))

	defer func() {
		tt.NotNil(t, recover())
	}()
	MatchedBy(func() {})
}

func TestRun(t *testing.T) {
	s := new(store)
	var seen []int
	s.On("Get", Anything).Return("", nil).Run(func(args Arguments) {
		seen = append(seen, args.Int(0))
	})
	s.Get(7)
	s.Get(8)
	tt.Equal(t, []int{7, 8}, seen)
}

func TestAssertFailures(t *testing.T) {
	s := new(store)
	s.On("Get", 1).Return("one", nil)

	mt := new(T)
	tt.False(t, s.AssertExpectations(mt))
	tt.False(t, s.AssertCalled(mt, "Get", 1))
	tt.False(t, s.AssertNumberOfCalls(mt, "Get", 1))
	tt.Equal(t, 3, len(mt.Errors))
	tt.True(t, strings.Contains(mt.Errors[0], "expected call not made: Get(1)"))
	tt.True(t, strings.Contains(mt.Errors[1], "calls: none"))

	s.Get(1)
	mt.Reset()
	tt.False(t, s.AssertNotCalled(mt, "Get", 1))
	tt.True(t, mt.Failed())

	s.Reset()
	tt.Equal(t, 0, len(s.Calls))
	tt.Equal(t, 0, len(s.ExpectedCalls))
}

func TestArguments(t *testing.T) {
	err := errors.New("e")
	args := Arguments{"s", 1, true, err, nil, []byte("b")}
	tt.Equal(t, "s", args.String(0))
	tt.Equal(t, 1, args.Int(1))
	tt.True(t, args.Bool(2))
	tt.Equal(t, err, args.Error(3))
	tt.Nil(t, args.Error(4))
	tt.Equal(t, []byte("b"), args.Get(5))

	tt.True(t, args.Is("s", 1, true, err, nil, []byte("b")))
	tt.False(t, args.Is("s"))
	tt.False(t, Arguments{1}.Is("1"))
	tt.False(t, Arguments{nil}.Is(1))

	for _, fn := range []func(){
		func() { args.Get(9) },
		func() { args.String(1) },
		func() { args.Int(0) },
		func() { args.Bool(0) },
		func() { args.Error(0) },
	} {
		func() {
			defer func() { tt.NotNil(t, recover()) }()
			fn()
		}()
	}
}
