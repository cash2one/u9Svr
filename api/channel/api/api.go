package channelApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"u9/api/channel/api/createOrder"
	"u9/api/channel/api/loginRequest"
	"u9/api/common"
	"u9/models"
)

const (
	testChannelId = 100
)

func CallLoginRequest(channelId, productId int, channelUserId, token string) (ret *common.BasicRet) {
	ret = new(common.BasicRet).Init()
	beego.Trace("callLoginRequest")
	var jsonParam *map[string]interface{}
	var err error
	if jsonParam, err = checkPackageParam(channelId, productId); err != nil {
		code, _ := strconv.Atoi(err.Error())
		ret.SetCode(code)
	}

	var lr loginRequest.LoginRequest
	switch channelId {
	case testChannelId: //test
		fallthrough
	case 122: //6YGame
		fallthrough
	case 127:
		fallthrough
	case 134:
		fallthrough
	case 136:
		fallthrough
	case 123: //熊猫玩
		beego.Trace(channelId)
		ret.SetCode(0)
		return
	case 101:
		lr = loginRequest.LrNewDangle(channelUserId, token, jsonParam)
	case 103:
		lr = loginRequest.LrNewAnZhi(channelUserId, token, jsonParam)
	case 104:
		lr = loginRequest.LrNewCCPay(channelUserId, token, jsonParam)
	case 106:
		lr = loginRequest.LrNewGFan(channelUserId, token, jsonParam)
	case 107:
		lr = loginRequest.LrNewGuopan(channelUserId, token, jsonParam)
	case 109:
		lr = loginRequest.LrNewM4399(channelUserId, token, jsonParam)
	case 110:
		lr = loginRequest.LrNewOppo(channelUserId, token, jsonParam)
	case 112:
		lr = loginRequest.LrNewMeiZu(channelUserId, token, jsonParam)
	case 117:
		lr = loginRequest.LrNewWandoujia(channelUserId, token, jsonParam)
	case 118:
		lr = loginRequest.LrNewXiaomi(channelUserId, token, jsonParam)
	case 120:
		lr = loginRequest.LrNewAmigo(channelUserId, token, jsonParam)
	case 125:
		lr = loginRequest.LrNewHaiMaWan(channelUserId, token, jsonParam)
	case 126:
		lr = loginRequest.LrNewLeTV(channelUserId, token, jsonParam)
	case 128:
		lr = loginRequest.LrNewZhuoYi(channelUserId, token, jsonParam)
	case 129:
		lr = loginRequest.LrNewShouMeng(channelUserId, token, jsonParam)
	case 130:
		lr = loginRequest.LrNewYYH(channelUserId, token, jsonParam)
	// case 131:

	case 132:
		lr = loginRequest.LrNewYiJie(channelUserId, token, jsonParam)
	case 133:
		lr = loginRequest.LrNewYouLong(channelUserId, token, jsonParam)
	default:
		ret.SetCode(3004)
		return
	}

	lr.InitParam()

	if err := lr.Response(); err != nil {
		beego.Trace(err)
		ret = lr.SetCode(3002)
		return
	}

	if err := lr.ParseChannelRet(); err != nil {
		beego.Trace(err)
		ret = lr.SetCode(9002)
		return
	}

	if !lr.CheckChannelRet() {
		beego.Trace(errors.New("channelRet code is fail."))
		ret = lr.SetCode(3003)
		return
	}
	ret = lr.SetCode(0)
	return
}

func checkPackageParam(channelId, productId int) (jsonParam *map[string]interface{}, err error) {
	pp := new(models.PackageParam)
	if err = pp.Query().Filter("channelId", channelId).Filter("productId", productId).One(pp); err != nil {
		msg := fmt.Sprintf("channelId=%d and productId=%d", channelId, productId)
		err = errors.New("1005")
		beego.Trace(err, ":", msg)
		return nil, err
	}

	jsonParam = new(map[string]interface{})
	beego.Trace(pp.XmlParam)
	if err = json.Unmarshal([]byte(pp.XmlParam), jsonParam); err != nil {
		beego.Trace(err, ":", err.Error())
		return nil, errors.New("9002")
	}
	return jsonParam, nil
}

func CallCreateOrder(lr *models.LoginRequest, orderId, host, ext string) (channelOrderId, ret string, err error) {
	var jsonParam *map[string]interface{}
	if testChannelId != lr.ChannelId {
		if jsonParam, err = checkPackageParam(lr.ChannelId, lr.ProductId); err != nil {
			beego.Error(err)
			return
		}
	}

	var co createOrder.CreateOrder
	switch lr.ChannelId {
	case 112:
		co = createOrder.CoNewMeizu(lr, orderId, host, ext, jsonParam)
	case 120:
		co = createOrder.CoNewAmigo(lr, orderId, host, ext, jsonParam)
	case 123: //熊猫玩
		co = createOrder.CoNewXmw(lr, orderId, host, ext, jsonParam)
	case 136:
		co = createOrder.CoNewCaishen(lr, orderId, host, ext, jsonParam)
	default:
		return
	}

	if err = co.InitParam(); err != nil {
		beego.Trace(err)
		return
	}

	if err = co.Response(); err != nil {
		beego.Trace(err)
		return
	}

	if err = co.ParseChannelRet(); err != nil {
		beego.Trace(err)
		return
	}

	ret = co.GetResult()
	channelOrderId = co.GetChannelOrderId()
	return
}
