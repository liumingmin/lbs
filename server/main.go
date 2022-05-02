package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/middleware"
)

func main() {
	InitOps()

	g := gin.Default()

	defaultResp := &middleware.DefaultServiceResponse{}
	g.POST("/query", middleware.ServiceHandler(QueryByLoc, GeoLocation{}, defaultResp))
	g.Run(":12800")
}
