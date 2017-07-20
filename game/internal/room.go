package internal

import (
	"github.com/name5566/leaf/go"
	"time"
	"fmt"
	"github.com/name5566/leaf/timer"
	"sort"
)

type Room struct {
	RoomId	string //房间id
	RoomOwner int //房主 的id
	RoomData *RoomData

	RoomUserId map[int]string //存放玩家座位 以及玩家accID
	Record *Record //记录 结算
	RoomState int //房间 状态 是否已开局

	CardsBase []int //记录下这副牌
	DealCards []int //记录已发的牌
	OutCards int //记录出的牌


	GameStartVote int //游戏开始需要 房间人的 准备 记数
	PlayerSignal	chan int//玩家传过 来的信息
	StartGameChan chan int
	RoomTurn chan int	//玩家 轮环
	PlayerTurn int //是否是 该玩家能打牌 或者 其他傍边等待的
	*g.LinearContext
	roomTimer *timer.Timer //进行房间的一个监督 是否人数满
	gameTimer *timer.Timer //游戏开始 监督
	//RoomMapUser	map[int]map[]
	//CardMap	map[int]int //用来存放手牌的 用map 利于删除(打出) key 为牌的值 value为 张数
	//PenCards map[int]int //用于存放碰 用map 利于升级为 杠 同时利于删除操作 key为牌的值 value暂时未定
	//GangCards []int


}

var rooms = make(map[string]*Room)//key为accID

func (r *Room)initRoom()  {
	r.Record = nil
	r.RoomState = 0
	r.RoomUserId = make(map[int]string)
	r.LinearContext = skeleton.NewLinearContext()

}


func (r *Room)isFull() bool {
	return len(r.RoomUserId) == r.RoomData.RoomVolume
}

func AddCreateRoom(newRoom *Room)  {
	rooms[newRoom.RoomData.RoomAccID] = newRoom
//	当有一个房间创建时 进行房间的一个监督 是否人数满 或者 没人这删掉房间
//	自动检测 房间是否可以开始游戏
	newRoom.autoCheckRoomState()
	newRoom.autoStartGame()
}

func (r *Room)deleteRoom()  {
	if len(rooms) >0 {
		r.roomTimer.Stop()
		delete(rooms,r.RoomData.RoomAccID)
	}
}

//每隔3分钟 自动检测房间是否还需要存在
func (r *Room)autoCheckRoomState()  {
	const duration  = 3 * 60 * time.Second
	r.roomTimer = skeleton.AfterFunc(duration, func() {
		r.Go(func() {
			length := len(r.RoomUserId)
			if length < 1{//房间没有玩家 则自动删除房间
				r.gameTimer.Stop()//停止对游戏 是否可以开始的自动检测
				r.deleteRoom()
			}
			fmt.Println("检查房间是否要关闭 3m")
		}, func() {
			r.autoCheckRoomState()
		})

	})
}

//每隔 5s 自动检测 当人数 齐时 将3s后 自动开始游戏
//除非退出房间 否则其他退出 无效
func (r *Room)autoStartGame()  {
	const duration  = 5 * time.Second
	r.gameTimer = skeleton.AfterFunc(duration, func() {
		r.Go(func() {
			fmt.Println("5s 检查游戏是否可以开始")
			length := len(r.RoomUserId)
			if length == r.RoomData.RoomVolume && r.RoomState == 0{
				fmt.Println("进入 10 等待游戏开始")
				//房间 满人后将10s后开始游戏
				select {
				case flag :=<- r.StartGameChan://玩家都点了准备 马上开始
					if flag == 1{
						r.RoomState = 1
						fmt.Println("信号 游戏开始")
						//r.startGame()
					}else {
					//	游戏被某个玩家取消了
						fmt.Println("游戏被某个玩家因退出而取消了")

					}
				case <- time.After(time.Second * 10)://10s准备 自动开始
					r.RoomState = 1
					//r.startGame()
					fmt.Println("超时 游戏开始")
				}
			}
		}, func() {
			r.autoStartGame()
		})
	})
}

func (r *Room)startGame()  {
	//val := <- r.StartGameChan
	//if val == 1{
	//	游戏开始
		r.CardsBase = ShuffleCards() //洗牌
	//	发牌
		fmt.Println("洗了牌")
		i := 0
		for _,v := range r.RoomUserId{
			user := accIDUsers[v]
			user.Cardings = r.CardsBase[(0 + i * 13) : (13 + i * 13)]
			//向前端 发送牌型
			sort.Ints(user.Cardings)
			user.rspAllCards()
			i ++
		}

	//	记录已发的牌
	//	r.DealCards = append(r.DealCards,r.CardsBase[:4*13])
	//	r.DealCards = make([]int,4 * 13)
	//	copy(r.DealCards,r.CardsBase[:4*13])
		r.DealCards = r.CardsBase[:4*13]
	//
	//}
}
//开始打牌
func (r *Room)playCard()  {
	//cp := make(chan int,1)
	//chance := <- cp
	for i := 4 * 13; i<136 ; i++{
		pCards := r.CardsBase[i]
		r.DealCards = append(r.DealCards,pCards)
		r.OutCards = pCards
		//var turn int
		for j := 0; j < 4; j++{//四个人进行循环
			select {
			case r.PlayerTurn = <- r.PlayerSignal://房间 操作


			case <-time.After(5 * time.Second)://房间等待玩家 操作 已超时
				fmt.Println("time over")
				//userLines[r.RoomUserId[r.PlayerTurn]].RoomSignal <- 1

			}
		}
	}

}

