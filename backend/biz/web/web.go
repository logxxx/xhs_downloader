package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/xhs_downloader/config"
	"time"
)

func InitWeb() {
	g := gin.Default()
	g.GET("/now", func(c *gin.Context) {
		c.String(200, time.Now().Format("2006-01-02 15:04:05"))
	})
	g.Run(fmt.Sprintf(":%v", config.GetConfig().Port))
}
