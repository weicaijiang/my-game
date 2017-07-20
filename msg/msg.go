package msg

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&Hello{})
	Processor.Register(&WUser{})
	Processor.Register(&LoginUser{})
//	微信登录
	Processor.Register(&WeChatLogin{})
// 各种信息注册
	Processor.Register(&CodeState{})

//	房间
	Processor.Register(&RoomBase{})
	Processor.Register(&JoinRoom{})
	Processor.Register(&RoomDataInfo{})

	Processor.Register(&QuitRoom{})
	Processor.Register(&ReadyGame{})

//	牌
	Processor.Register(&Cards{})
}

const  (

	SUCCESS_DONE = 1
	FAILURE_DONE = 0

	ERROR_Register	=	111

	ERROR_Params	= 113 //数据格式有误

	ERRO_LoginRepeated = 118 //重复登录

	ERRO_Unkonw	= 110
	ERRO_InitValue = 121 //初始值 initValue 报错

	ERRO_LoginFailed = 119 //登录失败

	ERRO_ONLYROOM	= 333//已有房间

	ERRO_NOTEXISITED = 444 //不存在操作

	SUCCESS_Register = 222




)

type Hello struct {
	Name	string
}

//用户x信息
type WUser struct {
	Name string	`json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

//用户登录信息
type LoginUser struct {
	Name string
	Password string
}

//微信登录
type WeChatLogin struct {
	Union	string //微信账号
}

//错误信息
type CodeState struct {
	MSG_STATE int
	Message string
}

//房间的基本信息
type RoomBase struct {
	Volume int
}
//返回创建了的房间信息
type RoomDataInfo struct {
	RoomID int "_id"//房间类型 确实的
	RoomAccID	string `json:"room_acc_id,omitempty"`
	RoomType int// 房间类型 即是什么类型的麻将
	RoomVolume int	//房间的容量
	RoomPay int	//房卡 需消耗
	RoomBaseMoney	int	//最低的 进房间 资金
	CreatedTime	int	//创建的时间
}

//加入房间
type JoinRoom struct {
	RoomAccID string
}

//牌的内容
type Cards struct {
	Cards []int//手牌
	PengCards []int//碰
	GangCards []int //杠
}

type QuitRoom  struct {
	Flag int
}

type ReadyGame struct {
	Flag int
}

