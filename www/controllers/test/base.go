package test

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Index() {

	var err error
	var num int64
	if num, err = this.statLoginRequest(0, 0, "", ""); err != nil {
		beego.Error(err)
	}

	this.Data["state"] = num
	this.Layout = "test/layout.html"
	this.TplName = "test/test.html"
}

func (this *BaseController) statLoginRequest(
	productId, channelId int,
	startTime, endTime string) (num int64, err error) {

	format := "CALL statLoginRequest(%d,%d,'%s','%s')"
	sqlText := fmt.Sprintf(format, productId, channelId, startTime, endTime)

	var rawPreparer orm.RawPreparer
	if rawPreparer, err = orm.NewOrm().Raw(sqlText).Prepare(); err != nil {
		return
	}
	defer func() {
		if err = rawPreparer.Close(); err != nil {
			return
		}
	}()

	var result sql.Result
	if result, err = rawPreparer.Exec(); err != nil {
		return
	}

	if num, err = result.RowsAffected(); err != nil {
		return
	}

	return
}
