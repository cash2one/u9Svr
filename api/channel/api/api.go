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

func CallLoginRequest(mlr *models.LoginRequest) (ret *common.BasicRet) {
	ret = new(common.BasicRet).Init()
	var jsonParam *map[string]interface{}
	var err error
	if jsonParam, err = checkPackageParam(mlr); err != nil {
		code, _ := strconv.Atoi(err.Error())
		beego.Error(err)
		ret.SetCode(code)
		return
	}

	var llr loginRequest.LoginRequest
	switch mlr.ChannelId {
	case testChannelId: //test
		fallthrough
	case 113:
		fallthrough
	case 122: //6YGame
		fallthrough
	case 127:
		fallthrough
	case 134:
		fallthrough
	case 136:
		fallthrough
	case 142: //朋友玩
		fallthrough
	case 143: //全民游戏
		fallthrough
	case 123: //熊猫玩
		beego.Trace(fmt.Sprintf("channelId:%d", mlr.ChannelId))
		ret.SetCode(0)
		return
	case 101:
		llr = loginRequest.LrNewDangle(mlr, jsonParam)
	case 103:
		llr = loginRequest.LrNewAnZhi(mlr, jsonParam)
	case 104:
		llr = loginRequest.LrNewCCPay(mlr, jsonParam)
	case 106:
		llr = loginRequest.LrNewGFan(mlr, jsonParam)
	case 107:
		llr = loginRequest.LrNewGuopan(mlr, jsonParam)
	case 108:
		llr = loginRequest.LrNewKaoPu(mlr, jsonParam)
	case 109:
		llr = loginRequest.LrNewM4399(mlr, jsonParam)
	case 110:
		llr = loginRequest.LrNewOppo(mlr, jsonParam)
	case 111:
		llr = loginRequest.LrNewMuMaYi(mlr, jsonParam)
	case 112:
		llr = loginRequest.LrNewMeiZu(mlr, jsonParam)
	case 114:
		llr = loginRequest.LrNewJiuYou(mlr, jsonParam)
	case 115:
		llr = loginRequest.LrNewSogou(mlr, jsonParam)
	case 117:
		llr = loginRequest.LrNewWandoujia(mlr, jsonParam)
	case 118:
		llr = loginRequest.LrNewXiaomi(mlr, jsonParam)
	case 120:
		llr = loginRequest.LrNewAmigo(mlr, jsonParam)
	case 125:
		llr = loginRequest.LrNewHaiMaWan(mlr, jsonParam)
	case 126:
		llr = loginRequest.LrNewLeTV(mlr, jsonParam)
	case 128:
		llr = loginRequest.LrNewZhuoYi(mlr, jsonParam)
	case 129:
		llr = loginRequest.LrNewShouMeng(mlr, jsonParam)
	case 130:
		llr = loginRequest.LrNewYYH(mlr, jsonParam)
	case 131:
		llr = loginRequest.LrNewSnail(mlr, jsonParam)
	case 132:
		llr = loginRequest.LrNewYiJie(mlr, jsonParam)
	case 133:
		llr = loginRequest.LrNewYouLong(mlr, jsonParam)
	case 135:
		llr = loginRequest.LrNewQikQik(mlr, jsonParam)
	case 137:
		llr = loginRequest.LrNewPPTV(mlr, jsonParam)
	case 139:
		llr = loginRequest.LrNewTencent(mlr, jsonParam)
	case 140:
		llr = loginRequest.LrNewTT(mlr, jsonParam)
	case 141:
		llr = loginRequest.LrNewC07073(mlr, jsonParam)
	default:
		ret.SetCode(3004)
		return
	}

	llr.InitParam()

	if err := llr.Response(); err != nil {
		beego.Error(err)
		ret = llr.SetCode(3002)
		return
	}

	if err := llr.ParseChannelRet(); err != nil {
		beego.Error(err)
		ret = llr.SetCode(9002)
		return
	}

	if !llr.CheckChannelRet() {
		beego.Error(errors.New("channelRet code is fail."))
		ret = llr.SetCode(3003)
		return
	}
	ret = llr.SetCode(0)
	return
}

func checkPackageParam(mlr *models.LoginRequest) (jsonParam *map[string]interface{}, err error) {
	pp := new(models.PackageParam)
	if mlr.ChannelId != testChannelId {
		if err = pp.Query().Filter("channelId", mlr.ChannelId).Filter("productId", mlr.ProductId).One(pp); err != nil {
			msg := fmt.Sprintf("1005:channelId=%d and productId=%d", mlr.ChannelId, mlr.ProductId)
			beego.Error(msg)
			return nil, errors.New("1005")
		}

		jsonParam = new(map[string]interface{})
		if err = json.Unmarshal([]byte(pp.XmlParam), jsonParam); err != nil {
			beego.Error(err)
			return nil, errors.New("9002")
		}
	}
	//beego.Trace(pp.XmlParam)
	//beego.Trace(jsonParam)
	return jsonParam, nil
}

func CallCreateOrder(mlr *models.LoginRequest, orderId, host, ext string) (channelOrderId, ret string, err error) {
	var jsonParam *map[string]interface{}
	if testChannelId != mlr.ChannelId {
		if jsonParam, err = checkPackageParam(mlr); err != nil {
			beego.Error("checkPackageParam is error.")
			return
		}
	}

	var co createOrder.CreateOrder
	switch mlr.ChannelId {
	case 112:
		co = createOrder.CoNewMeizu(mlr, orderId, host, ext, jsonParam)
	case 120:
		co = createOrder.CoNewAmigo(mlr, orderId, host, ext, jsonParam)
	case 123: //熊猫玩
		co = createOrder.CoNewXmw(mlr, orderId, host, ext, jsonParam)
	case 136:
		co = createOrder.CoNewCaishen(mlr, orderId, host, ext, jsonParam)
	default:
		return
	}

	if err = co.InitParam(); err != nil {
		beego.Error(err)
		return
	}

	if err = co.Response(); err != nil {
		beego.Error(err)
		return
	}

	if err = co.ParseChannelRet(); err != nil {
		beego.Error(err)
		return
	}

	ret = co.GetResult()
	channelOrderId = co.GetChannelOrderId()
	return
}
