package main

import (
	"github.com/zhaocloud/beego"
	"github.com/zhaocloud/beego/example/chat/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WSController{})
	beego.Run()
}
