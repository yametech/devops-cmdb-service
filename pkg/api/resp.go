package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type message struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func RequestErr(g *gin.Context, err error) {
	g.JSON(http.StatusBadRequest, &message{Data: err.Error(), Msg: "request not match"})
	g.Abort()
}

