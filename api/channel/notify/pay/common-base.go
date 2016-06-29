package channelPayNotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
	"u9/api/common"
	"u9/models"
	"u9/tool"
)

const ChannelApiVerKeyName = "channelAPIVerion"

type Base struct {
	//channel http input/output
	ctx               *context.Context
	urlParamCheckKeys *[]string //选项参数:URL参数KEY检查列表
	urlParams         *url.Values
	body              string

	channelTradeContent string      //中间变量:渠道业务数据的字符串内容，可以选择性使用
	channelTradeData    interface{} //选项参数:渠道业务数据接口
	channelRetData      interface{} //选项参数:渠道返回数据接口

	//channel parameter
	channelParamKeys map[string]string //选项参数:渠道参数KEY清单
	channelParams    map[string]string //渠道参数

	//game product
	productKey         string //游戏产品密钥，用于数据验证
	productSignContent string
	ProductRet         common.BasicRet //游戏服务器返回数据接口

	//trade input parameter
	channelId int //初始化参数:渠道ID
	productId int //初始化参数:产品ID

	//trade store data
	orderId              string //U9订单ID,必须存在
	channelUserId        string //渠道用户ID，根据requireChannelUserId确定是否检查为空
	requireChannelUserId bool   //default=true,根据渠道是否存在channelUserId设定是否检查报错
	channelOrderId       string //部分渠道不存在该值

	payAmount     int     //支付金额（分）,必须存在
	payDiscount   int     //支付折扣(分)，默认为默认为0
	exChangeRatio float64 //选项参数:对换比例，默认为100.00，即假设渠道单位为元

	//DB models
	packageParam models.PackageParam
	loginRequest models.LoginRequest
	orderRequest models.OrderRequest
	payOrder     models.PayOrder

	//other
	lastError     int
	gameNotifyUrl string
}

func (this *Base) Init(params ...interface{}) (err error) {
	this.ProductRet.Init()
	this.lastError = err_noerror
	this.requireChannelUserId = true
	this.urlParamCheckKeys = &emptyUrlKeys
	this.exChangeRatio = 100.00

	this.channelId = params[0].(int)
	this.productId = params[1].(int)
	this.ctx = params[2].(*context.Context)

	if this.ctx.Request.Form == nil {
		this.ctx.Request.ParseForm()
	}
	this.urlParams = &(this.ctx.Request.Form)

	this.channelParams = map[string]string{}
	this.channelParams["_gameId"] = ""
	this.channelParams["_gameKey"] = ""
	this.channelParams["_payKey"] = ""
	this.channelParams["_version"] = ""

	this.channelParamKeys = map[string]string{}
	this.channelParamKeys["_gameId"] = ""
	this.channelParamKeys["_gameKey"] = ""
	this.channelParamKeys["_payKey"] = ""
	this.channelParamKeys["_version"] = ChannelApiVerKeyName

	if err = this.initProductKey(); err != nil {
		return
	}

	if err = this.initChannelParam(); err != nil {
		return
	}
	return
}

func (this *Base) checkUrlParam() (err error) {
	for _, urlKey := range *(this.urlParamCheckKeys) {
		if this.urlParams.Get(urlKey) == "" {
			this.lastError = err_checkUrlParam
			return
		}
	}
	return
}

func (this *Base) getChannelParam(channelParamKey string) (ret string, err error) {
	ok := false
	if ret, ok = this.channelParams[channelParamKey]; !ok {
		format := `getChannelParam: channelParam %s isn't exist`
		msg := fmt.Sprintf(format, channelParamKey)
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	return
}

func (this *Base) initChannelParam() (err error) {
	defer func() {
		if err != nil {
			this.lastError = err_initChannelParam
		}
	}()

	if err = this.packageParam.Query().
		Filter("channelId", this.channelId).
		Filter("productId", this.productId).
		One(&this.packageParam); err != nil {
		beego.Error(err)
		return
	}

	if err = json.Unmarshal([]byte(this.packageParam.XmlParam), &this.channelParams); err != nil {
		beego.Error(err)
		return
	}

	return
}

func (this *Base) parseChannelParam() (err error) {
	key := "_gameId"
	if this.channelParamKeys[key] != "" {
		if this.channelParams[key], err = this.getChannelParam(this.channelParamKeys[key]); err != nil {
			this.lastError = err_initChannelGameId
			return
		}
	}

	key = "_gameKey"
	if this.channelParamKeys[key] != "" {
		if this.channelParams[key], err = this.getChannelParam(this.channelParamKeys[key]); err != nil {
			this.lastError = err_initChannelGameKey
			return
		}
	}

	key = "_payKey"
	if this.channelParamKeys[key] != "" {
		if this.channelParams[key], err = this.getChannelParam(this.channelParamKeys[key]); err != nil {
			this.lastError = err_initChannelPayKey
			return
		}
	}

	key = "_version"
	this.channelParams[key], _ = this.getChannelParam(this.channelParamKeys[key])
	return
}

func (this *Base) ParseInputParam(params ...interface{}) (err error) {
	if err = this.checkUrlParam(); err != nil {
		return
	}

	if err = this.parseChannelParam(); err != nil {
		return
	}

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, this.ctx.Request.Body); err != nil {
		err = nil
	} else {
		this.body = string(buffer.Bytes())
	}
	return
}

func (this *Base) initProductKey() (err error) {
	defer func() {
		if err != nil {
			this.lastError = err_initProductKey
		}
	}()

	product := new(models.Product)
	if err = product.Query().
		Filter("Id", this.productId).
		One(product); err != nil {
		format := "initProductKey: err:%v, product:%v"
		msg := fmt.Sprintf(format, err, product)
		beego.Error(msg)
		return
	}

	this.productKey = product.AppKey
	return
}

func (this *Base) prepareOrderRequest() (err error) {
	defer func() {
		if err != nil {
			this.lastError = err_prepareOrderRequest
		}
	}()

	orqs := this.orderRequest.Query().
		Filter("ChannelId", this.channelId).
		Filter("ProductId", this.productId).
		Filter("OrderId", this.orderId)
	if orqs.Exist() == false {
		format := "prepareOrderRequest: err:orderRequest is not exist, channelId:%d, productId:%d, orderId:%s"
		msg := fmt.Sprintf(format, this.channelId, this.productId, this.orderId)
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	if err = orqs.One(&this.orderRequest); err != nil {
		format := "prepareOrderRequest: err:%v, channelId:%d, productId:%d, orderId:%s"
		msg := fmt.Sprintf(format, err, this.channelId, this.productId, this.orderId)
		err = errors.New(msg)
		beego.Error(msg)
		return
	}

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(this.urlParams); err != nil {
		format := "prepareOrderRequest: err:%v"
		msg := fmt.Sprintf(format, err)
		err = errors.New(msg)
		beego.Error(err)
		return
	}

	this.orderRequest.ChannelLog = string(jsonBytes)
	if err = this.orderRequest.Update("ChannelLog"); err != nil {
		format := "prepareOrderRequest: err:%v"
		msg := fmt.Sprintf(format, err)
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	return
}

func (this *Base) prepareLoginRequest() (err error) {
	err = this.loginRequest.Query().
		Filter("Userid", this.orderRequest.UserId).
		One(&this.loginRequest)
	if err != nil {
		this.lastError = err_prepareLoginRequest

		format := "prepareLoginRequest: err:%v, UserId:%s, loginRequest:%+v"
		msg := fmt.Sprintf(format, err, this.orderRequest.UserId, this.loginRequest)
		err = errors.New(msg)
		beego.Error(msg)
		return
	}
	return
}

func (this *Base) PrepareTradeData() (err error) {
	if err = this.prepareOrderRequest(); err != nil {
		return
	}
	if err = this.prepareLoginRequest(); err != nil {
		return
	}
	return
}

func (this *Base) CheckSign(params ...interface{}) (err error) {
	signState := params[0].(bool)
	signMethod := params[1].(string)
	signMsg := params[2].(string)

	if !signState {
		this.lastError = err_checkSign

		format := "CheckSign: %s sign is invalid, err:%+v, %s"
		msg := fmt.Sprintf(format, signMethod, err, signMsg)
		err = errors.New(msg)
		beego.Error(err)
		return
	}
	return
}

func (this *Base) CheckChannelRet(params ...interface{}) (err error) {
	if this.orderId != this.orderRequest.OrderId {
		this.lastError = err_orderIsNotExist
		format := "CheckChannelRet: err:orderId!=orderRequest.OrderId, channelId:%s, productId:%s orderRequest:%+v, orderId:%s"
		msg := fmt.Sprintf(format, this.channelId, this.productId, this.orderRequest, this.orderId)
		err = errors.New(msg)
		beego.Error(err)
		return
	}

	if this.orderRequest.ReqAmount != this.payAmount {
		this.lastError = err_payAmountError
		format := "CheckChannelRet: err:orderRequest.ReqAmount!=payAmount, orderRequest:%+v, payAmount:%d"
		msg := fmt.Sprintf(format, this.orderRequest, this.payAmount)
		err = errors.New(msg)
		beego.Error(msg)
		return
	}

	if this.requireChannelUserId && this.channelUserId != this.loginRequest.ChannelUserid {
		this.lastError = err_channelUserIsNotExist
		format := "CheckChannelRet: channelUserId!=loginRequest.ChannelUserid, loginRequest:%+v, channelUserId:%s"
		msg := fmt.Sprintf(format, this.loginRequest, this.channelUserId)
		err = errors.New(msg)
		beego.Error(err)
		return
	}

	if len(params) > 1 {
		tradeState := params[0].(bool)

		if !tradeState {
			this.lastError = err_tradeFail

			format := "CheckChannelRet: err:%s, channelParamKeys:%+v, channelParam:%+v, urlParam:%+v, body:%s, channelTradeData:%+v"
			tradeFailDesc := params[1].(string)
			errMsg := fmt.Sprintf(format,
				tradeFailDesc,
				this.channelParamKeys, this.channelParams,
				this.urlParams, this.body, this.channelTradeData)
			beego.Warn(errMsg)
		}
	}
	return
}

func (this *Base) notifyProductSvr() (err error) {
	defer func() {
		if err != nil {
			this.lastError = err_notifyProductSvr
		}
	}()

	format := this.orderRequest.CallbackUrl + "?" +
		"UserId=%s" + "&" + "OrderId=%s" + "&" + "ChannelId=%d" + "&" +
		"ChannelOrderId=%s" + "&" + "ChannelUserId=%s" + "&" + "ProductId=%d" + "&" +
		"ProductOrderId=%s" + "&" + "ReqAmount=%d" + "&" + "AppExt=%s" + "&" +
		"PayAmount=%d" + "&" + "Sign=%s" + "&" + "Code=%d" + "&" + "ErrDesc=%s"

	this.productSignContent = fmt.Sprintf("%s%s%s%s",
		this.orderRequest.ProductOrderId,
		this.orderId,
		this.channelOrderId,
		this.productKey)

	this.gameNotifyUrl = fmt.Sprintf(format,
		this.orderRequest.UserId,
		this.orderId,
		this.channelId,
		this.channelOrderId,
		this.loginRequest.ChannelUserid,
		this.productId,
		this.orderRequest.ProductOrderId,
		this.orderRequest.ReqAmount,
		this.orderRequest.AppExt,
		this.payAmount,
		tool.Md5([]byte(this.productSignContent)),
		this.lastError,
		errorDescList[this.lastError])

	if this.orderRequest.CallbackUrl != "" {
		req := httplib.Get(this.gameNotifyUrl)

		if _, err = req.Response(); err != nil {
			format := "notifyProductSvr: err:%+v, gameNotifyUrl:%s"
			msg := fmt.Sprintf(format, err, this.gameNotifyUrl)
			beego.Error(msg)
			return
		}

		if err = req.ToJSON(&this.ProductRet); err != nil {
			format := "notifyProductSvr: err:%+v, gameNotifyUrl:%s ,gameServerRet:%s"
			gameServerRet, _ := req.String()
			msg := fmt.Sprintf(format, err, this.gameNotifyUrl, gameServerRet)
			beego.Error(msg)
			return
		}
	} else {
		this.ProductRet.SetCode(0)
	}
	beego.Trace(this.gameNotifyUrl)

	this.orderRequest.ProductCode = this.ProductRet.Code
	this.orderRequest.ProductMessage = this.ProductRet.Message
	this.orderRequest.State = 1
	if err = this.orderRequest.Update("State", "ProductMessage", "ProductCode"); err != nil {
		format := "notifyProductSvr: err:%+v, orderRequest:%+v"
		msg := fmt.Sprintf(format, err, this.orderRequest)
		err = errors.New(msg)
		beego.Error(msg)
	}

	if this.ProductRet.Code != 0 {
		format := "notifyProductSvr: err:this.ProductRet.Code!=0, gameNotifyUrl:%s, ProductRet:%+v"
		msg := fmt.Sprintf(format, this.gameNotifyUrl, this.ProductRet)
		err = errors.New(msg)
		beego.Error(msg)
		return
	}
	return
}

func (this *Base) handleOrder() (err error) {
	defer func() {
		if err != nil {
			this.lastError = err_handleOrder
		}
	}()

	created := false
	this.payOrder = models.PayOrder{OrderId: this.orderId, ChannelOrderId: this.channelOrderId}
	if created, _, err = orm.NewOrm().ReadOrCreate(&this.payOrder, "OrderId", "ChannelOrderId"); err != nil {
		format := "notifyProductSvr: err:%+v, payOrder:%+v"
		msg := fmt.Sprintf(format, err, this.payOrder)
		err = errors.New(msg)
		beego.Error(err)
		return
	}

	if !created {
		beego.Warn("notifyProductSvr: PayOrder:order is exist.")
		return
	}

	this.payOrder.PayAmount = this.payAmount
	this.payOrder.PayDiscount = this.payDiscount
	this.payOrder.PayTime = time.Now()
	if err = this.payOrder.Update("PayAmount", "PayDiscount"); err != nil {
		format := "notifyProductSvr: err:%+v, payOrder:%+v"
		msg := fmt.Sprintf(format, err, this.payOrder)
		err = errors.New(msg)
		beego.Error(msg)
		return
	}

	this.orderRequest.State = 2
	err = this.orderRequest.Update("State")
	if err != nil {
		format := "notifyProductSvr: err:%+v, orderRequest:%+v"
		msg := fmt.Sprintf(format, err, this.orderRequest)
		err = errors.New(msg)
		beego.Error(msg)
		return
	}
	return
}

func (this *Base) Handle() (err error) {
	switch this.orderRequest.State {
	case 0, 1:
		if err = this.notifyProductSvr(); err != nil {
			return err
		}
		if err = this.handleOrder(); err != nil {
			return err
		}
	case 2:
		return nil
	default:
		err = errors.New("handle: It isn't exist orderRequest state")
		beego.Error(err)
		return err
	}
	return
}

func (this *Base) GetResult(params ...interface{}) (ret string) {
	msg := this.Dump()
	if this.lastError != err_noerror {
		beego.Error(msg)
	} else {
		beego.Trace(msg)
	}

	paramLen := len(params)
	if paramLen < 1 {
		return
	}
	retFormat := params[0].(string)

	var succMsgParams []interface{}
	var failMsgParams []interface{}
	if paramLen > 2 {
		spiltParams := strings.Split(params[1].(string), ",")
		for _, spiltParam := range spiltParams {
			succMsgParams = append(succMsgParams, interface{}(spiltParam))
		}

		spiltParams = strings.Split(params[2].(string), ",")
		for _, spiltParam := range spiltParams {
			failMsgParams = append(failMsgParams, interface{}(spiltParam))
		}
	}

	if this.lastError == err_noerror {
		ret = fmt.Sprintf(retFormat, succMsgParams...)
	} else {
		ret = fmt.Sprintf(retFormat, failMsgParams...)
	}

	return
}

func (this *Base) parsePayAmount(amount, discount string) (err error) {
	defer func() {
		if err != nil {
			format := "err:%v, "
			msg := fmt.Sprintf(format, err)
			msg = msg + this.Dump()
			beego.Error(msg)
		}
	}()

	payAmount := 0.0
	if payAmount, err = strconv.ParseFloat(amount, 64); err != nil {
		this.lastError = err_parsePayAmount
		return
	}
	payDiscount := 0.0
	if discount != "" {
		if payDiscount, err = strconv.ParseFloat(discount, 64); err != nil {
			this.lastError = err_parsePayDiscount
			return
		}
	}

	this.payDiscount = int(payDiscount * this.exChangeRatio)
	this.payAmount = int(payAmount*this.exChangeRatio) + this.payDiscount
	return
}

func (this *Base) Dump() (ret string) {
	urlParamCheckKey := `[`
	for _, urlKey := range *(this.urlParamCheckKeys) {
		urlParamCheckKey = urlParamCheckKey + `,` + urlKey
	}
	urlParamCheckKey = urlParamCheckKey + `]`

	format := `productKey:%s productSignContent:%s` +
		"\nchannelParamKey: %+v,\n" + "channelParam: %+v,\n\n" +
		"urlParamCheckKey: %s,\n" + "request: %+v,\n" + "body: %s,\n\n" +
		"channelTradeContent: %s,\n\n" + "channelTradeData: %+v, \n\n" +
		"orderId: %s,  channelUserId: %s,  requireChannelUserId: %v,  channelOrderId: %s" +
		" ,payAmount: %d,  payDiscount: %d,  exChangeRatio: %f,\n\n" +
		"loginRequest: %+v,\n\n" +
		"orderRequest: %+v,\n\n" +
		"payOrder: %+v,\n\n" +
		"gameNotifyUrl: %s,\n" +
		"channelRetData: %+v,\n" +
		"lastError: %d,  lastErrorDesc: %s"

	ret = fmt.Sprintf(format,
		this.productKey, this.productSignContent,
		this.channelParamKeys, this.channelParams,
		this.urlParamCheckKeys,
		this.ctx.Request, this.body,
		this.channelTradeContent, this.channelTradeData,
		this.orderId, this.channelUserId, this.requireChannelUserId, this.channelOrderId,
		this.payAmount, this.payDiscount, this.exChangeRatio,
		this.loginRequest, this.orderRequest, this.payOrder,
		this.gameNotifyUrl,
		this.channelRetData,
		this.lastError, errorDescList[this.lastError])

	return
}
