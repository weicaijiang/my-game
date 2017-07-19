package internal

import (
	"reflect"
	"my-game/msg"
	"leaf/gate"
	"my-game/game"
	"fmt"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.WeChatLogin{},handleLogin)
}

func handleLogin(args []interface{})  {
	m := args[0].(*msg.WeChatLogin)
	a := args[1].(gate.Agent)
	fmt.Println("mmmmm",m)
	if m.Union == ""||len(m.Union) < 0{
		a.WriteMsg(&msg.CodeState{msg.ERROR_Params,"账号不能为空"})
		return
	}

	game.ChanRPC.Go("LoginAgent",m.Union,a)
	a.WriteMsg(&msg.CodeState{msg.SUCCESS_Register,"one success!"})
	fmt.Println("login success")

}