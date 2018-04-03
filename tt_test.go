package tt

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt/example"
)

func TestAdd(t *testing.T) {
	fmt.Println(add.Add(1, 1))

	// Expect(t, "1", add.Add(1, 1))
	Expect(t, "2", add.Add(1, 1))

	// Equal(t, 1, add.Add(1, 1))
	Equal(t, 2, add.Add(1, 1))
}
