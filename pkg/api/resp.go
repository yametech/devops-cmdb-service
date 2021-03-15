package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type message struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Code int32       `json:"code"`
}

func RequestErr(g *gin.Context, err error) {
	g.JSON(http.StatusBadRequest, &message{Data: "request not match", Msg: err.Error(), Code: 400})
	g.Abort()
}

func RequestDataErr(g *gin.Context, data string, code int32) {
	g.JSON(http.StatusOK, &message{Data: data, Msg: "request not match", Code: code})
	g.Abort()
}

func RequestOK(g *gin.Context, data interface{}) {
	g.JSON(http.StatusOK, &message{Data: data, Msg: "request success", Code: 200})
	g.Abort()
}
