package mjlib

import "fmt"

//自己写的 调用mjlib的胡牌

//传入 手牌数组 鬼牌位置
//返回 是否胡
func IsHu(a []int,guiIndex,guiValue int)bool  {
	 as := []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}
	fmt.Println("初始as:",as)
	if len(a) == 2{//手牌只有2张
		if a[0] == guiValue || a[1] == guiValue{
			return true
		}else if a[0] == a[1]{
			return true
		}else {
			return false
		}
	}
	for i := 0; i< len(a); i++{
		v := a[i]
		if a[i] > 100 && a[i] <110{
			as[(v%100 - 1)] += 1
		}else if a[i] >200 && a[i] < 210{
			as[(v%200 + 8)] += 1
		}else if a[i] >300 && a[i] < 310{
			as[(v%300 + 17)] += 1
		}else if a[i] >400 && a[i] < 410{
			as[(v%400 + 26)] += 1
		}
	}
	fmt.Println("整理后的as",as)
	if MHuLib.GetHuInfo(as,34,guiIndex,34){
		return true
	}
	return false
}
