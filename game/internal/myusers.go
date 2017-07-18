package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/go"
	"my-game/msg"
	"gopkg.in/mgo.v2/bson"
)

var (
	users = make(map[bson.ObjectId]*MyLine)
)

type MyLine struct {
	gate.Agent
	*g.LinearContext

	MyUserData *UserData
}

func (u *MyLine)register()  {
	//user := new(UserData)
	user := u.MyUserData
	skeleton.Go(func() {
		//err := user.register()
		err := user.register()
		if err != nil{
			u.WriteMsg(&msg.CodeState{msg.ERROR_Register,"注册失败!"})
			u.Close()
			return
		}
	}, func() {

		//u.WriteMsg(&msg.CodeState{msg.SUCCESS_Register,"注册成功!"})
		//u.MyUserData = user
		//users[user.Id] = u

	})
	//err := u.MyUserData.register()
	//if err != nil{

}
