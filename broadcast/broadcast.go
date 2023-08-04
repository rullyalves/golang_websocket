package broadcast

import "sync"

func IsCOpen(ch <-chan interface{}) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}

type Broadcaster struct {
	Listeners sync.Map
}

type Listener struct {
	ID      string
	Channel chan interface{}
	OnData  func(message interface{})
}

func (b *Broadcaster) Publish(message interface{}) {
	b.Listeners.Range(func(_, value any) bool {
		listener := value.(Listener)

		if !IsCOpen(listener.Channel) {
			b.Unsubscribe(listener)
			return true
		}
		listener.OnData(message)
		return true
	})
}

func (b *Broadcaster) Subscribe(listener Listener) {
	b.Listeners.Store(listener.ID, listener)
}

func (b *Broadcaster) Unsubscribe(listener Listener) {
	b.Listeners.Delete(listener.ID)
}
