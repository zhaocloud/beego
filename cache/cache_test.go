// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("zhaocloud", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("zhaocloud") {
		t.Error("check err")
	}

	if v := bm.Get("zhaocloud"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("zhaocloud") {
		t.Error("check err")
	}

	if err = bm.Put("zhaocloud", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("zhaocloud"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("zhaocloud"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("zhaocloud"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("zhaocloud"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("zhaocloud")
	if bm.IsExist("zhaocloud") {
		t.Error("delete err")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"/cache","FileSuffix":".bin","DirectoryLevel":2,"EmbedExpiry":0}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("zhaocloud", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("zhaocloud") {
		t.Error("check err")
	}

	if v := bm.Get("zhaocloud"); v.(int) != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("zhaocloud"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("zhaocloud"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("zhaocloud"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("zhaocloud"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("zhaocloud")
	if bm.IsExist("zhaocloud") {
		t.Error("delete err")
	}
	//test string
	if err = bm.Put("zhaocloud", "author", 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("zhaocloud") {
		t.Error("check err")
	}

	if v := bm.Get("zhaocloud"); v.(string) != "author" {
		t.Error("get err")
	}
}
