package zhaoutils

import (
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base64"
    "errors"
    "github.com/zhaocloud/beego/utils"
    "net/http"
    "strings"
)

type ZhaoAuth struct {
    SecretKeyID    string
    Signature      string
    StringToSign   string
    ClientUniqueID string
    Path           string
    Date           string
    ExpectedSign   string
}

/* {{{ CheckZhaoAuth
 */
func (za *ZhaoAuth) CheckZhaoAuth(r *http.Request) error {
    return nil
    za.ClientUniqueID = r.Header.Get("X-Zhao-DeviceId") //不区分大小写
    za.Date = r.Header.Get("Date")
    za.Path = r.URL.Path
    //if za.ClientUniqueID == "" || za.Date == "" {
    //    return errors.New("not enough info")
    //}
    // todo: check date time

    authHeader := r.Header.Get("Authorization")

    //auth的形式为 "zhao secretkeyid:signature"
    as := strings.SplitN(authHeader, " ", 2)

    if len(as) != 2 || strings.ToLower(as[0]) != "zhao" {
        return errors.New("not zhaocloud authorization")
    }

    s := strings.SplitN(as[1], ":", 2)

    if len(s) != 2 {
        return errors.New("signature error")
    }

    za.SecretKeyID = s[0]
    signature, err := base64.StdEncoding.DecodeString(s[1])
    if err != nil {
        return errors.New("base64 decode error")
    }
    za.Signature = s[1]

    // StringToSign: HTTP-Verb + "\n" + Date + "\n" + HTTP-Path + "\n" + {ClientUniqueID};
    za.StringToSign = r.Method + "\n" + za.Date + "\n" + za.Path + "\n" + za.ClientUniqueID

    secretKey := za.getSecretKey()
    if secretKey == "" {
        return errors.New("get key error")
    }

    mac := hmac.New(sha1.New, []byte(secretKey))
    mac.Write([]byte(za.StringToSign))
    expectedSign := mac.Sum(nil)
    za.ExpectedSign = base64.StdEncoding.EncodeToString(expectedSign)

    if ok := hmac.Equal(signature, expectedSign); !ok {
        return errors.New("signature error")
    }

    return nil
}

/* }}} */

/* {{{ get secretkey
 */
func (za *ZhaoAuth) getSecretKey() (sk string) {
    //暂时用算法解决，之后需要完全随机数,从高速缓存中查询
    if za.SecretKeyID == "" {
        return
    }
    sk = utils.LengthenUUID(za.SecretKeyID)
    return
}

/* }}} */
