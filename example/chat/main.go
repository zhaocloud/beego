// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     Unknwon
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
