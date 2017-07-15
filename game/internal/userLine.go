package internal

import "sort"

type UserLine struct {
	UserData *User
	RoomId int
	ReadySign int //准备信号

	CardMap	map[int]int //用来存放手牌的 用map 利于删除(打出) key 为牌的值 value为 张数
	PenCards map[int]int //用于存放碰 用map 利于升级为 杠 同时利于删除操作 key为牌的值 value暂时未定
	GangCards []int //用于存放杠
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

//