// zhao-cloud.com
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package beego

import (
    beecontext "github.com/zhaocloud/beego/context"
    "github.com/zhaocloud/beego/utils"
    "reflect"
    "strings"
)

type zhaoControllerInfo struct {
    controllerType reflect.Type
    methods        map[string]string
}

//通过这个函数注册controller, 每个controller的init()都需要调用这个函数
func RegistZhaoController(c ControllerInterface) {
    p := BeeApp.Handlers
    reflectVal := reflect.ValueOf(c)
    rt := reflectVal.Type()
    ct := reflect.Indirect(reflectVal).Type()
    cName := strings.ToLower(strings.TrimSuffix(ct.Name(), "Controller"))
    if _, ok := p.zhaoRouter[cName]; ok {
        return
    } else {
        route := &zhaoControllerInfo{}
        route.controllerType = ct
        methods := make(map[string]string)
        for i := 0; i < rt.NumMethod(); i++ {
            if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
                methods[strings.ToLower(rt.Method(i).Name)] = rt.Method(i).Name
            }
        }
        route.methods = methods
        p.zhaoRouter[cName] = route
    }
}

func (p *ControllerRegistor) getZhaoRunMethod(method string, ctx *beecontext.Context, router *zhaoControllerInfo) string {
    method = strings.ToLower(method)
    if method == "post" && strings.ToLower(ctx.Input.Query("_method")) == "put" {
        method = "put"
    }
    if method == "post" && strings.ToLower(ctx.Input.Query("_method")) == "delete" {
        method = "delete"
    }
    if _, ok := router.methods[method]; ok {
        return strings.Title(method)
    } else {
        Critical("router has no match method: ", method)
        return ""
    }
}

//parse request path
func parseRequest(path string, ctx *beecontext.Context) {
    requestPieces := strings.Split(path, "/")
    for off, piece := range requestPieces {
        if piece != "" {
            if off == 1 {
                ctx.Input.SetData("_endpoint", piece)
            }
            if off == 2 && piece[0] != '@' { //@开头是selector
                ctx.Input.SetData("_rowkey", piece)
            }
            if off > 1 && piece[0] == '@' {
                ctx.Input.SetData("_selector", piece)
            }
        }
    }
}
