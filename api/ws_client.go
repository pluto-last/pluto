package api

import (
	"github.com/gin-gonic/gin"
	"pluto/model/reply"
)

type WSClientCtl struct{}

// @Tags WSClient
// @Summary 机器人建立websocket连接
// @Produce  application/json
// @Success 200
// @Router /ws/client [get]
func (ws *WSClientCtl) BotWebsocket(c *gin.Context) {
	WSClientService.BotWebsocket(c)
	reply.Ok(c)
}
