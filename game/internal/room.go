package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/go"
	"time"
	"fmt"
)

type Room struct {
	RoomId	int //房间id
	RoomOwner int //房主
	RoomData *RoomData

	RoomUserId map[int]int //存放玩家座位 以及玩家ID
	Record *Record //记录 结算
	RoomState int //房间 状态 是否已开局

	CardsBase []int //记录下这副牌
	DealCards []int //记录已发的牌
	OutCartds int //记录出的牌


	GameStartVote int //游戏开始需要 房间人的 准备 记数
	PlayerSignal	chan int//玩家传过 来的信息
	StartGameChan chan int
	RoomTurn chan int	//玩家 轮环
	PlayerTurn int //是否是 该玩家能打牌 或者 其他傍边等待的
	gate.Agent
	*g.LinearContext
	//RoomMapUser	map[int]map[]
	//CardMap	map[int]int //用来存放手牌的 用map 利于删除(打出) key 为牌的值 value为 张数
	//PenCards map[int]int //用于存放碰 用map 利于升级为 杠 同时利于删除操作 key为牌的值 value暂时未定
	//GangCards []int


}

var rooms map[int]*Room

func init()  {
	rooms = make(map[int]*Room)
}

func (r *Room)isFull() bool {
	return len(r.RoomUserId) == r.RoomData.RoomVolume
}

func AddCreateRoom(newRoom *Room)  {
	rooms[newRoom.RoomId] = newRoom
}

func (r *Room)startGame()  {
	val := <- r.StartGameChan
	if val == 1{
	//	游戏开始
		r.CardsBase = ShuffleCards() //洗牌
	//	发牌
		var userLine1 UserLine
		var userLine2 UserLine
		var userLine3 UserLine
		var userLine4 UserLine
		for _, userId := range r.RoomUserId{
			userLine1 = userLines[userId]
			userLine2 = userLines[userId]
			userLine3 = userLines[userId]
			userLine4 = userLines[userId]
		}
		for i:=0; i< 4* 13; i = i+4{
			if v,ok := userLine1.CardMap[r.CardsBase[i]]; ok{
				userLine1.CardMap[r.CardsBase[i]] = v + 1
			}else {
				userLine1.CardMap[r.CardsBase[i]] = 1
			}
			if v,ok := userLine2.CardMap[r.CardsBase[i + 1]]; ok{
				userLine2.CardMap[r.CardsBase[i + 1]] = v + 1
			}else {
				userLine2.CardMap[r.CardsBase[i + 1]] = 1
			}
			if v,ok := userLine3.CardMap[r.CardsBase[i + 2]]; ok{
				userLine3.CardMap[r.CardsBase[i + 2]] = v + 1
			}else {
				userLine3.CardMap[r.CardsBase[i + 2]] = 1
			}
			if v,ok := userLine4.CardMap[r.CardsBase[i + 3]]; ok{
				userLine4.CardMap[r.CardsBase[i + 3]] = v + 1
			}else {
				userLine4.CardMap[r.CardsBase[i + 3]] = 1
			}
		}

	//	记录已发的牌
	//	r.DealCards = append(r.DealCards,r.CardsBase[:4*13])
		r.DealCards = make([]int,4 * 13)
		copy(r.DealCards,r.CardsBase[:4*13])
	//
	}
}
//开始打牌
func (r *Room)playCard()  {
	//cp := make(chan int,1)
	//chance := <- cp
	for i := 4 * 13; i<136 ; i++{
		pCards := r.CardsBase[i]
		r.DealCards = append(r.DealCards,pCards)
		r.OutCartds = pCards
		//var turn int
		for j := 0; j < 4; j++{//四个人进行循环
			select {
			case r.PlayerTurn = <- r.PlayerSignal://房间 操作


			case <-time.After(5 * time.Second)://房间等待玩家 操作 已超时
				fmt.Println("time over")
				userLines[r.RoomUserId[r.PlayerTurn]].RoomSignal <- 1

			}
		}
	}

}

