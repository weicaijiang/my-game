package internal

import (
	"github.com/name5566/leaf/gate"
	"my-game/msg"
	"fmt"
	"leaf/log"
)

type AgentInfo struct {
	accID string
	userID int
}

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("RegisterAgent",rpcRegisterAgent)
	skeleton.RegisterChanRPC("LoginAgent",rpcLoginAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	//_ = a
	a.SetUserData(new(AgentInfo))
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	accID := a.UserData().(*AgentInfo).accID
	a.SetUserData(nil)

	user := accIDUsers[accID]
	if user == nil {
		return
	}

	log.Debug("acc %v logout", accID)

	// logout
	if user.State == userLogin {
		user.State = userLogout
	} else {
		user.State = userLogout
		user.logout(accID)
	}
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

func rpcLoginAgent(args []interface{})  {
	union := args[0].(string)
	//fmt.Println("union=",union)
	a := args[1].(gate.Agent)
	//if a.UserData() == nil{
	//	fmt.Println("userDate==nil")
	//	return
	//}

//	login repeated
	oldUser := accIDUsers[union]
	if oldUser != nil{
		a.WriteMsg(&msg.CodeState{msg.ERRO_LoginRepeated,"重复登录！"})
		oldUser.WriteMsg(&msg.CodeState{msg.ERRO_LoginRepeated,"重复登录!"})
		a.Close()
		oldUser.Close()
		log.Debug("repeated login name=",union)
		return
	}

	log.Debug("name %v login",union)

	newUserLine := new(UserLine)
	newUserLine.Agent = a
	newUserLine.LinearContext = skeleton.NewLinearContext()
	newUserLine.State = userLogin
	a.UserData().(*AgentInfo).accID = union
	accIDUsers[union] = newUserLine
	newUserLine.login(union)
}

func rpcCreateRoom(args []interface{})  {

}