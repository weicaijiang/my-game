package internal

import (
	"sort"
	"time"
	"my-game/aglorithm"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/go"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/log"
	"my-game/msg"
	"github.com/name5566/leaf/timer"
	"github.com/name5566/leaf/util"
	"gopkg.in/mgo.v2"
	"fmt"
	"my-game/mjlib"
	"math"
)

type UserLine struct {
	gate.Agent
	*g.LinearContext
	userData *User
	RoomId string
	ReadySign int //准备信号
	State int //玩家状态
	saveDBTimer *timer.Timer
	RoomPosition int //座位号


	MyTurn	chan int //轮到我出牌的信号 即值为 我出的牌的值
	SumChan chan string //总共信息号  用来确定玩家 操作了 碰(peng) 杠(gang) 吃(chi) 点炮(fire)等 超时未 none
	RoomSignal chan int //房间释放信号 该玩家操作

	CardMap	map[int]int //用来存放手牌的 用map 利于删除(打出) key 为牌的值 value为 张数
	PenCards map[int]int //用于存放碰 用map 利于升级为 杠 同时利于删除操作 key为牌的值 value暂时未定
	GangCards []int //用于存放杠

	Cardings []int //存放手牌
	PengCardings []int //存放碰的牌
	GangCardings []int //存放杠的牌
}

const  (
	userLogin = iota
	userLogout
	userGame
)

//var userLines map[int]UserLine //保存玩家 key为userid value 为userline

var (
	accIDUsers = make(map[string]*UserLine)
	userLines = make(map[int]*UserLine)
)

func GetLineUsers()int  {
	count := len(userLines)
	return count
}

//登录
func (u *UserLine)login(name string)  {
	fmt.Println("userLiners =",GetLineUsers())
	userData := new(User)
	skeleton.Go(func() {
		db := mongoDB.Ref()
		defer mongoDB.UnRef(db)

		err := db.DB(DBName).C(C_USERS).Find(bson.M{"accid":name}).One(userData)
		if err != nil{
			fmt.Println("err db")
			if err != mgo.ErrNotFound{
				log.Error("load acc %v data error: %v", name, err)
				userData = nil
				u.WriteMsg(&msg.CodeState{msg.ERRO_Unkonw,"未知错误!"})
				u.Close()
				return
			}

		//	new
			err := userData.initValue(name)

			if err != nil{
				log.Error("init acc %v data error: %v",name, err)
				userData = nil
				u.WriteMsg(&msg.CodeState{msg.ERRO_InitValue,"初始值报错initValue!"})
				u.Close()
				return
			}
		}
	}, func() {
		if u.State == userLogout{
			u.logout(name)
			return
		}

		u.State = userGame
		if userData == nil{
			return
		}
		userData.LastLoginTime = int(time.Now().Unix())
		u.userData = userData
		u.Cardings = make([]int,0,200)
		userLines[u.userData.Id]= u
		u.UserData().(*AgentInfo).userID = userData.Id
		fmt.Println("userline=%v",u)
		u.autoSaveDB()
	})
}

//退出
func (u *UserLine)logout(name string)  {
	if u.userData != nil{
		u.saveDBTimer.Stop()
		//delete(accIDUsers,name)
		delete(userLines,u.userData.Id)
	}

	data := util.DeepClone(u.userData)
	u.Go(func() {
		if data != nil {
			db := mongoDB.Ref()
			defer mongoDB.UnRef(db)
			userID := data.(*User).Id
			_, err := db.DB(DBName).C(C_USERS).
				UpsertId(userID, data)
			if err != nil {
				log.Error("save user %v data error: %v", userID, err)
			}
		}
	}, func() {
		delete(accIDUsers, name)
	})
}
//自动存储
func (u *UserLine)autoSaveDB()  {
	//const duration  = 5 * time.Minute
	const duration  = 15 * time.Second

	u.saveDBTimer = skeleton.AfterFunc(duration, func() {
		data := util.DeepClone(u.userData)
		u.Go(func() {
			db := mongoDB.Ref()
			defer mongoDB.UnRef(db)
			userID := data.(*User).Id
			_, err := db.DB(DBName).C(C_USERS).UpsertId(userID, data)
			if err != nil{
				log.Error("save user %v data error",data)
				//return
			}
			fmt.Println("insert db")

		}, func() {
			u.autoSaveDB()
		})
	})
}

//创建房间
func (u *UserLine)createRoom(m interface{})  {
	if u.RoomId != ""{//已有房间
		u.WriteMsg(&msg.CodeState{msg.ERRO_ONLYROOM,"你已有房间了，请退出当前房间后操作!"})
		return
	}
	newRoom := new(RoomData)
	err := newRoom.initValue()
	if err != nil{
		u.WriteMsg(&msg.CodeState{msg.ERRO_InitValue,"内部错误，暂时无法创建房间!"})
		return
	}
	newRoom.RoomVolume = m.(*msg.RoomBase).Volume
	room := new(Room)
	room.RoomData = newRoom
	room.RoomOwner = u.userData.Id
	room.LinearContext = skeleton.NewLinearContext()
	room.initRoom()
	room.RoomUserId[0] = u.userData.AccID
	u.RoomPosition = 0
	AddCreateRoom(room)
	u.RoomId = room.RoomData.RoomAccID
	u.WriteMsg(&msg.RoomDataInfo{RoomID:newRoom.RoomID,RoomAccID:newRoom.RoomAccID,
	RoomType:newRoom.RoomType,RoomVolume:newRoom.RoomVolume,RoomPay:newRoom.RoomPay,RoomBaseMoney:newRoom.RoomBaseMoney,
	CreatedTime:newRoom.CreatedTime})

}

//加入房间
func (u *UserLine)joinRoom(room *Room)  {
	if u.RoomId != ""{
		u.WriteMsg(&msg.CodeState{msg.ERRO_ONLYROOM,"你已有房间了，请退出当前房间后操作!"})
		return
	}else{
		if len(room.RoomUserId) >= room.RoomData.RoomVolume{
			u.WriteMsg(&msg.CodeState{msg.FAILURE_DONE, "房间满了!"})
			return
		}else {
			u.RoomId = room.RoomData.RoomAccID
			for i := 0; i < room.RoomData.RoomVolume; i++{
				if _,ok := room.RoomUserId[i];!ok{
					room.RoomUserId[i] = u.userData.AccID
					u.RoomPosition = i
					fmt.Println("玩家ID=",u.userData.AccID+" 加入房间ID=",u.RoomId+" 座位号为:",i)
					break
				}
			}
		}
	}
}




//退出房间
//注意：游戏中 是无法退出的 除非 自己掉线了 将是 变成托管
//房间信号 置为 0 即 r.StartGameChan <- 0
func (u *UserLine)quitRoom()  {
	fmt.Println("quit--room",u.RoomId)
	if u.RoomId == ""{

	}else {
		fmt.Println("退出第一步")
		room := rooms[u.RoomId]
		room.RoomState = 0
		room.StartGameChan <- 0
		//go func() {
		//	if len(room.RoomUserId) == room.RoomData.RoomVolume{//说明本来是满人的 某人中途退出
		//		fmt.Println("置为0 退出")
		//		room.StartGameChan <- 0
		//		fmt.Println("<-0 ok")
		//	}else if len(room.RoomUserId) == 1{//最后一人退出
		//
		//	}
		//}()
		delete(room.RoomUserId,u.RoomPosition)
		u.RoomId = ""
		u.RoomPosition = 0
		fmt.Println("发送了消息")

	}
}

//整理牌 全部的牌 返回给前端
func (u *UserLine)rspAllCards()  {
	//sort.Ints(u.PengCardings)
	//sort.Ints(u.GangCardings)
	sort.Ints(u.Cardings)
	u.WriteMsg(&msg.Cards{u.Cardings,u.PengCardings,u.GangCardings})
}

//func InitUserLine(a gate.Agent)(u *UserLine)  {
//	user := new(UserLine)
//
//}

//将手牌为map类型 生成数组 并排好序 返回给前端
func (u *UserLine)getRealCards(cardMap map[int]int) (a []int){
	if len(cardMap) <=0{
		return 
	}
	for i,v := range cardMap{
		if v > 1{//是2 3 4个同样的牌
			for j:=0;j<v;j++{
				a = append(a,i)
			}
		}else{
			a = append(a,i)
		}
	}
	sort.Ints(a)
	return a
}

//打出手牌 处理
// 参数 value 为牌的值
func (u *UserLine)outOneRealCard(value int) bool  {
	if v, ok := u.CardMap[value]; ok {
		left := v - 1
		if left <= 0{
			delete(u.CardMap,value)
		}else {
			u.CardMap[value] = left
		}
		return true //正常出牌
	}else {
		return false //不能正常出牌 因为 没有这张牌 出异常
	}
}

//摸牌 处理 添加一张手牌
//参数 value 为牌的值
func (u *UserLine)addOneRealCard(value int) bool  {
	if v, ok := u.CardMap[value]; ok{
		u.CardMap[value] = v + 1 //已有 数量加 一
	}else {
		u.CardMap[value] = 1
	}
	return true
}

//检查 手牌是否能杠 摸完牌或者别人打出一张牌
//参数 value 为牌的值
func (u *UserLine)canGangByRealCard(value int)  bool {
	if v, ok := u.CardMap[value]; ok && v == 4 {
		u.GangCards = append(u.GangCards,value)
		return true
	}else {
		return false
	}
}

//手牌杠 后变化 以及杠牌组 变化
func (u *UserLine)deleteRealCardByGang(value int)  {
	delete(u.CardMap, value)
	u.GangCards = append(u.GangCards,value)
}

//检查 碰的牌 是否能杠
////参数 value 为牌的值
func (u *UserLine)canGangByPenCards(value int)bool  {
	if _, ok := u.PenCards[value]; ok {
		return true
	}else {
		return false
	}
}

// 碰变成杠
//参数 value 为牌的值
func (u *UserLine)gangByPenCards(value int)  {
	delete(u.PenCards,value)
	u.GangCards = append(u.GangCards,value)
}

//判断 是否能碰
func (u *UserLine) canPen(value int) bool  {
	if v, ok := u.CardMap[value]; ok && v == 2{
		return true
	}
	return false
}

//碰后 手牌变化 以及碰牌增加
func (u *UserLine)penDealCards(value int)  {
	delete(u.CardMap,value)
	u.PenCards[value] = 1
}

//玩家 信号切换 准备 //准备开始游戏
//准备后 只会加快游戏开始 不会取消游戏不开始
func (u *UserLine)readGame()  {
	fmt.Println("ready")
	u.ReadySign = u.ReadySign ^ 1
	room := rooms[u.RoomId]
	if u.ReadySign == 1{
		room.GameStartVote ++
		if room.GameStartVote == room.RoomData.RoomVolume{
			room.StartGameChan <- 1
			fmt.Println("都发送了 准备")
		}
	}else{
		room.GameStartVote --

	}
}

//创建房间
//房间 创建成功 后 就一直在等待 开始
func (u *UserLine) CreateRoom() bool {
	if u.RoomId == "" {
		room := Room{}
		AddCreateRoom(&room)
		room.Go(func() {
			room.startGame()
		}, func() {
		//游戏开始 洗牌 发牌等操作

		})

		return true
	}else{
		return false
	}
}

//出牌
//value 为摸起牌的值 index 为下标
func (u *UserLine)playMyCards(index,value int)  {
	if len(u.Cardings) >= index{
		if u.Cardings[index] == value{//存在 这张牌
			u.Cardings = append(u.Cardings[:index],u.Cardings[index+1:]...)
			u.rspAllCards()
		//	告知打出了牌
			u.MyTurn <- value
		}
	}

}


//玩家操作
func (u *UserLine)userPlayCard()  {
	for {
		room := rooms[u.RoomId]
		dealCard := room.DealCards[len(room.DealCards)]
		select {
		case <- u.RoomSignal:
		//	玩家操作
		//	摸牌
			if u.addOneRealCard(dealCard){
			//	摸完牌后 打牌或者判断是否胡了
			//	首先 判断是否能胡
				if aglorithm.IsHu(u.getRealCards(u.CardMap),u.CardMap){
					
				}
			//	是否能杠
				if u.canGangByRealCard(dealCard){
				//	杠后 手牌 变化等操作

				}
			//	是否是 他人 若是 则进行 是否能碰
			//
				if room.PlayerTurn == u.userData.Id && u.canPen(dealCard){
				//	碰后 手牌 变化等操作
				}

			//	是否能出牌
				if room.PlayerTurn == u.userData.Id {
				//	玩家出牌
					var value int //接收 玩家要出的牌 从前端
					if u.outOneRealCard(value){
						room.OutCard = value
						//rooms[u.RoomId].PlayerSignal <- (room.PlayerTurn + 1) % 4
					}
				}
				rooms[u.RoomId].PlayerSignal <- (room.PlayerTurn + 1) % 4

			}
		//case <- time.After(5 * time.Second):
		//	rooms[u.RoomId].PlayerSignal <- (room.PlayerTurn + 1) % 4

		}

	}
}

//碰操作
func (u *UserLine)isPeng(value int)  bool{
	for i := 0; i< len(u.Cardings)-1; i++{
		if u.Cardings[i] == value && u.Cardings[i] == u.Cardings[i+1]{
		//	可以碰
			u.WriteMsg(&msg.Peng{i,value})
			return true
		}
	}
	return false
}

//确定碰操作
func (u *UserLine)pengOK(index int)  {
	u.PengCardings = append(u.PengCardings,u.Cardings[index])
	u.Cardings = append(u.Cardings[:index],u.Cardings[index+2:]...)

	//u.rspAllCards()

}

//杠操作 杠别人 不是自己摸的 放杠
func (u *UserLine)fangGang(value int) bool  {
	for i := 0; i <len(u.Cardings)-2; i++{
		if u.Cardings[i] == value && u.Cardings[i+1] == u.Cardings[i] && u.Cardings[i+2] == u.Cardings[i]{
			u.WriteMsg(&msg.Gang{i,value,113})
			return true
		}
	}
	return false
}

//暗杠 自己摸的 判断
func (u *UserLine)anGang()bool  {
	for i := 0; i < len(u.Cardings) - 3; i++{
		if u.Cardings[i] == u.Cardings[i+1] && u.Cardings[i] == u.Cardings[i+2] && u.Cardings[i] == u.Cardings[i+3]{
			u.WriteMsg(&msg.Gang{i,u.Cardings[i],111})
			return true
		}
	}
	return false
}

//明杠 判断
//value 位摸上的牌
func (u *UserLine)mingGang(value int) bool {
	for i:=0; i<len(u.PengCardings); i++{
		if u.PengCardings[i] == value{
			u.WriteMsg(&msg.Gang{i,value,112})
			return true
		}
	}
	return false
}

//玩家确定杠 处理
func (u *UserLine)gangOK(gangType,index int,pengIndex int)  {
	if gangType == 111 && pengIndex == 0{//暗杠
		u.GangCardings = append(u.GangCardings,u.Cardings[index])
		u.Cardings = append(u.Cardings[:index],u.Cardings[index+4:]...)
		u.MyTurn <- 111
	}else if gangType == 112 && pengIndex !=0 && pengIndex < len(u.PengCardings){//明杠
		if u.PengCardings[pengIndex] == u.Cardings[index]{
			u.GangCardings = append(u.GangCardings,u.Cardings[index])
			u.PengCardings = append(u.PengCardings[:pengIndex],u.PengCardings[pengIndex+1:]...)
			u.Cardings = append(u.Cardings[:index],u.Cardings[index+1:]...)
			u.MyTurn <- 112
		}else {
			u.WriteMsg(&msg.CodeState{msg.ERROR_Params,"参数有误!"})
		}

	}else if gangType == 113 && pengIndex ==0 {//放杠
		u.GangCardings = append(u.GangCardings,u.Cardings[index])
		u.Cardings = append(u.Cardings[:index],u.Cardings[index+3:]...)
		u.MyTurn <- 113
	}else {

	}
}

//获取手牌的 某张牌的具体的数量
func (u *UserLine)getMapCards() (mapCards map[int]int) {
	mapCards = make(map[int]int)
	for i:=0; i< len(u.Cardings); i++{
		if v,ok := mapCards[u.Cardings[i]]; ok{
			mapCards[u.Cardings[i]] = v +1
		}else {
			mapCards[u.Cardings[i]] = 1
		}
	}
	return

}

//吃牌
//转化为map 记录个数
func (u *UserLine)isChi(value int) bool {
	if len(u.Cardings) >= 4{ //吃必须手牌4张以上
		mapCards := u.getMapCards()
		for i, _:= range mapCards{
			if math.Abs(float64(i - value)) ==2{
				if v1,ok1 := mapCards[value + 1] ;ok1&&(v1<10||mapCards[i]<10){
				//有一个 3 5  ==>有4
					mapCards[value +1] += 10
					mapCards[i] += 10
				}
				if v2,ok2 := mapCards[value-1]; ok2&&(v2<10||mapCards[i]<10){
					//1 3
					mapCards[value -1] +=10
					mapCards[i] +=10
				}
			}else if math.Abs(float64(i- value)) == 1{
				// 3 4 value =3 i=4  3 4 5
 				if v1,ok1 := mapCards[value+2]; ok1&&(v1<10||mapCards[i]<10){
					mapCards[value +2 ] += 10
					mapCards[i] += 10
				}
				// 2 3 4
				if v2, ok2 := mapCards[value-1]; ok2&&(v2<10||mapCards[i]<10){
					mapCards[value-1] += 10
					mapCards[i] += 10
				}
				// 2 3 value=3 i=2    2 3 4
				if v3, ok3 := mapCards[value + 1]; ok3 &&(v3 <10 || mapCards[i]<10){
					mapCards[value+1] += 10
					mapCards[i] += 10
				}
				// 1 2 3
				if v4, ok4 := mapCards[value -2]; ok4&&(v4<10&&mapCards[i]<10){
					mapCards[value-2] += 10
					mapCards[i] += 10
				}
			}
		}

	}
	return false
}
//胡牌
// 参数 value为牌的值 wIndex 王牌的下标 wValue 王牌的值
func (u *UserLine)isChiHu(value int,wIndex,wValue int) bool {
	a := make([]int,len(u.Cardings))
	a = u.Cardings
	a = append(a,value)
	if mjlib.IsHu(a,wIndex,wValue){
		u.WriteMsg(&msg.MimeHu{0,wValue})
		return true
	}
	return false
}