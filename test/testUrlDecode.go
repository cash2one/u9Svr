package main

import (
	fm "fmt"
	"net/url"
)

func main() {
	var urlStr string = "notify_time=2016-03-17 14:18:41&appid=d7492b8602d1a65f5165a5d84a7b7202&out_trade_no=20160317141827780&total_fee=0.01&subject=钻石&body=海马玩游戏demo&trade_status=107ea9899e1e47704e62c0fd82c211826, sign:0e4301e15ecca8ae00db934d053fa674" //"%E7%BE%BD%E6%AF%9B"
	//test := "abc=1_%E7%BE%BD%E6%AF%9B"
	//abc, _ := url.QueryUnescape(urlStr)
	abc := url.QueryEscape(urlStr)
	fm.Println(abc)
	// l, err := url.ParseQuery(urlStr)
	// fm.Println(l, err)
	// l2, err2 := url.ParseRequestURI(urlStr)
	// fm.Println(l2, err2)

	// l3, err3 := url.Parse(urlStr)
	// fm.Println(l3, err3)
	// fm.Println(l3.Path)
	// fm.Println(l3.RawQuery)
	// fm.Println(l3.Query())
	// fm.Println(l3.Query().Encode())

	// fm.Println(l3.RequestURI())
}
