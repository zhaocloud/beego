package controllers

import (
	"github.com/zhaocloud/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["host"] = this.Ctx.Request.Host
	this.TplNames = "index.tpl"
}
