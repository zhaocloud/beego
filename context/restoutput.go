package context

import (
    "compress/flate"
    "compress/gzip"
    "encoding/json"
    "io"
    "net/http"
    "strconv"
    "strings"
)

func (output *BeegoOutput) RESTJson(status int, data interface{}, hasIndent bool, coding bool) error {
    output.Header("Content-Type", "application/json; charset=UTF-8")
    var content []byte
    var err error
    if hasIndent {
        content, err = json.MarshalIndent(data, "", "  ")
    } else {
        content, err = json.Marshal(data)
    }
    if err != nil {
        http.Error(output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
        return err
    }
    if coding {
        content = []byte(stringsToJson(string(content)))
    }
    output_writer := output.Context.ResponseWriter.(io.Writer)
    if output.EnableGzip == true && output.Context.Input.Header("Accept-Encoding") != "" {
        splitted := strings.SplitN(output.Context.Input.Header("Accept-Encoding"), ",", -1)
        encodings := make([]string, len(splitted))

        for i, val := range splitted {
            encodings[i] = strings.TrimSpace(val)
        }
        for _, val := range encodings {
            if val == "gzip" {
                output.Header("Content-Encoding", "gzip")
                output_writer, _ = gzip.NewWriterLevel(output.Context.ResponseWriter, gzip.BestSpeed)

                break
            } else if val == "deflate" {
                output.Header("Content-Encoding", "deflate")
                output_writer, _ = flate.NewWriter(output.Context.ResponseWriter, flate.BestSpeed)
                break
            }
        }
    } else {
        output.Header("Content-Length", strconv.Itoa(len(content)))
    }
    output.SetStatus(status)
    //output.Context.ResponseWriter.WriteHeader(status)
    output_writer.Write(content)
    switch output_writer.(type) {
    case *gzip.Writer:
        output_writer.(*gzip.Writer).Close()
    case *flate.Writer:
        output_writer.(*flate.Writer).Close()
    }
    return nil
}

// 获取回应的总长度,包括Header, 以及 response Header的内容
func (output *BeegoOutput) GetOutputInfo(code int) (l int, hs []byte) {
    //status line
    codestring := strconv.Itoa(code)
    text, ok := StatusText[code]
    if !ok {
        text = "status code " + codestring
    }
    l = len("HTTP/1.1 "+codestring+" "+text) + len(CRLF) //HTTP/1.1 200 OK\r\n
    //headers
    headers := output.Context.ResponseWriter.Header()
    for key, vals := range headers {
        for _, val := range vals {
            l += len(key) + len(ColonSpace) + len(val) + len(CRLF)
        }
    }
    l += len(CRLF) * 2 //最后连续两个\r\n才到body
    cl, _ := strconv.Atoi(headers.Get("Content-Length"))
    l += cl
    hs, _ = json.Marshal(headers)

    return l, hs
}

func (input *BeegoInput) GetInputInfo() (l int, hs []byte) {
    l = len(input.Method()+" "+input.Uri()+" "+input.Protocol()) + len(CRLF)

    headers := input.Request.Header
    for key, vals := range headers {
        for _, val := range vals {
            l += len(key) + len(ColonSpace) + len(val) + len(CRLF)
        }
    }
    l += len(CRLF) * 2 //最后连续两个\r\n才到body
    cl, _ := strconv.Atoi(headers.Get("Content-Length"))
    l += cl
    hs, _ = json.Marshal(headers)

    return l, hs
}
