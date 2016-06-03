package channelApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
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
	case 102: //qihoo360
		fallthrough
	case 113: //拇指玩
		fallthrough
	case 122: //6YGame
		fallthrough
	case 123: //熊猫玩
		fallthrough
	case 127: //htc
		fallthrough
	case 134: //芒果玩
		fallthrough
	case 136: //小财神
		fallthrough
	case 142: //朋友玩
		fallthrough
	case 143: //全民游戏
		fallthrough
	case 145: //huawei
		fallthrough
	case 146: //lenovo
		beego.Trace(fmt.Sprintf("channelId:%d", mlr.ChannelId))
		ret.SetCode(0)
		return
	case 101: //当乐
		llr = loginRequest.LrNewDangle(mlr, jsonParam)
	case 103: //安智
		llr = loginRequest.LrNewAnZhi(mlr, jsonParam)
	case 104: //虫虫
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
	case 144:
		llr = loginRequest.LrNewVivo(mlr, jsonParam)
	case 147:
		llr = loginRequest.LrNewBaidu(mlr, jsonParam)
	default:
		ret.SetCode(3004)
		return
	}

	if err := llr.InitParam(); err != nil {
		beego.Error(err)
		ret = llr.SetCode(9001)
		return
	}

	if err := llr.GetResponse(); err != nil {
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

func CallCreateOrder(mlr *models.LoginRequest, orderId,
	extParamStr string, ctx *context.Context) (channelOrderId, ret string, err error) {
	var channelParams *map[string]interface{}
	if testChannelId != mlr.ChannelId {
		if channelParams, err = checkPackageParam(mlr); err != nil {
			beego.Error("checkPackageParam is error.")
			return
		}
	}
	var co createOrder.CreateOrder
	switch mlr.ChannelId {
	case 112: //魅族游戏
		co = new(createOrder.Meizu)
	case 120: //金立
		co = new(createOrder.Amigo)
	case 123: //熊猫玩
		co = new(createOrder.Xmw)
	case 136: //小财神
		co = new(createOrder.Caishen)
	case 139: //tencent
		co = new(createOrder.Tencent)
	case 144: //vivo
		co = new(createOrder.Vivo)
	case 145: //huawei
		co = new(createOrder.Huawei)
	default:
		return
	}

	beego.Trace("1:Prepare")
	if err = co.Prepare(mlr, orderId, extParamStr, channelParams, ctx); err != nil {
		return
	}

	beego.Trace("2:InitParam")
	if err = co.InitParam(); err != nil {
		return
	}

	beego.Trace("3:Response")
	if err = co.GetResponse(); err != nil {
		return
	}

	beego.Trace("4:ParseChannelRet")
	if err = co.ParseChannelRet(); err != nil {
		return
	}

	beego.Trace("5:GetResult")
	ret = co.GetResult()
	channelOrderId = co.GetChannelOrderId()
	return
}
