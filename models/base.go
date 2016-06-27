package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	if dbport == "" {
		dbport = "3306"
	}

	dburl := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname +
		"?charset=utf8&loc=Asia%2FShanghai"
	orm.RegisterDataBase("default", "mysql", dburl)

	orm.RegisterModel(new(Channel), new(Product), new(ProductVersion),
		new(LoginRequest), new(OrderRequest), new(PayOrder),
		new(Manager), new(Cp),
		new(PackageParam), new(ChannelPackageParam),
		new(PackageTask), new(PackageTaskList), new(StatChannelPay),
		new(LoginRequestLog))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}
