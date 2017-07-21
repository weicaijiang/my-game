package internal

import (
	"github.com/name5566/leaf/go"
	"time"
	"fmt"
	"github.com/name5566/leaf/timer"
	"sync"
)

type Room struct {
	RoomId	string //房间id acc可以给玩家看到的
	RoomOwner int //房主 的id
	RoomData *RoomData

	RoomUserId map[int]string //存放玩家座位 以及玩家accID
	Record *Record //记录 结算
	RoomState int //房间 状态 是否已开局

	CardsBase []int //记录下这副牌
	DealCards []int //记录已发的牌
	OutCard int //记录出的牌

	pengFlag bool//碰的标记 一副牌中 一次只能一次碰
	gangFlag bool //杠的标记 一副牌中 一次只能一次杠

	chiHu []string //记录玩家 点击吃胡 的ACCID
	pengID []string //记录玩家 点击碰 的ACCID
	gangID []string //记录玩家 点击杠 的ACCID
	chiID []string //记录玩家 点击吃 的ACCID

	userWant map[string]map[int]string  //记录玩家 点击了 吃胡 吃 碰 杠等按钮  k0-{fire chi peng gang} k1(int)位座位 v为accID

	Playing string //正在操作的玩家 可以出牌的玩家
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
	r.pengFlag = false
	r.gangFlag = false
	r.StartGameChan = make(chan int,1)
	r.RoomUserId = make(map[int]string)
	//r.LinearContext = skeleton.NewLinearContext()

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
					if flag == 1 && len(r.RoomUserId) == r.RoomData.RoomVolume{
						r.RoomState = 1
						fmt.Println("信号 游戏开始")
						r.startGame()
					}else {
					//	游戏被某个玩家取消了
						fmt.Println("游戏被某个玩家因退出而取消了")

					}
				case <- time.After(time.Second * 30)://30s准备 自动开始
					if len(r.RoomUserId) == r.RoomData.RoomVolume{
						r.RoomState = 1
						fmt.Println("超时 游戏开始")
						r.startGame()
					}else{
						fmt.Println("什么都没做！")
					}
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
			//sort.Ints(user.Cardings)
			user.rspAllCards()
			i ++
		}

	//	记录已发的牌
	//	r.DealCards = append(r.DealCards,r.CardsBase[:4*13])
	//	r.DealCards = make([]int,4 * 13)
	//	copy(r.DealCards,r.CardsBase[:4*13])
		r.DealCards = r.CardsBase[:4*13]
	fmt.Println("已发的牌有:",r.DealCards)
	//开始 打牌
	r.Go(func() {
		r.playCard()
	}, func() {//游戏结束 一局
		//初始化 一些资源
		r.CardsBase = *new([]int)
		r.DealCards = *new([]int)
		r.RoomState = 0
	})
	//
	//}
}
//开始打牌
func (r *Room)playCard()  {
	//cp := make(chan int,1)
	//chance := <- cp
	//获取庄家的 座位号 即目前 谁是庄家
	//根据座位号进行 轮环
	var players map[int]UserLine
	players =make(map[int]UserLine)
	length := len(r.RoomUserId)
	for i:=0; i< length; i++{
		players[i] = *accIDUsers[r.RoomUserId[i]]
	}
	startPosition := userLines[r.RoomOwner].RoomPosition
	for i := 4 * 13; i<136 ; i++{
		pCard := r.CardsBase[i]
		r.DealCards = append(r.DealCards,pCard)
		r.OutCard = pCard
		//var turn int
		player := players[startPosition]
		player.Cardings = append(player.Cardings,pCard)
		player.rspAllCards()
		select {//自己牌的检测 是否可以胡 或者杠
			case r.OutCard = <- player.MyTurn:
			//告知 其他玩家 是某张牌 让其他玩家 做出反应
			//for j:=0; j< length; j++{
			//	if j == startPosition {//刚出了牌的玩家
			//	}else{//其他玩家看牌
			//	//	是否能碰 能杠 能吃 能点炮....
			//		if !r.pengFlag{//有一个玩家符合即可 其他玩家都不用再操作了
			//			r.pengFlag = players[j].isPeng(r.OutCard)
			//		}
			//		if !r.gangFlag{//有一个玩家符合即可 其他玩家都不用再操作了
			//			r.gangFlag = players[j].isGang(r.OutCard)
			//		}
			//		players[j].isChi(r.OutCard)
			//	}
			//}

			case <- time.After( 10 * time.Second)://超时 玩家没有动静 则出刚摸到的手牌 即 pCards == r.OutCard

		}


		//其他玩家 检测自己的手牌  /// r.OutCard 为玩家出的牌
		// 返回 信息给玩家
		wg := new(sync.WaitGroup)
		for j:=0; j< length; j++{
			if j == startPosition {//刚出了牌的玩家
			}else{//其他玩家看牌
				//	是否能碰 能杠 能吃 能点炮....
				wg.Add(1)
				go func(player UserLine,r *Room) {
					if !r.pengFlag{//有一个玩家符合即可 其他玩家都不用再操作了
						r.pengFlag = player.isPeng(r.OutCard)
					}
					if !r.gangFlag{//有一个玩家符合即可 其他玩家都不用再操作了
						r.gangFlag = player.isGang(r.OutCard)
					}
					player.isChi(r.OutCard)
					player.isChiHu(r.OutCard)
					wg.Done()
				}(players[j],r)
			}
		}
		wg.Wait()

		//玩家给回信息
		wgPlayer := new(sync.WaitGroup)
		for j := 0; j< length; j++{
			if j == startPosition{

			}else {
				wgPlayer.Add(1)
				go func(player UserLine, r *Room) {
					select {
					case flag := <- player.SumChan:
						switch flag {
						case "peng":
							fmt.Println("玩家ID：",player.userData.AccID," 叫碰")
							(r.userWant["peng"])[player.RoomPosition] = player.userData.AccID
						case "gang":
							fmt.Println("玩家ID：",player.userData.AccID," 叫杠")
							(r.userWant["gang"])[player.RoomPosition] = player.userData.AccID
						case "chi":
							fmt.Println("玩家ID：",player.userData.AccID," 叫吃")
							(r.userWant["chi"])[player.RoomPosition] = player.userData.AccID
						case "fire":
							fmt.Println("玩家ID：",player.userData.AccID," 叫胡")
							(r.userWant["fire"])[player.RoomPosition] = player.userData.AccID
						default:
							fmt.Println("玩家ID：",player.userData.AccID," pass")
						}
						wgPlayer.Done()
					case <- time.After(5 * time.Second)://5s 反应
						fmt.Println("玩家ID：",player.userData.AccID," 超时反应")
					}
					wgPlayer.Done()
				}(players[j],r)
			}
		}
		wgPlayer.Wait()
		//处理 玩家给回的信息 需求的操作
		if len(r.userWant) == 0{
		//	所有玩家 没有需求
		} else {
			if len(r.userWant["fire"]) >0{
				if len(r.userWant["fire"]) == 1{//只有一个玩家吃胡
				//	初始化
					r.initRoom()
				}else{//多个玩家吃胡
					for i:=1 ; i <= 3; i++{//
						if id, ok :=(r.userWant["fire"])[startPosition+i]; ok{//是下家胡?
							fmt.Println("玩家iD",id,"吃胡了")
							r.initRoom()
							break
						//	初始化 游戏结束
						}
					}
					break
				}
			}else if len(r.userWant["gang"]) == 1{
			//

			}else if len(r.userWant["peng"]) == 1{
			//	碰优先于吃

			}else if len(r.userWant["chi"]) >0 {
				if len(r.userWant["chi"]) == 1{//一家吃

				}else{//多家吃
					for i:=1; i <= 3; i++{
						if id, ok :=(r.userWant["chi"])[startPosition+i]; ok{//是下家吃?
							fmt.Println("玩家iD",id,"吃")
						}
					}
				}
			}else{

			}
		}


		startPosition = (startPosition + 1) % r.RoomData.RoomVolume
	}

}

