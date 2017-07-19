package internal

import (
	"reflect"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"my-game/msg"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.Hello{}, handleHello)

	handler(&msg.RoomBase{},handleCreateRoom)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleHello(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.Hello)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("hello %v", m.Name)

	// 给发送者回应一个 Hello 消息
	a.WriteMsg(&msg.Hello{
		Name: "client",
	})
}

func handleCreateRoom(args []interface{})  {
	m := args[0].(*msg.RoomBase)

	a := args[1].(gate.Agent)
	log.Debug("room %v",m)
	log.Debug("a =%v",a.UserData().(*AgentInfo).accID)
	log.Debug("a =%v",a)
	a.WriteMsg(&msg.RoomBase{6})
}
