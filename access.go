package beego

import (
    "github.com/zhaocloud/beego/logs"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// logger references the used application logger.
var AccessLogger *logs.BeeLogger

func SaveAccess(start time.Time, w *responseWriter, r *http.Request) {
    sc := w.status
    if sc == 0 {
        //没有特别设置,则为http 200 ok
        sc = http.StatusOK
    }
    Access(strconv.Quote(start.Format(time.RFC3339)), strconv.Quote(r.RemoteAddr), sc, w.Header().Get("Content-Length"), strconv.Quote(w.Header().Get("Content-Type")), strconv.Quote(r.Method), strconv.Quote(r.RequestURI), strconv.Quote(r.Proto), strconv.FormatFloat(time.Since(start).Seconds(), 'f', 6, 64))
}

func Access(v ...interface{}) {
    format := generateAccessFmt(len(v))
    AccessLogger.Access(format, v...)
}

func generateAccessFmt(n int) string {
    return strings.Repeat("%v,", n)
}
