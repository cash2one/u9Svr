package createOrder

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"u9/models"
	"u9/tool"
)

type caishenExtParam struct {
	ProductId string `json:"product_id"`
	Price     int    `json:"price"`
	GameUid   string `json:"game_uid"`
	UId       string `json:"u_id"`
	GameId    string `json:"game_id"`
}

type Caishen struct {
	Cr
}

func (this *Caishen) Prepare(lr *models.LoginRequest, orderId, extParamStr string,
	channelParams *map[string]interface{}, ctx *context.Context) (err error) {

	if err = this.Cr.Initial(lr, orderId, nil, new(caishenExtParam),
		extParamStr, channelParams, ctx); err != nil {
		beego.Error(err)
		return err
	}
	this.parseAppKey("PAYKEY")
	return nil
}

func (this *Caishen) InitParam() (err error) {
	return nil
}

func (this *Caishen) GetResponse() (err error) {
	return nil
}

func (this *Caishen) ParseChannelRet() (err error) {
	return nil
}

func (this *Caishen) GetResult() (ret string) {

	format := `%s_%s_%d_%s_%s_%s_%s`

	extParam := this.extParam.(*caishenExtParam)
	context := fmt.Sprintf(format, this.orderId, extParam.ProductId, extParam.Price,
		extParam.GameUid, extParam.UId, extParam.GameId, this.appKey)

	sign := tool.Md5([]byte(context))

	ret = sign
	return ret
}
