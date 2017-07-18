package msg

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&Hello{})
	Processor.Register(&WUser{})
// 各种信息注册
	Processor.Register(&CodeState{})
}

const  (
	ERROR_Register	=	111

	ERROR_Params	= 113 //数据格式有误


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

type CodeState struct {
	MSG_STATE int
	Message string
}