package tt

import (
	"fmt"
	"testing"

	add "github.com/vcaesar/tt/example"
)

func TestTT(t *testing.T) {
	mockT := new(testing.T)
	fmt.Println(Expect(mockT, "1", add.Add(1, 1)))

	if Expect(mockT, "1", add.Add(1, 1)) {
		t.Error("Equal should return false")
	}
	if Equal(mockT, "1", add.Add(1, 1)) {
		t.Error("Equal should return false")
	}

	if !Expect(mockT, "Hello World", "Hello World") {
		t.Error("Equal should return true")
	}
	if !Equal(mockT, "Hello World", "Hello World") {
		t.Error("Equal should return true")
	}

	fmt.Println(add.Add(1, 1))

	// Expect(t, "1", add.Add(1, 1))
	Expect(t, "2", add.Add(1, 1))
	// Equal(t, 1, add.Add(1, 1))
	Equal(t, 2, add.Add(1, 1))

	Nil(t, nil)
	Empty(t, "")
	Zero(t, 0)
	Bool(t, 1 == 1)
	True(t, 1 == 1)
	False(t, 1 != 1)

	at := New(t)
	at.Expect("2", add.Add(1, 1))
	at.Equal(2, add.Add(1, 1))

	at.Nil(nil)
	at.Empty("")
	at.Zero(0)
	at.Bool(1 == 1)
	at.True(1 == 1)
	at.False(1 != 1)

	True(t, Pprof(1))
}

func TestNot(t *testing.T) {
	mockT := new(testing.T)
	fmt.Println(Expect(mockT, "1", add.Add(1, 1)))

	if !Not(mockT, "1", add.Add(1, 1)) {
		t.Error("Equal should return false")
	}
	if !NotExpect(mockT, "1", add.Add(1, 1)) {
		t.Error("Equal should return false")
	}
	if !NotEqual(mockT, "1", add.Add(1, 1)) {
		t.Error("Equal should return false")
	}

	if Not(mockT, "Hello World", "Hello World") {
		t.Error("Equal should return true")
	}
	if NotExpect(mockT, "Hello World", "Hello World") {
		t.Error("Equal should return true")
	}
	if NotEqual(mockT, "Hello World", "Hello World") {
		t.Error("Equal should return true")
	}

	fmt.Println(add.Add(1, 1))

	Not(t, "3", add.Add(1, 1))
	NotExpect(t, "3", add.Add(1, 1))
	NotEqual(t, 3, add.Add(1, 1))
	//
	NotNil(t, "nil", "", "tt_test.go:86")
	NotEmpty(t, 1)
	NotZero(t, "1")

	at := New(t)
	at.Not("3", add.Add(1, 1))
	at.NotExpect("3", add.Add(1, 1))
	at.NotEqual(3, add.Add(1, 1))
	//
	at.NotNil("nil", "", "tt_test.go:95")
	at.NotEmpty(1)
	at.NotZero("1")
}

func TestArgs(t *testing.T) {
	Expect(t, "2", add.Add(1, 1), "must be expect", 3)
	Equal(t, 2, add.Add(1, 1), "", 4)

	NotExpect(t, "3", add.Add(1, 1), "", 3)
	Not(t, 3, add.Add(1, 1), "", 4)
	NotEqual(t, 3, add.Add(1, 1), "", 4)

	Nil(t, nil, "", 3)
	Empty(t, "", "", 3)
	Bool(t, 1 == 1, "", 4)
	True(t, 1 == 1, "", 4)
	False(t, 1 != 1, "", 4)

	at := New(t)
	at.Expect("2", add.Add(1, 1), "", 4)
	at.Equal(2, add.Add(1, 1), "", 5)

	at.NotExpect("3", add.Add(1, 1), "", 4)
	at.Not("3", add.Add(1, 1), "", 5)
	at.NotEqual(3, add.Add(1, 1), "", 5)

	at.Nil(nil, "", 4)
	at.Empty("", "", 4)
	at.Bool(1 == 1, "", 4)
	at.True(1 == 1, "", 4)
	at.False(1 != 1, "", 4)
}

func TestType(t *testing.T) {
	Type = true
	Equal(t, 1, 1)
	Nil(t, nil)

	mockT := new(testing.T)
	if Equal(mockT, 1, "1", "must be equal", "tt/tt_test.go:125") {
		t.Error("Equal should return false")
	}

	if NotEqual(mockT, 1, 1, "must be not equal") {
		t.Error("Equal should return false")
	}

	if !NotEqual(mockT, 1, "1") {
		t.Error("Equal should return true")
	}

	DEqual(t, 1, 1, "", "tt/tt_test.go:137")
}

func TestDbg(t *testing.T) {
	Nil(t, Log("log test"))
	// Nil(t, Err("err test"))
}

func Benchmark1(b *testing.B) {
	at := New(b)
	at.Bool(true)
	for i := 0; i < b.N; i++ {
		// do something
		// fmt.Println("do something")
		at.Equal(2, add.Add(1, 1))
	}
}

func Benchmark2(b *testing.B) {
	fn := func() {
		Equal(b, 2, add.Add(1, 1))
	}

	BM(b, fn)
}

func Benchmark3(b *testing.B) {
	at := New(b)
	fn := func() {
		at.Equal(2, add.Add(1, 1))
	}

	at.BM(b, fn)
}
