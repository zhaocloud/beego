package beego

// ServeJson sends a json response with encoding charset.
func (c *Controller) ServeREST(err error, encoding ...bool) {
    var hasIndent bool
    var hasencoding bool
    if RunMode == "prod" {
        hasIndent = false
    } else {
        hasIndent = true
    }
    if len(encoding) > 0 && encoding[0] == true {
        hasencoding = true
    }
    if err == nil {
        c.Ctx.Output.Json(c.Data["json"], hasIndent, hasencoding)
    } else {
        c.Ctx.Output.RESTBadRequest(err)
    }
}
