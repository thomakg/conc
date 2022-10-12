package util

import (
	"sync"
)

type Concurrent[T any] struct {
	s  chan T
	e  chan error
	wg *sync.WaitGroup
}

func NewConcurrent[T any]() *Concurrent[T] {
	return &Concurrent[T]{
		s:  make(chan T),
		e:  make(chan error),
		wg: &sync.WaitGroup{},
	}
}

func (c *Concurrent[T]) RunFn(fn func() (T, error)) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()
		s, err := fn()
		if err != nil {
			c.e <- err
			return
		}
		c.s <- s
	}()
}

func (c *Concurrent[T]) WaitAndReturn() ([]T, []error) {
	res := []T{}
	errs := []error{}
	go func() {
		for s := range c.s {
			res = append(res, s)
		}
	}()
	go func() {
		for e := range c.e {
			errs = append(errs, e)
		}
	}()
	c.wg.Wait()
	close(c.s)
	close(c.e)
	return res, errs
}
