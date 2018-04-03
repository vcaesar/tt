# tt
Simple and colorful test tools

## Installation/Update

```
go get -u github.com/vcaesar/tt
```

## Usage:

#### [Look at an example](/examples/)

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
}
```
## Thanks

[Testify](https://github.com/stretchr/testify), the code has some inspiration.