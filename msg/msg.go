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
	Processor.Register(&Card{})

	Processor.Register(&Peng{})
	Processor.Register(&Gang{})

//	胡
	Processor.Register(&MimeHu{})

//	吃
	Processor.Register(&ChiPai{})
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

type Card struct {
	Index int //要出的牌的index
	Value int //牌的值
}

type Peng struct {
	Index int //要 碰的牌的下标
	Value int
}

type Gang struct {
	Index int //要 杠的牌的下标
	Value int // value为 说明是自个摸起的杠 111为暗杠 112为明杠 113为放杠
	GangType int //杠的类型 是否是自己摸起来杠的，是为1则为暗杠，2为明杠， 否则为3为别人放杠
}
//胡牌
type MimeHu struct {
	HuType int //胡牌类型 1为自摸，0为吃胡
	CardValue int //胡的牌值
}

//吃
type ChiPai struct {
	Index int
	Value int
	Array [][2]int
}

//确定碰
type PengOK struct {
	Index int
	Value int
}

//确定杠
//放杠
type GangOthers struct {
	Index int
	Value int
	GangType int
}