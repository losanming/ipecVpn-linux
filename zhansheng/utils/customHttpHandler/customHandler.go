package customHttpHandler

import (
	"danfwing.com/m/zhansheng/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CustomHTTPErrorHandler 捕捉所有未处理的异常
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if h, ok := err.(*echo.HTTPError); ok {
		if h.Code == 401 {
			err = c.NoContent(h.Code)
			return
		} else if h.Code == 404 {
			err = c.NoContent(http.StatusOK)
			return
		} else {
			err = BadRequest(c, fmt.Sprintf("发生内部错误(%d)", h.Code))
			return
		}
	} else {
		err = BadRequest(c, "未知错误")
		return
	}
}

// JsonResult API返回结果
type JsonResult struct {
	Code    int32       `json:"code"`
	Api     string      `json:"api"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"` // maybe map or slice
}

// BadRequest 异常返回
func BadRequest(c echo.Context, msg string) error {
	j := &JsonResult{}
	j.Code = http.StatusBadRequest
	j.Api = c.Request().URL.Path
	j.Message = msg
	j.Data = ""
	info := utils.GetInfo(2)
	utils.PrintFileLog("server_err", msg+" "+info)
	return c.JSON(http.StatusBadRequest, j)
}

// SuccessRequest 成功返回
func SuccessRequest(c echo.Context, data interface{}) error {
	j := &JsonResult{}
	j.Code = http.StatusOK
	j.Api = c.Request().URL.Path
	j.Message = "执行成功"
	j.Data = data
	return c.JSON(http.StatusOK, j)
}
