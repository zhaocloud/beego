// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package beego

import (
    "bufio"
    "errors"
    "fmt"
    "net"
    "net/http"
    "net/url"
    "reflect"
    "regexp"
    "runtime"
    "strings"
    "strconv"
    "time"

    beecontext "github.com/zhaocloud/beego/context"
    "github.com/zhaocloud/beego/middleware"
    "github.com/zhaocloud/beego/toolbox"
    "github.com/zhaocloud/beego/utils"
    "github.com/zhaocloud/beego/zhaoutils"
)

const (
    // default filter execution points
    BeforeRouter = iota
    AfterStatic
    BeforeExec
    AfterExec
    FinishRouter
)

const (
    routerTypeBeego = iota
    routerTypeRESTFul
    routerTypeHandler
)

var (
    //zhao auth
    ZhaoAuth *zhaoutils.ZhaoAuth
    // supported http methods.
    HTTPMETHOD = []string{"get", "post", "put", "delete", "patch", "options", "head", "trace", "connect"}
    // these beego.Controller's methods shouldn't reflect to AutoRouter
    exceptMethod = []string{"Init", "Prepare", "Finish", "Render", "RenderString",
        "RenderBytes", "Redirect", "Abort", "StopRun", "UrlFor", "ServeJson", "ServeJsonp",
        "ServeXml", "Input", "ParseForm", "GetString", "GetStrings", "GetInt", "GetBool",
        "GetFloat", "GetFile", "SaveToFile", "StartSession", "SetSession", "GetSession",
        "DelSession", "SessionRegenerateID", "DestroySession", "IsAjax", "GetSecureCookie",
        "SetSecureCookie", "XsrfToken", "CheckXsrfCookie", "XsrfFormHtml",
        "GetControllerAndAction"}
)

// To append a slice's value into "exceptMethod", for controller's methods shouldn't reflect to AutoRouter
func ExceptMethodAppend(action string) {
    exceptMethod = append(exceptMethod, action)
}

type controllerInfo struct {
    pattern        string
    regex          *regexp.Regexp
    params         map[int]string
    controllerType reflect.Type
    methods        map[string]string
    hasMethod      bool
    handler        http.Handler
    runfunction    FilterFunc
    routerType     int
    isPrefix       bool
}

// ControllerRegistor containers registered router rules, controller handlers and filters.
type ControllerRegistor struct {
    routers      []*controllerInfo // regexp router storage
    fixrouters   []*controllerInfo // fixed router storage
    enableFilter bool
    filters      map[int][]*FilterRouter
    enableAuto   bool
    autoRouter   map[string]map[string]reflect.Type //key:controller key:method value:reflect.type
    zhaoRouter   map[string]*zhaoControllerInfo     //半auto模式, added by odin
}

// NewControllerRegistor returns a new ControllerRegistor.
func NewControllerRegistor() *ControllerRegistor {
    return &ControllerRegistor{
        routers:    make([]*controllerInfo, 0),
        autoRouter: make(map[string]map[string]reflect.Type),
        filters:    make(map[int][]*FilterRouter),
        zhaoRouter: make(map[string]*zhaoControllerInfo),
    }
}

// Add controller handler and pattern rules to ControllerRegistor.
// usage:
//	default methods is the same name as method
//	Add("/user",&UserController{})
//	Add("/api/list",&RestController{},"*:ListFood")
//	Add("/api/create",&RestController{},"post:CreateFood")
//	Add("/api/update",&RestController{},"put:UpdateFood")
//	Add("/api/delete",&RestController{},"delete:DeleteFood")
//	Add("/api",&RestController{},"get,post:ApiFunc")
//	Add("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
func (p *ControllerRegistor) Add(pattern string, c ControllerInterface, mappingMethods ...string) {
    j, params, parts := p.splitRoute(pattern)
    reflectVal := reflect.ValueOf(c)
    t := reflect.Indirect(reflectVal).Type()
    methods := make(map[string]string)
    if len(mappingMethods) > 0 {
        semi := strings.Split(mappingMethods[0], ";")
        for _, v := range semi {
            colon := strings.Split(v, ":")
            if len(colon) != 2 {
                panic("method mapping format is invalid")
            }
            comma := strings.Split(colon[0], ",")
            for _, m := range comma {
                if m == "*" || utils.InSlice(strings.ToLower(m), HTTPMETHOD) {
                    if val := reflectVal.MethodByName(colon[1]); val.IsValid() {
                        methods[strings.ToLower(m)] = colon[1]
                    } else {
                        panic(colon[1] + " method doesn't exist in the controller " + t.Name())
                    }
                } else {
                    panic(v + " is an invalid method mapping. Method doesn't exist " + m)
                }
            }
        }
    }
    if j == 0 {
        //now create the Route
        route := &controllerInfo{}
        route.pattern = pattern
        route.controllerType = t
        route.methods = methods
        route.routerType = routerTypeBeego
        if len(methods) > 0 {
            route.hasMethod = true
        }
        p.fixrouters = append(p.fixrouters, route)
    } else { // add regexp routers
        //recreate the url pattern, with parameters replaced
        //by regular expressions. then compile the regex
        pattern = strings.Join(parts, "/")
        regex, regexErr := regexp.Compile(pattern)
        if regexErr != nil {
            //TODO add error handling here to avoid panic
            panic(regexErr)
        }

        //now create the Route

        route := &controllerInfo{}
        route.regex = regex
        route.params = params
        route.pattern = pattern
        route.methods = methods
        route.routerType = routerTypeBeego
        if len(methods) > 0 {
            route.hasMethod = true
        }
        route.controllerType = t
        p.routers = append(p.routers, route)
    }
}

// add get method
// usage:
//    Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Get(pattern string, f FilterFunc) {
    p.AddMethod("get", pattern, f)
}

// add post method
// usage:
//    Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Post(pattern string, f FilterFunc) {
    p.AddMethod("post", pattern, f)
}

// add put method
// usage:
//    Put("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Put(pattern string, f FilterFunc) {
    p.AddMethod("put", pattern, f)
}

// add delete method
// usage:
//    Delete("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Delete(pattern string, f FilterFunc) {
    p.AddMethod("delete", pattern, f)
}

// add head method
// usage:
//    Head("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Head(pattern string, f FilterFunc) {
    p.AddMethod("head", pattern, f)
}

// add patch method
// usage:
//    Patch("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Patch(pattern string, f FilterFunc) {
    p.AddMethod("patch", pattern, f)
}

// add options method
// usage:
//    Options("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Options(pattern string, f FilterFunc) {
    p.AddMethod("options", pattern, f)
}

// add all method
// usage:
//    Any("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) Any(pattern string, f FilterFunc) {
    p.AddMethod("*", pattern, f)
}

// add http method router
// usage:
//    AddMethod("get","/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegistor) AddMethod(method, pattern string, f FilterFunc) {
    if method != "*" && !utils.InSlice(strings.ToLower(method), HTTPMETHOD) {
        panic("not support http method: " + method)
    }
    route := &controllerInfo{}
    route.routerType = routerTypeRESTFul
    route.runfunction = f
    methods := make(map[string]string)
    if method == "*" {
        for _, val := range HTTPMETHOD {
            methods[val] = val
        }
    } else {
        methods[method] = method
    }
    route.methods = methods
    paramnums, params, parts := p.splitRoute(pattern)
    if paramnums == 0 {
        //now create the Route
        route.pattern = pattern
        p.fixrouters = append(p.fixrouters, route)
    } else {
        //recreate the url pattern, with parameters replaced
        //by regular expressions. then compile the regex
        pattern = strings.Join(parts, "/")
        regex, regexErr := regexp.Compile(pattern)
        if regexErr != nil {
            panic(regexErr)
        }
        //now create the Route
        route.regex = regex
        route.params = params
        route.pattern = pattern
        p.routers = append(p.routers, route)
    }
}

func (p *ControllerRegistor) Handler(pattern string, h http.Handler, options ...interface{}) {
    paramnums, params, parts := p.splitRoute(pattern)
    route := &controllerInfo{}
    route.routerType = routerTypeHandler
    route.handler = h
    if len(options) > 0 {
        if v, ok := options[0].(bool); ok {
            route.isPrefix = v
        }
    }
    if paramnums == 0 {
        route.pattern = pattern
        p.fixrouters = append(p.fixrouters, route)
    } else {
        //recreate the url pattern, with parameters replaced
        //by regular expressions. then compile the regex
        pattern = strings.Join(parts, "/")
        regex, regexErr := regexp.Compile(pattern)
        if regexErr != nil {
            panic(regexErr)
        }
        //now create the Route
        route.regex = regex
        route.params = params
        route.pattern = pattern
        p.routers = append(p.routers, route)
    }
}

// analisys the patter to params & parts
func (p *ControllerRegistor) splitRoute(pattern string) (paramnums int, params map[int]string, parts []string) {
    parts = strings.Split(pattern, "/")
    j := 0
    params = make(map[int]string)
    for i, part := range parts {
        if strings.HasPrefix(part, ":") {
            expr := "(.*)"
            //a user may choose to override the defult expression
            // similar to expressjs: ‘/user/:id([0-9]+)’
            if index := strings.Index(part, "("); index != -1 {
                expr = part[index:]
                part = part[:index]
                //match /user/:id:int ([0-9]+)
                //match /post/:username:string	([\w]+)
            } else if lindex := strings.LastIndex(part, ":"); lindex != 0 {
                switch part[lindex:] {
                case ":int":
                    expr = "([0-9]+)"
                    part = part[:lindex]
                case ":string":
                    expr = `([\w]+)`
                    part = part[:lindex]
                }
                //marth /user/:id! non-empty value
            } else if part[len(part)-1] == '!' {
                expr = `(.+)`
                part = part[:len(part)-1]
            }
            params[j] = part
            parts[i] = expr
            j++
        }
        if strings.HasPrefix(part, "*") {
            expr := "(.*)"
            if part == "*.*" {
                params[j] = ":path"
                parts[i] = "([^.]+).([^.]+)"
                j++
                params[j] = ":ext"
                j++
            } else {
                params[j] = ":splat"
                parts[i] = expr
                j++
            }
        }
        //url like someprefix:id(xxx).html
        if strings.Contains(part, ":") && strings.Contains(part, "(") && strings.Contains(part, ")") {
            var out []rune
            var start bool
            var startexp bool
            var param []rune
            var expt []rune
            for _, v := range part {
                if start {
                    if v != '(' {
                        param = append(param, v)
                        continue
                    }
                }
                if startexp {
                    if v != ')' {
                        expt = append(expt, v)
                        continue
                    }
                }
                if v == ':' {
                    param = make([]rune, 0)
                    param = append(param, ':')
                    start = true
                } else if v == '(' {
                    startexp = true
                    start = false
                    params[j] = string(param)
                    j++
                    expt = make([]rune, 0)
                    expt = append(expt, '(')
                } else if v == ')' {
                    startexp = false
                    expt = append(expt, ')')
                    out = append(out, expt...)
                } else {
                    out = append(out, v)
                }
            }
            parts[i] = string(out)
        }
    }
    return j, params, parts
}

// Add auto router to ControllerRegistor.
// example beego.AddAuto(&MainContorlller{}),
// MainController has method List and Page.
// visit the url /main/list to execute List function
// /main/page to execute Page function.
func (p *ControllerRegistor) AddAuto(c ControllerInterface) {
    p.enableAuto = true
    reflectVal := reflect.ValueOf(c)
    rt := reflectVal.Type()
    ct := reflect.Indirect(reflectVal).Type()
    firstParam := strings.ToLower(strings.TrimSuffix(ct.Name(), "Controller"))
    if _, ok := p.autoRouter[firstParam]; ok {
        return
    } else {
        p.autoRouter[firstParam] = make(map[string]reflect.Type)
    }
    for i := 0; i < rt.NumMethod(); i++ {
        if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
            p.autoRouter[firstParam][rt.Method(i).Name] = ct
        }
    }
}

// Add auto router to ControllerRegistor with prefix.
// example beego.AddAutoPrefix("/admin",&MainContorlller{}),
// MainController has method List and Page.
// visit the url /admin/main/list to execute List function
// /admin/main/page to execute Page function.
func (p *ControllerRegistor) AddAutoPrefix(prefix string, c ControllerInterface) {
    p.enableAuto = true
    reflectVal := reflect.ValueOf(c)
    rt := reflectVal.Type()
    ct := reflect.Indirect(reflectVal).Type()
    firstParam := strings.Trim(prefix, "/") + "/" + strings.ToLower(strings.TrimSuffix(ct.Name(), "Controller"))
    if _, ok := p.autoRouter[firstParam]; ok {
        return
    } else {
        p.autoRouter[firstParam] = make(map[string]reflect.Type)
    }
    for i := 0; i < rt.NumMethod(); i++ {
        if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
            p.autoRouter[firstParam][rt.Method(i).Name] = ct
        }
    }
}

// [Deprecated] use InsertFilter.
// Add FilterFunc with pattern for action.
func (p *ControllerRegistor) AddFilter(pattern, action string, filter FilterFunc) error {
    mr, err := buildFilter(pattern, filter)
    if err != nil {
        return err
    }

    switch action {
    case "BeforeRouter":
        p.filters[BeforeRouter] = append(p.filters[BeforeRouter], mr)
    case "AfterStatic":
        p.filters[AfterStatic] = append(p.filters[AfterStatic], mr)
    case "BeforeExec":
        p.filters[BeforeExec] = append(p.filters[BeforeExec], mr)
    case "AfterExec":
        p.filters[AfterExec] = append(p.filters[AfterExec], mr)
    case "FinishRouter":
        p.filters[FinishRouter] = append(p.filters[FinishRouter], mr)
    }
    p.enableFilter = true
    return nil
}

// Add a FilterFunc with pattern rule and action constant.
func (p *ControllerRegistor) InsertFilter(pattern string, pos int, filter FilterFunc) error {
    mr, err := buildFilter(pattern, filter)
    if err != nil {
        return err
    }
    p.filters[pos] = append(p.filters[pos], mr)
    p.enableFilter = true
    return nil
}

// UrlFor does another controller handler in this request function.
// it can access any controller method.
func (p *ControllerRegistor) UrlFor(endpoint string, values ...string) string {
    paths := strings.Split(endpoint, ".")
    if len(paths) <= 1 {
        Warn("urlfor endpoint must like path.controller.method")
        return ""
    }
    if len(values)%2 != 0 {
        Warn("urlfor params must key-value pair")
        return ""
    }
    urlv := url.Values{}
    if len(values) > 0 {
        key := ""
        for k, v := range values {
            if k%2 == 0 {
                key = v
            } else {
                urlv.Set(key, v)
            }
        }
    }
    controllName := strings.Join(paths[:len(paths)-1], ".")
    methodName := paths[len(paths)-1]
    for _, route := range p.fixrouters {
        if route.controllerType.Name() == controllName {
            var finded bool
            if utils.InSlice(strings.ToLower(methodName), HTTPMETHOD) {
                if route.hasMethod {
                    if m, ok := route.methods[strings.ToLower(methodName)]; ok && m != methodName {
                        finded = false
                    } else if m, ok = route.methods["*"]; ok && m != methodName {
                        finded = false
                    } else {
                        finded = true
                    }
                } else {
                    finded = true
                }
            } else if route.hasMethod {
                for _, md := range route.methods {
                    if md == methodName {
                        finded = true
                    }
                }
            }
            if !finded {
                continue
            }
            if len(values) > 0 {
                return route.pattern + "?" + urlv.Encode()
            }
            return route.pattern
        }
    }
    for _, route := range p.routers {
        if route.controllerType.Name() == controllName {
            var finded bool
            if utils.InSlice(strings.ToLower(methodName), HTTPMETHOD) {
                if route.hasMethod {
                    if m, ok := route.methods[strings.ToLower(methodName)]; ok && m != methodName {
                        finded = false
                    } else if m, ok = route.methods["*"]; ok && m != methodName {
                        finded = false
                    } else {
                        finded = true
                    }
                } else {
                    finded = true
                }
            } else if route.hasMethod {
                for _, md := range route.methods {
                    if md == methodName {
                        finded = true
                    }
                }
            }
            if !finded {
                continue
            }
            var returnurl string
            var i int
            var startreg bool
            for _, v := range route.regex.String() {
                if v == '(' {
                    startreg = true
                    continue
                } else if v == ')' {
                    startreg = false
                    returnurl = returnurl + urlv.Get(route.params[i])
                    i++
                } else if !startreg {
                    returnurl = string(append([]rune(returnurl), v))
                }
            }
            if route.regex.MatchString(returnurl) {
                return returnurl
            }
        }
    }
    if p.enableAuto {
        for cName, methodList := range p.autoRouter {
            if strings.ToLower(strings.TrimSuffix(paths[len(paths)-2], "Controller")) == cName {
                if _, ok := methodList[methodName]; ok {
                    if len(values) > 0 {
                        return "/" + strings.TrimSuffix(paths[len(paths)-2], "Controller") + "/" + methodName + "?" + urlv.Encode()
                    } else {
                        return "/" + strings.TrimSuffix(paths[len(paths)-2], "Controller") + "/" + methodName
                    }
                }
            }
        }
    }
    return ""
}

// Implement http.Handler interface.
func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    starttime := time.Now()
    requestPath := r.URL.Path
    var runrouter reflect.Type
    var runMethod string
    var endpoint string
    params := make(map[string]string)
    ZhaoAuth = &zhaoutils.ZhaoAuth{}

    w := &responseWriter{writer: rw}
    w.Header().Set("Server", BeegoServerName)
    //w.Header().Set("Date", starttime.Format(time.RFC1123))

    // init context
    context := &beecontext.Context{
        ResponseWriter: w,
        Request:        r,
        Input:          beecontext.NewInput(r),
        Output:         beecontext.NewOutput(),
    }
    context.Output.Context = context
    context.Output.EnableGzip = EnableGzip

    if context.Input.IsWebsocket() {
        context.ResponseWriter = rw
    }

    defer func() {
        if err := recover(); err != nil {
            var stack string
            Critical("the request url is ", r.URL.Path)
            Critical("Handler crashed with error", err)
            for i := 1; ; i++ {
                _, file, line, ok := runtime.Caller(i)
                if !ok {
                    break
                }
                Critical(file, line)
                stack = stack + fmt.Sprintln(file, line)
            }
            context.Output.RESTPanic(err)
        }
        //save access log
        SaveAccess(starttime, context)
    }()

    //解析请求
    parseRequest(requestPath, context)

    // session init
    if SessionOn {
        context.Input.CruSession = GlobalSessions.SessionStart(w, r)
        defer func() {
            context.Input.CruSession.SessionRelease(w)
        }()
    }

    if !utils.InSlice(strings.ToLower(r.Method), HTTPMETHOD) {
        context.Output.RESTMethodNotAllowed(errors.New("Method Not Allowed"))
        goto Admin
    }
    //static file server
    if serverStaticRouter(context) {
        goto Admin
    }

    if !context.Input.IsGet() && !context.Input.IsHead() {
        if CopyRequestBody && !context.Input.IsUpload() {
            context.Input.CopyBody()
        }
        context.Input.ParseFormOrMulitForm(MaxMemory)
    }

    if ep := context.Input.GetData("_endpoint"); ep != nil {
        endpoint = ep.(string)
    }

    //zhao auth
    if endpoint != "" && !utils.InSlice(endpoint, SkipAuth) {
        if err := ZhaoAuth.CheckZhaoAuth(r); err != nil {
            Critical("zhao auth failed: ", err)
            Critical("zhao auth failed: signature:", ZhaoAuth.Signature, ", expectedSign:", ZhaoAuth.ExpectedSign, ", keyid:",
            ZhaoAuth.SecretKeyID, ", key:", ZhaoAuth.SecretKey, ", date:", ZhaoAuth.Date, ", path:", ZhaoAuth.Path, ", uid:", ZhaoAuth.ClientUniqueID, ", uid:", strconv.Quote(ZhaoAuth.StringToSign))
            context.Output.RESTUnauthorized(err)
            goto Admin
        }
    }

    if endpoint != "" {
        //Debug("endpoint: ", endpoint)
        if route, ok := p.zhaoRouter[endpoint]; ok {
            runMethod = p.getZhaoRunMethod(r.Method, context, route)
            if runMethod != "" {
                context.Input.Params = params
                //routerInfo = route
                vc := reflect.New(route.controllerType)
                execController, ok := vc.Interface().(ControllerInterface)
                if !ok {
                    panic("controller is not ControllerInterface")
                }
                //call the controller init function
                execController.Init(context, route.controllerType.Name(), runMethod, vc.Interface())
                //call prepare function
                execController.Prepare()
                if !w.started {
                    //exec main logic
                    switch runMethod {
                    case "Get":
                        execController.Get()
                    case "Post":
                        execController.Post()
                    case "Delete":
                        execController.Delete()
                    case "Put":
                        execController.Put()
                    case "Head":
                        execController.Head()
                    case "Patch":
                        execController.Patch()
                    case "Options":
                        execController.Options()
                    default:
                        in := make([]reflect.Value, 0)
                        method := vc.MethodByName(runMethod)
                        method.Call(in)
                    }

                }

                // finish all runrouter. release resource
                execController.Finish()
            } else {
                context.Output.RESTMethodNotAllowed(errors.New("Method Not Allowed"))
            }
        } else {
            //Debug("not found endpoint: ", endpoint)
            context.Output.RESTBadRequest(errors.New("Bad Request"))
        }
    } else {
        //root
        context.Output.RESTForbidden(errors.New("Forbidden"))
    }

Admin:
    //admin module record QPS
    if EnableAdmin {
        timeend := time.Since(starttime)
        if FilterMonitorFunc(r.Method, requestPath, timeend) {
            if runrouter != nil {
                go toolbox.StatisticsMap.AddStatistics(r.Method, requestPath, runrouter.Name(), timeend)
            } else {
                go toolbox.StatisticsMap.AddStatistics(r.Method, requestPath, "", timeend)
            }
        }
    }
}

// there always should be error handler that sets error code accordingly for all unhandled errors.
// in order to have custom UI for error page it's necessary to override "500" error.
func (p *ControllerRegistor) getErrorHandler(errorCode string) func(rw http.ResponseWriter, r *http.Request) {
    handler := middleware.SimpleServerError
    ok := true
    if errorCode != "" {
        handler, ok = middleware.ErrorMaps[errorCode]
        if !ok {
            handler, ok = middleware.ErrorMaps["500"]
        }
        if !ok || handler == nil {
            handler = middleware.SimpleServerError
        }
    }

    return handler
}

// returns method name from request header or form field.
// sometimes browsers can't create PUT and DELETE request.
// set a form field "_method" instead.
func (p *ControllerRegistor) getRunMethod(method string, context *beecontext.Context, router *controllerInfo) string {
    method = strings.ToLower(method)
    if method == "post" && strings.ToLower(context.Input.Query("_method")) == "put" {
        method = "put"
    }
    if method == "post" && strings.ToLower(context.Input.Query("_method")) == "delete" {
        method = "delete"
    }
    if router.hasMethod {
        if m, ok := router.methods[method]; ok {
            return m
        } else if m, ok = router.methods["*"]; ok {
            return m
        } else {
            return ""
        }
    } else {
        return strings.Title(method)
    }
}

//responseWriter is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type responseWriter struct {
    writer  http.ResponseWriter
    started bool
    status  int
}

// Header returns the header map that will be sent by WriteHeader.
func (w *responseWriter) Header() http.Header {
    return w.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true.
// started means the response has sent out.
func (w *responseWriter) Write(p []byte) (int, error) {
    w.started = true
    return w.writer.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (w *responseWriter) WriteHeader(code int) {
    w.status = code
    w.started = true
    w.writer.WriteHeader(code)
}

// hijacker for http
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
    hj, ok := w.writer.(http.Hijacker)
    if !ok {
        println("supported?")
        return nil, nil, errors.New("webserver doesn't support hijacking")
    }
    return hj.Hijack()
}
