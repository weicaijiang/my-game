package internal

import (
	"github.com/name5566/leaf/gate"
	"my-game/msg"
	"fmt"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("RegisterAgent",rpcRegisterAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcRegisterAgent(args []interface{})  {
	m := args[0].(*msg.WUser)
	fmt.Println("m=",m)
	a := args[1].(gate.Agent)
	fmt.Println("a=%v",a)
	userdata := new(UserData)
	userdata.Name = m.Name
	userdata.Password = m.Password
	user := new(MyLine)
	user.MyUserData = userdata

	user.register()




}
