# func-pool

A simple function pool executor for Go. Allows you to scatter 
func calls go routines and gather results from your functions. 
Look at the tests in this repository for examples on how to
use this package.

Example:
```go
package main

import (
	"context"
	"errors"
  "fmt"

  "github.com/alekLukanen/func-pool
)

var (
	ErrMock = errors.New("mock error")
)

type FuncMock struct {
	Val1 int

	ReturnValue string
	Done        bool
	Error       error
}

func (obj *FuncMock) Call(ctx context.Context) struct{} {
	defer func() {
		obj.Done = true
	}()
	if obj.Val1%2 == 0 {
		obj.ReturnValue = "even"
		return struct{}{}
	} else {
		obj.ReturnValue = "odd"
		obj.Error = ErrMock
		return struct{}{}
	}
}

func main() {
  size := 10

	ctx := context.Background()
	pool := funcpool.NewFuncPool(ctx, 2, size)
	pool.Start()

	funcs := make([]FuncMock, 0, size)
	for i := 0; i < size; i++ {
		funcs = append(funcs, FuncMock{Val1: i})
		pool.AddFunc(&funcs[i])
	}

	countReceived := 0
	for pool.HasResults() {
		countReceived++
		<-pool.Results
	}

	pool.Stop()
 
  fmt.Println(countReceived)
  for _, f := range funcs {
    fmt.Println(f.ReturnValue, f.Done, f.Error)
  }

}

```

Output:
```shell
10
even true <nil>
odd true mock error
even true <nil>
odd true mock error
even true <nil>
odd true mock error
even true <nil>
odd true mock error
even true <nil>
odd true mock error
```

