package login

import (
	"u9/api/controllers"
	"u9/api/controllers/login/loginRequestHandle"
)

type LoginController struct {
	base.BaseController
	lrParam loginRequestHandle.Param
	lrRet   LoginRequestRet
}
