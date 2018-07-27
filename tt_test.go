package tt

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt/example"
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
	Bool(t, 1 == 1)
	True(t, 1 == 1)
	False(t, 1 != 1)

	at := New(t)
	at.Expect("2", add.Add(1, 1))
	at.Equal(2, add.Add(1, 1))

	at.Nil(nil)
	at.Bool(1 == 1)
	at.True(1 == 1)
	at.False(1 != 1)
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

	at := New(t)
	at.Not("3", add.Add(1, 1))
	at.NotExpect("3", add.Add(1, 1))
	at.NotEqual(3, add.Add(1, 1))
}

func TestArgs(t *testing.T) {
	Expect(t, "2", add.Add(1, 1), 3)
	Equal(t, 2, add.Add(1, 1), 4)

	NotExpect(t, "3", add.Add(1, 1), 3)
	Not(t, 3, add.Add(1, 1), 4)
	NotEqual(t, 3, add.Add(1, 1), 4)

	Nil(t, nil, 3)
	Bool(t, 1 == 1, 4)
	True(t, 1 == 1, 4)
	False(t, 1 != 1, 4)

	at := New(t)
	at.Expect("2", add.Add(1, 1), 4)
	at.Equal(2, add.Add(1, 1), 5)

	at.NotExpect("3", add.Add(1, 1), 4)
	at.Not("3", add.Add(1, 1), 5)
	at.NotEqual(3, add.Add(1, 1), 5)

	at.Nil(nil, 4)
	at.Bool(1 == 1, 4)
	at.True(1 == 1, 4)
	at.False(1 != 1, 4)
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
	at := New(b)
	fn := func() {
		at.Equal(2, add.Add(1, 1))
	}

	BM(b, fn)
	at.BM(b, fn)
}
