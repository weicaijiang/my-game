package gate

import (
	"my-game/msg"
	"my-game/game"
	"my-game/login"
)

func init() {
	msg.Processor.SetRouter(&msg.Hello{},game.ChanRPC)

	//msg.Processor.SetRouter(&msg.WUser{},login.ChanRPC)
	msg.Processor.SetRouter(&msg.WeChatLogin{},login.ChanRPC)

	msg.Processor.SetRouter(&msg.RoomBase{},game.ChanRPC)

	//msg.Processor.SetRouter(&msg.LoginUser{},)
}
