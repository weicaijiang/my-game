package toJson

import (
	"net"
	"github.com/name5566/leaf/network"
	"github.com/name5566/leaf/gate"
	"reflect"
	"github.com/name5566/leaf/log"
	"fmt"
)

type WriteToJSON struct {
	conn network.Conn
	gate *gate.Gate
	userData interface{}
}

func (a *WriteToJSON) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *WriteToJSON) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Call0("CloseAgent", a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *WriteToJSON) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		fmt.Println("data---w",data)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *WriteToJSON) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *WriteToJSON) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *WriteToJSON) Close() {
	a.conn.Close()
}

func (a *WriteToJSON) Destroy() {
	a.conn.Destroy()
}

func (a *WriteToJSON) UserData() interface{} {
	return a.userData
}

func (a *WriteToJSON) SetUserData(data interface{}) {
	a.userData = data
}