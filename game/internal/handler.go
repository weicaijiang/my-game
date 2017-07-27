package internal

import (
	"reflect"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"my-game/msg"
	"strconv"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.Hello{}, handleHello)

	handleMsg(&msg.RoomBase{},handleCreateRoom)
	handleMsg(&msg.JoinRoom{},handleJoinRoom)

	handleRoom(&msg.ReadyGame{},handleReadyGame)
	handleRoom(&msg.QuitRoom{},handleQuitRoom)

//	出牌
//	handleRoom(&msg.Card{},handleInRoomMyTime)
	handleInRoomMyTime(&msg.Card{},handleOneCardByIndex)

//	杠操作
	handleRoom(&msg.Gang{},handleGang)
	handleRoom(&msg.FangGang{},handleFangGang)

//	胡
	handleRoom(&msg.MimeHu{},handleMimeHu)
	handleRoom(&msg.FireHu{},handleFireHu)



//	碰
	handleRoom(&msg.Peng{},handlePeng)

//


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
			//fmt.Println("房间判断在---------")
			args[1] = user
			h.(func([]interface{}))(args)
		}else {//报错 该房间不存在
			user.WriteMsg(&msg.CodeState{msg.ERRO_NOTEXISITED,"房间不存在"})
			return 
		}
	})
}

func handleInRoomMyTime(m interface{}, h interface{})  {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), func(args []interface{}) {
		a := args[1].(gate.Agent)
		user := userLines[a.UserData().(*AgentInfo).userID]
		if user == nil ||user.State == userLogout{
			return
		}
		if user.RoomId == ""{//都没有加入任何房间
			return
		}
		if v, ok := rooms[user.RoomId]; ok{//在房间里 是否该我打
			if v.Playing == user.userData.AccID{
				args[1] = user
				h.(func([]interface{}))(args)
			}else{
				user.WriteMsg(&msg.CodeState{msg.FAILURE_DONE,"还没轮到你出牌"})
				return
			}

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
	user.WriteMsg(&msg.QuitRoom{2})
	//user.WriteMsg(&msg.CodeState{msg.SUCCESS_DONE,"退出成功!"})
}

//准备 按钮
func handleReadyGame(args []interface{})  {

	user := args[1].(*UserLine)
	user.readGame()
	user.WriteMsg(&msg.ReadyGame{2})
	//user.WriteMsg(&msg.CodeState{msg.SUCCESS_DONE,"退出成功!"})
}

//手牌处理 即摸到的牌 要打出的牌
func handleOneCardByIndex(args []interface{})  {
	m := args[0].(*msg.Card)
	//fmt.Println("打出的手牌有:",m.Index,m.Value)
	user := args[1].(*UserLine)
	user.playMyCards(m.Index,m.Value)
}

//杠处理 即点击确定杠 摸牌的人 明 暗杠
func handleGang(args []interface{})  {

	m := args[0].(*msg.Gang)
	user := args[1].(*UserLine)
	user.gangOK(m.GangType,m.Index,m.Value)
}


//碰处理 点击确定碰
func handlePeng(args []interface{})  {
	m := args[0].(*msg.Peng)
	user := args[1].(*UserLine)
	user.SumChan <- "peng+" + strconv.Itoa(m.Index)+"+"+strconv.Itoa(m.Value)
}

//杠处理 点击确定杠 明杠处理
func handleFangGang(args []interface{})  {
	m := args[0].(*msg.FangGang)
	//fmt.Println("放杠",m.Index,m.Value,m.GangType)
	user := args[1].(*UserLine)
	user.SumChan <- "gang+" + strconv.Itoa(m.Index) + "+" + strconv.Itoa(m.Value)

}

//点击碰牌后 需出一张牌
func handlePengOut(args []interface{})  {
	m := args[0].(*msg.Card)
	user := args[1].(*UserLine)
	user.MyTurn <- m.Value
	user.Cardings = append(user.Cardings[:m.Index],user.Cardings[m.Index+1:]...)
	user.rspAllCards()

}

//处理吃
func handleChi(args []interface{})  {
	//m := args[0].(msg.ChiPai)
	//user := args[1].(*UserLine)


}

//胡 自摸
func handleMimeHu(args []interface{})  {
	m := args[0].(*msg.MimeHu)
	user := args[1].(*UserLine)
	user.MyTurn <- m.HuType
}

//点炮
func handleFireHu(args []interface{})  {
	m := args[0].(*msg.FireHu)
	user := args[1].(*UserLine)
	user.SumChan <- "fire+" + strconv.Itoa(m.CardValue)
}
