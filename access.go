package beego

import (
    "github.com/zhaocloud/beego/logs"
)

// logger references the used application logger.
var AccessLogger *logs.BeeLogger

func SaveAccess(v ...interface{}) {
    AccessLogger.Access(generateFmtStr(len(v)), v...)
}
