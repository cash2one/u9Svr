package login

import (
	"u9/api/controllers"
	"u9/api/controllers/login/lrBefore"
	"u9/models"
)

type LoginController struct {
	base.BaseController
	lr      models.LoginRequest
	lrParam lrBefore.Param
	lrRet   LoginRequestRet
}
