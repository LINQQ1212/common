package vhost

import (
	"github.com/LINQQ1212/common/utils"
	"github.com/jxskiss/base62"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

type Site struct {
	Domain string
	Index  int
	Table  string // 字符串 表
	Key    string // 字符串 表
	//TableMap   map[string]string `json:"-" yaml:"-"`
	ListTemp   string
	TempsLen   int
	Temps      []string // 模板还是要一次性使用
	Num        int      // 随机数
	Num2       uint64   // 随机数2
	Path       string   // 伪装路径
	CPath      string   // 分类伪装路径
	Suffix     string   // 后缀
	SA         string   // html 前
	SE         string   // html 后
	SearchPath string   `json:"search-path" yaml:"search-path"`
	IsParam    bool
	Type       int              // 方式
	HashID     *base62.Encoding `json:"-" yaml:"-"`
	MinLength  int
}

func (s *Site) BuildCUri(host, f string, id uint64) string {
	if s.IsParam {
		return host + "?" + s.IdToStr(id) + s.Suffix
	}
	return host + s.IdToStr(id) + s.Suffix
}

func (s *Site) BuildImageLink(host, f, path string) string {
	if s.IsParam {
		return host + "?images=" + strings.TrimPrefix(path, "/images/")
	}
	return host + path
}
func (s *Site) BuildPath(host, f, path string) string {
	if s.IsParam {
		return host + "?" + path
	}
	return host + strings.TrimPrefix(path, "/")
}

func (s *Site) BuildPath2(host, f, path string) string {
	if s.IsParam {
		return host + "?" + path + s.Suffix
	}
	return host + path + s.Suffix
}

func (s *Site) BuildProUri(host, f string, id uint64) string {
	if s.IsParam {
		return host + "?" + s.Path + "/" + s.Base32Id(id) + s.Suffix
	}
	return host + s.Path + "/" + s.Base32Id(id) + s.Suffix
}

func (s *Site) BuildProUri2(host, f string, cid, id uint64) string {
	if s.IsParam {
		return host + "?" + s.IdToStr(cid) + "/" + s.IdToStr(id) + s.Suffix
	}
	return host + s.IdToStr(cid) + "/" + s.IdToStr(id) + s.Suffix
}

func (s *Site) BuildProUriCName(host, f, cname string, id uint64) string {
	if s.IsParam {
		return host + "?" + url.PathEscape(cname) + "/" + s.IdToStr(id) + s.Suffix
	}
	return host + url.PathEscape(cname) + "/" + s.IdToStr(id) + s.Suffix
}

func (s *Site) BuildProUriTPath(host, f, tpName string, id uint64) string {
	if s.IsParam {
		return host + "?" + tpName + "/" + s.IdToStr(id) + s.Suffix
	}
	return host + tpName + "/" + s.IdToStr(id) + s.Suffix
}

func GetSuffix() string {
	//arr := []string{"", ".html", ".shtml", ".jsp", ".asp", ".htm"}
	return utils.Suffixs[rand.Int()%len(utils.Suffixs)]
}
func GetSearch() string {
	arr := []string{"q", "query", "s", "keys", "search"}
	return arr[rand.Int()%5]
}

func (s *Site) InitTableMap() {
	/*i := 0
	s.TableMap = map[string]string{}
	for _, r := range s.Table {
		s.TableMap[string(r)] = strconv.Itoa(i)
		i++
		if i > 9 {
			i = 0
		}
	}*/
	s.HashID = base62.NewEncoding(s.Table)
}

func (s *Site) InitHtml() {
	switch rand.Int() % 8 {
	case 0:
		s.SA = "<p>"
		s.SE = "</p>"
	case 1:
		s.SA = "<span>"
		s.SE = "</span>"
	case 2:
		s.SA = ""
		s.SE = "</br>"
	case 3:
		s.SA = "<div><p>"
		s.SE = "</p></div>"
	case 4:
		s.SA = "<P>"
		s.SE = "</P>"
	case 5:
		s.SA = "<div><span>"
		s.SE = "</span></div>"
	case 6:
		s.SA = "<div>"
		s.SE = "</div>"
	case 7:
		s.SA = "<li>"
		s.SE = "</li>"
	}
}

func (s *Site) Str2Id(str string) uint64 {
	if str == "" {
		return 0
	}
	var id uint64
	for i, i2 := range str {
		id += uint64(i2 << i)
	}
	if id < 0 {
		id = ^id
	}
	return id * uint64(len(str)+s.Index)
}

func (s *Site) IdToStr(id uint64) string {
	return strconv.FormatUint(id+s.Num2, 10)
}
func (s *Site) StrToId(str string) uint64 {
	isNum := true
	for _, i2 := range str {
		if !unicode.IsNumber(i2) {
			isNum = false
			break
		}
	}
	if isNum {
		num, _ := strconv.ParseUint(str, 10, 64)
		if num <= s.Num2 {
			return num
		}
		return num - s.Num2
	}
	return s.Str2Id(str)
}

func (s *Site) Base32Id(id uint64) string {
	b := s.HashID.Encode([]byte(strconv.FormatUint(id, 10)))
	return *(*string)(unsafe.Pointer(&b))
}

func (s *Site) HexId(str string) string {
	b := s.HashID.Encode([]byte(str))
	return *(*string)(unsafe.Pointer(&b))
}

func (s *Site) Base32ToId(str string) uint64 {
	b, err := s.HashID.Decode([]byte(str))
	if err == nil {
		if i, err := strconv.ParseUint(*(*string)(unsafe.Pointer(&b)), 10, 64); err == nil {
			return i
		}
	}
	return s.Str2Id(str)
}

func (s *Site) Base32ToIdIsMy(str string) (uint64, bool) {
	b, err := s.HashID.Decode([]byte(str))
	if err == nil {
		if i, err := strconv.ParseUint(*(*string)(unsafe.Pointer(&b)), 10, 64); err == nil {
			return i, true
		}
		return s.Str2Id(str), true
	}
	return s.Str2Id(str), false
}

func (s *Site) IdToRandStr(id uint64) string {
	var b []byte
	cache := id
	v := false
	for i := 0; i < 10; i++ {
		if v {
			break
		}
		if idx := int(cache & 63); idx < len(s.Table) {
			if idx == 0 {
				v = true
			}
			b = append(b, s.Table[idx])
			cache >>= 6
		}
	}
	return *(*string)(unsafe.Pointer(&b))
}
