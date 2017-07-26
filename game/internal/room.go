package internal

import (
	"github.com/name5566/leaf/go"
	"time"
	"fmt"
	"github.com/name5566/leaf/timer"
	"sync"
	"my-game/mjlib"
	"my-game/msg"
	"strings"
	"strconv"
	"sort"
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
	WValueCard int //记录王牌的值
	WIndexCard int //记录王牌的 index

	pengFlag bool//碰的标记 一副牌中 一次只能一次碰
	gangFlag bool //杠的标记 一副牌中 一次只能一次杠
	pengUserID	string //需要碰的玩家的ACCID
	gangUserID string //需要杠的玩家 ACCID

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
				case <- time.After(time.Second * 10)://30s准备 自动开始
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

//获取王牌的index
func (r *Room)getWIndex()  {
	i := r.WValueCard / 100
	w := (i - 1) * 9 + (r.WValueCard % 100) -1
	r.WIndexCard = w

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
			sort.Ints(user.Cardings)
			fmt.Println("玩家id:",v,"长度有:",len(user.Cardings),"开始手牌有:",user.Cardings)

			//向前端 发送牌型
			//sort.Ints(user.Cardings)
			user.rspAllCards()
			i ++
		}

	//	记录已发的牌
	//	r.DealCards = append(r.DealCards,r.CardsBase[:4*13])
	//	r.DealCards = make([]int,4 * 13)
	//	copy(r.DealCards,r.CardsBase[:4*13])
	//	r.DealCards = r.CardsBase[:4*13]
	//fmt.Println("已发的牌有:",r.DealCards)
	//fmt.Println("已发的牌的长度:",len(r.DealCards))
	fmt.Println("原始牌有:",r.CardsBase)
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
	//var players map[int]*UserLine
	players := make(map[int]*UserLine)
	//players := new(map[int]*UserLine)
	length := len(r.RoomUserId)
	//var mutex sync.Mutex
	for i:=0; i< length; i++{
		players[i] = accIDUsers[r.RoomUserId[i]]
		fmt.Println("players[i]",players[i].Cardings)
		fmt.Println("lenp------",len(players[i].Cardings))
		fmt.Println("lenp------",cap(players[i].Cardings))
	}
	fmt.Println("players=len=",len(players))
	startPosition := userLines[r.RoomOwner].RoomPosition
	fmt.Println("起始位置开始:",startPosition)
	for i := 4 * 13; i<136 ; i++{
		r.OutCard = r.CardsBase[i]
		fmt.Println("摸到的牌是pCards:",r.OutCard)
		fmt.Println("摸到的牌的前一张:",r.CardsBase[i-1])
		//r.DealCards = append(r.DealCards,pCard)
		//var turn int
		//player := players[startPosition]
		player := new(UserLine)
		player = players[startPosition]
		//mutex.Lock()
		players[startPosition].Cardings = append(players[startPosition].Cardings,r.OutCard)
		//mutex.Unlock()
		fmt.Println("玩家id:",player.userData.AccID,"的手牌为:",player.Cardings)
		for j:=0; j<length; j++{
			fmt.Println("第ddttttddvvdjj",j,"次")
			fmt.Println("第ddttttvvdddjj",players[j].userData.AccID)
			fmt.Println("第ddtttvvvdddjj",players[j].Cardings,"次")
			fmt.Println("lenp--dddd----",len(players[j].Cardings))
			fmt.Println("lenp--dddd----",cap(players[j].Cardings))
		}
		r.Playing = player.userData.AccID
		player.rspAllCards()
		//for j:=0; j<length; j++{
		//	fmt.Println("第ddddvvdddd",j,"次")
		//	fmt.Println("第ddvvddd",players[j].userData.AccID)
		//	fmt.Println("第ddvvvddd",players[j].Cardings,"次")
		//}

		//go func(p *UserLine,room *Room) {
			if mjlib.IsHu(player.Cardings,r.WIndexCard,r.WValueCard){
				player.WriteMsg(&msg.MimeHu{1,r.OutCard})
			}
		for j:=0; j<length; j++{
			fmt.Println("第ddtttddd",j,"次")
			fmt.Println("第ddttddd",players[j].userData.AccID)
			fmt.Println("第ddttddd",players[j].Cardings,"次")
		}
			fmt.Println("1玩家iD:",player.userData.AccID,"进入检验胡牌后的牌:",player.Cardings)
			if player.anGang()|| player.mingGang(r.OutCard){

			}
		for j:=0; j<length; j++{
			fmt.Println("第ddd",j,"次")
			fmt.Println("第ddd",players[j].userData.AccID)
			fmt.Println("第ddd",players[j].Cardings,"次")
		}
			fmt.Println("2玩家iD:",player.userData.AccID,"进入检验胡牌后的牌:",player.Cardings)
		//}(player,r)
		select {//自己牌的检测 是否可以胡 或者杠
			case flag  := <- player.MyTurn:
				if flag == 100 {//自个胡牌
					player.WriteMsg(&msg.MimeHu{1,r.OutCard})
					break
				}else if flag == 111{//自个杠 为暗杠

					continue
				}else if flag == 112{//明杠
					continue
				}else {//没有胡与杠 出牌 放杠不在这里处理
					r.OutCard = flag
				}
			case <- time.After( 15 * time.Second)://超时 玩家没有动静 则出最后一张牌
				r.OutCard = player.Cardings[len(player.Cardings)-1]
				player.Cardings = player.Cardings[:(len(player.Cardings)-1)]
				fmt.Println("玩家ID为:",player.userData.AccID,"超时打出的牌为:",r.OutCard,"此时的手牌有:",player.Cardings)
				player.rspAllCards()
		}


		//其他玩家 检测自己的手牌  /// r.OutCard 为玩家出的牌
		// 返回 信息给玩家
		//Again: wg := new(sync.WaitGroup)
		//for j:=0; j< length; j++{
		//	if j == startPosition {//刚出了牌的玩家
		//
		//	}else{//其他玩家看牌
		//		//	是否能碰 能杠 能吃 能点炮....
		//		wg.Add(1)
		//		go func(player UserLine,r *Room) {
		//			fmt.Println("kaishi:",player.Cardings)
		//			player.isChiHu(r.OutCard,r.WIndexCard,r.WValueCard)//吃胡 判断
		//			//吃胡后的牌
		//			fmt.Println("吃胡后 判断:",player.Cardings)
		//			if !r.pengFlag{//有一个玩家符合即可 其他玩家都不用再操作了
		//				r.pengFlag = player.isPeng(r.OutCard)
		//			}
		//			if !r.gangFlag{//有一个玩家符合即可 其他玩家都不用再操作了
		//				r.gangFlag = player.fangGang(r.OutCard)
		//			}
		//			//player.isChi(r.OutCard)
		//			wg.Done()
		//		}(*players[j],r)
		//	}
		//}
		//wg.Wait()



		//将碰 杠开启 下一轮
		r.pengFlag = false
		r.gangFlag = false

		//玩家给回信息
		wgPlayer := new(sync.WaitGroup)
		for j := 0; j< length; j++{
			if j == startPosition{

			}else {
				fmt.Println("第",j,"次")
				fmt.Println("第",players[j].userData.AccID)
				fmt.Println("第",players[j].Cardings,"次")
				wgPlayer.Add(1)
				go func(player UserLine, r *Room) {
					fmt.Println("玩家id",player.userData.AccID)
					select {
					case flag := <- player.SumChan:
						strFlagArray := strings.Split(flag,"+")
						switch strFlagArray[0] {
						case "peng":
							fmt.Println("玩家ID：",player.userData.AccID," 叫碰")
							(r.userWant["peng"])[player.RoomPosition] = player.userData.AccID + "+" + strFlagArray[1]
							r.pengUserID = player.userData.AccID
						case "gang":
							fmt.Println("玩家ID：",player.userData.AccID," 叫杠")
							(r.userWant["gang"])[player.RoomPosition] = player.userData.AccID + "+" + strFlagArray[1]
							r.gangUserID = player.userData.AccID
						case "chi":
							fmt.Println("玩家ID：",player.userData.AccID," 叫吃")
							(r.userWant["chi"])[player.RoomPosition] = player.userData.AccID
						case "fire":
							fmt.Println("玩家ID：",player.userData.AccID," 叫胡")
							(r.userWant["fire"])[player.RoomPosition] = player.userData.AccID
						default:
							fmt.Println("玩家ID：",player.userData.AccID," pass")
						}
						//wgPlayer.Done()
					case <- time.After(15 * time.Second)://5s 反应
						fmt.Println("玩家ID：",player.userData.AccID," 超时反应","手牌为:",player.Cardings)
					}
					fmt.Println("完成一次:")
					wgPlayer.Done()
				}(*players[j],r)
			}
		}
		wgPlayer.Wait()
		//处理 玩家给回的信息 需求的操作
		if len(r.userWant) == 0{
			//所有玩家 没有需求
			//顺序 进行游戏
			fmt.Println("len===",startPosition)
			startPosition = (startPosition + 1) % r.RoomData.RoomVolume
			fmt.Println("lenPPPP===",startPosition)
		} else {
			if len(r.userWant["fire"]) >0{
				if len(r.userWant["fire"]) == 1{//只有一个玩家吃胡
				//	初始化
				//	r.Playing
					r.initRoom()
				}else{//多个玩家吃胡
					for i:=1 ; i <= 3; i++{//
						if id, ok :=(r.userWant["fire"])[(startPosition+i)%4]; ok{//是下家胡?
							fmt.Println("玩家iD",id,"吃胡了")
							r.initRoom()
							break
						//	初始化 游戏结束
						}
					}
				}
				break
			}else if len(r.userWant["gang"]) == 1{
			//有一家杠 放杠处理
				user := accIDUsers[r.pengUserID]
				r.Playing = user.userData.AccID
				ugang := r.userWant["gang"]
				if v, ok := ugang[user.RoomPosition]; ok{
					as := strings.Split(v,"+")
					index, err := strconv.Atoi(as[1])
					if err !=nil{

					}
					user.gangOK(113,index,0)
					user.rspAllCards()
				}
				startPosition = user.RoomPosition
				continue

			}else if len(r.userWant["peng"]) == 1{
			//	碰优先于吃
			//	碰返回 去掉碰的 返回手牌
				um := r.userWant["peng"]
				r.Playing = r.pengUserID
				user := accIDUsers[r.pengUserID]
				if v, ok := um[user.RoomPosition]; ok{
					as := strings.Split(v,"+")
					//user := accIDUsers[as[0]]
					index, err := strconv.Atoi(as[1])
					if err !=nil{

					}
					user.pengOK(index)
					user.rspAllCards()
				}
				select {
				case r.OutCard = <- user.MyTurn:
				//	玩家有操作
				case <- time.After(5 * time.Second):
					//碰了 但是超时还没出牌
				//	将出手牌的最后一张 按一排顺序的
					length := len(user.Cardings)
					r.OutCard = user.Cardings[length-1]
					user.Cardings = user.Cardings[:length-1]
					//u.Cardings = append(u.Cardings[:index],u.Cardings[index+2:]...)

				}
				startPosition = user.RoomPosition
				//goto Again

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
		//startPosition = (startPosition + 1) % r.RoomData.RoomVolume
	}

}

