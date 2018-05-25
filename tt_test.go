package tt

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt/example"
)

func TestAdd(t *testing.T) {
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
}
