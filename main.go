package main

import (
	"jiratool/conf"
	"jiratool/fynetool"
)

func main() {
	// 初始化設定檔
	conf.ConfigInit()
	fynetool.InitDataList()

	// 初始化視窗元件
	w := fynetool.InitFyneApp()

	// 設定視窗內容元件
	w = fynetool.SetttingWidget(w)

	// 啟動視窗
	w.ShowAndRun()
}
