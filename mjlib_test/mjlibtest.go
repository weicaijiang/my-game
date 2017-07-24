package main

import (
	"fmt"
	"qipai/mjlib_go/src/mjlib"
)

func print_cards(cards []int) {
	for i := 0; i < 9; i++ {
		fmt.Printf("%d,", cards[i])
	}
	fmt.Printf("\n")

	for i := 9; i < 18; i++ {
		fmt.Printf("%d,", cards[i])
	}
	fmt.Printf("\n")

	for i := 18; i < 27; i++ {
		fmt.Printf("%d,", cards[i])
	}
	fmt.Printf("\n")

	for i := 27; i < 34; i++ {
		fmt.Printf("%d,", cards[i])
	}
	fmt.Printf("\n")
}

var tested = map[int]bool{}

func check_hu(cards []int, max int) {
	for i := 0; i < max; i++ {
		if cards[i] > 4 {
			return
		}
	}

	num := 0
	for i := 0; i < 9; i++ {
		num = num*10 + cards[i]
	}

	_, ok := tested[num]
	if ok {
		return
	}

	tested[num] = true

	for i := 0; i < max; i++ {
		if !mjlib.MHuLib.GetHuInfo(cards, 34, 34, 34) {
			fmt.Printf("测试失败 i=%d\n", i)
			print_cards(cards)
		}
	}
}

func gen_auto_table_sub(cards []int, level int) {
	for i := 0; i < 32; i++ {
		index := -1
		if i <= 17 {
			cards[i] += 3
		} else if i <= 24 {
			index = i - 18
		} else {
			index = i - 16
		}

		if index >= 0 {
			cards[index] += 1
			cards[index+1] += 1
			cards[index+2] += 1
		}

		if level == 4 {
			check_hu(cards, 18)
		} else {
			gen_auto_table_sub(cards, level+1)
		}

		if i <= 17 {
			cards[i] -= 3
		} else {
			cards[index] -= 1
			cards[index+1] -= 1
			cards[index+2] -= 1
		}
	}
}

func test_two_color() {
	fmt.Println("测试两种花色")
	cards := []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}

	for i := 0; i < 18; i++ {
		cards[i] = 2
		fmt.Printf("将 %d\n", i+1)
		gen_auto_table_sub(cards, 1)
		cards[i] = 0
	}
}

func test_one_success() {
	cards := []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 2, 3, 0, 0, 0, 0,
		0, 0, 0, 2, 2, 2, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}

	//fmt.Println("测试1种能胡的牌型")
	//print_cards(cards)
	mjlib.MHuLib.GetHuInfo(cards, 34, 34, 34)
	//if mjlib.MHuLib.GetHuInfo(cards, 34, 34, 34) {
	//	//fmt.Println("测试通过：胡牌")
	//} else {
	//	//fmt.Println("测试失败：能胡的牌型判断为不能胡牌")
	//}
}

func test_one_fail() {
	cards := []int{
		0, 1, 1, 1, 0, 0, 1, 0, 1,
		0, 1, 1, 1, 0, 0, 2, 2, 2,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}

	fmt.Println("测试1种不能胡的牌型")
	print_cards(cards)
	if !mjlib.MHuLib.GetHuInfo(cards, 34, 34, 34) {
		fmt.Println("测试通过：不能胡牌")
	} else {
		fmt.Println("测试失败：不能胡牌的牌型判断为胡了")
	}
}

func main() {
	fmt.Println("test hulib begin...")

	mjlib.Init()
	//fmt.Println("init")
	mjlib.MTableMgr.LoadTable()
	mjlib.MTableMgr.LoadFengTable()
	//fmt.Println("success")

	//cards := []int{
	//	0, 1, 1, 1, 0, 0, 1, 0, 1,
	//	0, 1, 1, 1, 0, 0, 2, 2, 2,
	//	0, 0, 0, 0, 0, 0, 0, 0, 0,
	//	0, 0, 0, 0, 0, 0, 0,
	//}
	//cards := []int{
	//	0, 0, 0, 0, 0, 2, 0, 0, 0,
	//	0, 0, 0, 0, 0, 0, 0, 0, 0,
	//	1, 1, 1, 1, 1, 1, 3, 2, 1,
	//	0, 0, 0, 0, 0, 0, 0,
	//}
	cards := []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 2, 3, 0, 0, 0, 0,
		0, 0, 0, 2, 2, 2, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0,
	}
	//time1 := time.Now().UnixNano()
	//for i:=0; i<1000000; i++{
		//test_one_success()
		if mjlib.MHuLib.GetHuInfo(cards, 34, 34, 34){
			fmt.Println("hu")
		}else {
			fmt.Println("no hu")
		}
	//}
	//time2 := time.Now().UnixNano()
	//fmt.Println("time=",time2-time1)
	//test_one_fail()

	//    test_two_color()
}
