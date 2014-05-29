package beego

import (
    "github.com/zhaocloud/beego/logs"
    "strings"
)

// logger references the used application logger.
var AccessLogger *logs.BeeLogger

func Access(v ...interface{}) {
    format := generateAccessFmt(len(v))
    AccessLogger.Access(format, v...)
}

func generateAccessFmt(n int) string {
    return strings.Repeat("%v,", n)
}
