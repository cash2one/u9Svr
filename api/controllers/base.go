package base

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"strconv"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Validate(obj interface{}) (code int) {
	if err := this.ParseForm(obj); err != nil {
		beego.Error(err)
		return
	}

	valid := validation.Validation{}
	if valid.Valid(obj); valid.HasErrors() {
		for _, verr := range valid.Errors {
			if code, cerr := strconv.Atoi(verr.Key); cerr == nil {
				//beego.Error(verr.Message)
				return code
			}
		}
	}
	return 0
}
