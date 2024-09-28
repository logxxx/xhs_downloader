package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/config"
	"time"
)

func InitWeb() {
	g := gin.Default()
	g.GET("/now", func(c *gin.Context) {
		c.String(200, time.Now().Format("2006-01-02 15:04:05"))
	})
	g.GET("/all_note", func(c *gin.Context) {
		resp := storage.GetStorage().GetAllNotes(c)
		c.JSON(200, resp)
	})
	g.GET("/all_uper", func(c *gin.Context) {
		resp := storage.GetStorage().GetAllNotes(c)
		c.JSON(200, resp)
	})
	port := config.GetConfig().Port
	if port <= 0 {
		port = 8080
	}
	g.Run(fmt.Sprintf(":%v", port))
}
