package internal

import (
	"fmt"
	"time"
)

type RoomData struct {
	RoomID int "_id"//房间类型 确实的
	RoomAccID	string
	RoomType int// 房间类型 即是什么类型的麻将
	RoomVolume int	//房间的容量
	RoomPay int	//房卡 需消耗
	RoomBaseMoney	int	//最低的 进房间 资金
	CreatedTime	int	//创建的时间
}

func (r *RoomData)initValue()error  {
	roomID, err := mongoDBNextSeq(C_ROOMS)
	if err != nil{
		return fmt.Errorf("get next rooms id error: %v", err)
	}
	r.RoomID = roomID
	r.RoomAccID = fmt.Sprintf("%06d", roomID)
	r.CreatedTime = int(time.Now().Unix())
	return nil
}