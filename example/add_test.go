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

func Benchmark1(b *testing.B) {
	at := tt.New(b)
	for i := 0; i < b.N; i++ {
		at.Equal(2, Add(1, 1))
	}
}

func Benchmark2(b *testing.B) {
	at := tt.New(b)
	fn := func() {
		at.Equal(2, Add(1, 1))
	}

	tt.BM(b, fn)
	// at.BM(b, fn)
}
