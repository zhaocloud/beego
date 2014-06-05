// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package utils

import (
	"testing"
)

func TestInSlice(t *testing.T) {
	sl := []string{"A", "b"}
	if !InSlice("A", sl) {
		t.Error("should be true")
	}
	if InSlice("B", sl) {
		t.Error("should be false")
	}
}
