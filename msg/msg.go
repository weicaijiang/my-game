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
}

const  (
	ERROR_Register	=	111

	ERROR_Params	= 113 //数据格式有误

	ERRO_LoginRepeated = 118 //重复登录

	ERRO_Unkonw	= 110
	ERRO_InitValue = 121 //初始值 initValue 报错

	ERRO_LoginFailed = 119 //登录失败

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