package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils/reqresp"
	"github.com/logxxx/xhs_downloader/biz/blog"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	cookie2 "github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/queue"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"time"
)

var ()

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
		log.Infof("[%v]parse_blog reqURL:%v reqCookie:%v", traceID, blogURL, cookie2.GetCookieName(cookie))
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
	g.POST("/send_work", func(c *gin.Context) {
		work := &model.Work{}
		log.Infof("send work:%+v", work)

		err = reqresp.ParseReq(c, work)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		if work.NoteID == "" {
			reqresp.MakeErrMsg(c, errors.New("empty note id"))
			return
		}

		err = queue.Push("work", work)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		reqresp.MakeRespOk(c)

	})
	g.GET("/recv_work", func(c *gin.Context) {
		work := &model.Work{}
		err = queue.Pop("work", work, false)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}
		log.Infof("recv work:%+v", work)

		reqresp.MakeResp(c, work)
	})
	g.POST("/send_work_result", func(c *gin.Context) {
		result := &blogmodel.ParseBlogResp{}
		err = reqresp.ParseReq(c, result)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		if result.NoteID == "" {
			reqresp.MakeErrMsg(c, errors.New("empty note id"))
			return
		}

		err = queue.Push("work_result", result)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		reqresp.MakeRespOk(c)
	})
	g.GET("/recv_work_result", func(c *gin.Context) {
		result := &blogmodel.ParseBlogResp{}
		err = queue.Pop("work_result", result, false)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		reqresp.MakeResp(c, result)
	})

	g.Run("0.0.0.0:8088")
}
