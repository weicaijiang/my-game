package main

import (
	"testing"
	//"fmt"
	"my-game/mjlib"
)


func Benchmark_mj(b *testing.B)  {
	mjlib.Init()
	mjlib.MTableMgr.LoadTable()
	mjlib.MTableMgr.LoadFengTable()
	for i := 0; i < b.N; i++ {
		cards := []int{
			0, 0, 0, 0, 0, 0, 0, 0, 0,
			1, 1, 1, 2, 3, 0, 0, 0, 0,
			0, 0, 0, 2, 2, 2, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0,
		}
		mjlib.MHuLib.GetHuInfo(cards, 34, 34, 34)
	}

	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	//log.Fatal(err)
	//	fmt.Println("err",err)
	//	return
	//}
	//fmt.Println(strings.Replace(dir, "\\", "/", -1))
}
