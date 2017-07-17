package internal

import (
	"sort"
	"time"
	"my-game/aglorithm"
)

type UserLine struct {
	UserData *User
	RoomId int
	ReadySign int //准备信号

	MyTurn	chan int //轮到我 的信号
	RoomSignal chan int //房间释放信号 该玩家操作

	CardMap	map[int]int //用来存放手牌的 用map 利于删除(打出) key 为牌的值 value为 张数
	PenCards map[int]int //用于存放碰 用map 利于升级为 杠 同时利于删除操作 key为牌的值 value暂时未定
	GangCards []int //用于存放杠
}

var userLines map[int]UserLine //保存玩家 key为userid value 为userline

func init()  {
	userLines = make(map[int]UserLine)
}

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
func (u *UserLine)addOneRealCard(value int)  {
	if v, ok := u.CardMap[value]; ok{
		u.CardMap[value] = v + 1 //已有 数量加 一
	}else {
		u.CardMap[value] = 1
	}
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

//玩家 信号切换 准备
func (u *UserLine)readGame()  {
	u.ReadySign = u.ReadySign ^ 1
	if u.ReadySign{
		room := rooms[u.RoomId]
		room.GameStartVote ++
		if room.GameStartVote == room.RoomData.RoomVolume{
			room.StartGameChan <- 1
		}
	}else{
		rooms[u.RoomId].GameStartVote --
	}
}

//创建房间
//房间 创建成功 后 就一直在等待 开始
func (u *UserLine) CreateRoom() bool {
	if u.RoomId == 0 {
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
//value 为摸起牌的值
func (u *UserLine)playMyCards(value int)  {

	var room Room
	select {
	case room = rooms[<-u.RoomSignal]:
	//	玩家操作

	case <-time.After( 5 * time.Second)://玩家超时 还没打 将自动打出 摸起的牌
		u.outOneRealCard(value)
		//room := rooms[u.RoomId]
	}
	room.RoomTurn <- (room.RoomUserId[u.UserData.Id] + 1) % 4 //下一家 摸牌等 权限
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
				if room.PlayerTurn == u.UserData.Id && u.canPen(dealCard){
				//	碰后 手牌 变化等操作
				}

			//	是否能出牌
				if room.PlayerTurn == u.UserData.Id {
				//	玩家出牌
					var value int //接收 玩家要出的牌 从前端
					if u.outOneRealCard(value){
						room.OutCartds = value
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