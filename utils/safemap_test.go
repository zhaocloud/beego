// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package utils

import (
	"testing"
)

func Test_beemap(t *testing.T) {
	bm := NewBeeMap()
	if !bm.Set("zhaocloud", 1) {
		t.Error("set Error")
	}
	if !bm.Check("zhaocloud") {
		t.Error("check err")
	}

	if v := bm.Get("zhaocloud"); v.(int) != 1 {
		t.Error("get err")
	}

	bm.Delete("zhaocloud")
	if bm.Check("zhaocloud") {
		t.Error("delete err")
	}
}
