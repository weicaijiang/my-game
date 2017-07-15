package gate

import (
	"my-game/msg"
	"my-game/game"
)

func init() {
	msg.Processor.SetRouter(&msg.Hello{},game.ChanRPC)
}
