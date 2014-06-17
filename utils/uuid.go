// @description make uuid, support short uuid
// @reference   https://github.com/stochastic-technologies/shortuuid
// @authors     Odin

package utils

import (
    "code.google.com/p/go-uuid/uuid"
    "fmt"
    "math/big"
    "strings"
)

const (
    alphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
    shortlen = 22 //这个长度是根据alphabet得来的,省略计算步骤
)

func shortenUUID(s string) (ss string) {
    uuInt, tr, tm, off := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
    //remove "-"
    s = strings.ToLower(strings.Replace(s, "-", "", -1))
    fmt.Sscan("0x"+s, uuInt)

    alphaLen := big.NewInt(int64(len(alphabet)))
    for uuInt.Cmp(big.NewInt(0)) > 0 {
        uuInt, off = tr.DivMod(uuInt, alphaLen, tm)
        ss += string(alphabet[off.Int64()])
    }
    //如果不足22位,用第一个字符补全
    if diff := shortlen - len(ss); diff > 0 {
        ss += strings.Repeat(string(alphabet[0]), diff)
    }

    return
}

//default, version 4
func NewUUID() string {
    return uuid.New()
}

func NewShortUUID() string {
    newUUID := NewUUID()
    return shortenUUID(newUUID)
}

// generate uuid v5
// namespace直接写域名或者url
func NewUUID5(namespace, data string) string {
    //以http开头的为URL, 其余都为DNS
    if strings.HasPrefix(strings.ToLower(namespace), "http") {
        return uuid.NewSHA1(uuid.NameSpace_URL, []byte(namespace+data)).String()
    } else {
        return uuid.NewSHA1(uuid.NameSpace_DNS, []byte(namespace+data)).String()
    }
}

func NewShortUUID5(namespace, data string) string {
    newUUID := NewUUID5(namespace, data)
    return shortenUUID(newUUID)
}
