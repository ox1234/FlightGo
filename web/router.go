package web

import (
	"github.com/gin-gonic/gin"
	"pentestplatform/web/controller"
)

func SetRouter(engine *gin.Engine){
	// 设置信息收集的路由
	gather := engine.Group("/gather")
	{
		gather.GET("/subdomain", controller.SubDomain)
		gather.GET("/port", controller.PortScan)
		gather.GET("/dir", controller.DirScan)
		gather.POST("/dir", controller.DirScan)
		gather.GET("/basic", controller.BasicScan)
		gather.POST("/start", controller.Start)
		gather.GET("/vt", controller.VtDomain)
		gather.GET("/rapiddns", controller.RapidDnsDomain)
		gather.GET("/alldomain", controller.AllDomain)
		gather.GET("/dump", controller.DumpData)
	}

	attack := engine.Group("/attack")
	{
		attack.GET("/show", controller.ShowPayload)
		attack.POST("/attack", controller.DoExploit)
	}
}
