package common

import (
	"fmt"
	"github.com/asaskevich/EventBus"
)

var bus = EventBus.New()

func Subscribe(mainCmd string, subCmd string, fn interface{}) {
	bus.Subscribe(fmt.Sprintf("%v_%v", mainCmd, subCmd), fn)
}

func Publish(mainCmd string, subCmd string, args ...interface{}) {
	bus.Publish(fmt.Sprintf("%v_%v", mainCmd, subCmd), args...)
}
