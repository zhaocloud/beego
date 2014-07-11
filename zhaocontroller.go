package beego

import (
    "net/http"
    "strings"
)

var SUCCODE = map[string]int{
    "get":    http.StatusOK,
    "delete": http.StatusNoContent,
    "put":    http.StatusCreated,
    "post":   http.StatusCreated,
    "patch":  http.StatusResetContent,
    "head":   http.StatusOK,
}

// ServeJson sends a json response with encoding charset.
func (c *Controller) ServeREST(err error, encoding ...bool) {
    var hasIndent bool
    var hasencoding bool
    var status int
    if RunMode == "prod" {
        hasIndent = false
    } else {
        hasIndent = true
    }
    if len(encoding) > 0 && encoding[0] == true {
        hasencoding = true
    }
    if err == nil {
        method := strings.ToLower(c.Ctx.Input.Method())
        if _, ok := SUCCODE[method]; !ok {
            status = SUCCODE[method]
        } else {
            status = http.StatusOK
        }
        //c.Ctx.Output.Json(c.Data["json"], hasIndent, hasencoding)
        c.Ctx.Output.RESTJson(status, c.Data["json"], hasIndent, hasencoding)
    } else {
        c.Ctx.Output.RESTBadRequest(err)
    }
}
