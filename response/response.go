package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func Data404WithMessage(message string, c *gin.Context) {
	c.Writer.WriteString("4404 ")
	c.Writer.WriteString(message)
	c.Status(http.StatusNotFound)
	c.Abort()
}

func DataMain404(c *gin.Context) {
	c.Status(http.StatusNotFound)
	c.Abort()
}

func DataMain500(c *gin.Context) {
	c.Status(http.StatusInternalServerError)
	c.Abort()
}

func DataMain302(c *gin.Context, link string) {
	RedirectHtml(c, link)
	return
	c.Redirect(http.StatusFound, link)
	c.Abort()
}

func Data404(c *gin.Context) {
	c.Writer.WriteString("4404")
	c.Status(http.StatusNotFound)
	c.Abort()
}

func Data500(c *gin.Context) {
	c.Writer.WriteString("5500")
	c.Status(http.StatusInternalServerError)
	c.Abort()
}

func Data302(c *gin.Context, link string) {
	RedirectHtml(c, link)
	return
	c.Writer.WriteString("3" + link)
	c.Status(http.StatusOK)
	c.Abort()
}

func RedirectHtml(c *gin.Context, link string) {
	c.Writer.WriteString(`<html><head><meta charset="utf-8"><meta http-equiv="refresh" content="0; url=` + link + `" /><script type="text/javascript">window.location.href="` + link + `";</script></head></html>`)
	c.Status(http.StatusOK)
	c.Abort()
}
