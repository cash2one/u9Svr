package main

import (
	"fmt"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"u9/models"
)

func main() {

	created := false
	var err error
	payOrder := models.PayOrder{OrderId: "20160316150010665", ChannelOrderId: "16031615001270000011"}
	if created, _, err = orm.NewOrm().ReadOrCreate(&payOrder, "OrderId", "ChannelOrderId"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(created)
}
