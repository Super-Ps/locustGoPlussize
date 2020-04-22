package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func (s *ServerEntry) WebPage(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	ctx.HTML(http.StatusOK, "index.html", gin.H{"sys": template.HTML(GetSysMonitorTemplate()),
														"pro": template.HTML(GetProMonitorTemplate()),
														"info": template.HTML(GetOsInfoTemplate()),
														"startRoute": template.HTML(fmt.Sprintf("%s/%s", Param.Route.Monitor, "list")),
														"restartRoute": template.HTML(Param.Route.Restart),
														"stopRoute": template.HTML(Param.Route.Stop),
														"exitRoute": template.HTML(Param.Route.Quit),
														"staticRoute": template.HTML(Param.Route.Static)})
}