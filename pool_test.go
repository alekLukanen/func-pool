package funcpool

import (
	"context"
	"errors"
	"testing"
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

func TestFuncPool(t *testing.T) {
	ctx := context.Background()

	pool := NewFuncPool(ctx, 2, 0)
	pool.Start()

	f1 := &FuncMock{Val1: 0}
	pool.AddFunc(f1)
	<-pool.Results

	f2 := &FuncMock{Val1: 1}
	pool.AddFunc(f2)
	<-pool.Results

	pool.Stop()

	if f1.ReturnValue != "even" && f1.Done != true && f1.Error != nil {
		t.Errorf("Expected nil, got %v", f1)
	}
	if f2.ReturnValue != "odd" && f2.Done != true && f2.Error != ErrMock {
		t.Errorf("Expected ErrMock, got %v", f2)
	}

}

func TestSetOfFuncs(t *testing.T) {
	size := 10_000

	ctx := context.Background()
	pool := NewFuncPool(ctx, 2, size)
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

	evenCount := 0
	oddCount := 0
	doneCount := 0
	errorCount := 0
	for _, f := range funcs {
		if f.ReturnValue == "even" {
			evenCount++
		} else if f.ReturnValue == "odd" {
			oddCount++
		}
		if f.Done {
			doneCount++
		}
		if f.Error != nil {
			errorCount++
		}
	}

	if evenCount != size/2 {
		t.Errorf("Unexpected even count, got %d", evenCount)
	}
	if oddCount != size/2 {
		t.Errorf("Unexpected odd count, got %d", oddCount)
	}
	if doneCount != size {
		t.Errorf("Unexpected done count, got %d", doneCount)
	}
	if errorCount != size/2 {
		t.Errorf("Unexpected error count, got %d", errorCount)
	}

}
