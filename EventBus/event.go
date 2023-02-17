package EventBus

import (
	"github.com/asaskevich/EventBus"
)

var bus EventBus.Bus

func init() {
	bus = EventBus.New()
}

func Subscribe(topic string, fn interface{}) {
	bus.Subscribe(topic, fn)
}

func Publish(topic string, args ...interface{}) {
	bus.Publish(topic, args...)
}
