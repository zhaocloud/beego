// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package toolbox

import (
	"os"
	"testing"
)

func TestProcessInput(t *testing.T) {
	ProcessInput("lookup goroutine", os.Stdout)
	ProcessInput("lookup heap", os.Stdout)
	ProcessInput("lookup threadcreate", os.Stdout)
	ProcessInput("lookup block", os.Stdout)
	ProcessInput("gc summary", os.Stdout)
}
