// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package beego

import (
	"testing"
)

func TestDefaults(t *testing.T) {
	if FlashName != "BEEGO_FLASH" {
		t.Errorf("FlashName was not set to default.")
	}

	if FlashSeperator != "BEEGOFLASH" {
		t.Errorf("FlashName was not set to default.")
	}
}
