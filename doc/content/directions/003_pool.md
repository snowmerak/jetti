---
title: 'Pool'
date: 2023-05-16T21:12:23+09:00
---

`Jetti` generate two type of pool, `sync.Pool` and `chan T`.

## sync pool

We can generate `sync.Pool` with `pool` and `sync:` direction.

```go
// jetti:pool sync:people
type Person struct {
    Name string
    Age  int
}
```

### Generated Code

Run `jetti generate` command, `jetti` generates the following code.

```go
package person

import "sync"
import "errors"
import "runtime"

var errPeopleCannotGet error = errors.New("cannot get people")

type PeoplePool struct {
	pool *sync.Pool
}

func (p *PeoplePool) Get() (*Person, error) {
	v := p.pool.Get()
	if v == nil {
		return nil, errPeopleCannotGet
	}
	return v.(*Person), nil
}

func (p *PeoplePool) GetWithFinalizer() (*Person, error) {
	v := p.pool.Get()
	if v == nil {
		return nil, errPeopleCannotGet
	}
	runtime.SetFinalizer(v, func(v interface{}) {
		p.pool.Put(v)
	})
	return v.(*Person), nil
}

func (p *PeoplePool) Put(v *Person) {
	p.pool.Put(v)
}

func NewPeoplePool() PeoplePool {
	return PeoplePool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(Person)
			},
		},
	}
}

func IsPeopleCannotGetErr(err error) bool {
	return errors.Is(err, errPeopleCannotGet)
}
```

## chan T

We can generate `chan T` pool with `pool` and `chan:` direction.

```go
// jetti:pool chan:candidate
type Person struct {
	Name string
	Age  int
}
```

### Generated Code

Run `jetti generate` command, `jetti` generates the following code.

```go
package person

import (
	"runtime"
	"time"
)

type CandidatePool struct {
	pool    chan *Person
	timeout time.Duration
}

func (c *CandidatePool) Get() *Person {
	after := time.After(c.timeout)
	select {
	case v := <-c.pool:
		return v
	case <-after:
		return new(Person)
	}
}

func (c *CandidatePool) GetWithFinalizer() *Person {
	after := time.After(c.timeout)
	resp := (*Person)(nil)
	select {
	case v := <-c.pool:
		resp = v
	case <-after:
		resp = new(Person)
	}
	runtime.SetFinalizer(resp, func(v interface{}) {
		c.pool <- v.(*Person)
	})
	return resp
}

func (c *CandidatePool) Put(v *Person) {
	select {
	case c.pool <- v:
	default:
	}
}

func NewCandidatePool(size int, timeout time.Duration) CandidatePool {
	pool := make(chan *Person, size)
	return CandidatePool{
		pool:    pool,
		timeout: timeout,
	}
}
```