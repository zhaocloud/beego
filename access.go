package beego

import (
    beecontext "github.com/zhaocloud/beego/context"
    "github.com/zhaocloud/beego/logs"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// logger references the used application logger.
var AccessLogger *logs.BeeLogger

func SaveAccess(start time.Time, context *beecontext.Context) {
    sc := context.Output.Status
    if sc == 0 {
        //没有特别设置,则为http 200 ok
        sc = http.StatusOK
    }
    Access(strconv.Quote(start.Format(time.RFC3339)), strconv.Quote(context.Input.IP()), sc, context.ResponseWriter.Header().Get("Content-Length"), strconv.Quote(context.ResponseWriter.Header().Get("Content-Type")), strconv.Quote(context.Input.Method()), strconv.Quote(context.Input.Uri()), strconv.Quote(context.Input.Protocol()), strconv.FormatFloat(time.Since(start).Seconds(), 'f', 6, 64))
}

func Access(v ...interface{}) {
    format := generateAccessFmt(len(v))
    AccessLogger.Access(format, v...)
}

func generateAccessFmt(n int) string {
    return strings.Repeat("%v,", n)
}
