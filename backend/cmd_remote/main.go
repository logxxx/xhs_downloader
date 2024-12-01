package main

import (
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils/reqresp"
	"github.com/logxxx/xhs_downloader/biz/blog"
	cookie2 "github.com/logxxx/xhs_downloader/biz/cookie"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "01/02 15:04:05",
	})

	log.Infof("hello world!")

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		time.Local = loc
	}

	g := gin.Default()
	g.GET("/ping", func(c *gin.Context) {
		log.Infof("recv PING msg")
		reqresp.MakeResp(c, "PONG")
	})
	g.GET("/parse_blog", func(c *gin.Context) {
		blogURL := c.Query("blog_url")
		cookie := c.Request.Header.Get("mycookie")
		traceID := c.Query("trace_id")
		log.Infof("[%v]parse_blog reqURL:%v reqCookie:%v", traceID, blogURL, cookie)
		if blogURL == "" {
			reqresp.MakeResp(c, "BLOG URL IS EMPTY")
			return
		}
		resp, err := blog.ParseBlogCore(blogURL, "")
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}
		if len(resp.Medias) > 0 {
			reqresp.MakeResp(c, resp)
			return
		}

		if cookie != "" {
			resp, err = blog.ParseBlogCore(blogURL, cookie)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
			resp.UseCookie = cookie2.GetCookieName(cookie)
		}

		reqresp.MakeResp(c, resp)
		return
	})
	g.Run("0.0.0.0:8088")
}
