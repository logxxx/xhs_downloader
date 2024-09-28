package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils/reqresp"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	"strconv"
	"time"
)

type GetNotesResp struct {
	Data  []model.Note `json:"data"`
	Token string       `json:"token"`
}

func InitWeb() {
	g := gin.Default()
	g.GET("/now", func(c *gin.Context) {
		c.String(200, time.Now().Format("2006-01-02 15:04:05"))
	})

	g.GET("/all_note_count", func(c *gin.Context) {
		resp := storage.GetStorage().GetNoteTotalCount()
		c.JSON(200, resp)
	})
	g.GET("/all_uper_count", func(c *gin.Context) {
		resp := storage.GetStorage().GetUperTotalCount()
		c.JSON(200, resp)
	})
	g.GET("/all_note", func(c *gin.Context) {
		resp := storage.GetStorage().GetAllNotes(c)
		c.JSON(200, resp)
	})
	g.GET("/all_uper", func(c *gin.Context) {
		resp := storage.GetStorage().GetAllNotes(c)
		c.JSON(200, resp)
	})
	g.GET("/uper", func(c *gin.Context) {
		uid := c.Query("uid")
		if uid == "" {
			reqresp.MakeErrMsg(c, errors.New("empty uid"))
			return
		}
		u := storage.GetStorage().GetUper(0, uid)
		if u.ID <= 0 {
			reqresp.MakeErrMsg(c, fmt.Errorf("user not found:%v", uid))
			return
		}
		reqresp.MakeResp(c, u)
	})
	g.GET("/update_uper", func(c *gin.Context) {

		uid := c.Query("uid")
		if uid == "" {
			reqresp.MakeErrMsg(c, errors.New("empty uid"))
			return
		}

		action := c.Query("action")
		if action == "" {
			reqresp.MakeErrMsg(c, errors.New("action"))
			return
		}

		u := storage.GetStorage().GetUper(0, uid)
		if u.ID <= 0 {
			reqresp.MakeErrMsg(c, fmt.Errorf("user not found:%v", uid))
			return
		}

		switch action {
		case "like":
			u.IsLike = true
		case "cancel_like":
			u.IsLike = false
		case "delete":
			u.IsDelete = true
		case "cancel_delete":
			u.IsDelete = false
		default:
			reqresp.MakeErrMsg(c, fmt.Errorf("unknown action:%v", action))
			return
		}

		_, err := storage.GetStorage().InsertOrUpdateUper(u)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}
		reqresp.MakeRespOk(c)
	})

	g.GET("/notes", func(c *gin.Context) {
		token := c.Query("token")
		limitStr := c.Query("limit")
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 10
		}
		notes, nextToken := storage.GetStorage().GetNotes(c, storage.GetUpersOpt{}, limit, token)
		resp := &GetNotesResp{
			Data:  notes,
			Token: nextToken,
		}
		c.JSON(200, resp)
	})
	port := config.GetConfig().Port
	if port <= 0 {
		port = 8080
	}
	g.Run(fmt.Sprintf(":%v", port))
}
