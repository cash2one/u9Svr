package common

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

const (
	CommonTimeLayout = "20060102150405"
)

type BasicRet struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
	Ext     string `json:"Ext"`
}

func (this *BasicRet) Init() *BasicRet {
	this.SetCode(9001)
	return this
}

func (this *BasicRet) SetCode(code int) {
	this.Code = code
	switch code {
	case 0:
		this.Message = "success"
	case 9001:
		this.Message = "Unknow execption"
	case 9002:
		this.Message = "Parse json execption"
	case 1001:
		this.Message = "LoginRequest:channelId exception"
	case 1002:
		this.Message = "LoginRequest:productId exception"
	case 1003:
		this.Message = "LoginRequest:token exception"
	case 1004:
		this.Message = "LoginRequest:userId exception"
	case 1005:
		this.Message = "PackageParam:channelId/productId exception"
	case 2001:
		this.Message = "OrderRequest:userId/productId exception"
	case 2002:
		this.Message = "ProductOrderId exception"
	case 2003:
		this.Message = "Amount exception"
	case 2004:
		this.Message = "OrderRequest:cp callbackUrl exception"
	case 2005:
		this.Message = "OrderRequest:order is already exist"
	case 3001:
		this.Message = "Channel api paramater exception"
	case 3002:
		this.Message = "Call channel api exception"
	case 3003:
		this.Message = "Channel api return failure"
	case 3004:
		this.Message = "ChannelId method is not implement"
	case 3005:
		this.Message = "ChannelId pay notify api is not implement"
	case 3006:
		this.Message = "Channel:channelId exception"
	case 3007:
		this.Message = "Product:productId exception"
	case 4001:
		this.Message = "pack:packId exception"
	default:
		this.Code = 9001
		this.Message = "Unknow execption"
	}
}

func DumpCtx(ctx *context.Context) (ret string) {
	format := "request:%+v"
	ret = fmt.Sprintf(format, ctx.Request)
	return
}
