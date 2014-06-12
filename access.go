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
var (
    AccessLogger *logs.BeeLogger
)

func SaveAccess(start time.Time, context *beecontext.Context) {
    logTime := strconv.Quote(start.Format(time.RFC3339))
    logIP := strconv.Quote(context.Input.IP())
    sc := context.Output.Status
    if sc == 0 {
        //没有特别设置,则为http 200 ok
        sc = http.StatusOK
    }
    requestSize, rhs := context.Input.GetInputInfo()
    rHeaderStr := strconv.Quote(string(rhs))
    //sentSize := context.ResponseWriter.Header().Get("Content-Length")
    sentSize, hs := context.Output.GetOutputInfo(sc)
    headerStr := strconv.Quote(string(hs))
    //contentType := strconv.Quote(context.ResponseWriter.Header().Get("Content-Type"))
    logMethod := strconv.Quote(context.Input.Method())
    logUri := strconv.Quote(context.Input.Uri())
    logProtocol := strconv.Quote(context.Input.Protocol())
    requestTime := strconv.FormatFloat(time.Since(start).Seconds(), 'f', 6, 64)
    //userAgent := strconv.Quote(context.Input.UserAgent())

    Access(logTime, logIP, requestSize, sc, sentSize, logMethod, logUri, logProtocol, requestTime, headerStr, rHeaderStr)
}

func Access(v ...interface{}) {
    format := generateAccessFmt(len(v))
    AccessLogger.Access(format, v...)
}

func generateAccessFmt(n int) string {
    return strings.Repeat("%v,", n)
}
