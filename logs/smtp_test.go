// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package logs

import (
	"testing"
	"time"
)

func TestSmtp(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("smtp", `{"username":"beegotest@gmail.com","password":"xxxxxxxx","host":"smtp.gmail.com:587","sendTos":["xiemengjun@gmail.com"]}`)
	log.Critical("sendmail critical")
	time.Sleep(time.Second * 30)
}
