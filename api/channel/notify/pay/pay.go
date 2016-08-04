package channelPayNotify

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func CallPayNotify(channelId, productId int, ctx *context.Context) (ret string, err error) {

	var pn PayNotify
	ret = "failure"

	defer func() {
		format := "callPayNotify: 7:GetResult:%s"
		if pn != nil {
			ret = pn.GetResult()
			msg := fmt.Sprintf(format, ret)
			beego.Trace(msg)
		} else {
			msg := fmt.Sprintf(format, ret)
			beego.Error(msg)
		}
	}()

	switch channelId {
	case 100:
		pn = new(Test)
	case 101:
		pn = new(Dangle)
	case 102:
		pn = new(Qihoo360)
	case 103:
		pn = new(AnZhi)
	case 104:
		pn = new(CCPay)
	case 106:
		pn = new(GFan)
	case 107:
		pn = new(Guopan)
	case 108:
		pn = new(KaoPu)
	case 109:
		pn = new(M4399)
	case 110:
		pn = new(Oppo)
	case 111:
		pn = new(MuMaYi)
	case 112:
		pn = new(Meizu)
	case 113:
		pn = new(MZW)
	case 114:
		pn = new(Jiuyou)
	case 115:
		pn = new(Sogou)
	case 117:
		pn = new(Wandoujia)
	case 118:
		pn = new(Xiaomi)
	case 120:
		pn = new(Amigo)
	case 122:
		pn = new(CYGame)
	case 123:
		pn = new(Xmw)
	case 125:
		pn = new(Haimawan)
	case 126:
		pn = new(LeTV)
	case 127:
		pn = new(HTC)
	case 128:
		pn = new(ZhuoYi)
	case 129:
		pn = new(ShouMeng)
	case 130:
		pn = new(YYH)
	case 131:
		pn = new(Snail)
	case 132:
		pn = new(YiJie)
	case 133:
		pn = new(YouLong)
	case 134:
		pn = new(Mango)
	case 135:
		pn = new(QikQik)
	case 136:
		pn = new(Xcs)
	case 137:
		pn = new(PPTV)
	case 139:
		pn = new(Tencent)
	case 140:
		pn = new(TT)
	case 141:
		pn = new(C07073)
	case 142:
		pn = new(Pengyouwan)
	case 143:
		pn = new(Qmyx)
	case 144:
		pn = new(Vivo)
	case 145:
		pn = new(Huawei)
	case 146:
		pn = new(Lenovo)
	case 147:
		pn = new(Baidu)
	case 148:
		pn = new(AnFeng)
	case 149:
		pn = new(CoolPad)
	case 150:
	 	pn = new(PaoJiao)
	case 151:
		pn = new(WeiUU)
	default:
		err = errors.New("channelId isn't exist.")
		return
	}

	beego.Trace("callPayNotify: 1:Init")
	if err = pn.Init(channelId, productId, ctx); err != nil {
		return
	}

	beego.Trace("callPayNotify: 2:ParseInputParam")
	if err = pn.ParseInputParam(); err != nil {
		return
	}
	beego.Trace("callPayNotify: 3:PrepareTradeData")
	if err = pn.PrepareTradeData(); err != nil {
		return
	}
	beego.Trace("callPayNotify: 4:CheckSign")
	if err = pn.CheckSign(); err != nil {
		return
	}
	beego.Trace("callPayNotify: 5:CheckChannelRet")
	if err = pn.CheckChannelRet(); err != nil {
		return
	}
	beego.Trace("callPayNotify: 6:Handle")
	if err = pn.Handle(); err != nil {
		return
	}
	return
}
