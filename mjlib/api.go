package mjlib

var (
    MTableMgr *TableMgr
    MHuLib *HuLib
    InitCards = []int{
        0, 0, 0, 0, 0, 0, 0, 0, 0,
        0, 0, 0, 0, 0, 0, 0, 0, 0,
        0, 0, 0, 0, 0, 0, 0, 0, 0,
        0, 0, 0, 0, 0, 0, 0,
    }
)

func Init(){
    MTableMgr = &TableMgr{}
    MTableMgr.Init()
    MHuLib = &HuLib{}
}
