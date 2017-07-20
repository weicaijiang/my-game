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

	handleMsg(&msg.RoomBase{},handleCreateRoom)
	handleMsg(&msg.JoinRoom{},handleJoinRoom)

	handleRoom(&msg.ReadyGame{},handleReadyGame)
	handleRoom(&msg.QuitRoom{},handleQuitRoom)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), func(args []interface{}) {
		// user
		a := args[1].(gate.Agent)
		user := userLines[a.UserData().(*AgentInfo).userID]
		//fmt.Printf("user=%v\n",user)
		if user == nil ||user.State == userLogout{
			return
		}

		// agent to user
		args[1] = user
		h.(func([]interface{}))(args)
	})
}

func handleRoom(m interface{}, h interface{})  {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), func(args []interface{}) {
		a := args[1].(gate.Agent)
		user := userLines[a.UserData().(*AgentInfo).userID]
		if user == nil ||user.State == userLogout{
			return
		}
		if user.RoomId == ""{//都没有加入任何房间
			return
		}
		if _, ok := rooms[user.RoomId]; ok{//在房间里
			args[1] = user
			h.(func([]interface{}))(args)

		}else {//报错 该房间不存在
			user.WriteMsg(&msg.CodeState{msg.ERRO_NOTEXISITED,"房间不存在"})
			return 
		}
	})
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

	user := args[1].(*UserLine)
	user.createRoom(m)
	user.WriteMsg(&msg.CodeState{msg.SUCCESS_DONE,"创建成功!"})

	//accID := a.UserData().(*AgentInfo).accID
	//log.Debug("room %v",m)
	//log.Debug("a =%v",a.UserData().(*AgentInfo).accID)
	//log.Debug("a =%v",a)
	//a.WriteMsg(&msg.RoomBase{6})
}

func handleJoinRoom(args []interface{})  {
	m := args[0].(*msg.JoinRoom)
	user := args[1].(*UserLine)
	//fmt.Println("rooms information ",len(rooms))
	if room, ok := rooms[m.RoomAccID]; ok {//房间存在
		user.joinRoom(room)
		user.WriteMsg(&msg.CodeState{msg.SUCCESS_DONE,"加入房间成功!"})
	}else{//房间不存在
		user.WriteMsg(&msg.CodeState{msg.ERROR_Params,m.RoomAccID + "房间不存在!"})
	}
}
//退出房间
func handleQuitRoom(args []interface{})  {

	user := args[1].(*UserLine)
	user.quitRoom()
	user.WriteMsg(&msg.CodeState{msg.SUCCESS_DONE,"退出成功!"})
}

//准备 按钮
func handleReadyGame(args []interface{})  {

	user := args[1].(*UserLine)
	user.readGame()
	user.WriteMsg(&msg.ReadyGame{2})
	//user.WriteMsg(&msg.CodeState{msg.SUCCESS_DONE,"退出成功!"})
}

