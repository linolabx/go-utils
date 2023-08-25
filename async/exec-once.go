package async

import (
	"sync"
)

type ExecOnce[T any] struct {
	mu        *sync.Mutex
	successed bool
	_func     func(p T) error
}

func ExecOnceWrap[T any](f func(p T) error) func(p T) error {
	return func(p T) error {
		return ExecOnceNew(f).Exec(p)
	}
}

func ExecOnceNew[T any](_func func(p T) error) *ExecOnce[T] {
	return &ExecOnce[T]{
		mu:        &sync.Mutex{},
		successed: false,
		_func:     _func,
	}
}

func (t *ExecOnce[T]) Exec(p T) error {
	if t.successed {
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.successed {
		return nil
	}

	err := t._func(p)
	if err != nil {
		return err
	}

	t.successed = true

	return nil
}
