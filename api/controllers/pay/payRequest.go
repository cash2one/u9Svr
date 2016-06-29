package pay

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
	"time"
	"u9/api/channel/api"
	"u9/api/common"
	"u9/models"
)

type PayRequestRet struct {
	common.BasicRet
	OrderId   string `json:"OrderId"`
	TransType string `json:"TransType"`
}

type PayRequestParam struct {
	ProductId      int    `json:"ProductId"`
	UserId         string `json:"UserId"`
	ProductOrderId string `json:"ProductOrderId"`
	Amount         int    `json:"Amount"`
	CallbackUrl    string `json:"CallbackUrl"`
	Ext            string `json:"Ext"`
	AppExt         string `json:"AppExt"`
}

func (this *PayRequestParam) Valid(v *validation.Validation) {
	switch {
	case this.ProductId <= 0:
		v.SetError("2001", "ProductId's range must not less and equal zero.")
		return
	case strings.TrimSpace(this.UserId) == "":
		v.SetError("2001", "Require userid.")
		return
	case strings.TrimSpace(this.ProductOrderId) == "":
		v.SetError("2002", "ProductOrderId is empty.")
		return
	case this.Amount <= 0:
		v.SetError("2003", "Amount's range must not less and equal zero.")
		return
	}

	isExistLoginRequest := new(models.LoginRequest).Query().
		Filter("UserId", this.UserId).Filter("ProductId", this.ProductId).Exist()

	if isExistLoginRequest == false {
		v.SetError("2001", "Query LoginRequest isn't exist.")
		return
	}

	var product models.Product
	if err := product.Query().Filter("id", this.ProductId).One(&product); err != nil {
		v.SetError("2001", "product isn't exist in database.")
		return
	}
	if this.CallbackUrl == "" {
		this.CallbackUrl = product.CallbackUrl
	}
}

func (this *PayController) PayRequest() {
	ret := PayRequestRet{OrderId: "", TransType: "CREATE_ORDER"}
	ret.Init()

	defer func() {
		this.Data["json"] = ret
		this.ServeJSON(true)
	}()

	msg := common.DumpCtx(this.Ctx)
	beego.Trace(msg)

	beego.Trace("payRequestParam: 1:validate")
	prp := new(PayRequestParam)
	if code := this.Validate(prp); code != 0 {
		ret.SetCode(code)
		return
	}

	beego.Trace("validateLogin: 2:query loginRequest with userId and productId")
	var lr models.LoginRequest
	lr.Query().
		Filter("UserId", prp.UserId).
		Filter("ProductId", prp.ProductId).
		One(&lr)

	beego.Trace("validateLogin: 3:create/update orderRequest")
	or := models.OrderRequest{
		UserId:         lr.Userid,
		ChannelId:      lr.ChannelId,
		ProductId:      lr.ProductId,
		ProductOrderId: prp.ProductOrderId,
		ReqAmount:      prp.Amount,
		AppExt:         prp.AppExt,
		Ext:            prp.Ext,
		ReqTime:        time.Now(),
		ProductCode:    -1,
		State:          0,
		CallbackUrl:    prp.CallbackUrl}
	create, _, err := orm.NewOrm().ReadOrCreate(&or,
		"UserId", "ChannelId", "ProductId", "ProductOrderId", "ReqAmount", "AppExt",
		"ReqTime", "State", "CallbackUrl")

	if create {
		or.OrderId = time.Now().Format(common.CommonTimeLayout) + strconv.Itoa(or.Id)
		if err = or.Update("OrderId"); err != nil {
			format := `validateLogin: err:%v`
			msg := fmt.Sprintf(format, err)
			beego.Error(msg)

			ret.SetCode(2005)
			return
		}
	} else {
		format := `validateLogin: record is exist, err:%v`
		msg := fmt.Sprintf(format, err)
		beego.Warn(msg)

		ret.SetCode(2005)
		return
	}

	beego.Trace("validateLogin: 4:callCreateOrder")
	channelExt := ""
	channelOrderId := ""
	channelOrderId, channelExt, err = channelApi.CallCreateOrder(&lr, or.OrderId, prp.Ext, this.Ctx)
	if err != nil {
		ret.SetCode(3002)
		if err = or.Delete(); err != nil {
			beego.Error(err)
		}
		return
	}

	beego.Trace("validateLogin: 5:update orderRequest with channelOrderId")
	or.ChannelOrderId = channelOrderId
	if err = or.Update("ChannelOrderId"); err != nil {
		ret.SetCode(3002)
		beego.Error(err)
		return
	}

	beego.Trace("validateLogin: 6:set ret")
	ret.SetCode(0)
	ret.Ext = channelExt
	ret.OrderId = or.OrderId
}

/*
  test url:
  http://192.168.0.185/api/gamePayRequest/?
  ProductId=1000&
  UserId=23c16f5323755132272fba79ab2e11d8&
  ProductOrderId=game20160114142841787&
  Amount=32&
  CallbackUrl=http://www.baidu.com
*/
