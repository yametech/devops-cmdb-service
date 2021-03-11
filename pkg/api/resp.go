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
	g.JSON(http.StatusBadRequest, &message{Data: err.Error(), Msg: "request not match", Code: 400})
	g.Abort()
}

func RequestOK(g *gin.Context, data interface{}) {
	g.JSON(http.StatusOK, &message{Data: data, Msg: "request success", Code: 200})
	g.Abort()
}
