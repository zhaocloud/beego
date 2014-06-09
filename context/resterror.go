//这里处理所有要output的REST错误
package context

import (
    "fmt"
    "net/http"
)

type RESTErrorData struct {
    Massage string            `json:"massage"`
    Errors  map[string]string `json:"errors"`
}

// 基本错误
func (output *BeegoOutput) RESTError(status int, err interface{}) {
    errors := make(map[string]string)
    errors["method"] = output.Context.Input.Method()
    errors["endpoint"] = output.Context.Input.GetData("_endpoint").(string)
    errors["code"] = "service_panic"

    re := RESTErrorData{
        Massage: fmt.Sprint(err),
        Errors:  errors,
    }
    output.RESTJson(status, re, true, true)
}

// 无法预知的错误,将来可以进一步封装,不把程序错误直接output出去
func (output *BeegoOutput) RESTPanic(err interface{}) {
    output.RESTError(http.StatusInternalServerError, err)
}

// NotFound
func (output *BeegoOutput) RESTNotFound(err interface{}) {
    output.RESTError(http.StatusNotFound, err)
}

//MethodNotAllowed
func (output *BeegoOutput) RESTMethodNotAllowed(err interface{}) {
    output.RESTError(http.StatusMethodNotAllowed, err)
}

// NotFound
func (output *BeegoOutput) RESTBadRequest(err interface{}) {
    output.RESTError(http.StatusBadRequest, err)
}
