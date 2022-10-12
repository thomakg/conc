package util_test

import (
	"fmt"
	"testing"

	"github.com/thomakg/conc/util"
)

var f = staticCounter()

func staticCounter() (f func() int) {
	var i int
	f = func() int {
		i++
		return i
	}
	return
}

func specialWait() (string, error) {
	i := f()
	if i%2 == 0 {
		return "", fmt.Errorf("failed at %d", i)
	}
	return fmt.Sprintf("%d", i), nil
}
func TestConc(t *testing.T) {
	it := 20000

	conc := util.NewConcurrent[string]()
	for y := 0; y < it; y++ {
		conc.RunFn(specialWait)
	}
	s, e := conc.WaitAndReturn()

	count := 0
	errCount := 0
	for range s {
		count++
	}

	for range e {
		errCount++
	}

	if count+errCount != it {
		t.Fatalf("expected count of %d, but got %d", it, count+errCount)
	}
}
