package broadcast

import (
	"sync"
	"sync/atomic"
)

type Dispose func()

type Listener[T any] func(T)

type Observer[T any] interface {
	Subscribe(listener Listener[T]) Dispose
}

type Publisher[T any] interface {
	Observer[T]
	Publish(message T)
}

func New[T any](workers int) Publisher[T] {
	ch := make(chan func())

	for i := 0; i < workers; i++ {
		go execute(ch)
	}

	return &observerImpl[T]{
		listeners: sync.Map{},
		counter:   atomic.Int32{},
		ch:        ch,
	}
}

type observerImpl[T any] struct {
	counter   atomic.Int32
	listeners sync.Map
	ch        chan func()
}

func execute(channel chan func()) {
	for {
		select {
		case action := <-channel:
			action()
		}
	}
}

func (b *observerImpl[T]) Publish(message T) {
	b.listeners.Range(func(_, value any) bool {
		listener := value.(Listener[T])

		fun := func() {
			listener(message)
		}

		b.ch <- fun
		return true
	})
}

func (b *observerImpl[T]) Subscribe(listener Listener[T]) Dispose {
	id := b.counter.Add(1)

	b.listeners.Store(id, listener)

	return func() {
		b.listeners.Delete(id)
	}
}

type Transformer[T any, S any] func(event T) S

type mapObserver[T any, S any] struct {
	source    Observer[T]
	transform Transformer[T, S]
}

func (b *mapObserver[T, S]) Subscribe(listener Listener[S]) Dispose {
	unsubscribe := b.source.Subscribe(func(s T) {
		newValue := b.transform(s)
		listener(newValue)
	})
	return unsubscribe
}

func Map[T any, S any](source Observer[T], transformer Transformer[T, S]) Observer[S] {
	return &mapObserver[T, S]{
		source:    source,
		transform: transformer,
	}
}

type Filter[T any] func(event T) bool

type filterObserver[T any] struct {
	source    Observer[T]
	transform Filter[T]
}

func (b *filterObserver[T]) Subscribe(listener Listener[T]) Dispose {
	unsubscribe := b.source.Subscribe(func(s T) {
		canNext := b.transform(s)
		if canNext {
			listener(s)
		}
	})
	return unsubscribe
}

func Where[T any](source Observer[T], transformer Filter[T]) Observer[T] {
	return &filterObserver[T]{
		source:    source,
		transform: transformer,
	}
}
