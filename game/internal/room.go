package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/go"
)

type Room struct {
	RoomId	int //房间id
	RoomOwner int //房主
	RoomData *RoomData

	RoomUserId map[int]int //存放玩家座位 以及玩家ID
	Record *Record //记录 结算
	RoomState int //房间 状态 是否已开局


	GameStartVote int //游戏开始需要 房间人的 准备 记数
	StartGameChan chan int
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
	}
}