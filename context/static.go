package context

import (
    "net/http"
)

func (ctx *Context) ServeFile(file string) {
    ctx.Output.SetStatus(http.StatusOK)
    http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
}
