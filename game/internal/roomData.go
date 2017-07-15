package internal

type RoomData struct {
	RoomType int// 房间类型 即是什么类型的麻将
	RoomVolume int	//房间的容量
	RoomPay int	//房卡 需消耗
	RoomBaseMoney	int	//最低的 进房间 资金
	CreatedTime	string	//创建的时间
}
