package channelPayNotify

type PayNotify interface {
	Init(params ...interface{}) (err error)
	ParseInputParam(params ...interface{}) (err error)
	PrepareTradeData() (err error)
	CheckSign(params ...interface{}) (err error)
	Handle() (err error)
	CheckChannelRet(params ...interface{}) (err error)
	GetResult(params ...interface{}) (ret string)
}

var emptyUrlKeys []string = []string{}

var errorDescList []string = []string{
	"err_noerror",

	/* init */
	"err_init",
	"err_initProductKey",

	/* parseInputParam */
	"err_parseInputParam",

	"err_checkUrlParam",

	"err_initChannelParam",
	"err_initChannelGameId",
	"err_initChannelGameKey",
	"err_initChannelPayKey",

	"err_parseOrderId",
	"err_parseChannelUserId",
	"err_parseChannelOrderId",
	"err_parsePayAmount",
	"err_parsePayDiscount",

	/* prepareTradeData */
	"err_prepareOrderRequest",
	"err_prepareLoginRequest",

	/* checkSign */
	"err_checkSign",
	"err_parseRsaPublicKey",
	"err_parseRsaPrivateKey",

	/* checkChannelRet */
	"err_checkChannelRet",
	"err_tradeFail",
	"err_payAmountError",
	"err_orderIsNotExist",
	"err_channelUserIsNotExist",

	/* handle */
	"err_handleOrder",
	"err_notifyProductSvr",
}

const (
	err_noerror = iota

	/* init */
	err_init
	err_initProductKey

	/* parseInputParam */
	err_parseInputParam

	err_checkUrlParam

	err_initChannelParam
	err_initChannelGameId
	err_initChannelGameKey
	err_initChannelPayKey

	err_parseOrderId
	err_parseChannelUserId
	err_parseChannelOrderId
	err_parsePayAmount
	err_parsePayDiscount

	/* prepareTradeData */
	err_prepareOrderRequest
	err_prepareLoginRequest

	/* checkSign */
	err_checkSign
	err_parseRsaPublicKey
	err_parseRsaPrivateKey

	/* checkChannelRet */
	err_checkChannelRet
	err_tradeFail
	err_payAmountError
	err_orderIsNotExist
	err_channelUserIsNotExist

	/* handle */
	err_handleOrder
	err_notifyProductSvr
)
