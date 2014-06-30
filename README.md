## beego

修改自github.com/astaxie/beego

## Original Features

* RESTful support
* MVC architecture
* Session support (store in memory, file, Redis or MySQL)
* Cache support (store in memory, Redis or Memcache)
* Global Config
* Intelligent routing
* Thread-safe map
* Friendly displaying of errors
* Useful template functions

## Added Features

* Daemonize, 可在app.conf中用 Daemonize={bool}配置, pidfile默认写到程序目录的run/<appname>.pid
* Access Logging, 默认写到程序目录的logs/access.log
* 简化路由功能, 只支持/{endpoint}/{rowkeys}/{selector}调度

## 使用
* main.go去除router.go,加上`_ "zhaoonline.com/applications/controllers"`

* controller对应endpoint需要如下代码(字符串对应,首字母大写,尖括号<>包含的内容可变):

```
type <Endpoint>Controller struct {
        beego.Controller
}

func init() {
    beego.RegistZhaoController(&<Endpoint>Controller{})
}
```

