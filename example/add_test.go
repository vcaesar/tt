package add

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt"
)

func TestAdd(t *testing.T) {
	fmt.Println(Add(1, 1))

	// tt.Expect(t, "1", Add(1, 1))
	tt.Expect(t, "2", Add(1, 1))
	// tt.Equal(t, 1, Add(1, 1))
	tt.Equal(t, 2, Add(1, 1))

	at := tt.New(t)
	at.Expect("2", Add(1, 1))
	at.Equal(2, Add(1, 1))
}
