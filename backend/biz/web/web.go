package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/reqresp"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	"github.com/logxxx/xhs_downloader/proto"
	"path/filepath"
	"strconv"
	"time"
)

type GetUpersResp struct {
	Data  []model.Uper `json:"data"`
	Token string       `json:"token"`
}

type GetNotesResp struct {
	Data  []model.Note `json:"data"`
	Token string       `json:"token"`
}

func InitWeb() {

	g := gin.Default()

	g.Use(reqresp.Cors())

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
		uid := c.Query("uid")
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 10
		}
		notes, nextToken := storage.GetStorage().GetNotes(c, storage.GetUpersOpt{Uid: uid}, limit, token)
		resp := &GetNotesResp{
			Data:  notes,
			Token: nextToken,
		}
		c.JSON(200, resp)
	})

	g.GET("/file", func(c *gin.Context) {

		id := c.Query("id")
		//log.Infof("get file:%v", id)
		isPreview := c.Query("is_preview")
		_ = isPreview

		if id == "" {
			reqresp.MakeErrMsg(c, errors.New("empty id"))
			return
		}

		filePath := utils.B64To(id)

		c.File(filePath)

	})

	g.GET("/uper_notes", func(c *gin.Context) {
		uid := c.Query("uid")
		if uid == "" {
			reqresp.MakeErrMsg(c, errors.New("empty uid"))
			return
		}

		dbUper := storage.GetStorage().GetUper(0, uid)
		if dbUper.ID <= 0 {
			reqresp.MakeErrMsg(c, errors.New("uper not found"))
			return
		}

		resp := &proto.ApiGetUperNotesResp{}

		for i, n := range dbUper.Notes {

			if i > 10 { //TODO: page
				break
			}

			dbNote := storage.GetStorage().GetNote(n)
			if dbNote.NoteID == "" {
				continue
			}
			note := proto.ApiUperNote{
				NoteID:  dbNote.NoteID,
				Poster:  utils.B64(GetNotePosterPath(dbNote.UperUID, dbNote.NoteID)),
				Title:   dbNote.Title,
				Content: dbNote.Content,
				Video:   "",  //todo
				Images:  nil, //todo
				Lives:   nil, //todo
			}
			resp.Data = append(resp.Data, note)
		}

		reqresp.MakeResp(c, resp)

	})

	g.GET("/upers", func(c *gin.Context) {
		token := c.Query("token")
		limitStr := c.Query("limit")
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 10
		}
		dbUpers, nextToken := storage.GetStorage().GetUpers(c, storage.GetUpersOpt{}, limit, token)
		resp := &proto.ApiGetUperInfoResp{
			Token: nextToken,
		}

		for _, db := range dbUpers {
			uper := proto.ApiUperInfo{
				UID:    db.UID,
				Name:   db.Name,
				Desc:   db.Desc,
				Tags:   db.Tags,
				Avatar: utils.B64(GetUperAvatarPath(db.UID)),
			}
			resp.Data = append(resp.Data, uper)
		}

		c.JSON(200, resp)
	})

	g.GET("/debug/upers", func(c *gin.Context) {
		token := c.Query("token")
		limitStr := c.Query("limit")
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 10
		}
		upers, nextToken := storage.GetStorage().GetUpers(c, storage.GetUpersOpt{}, limit, token)
		resp := &GetUpersResp{
			Data:  upers,
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

func GetUperAvatarPath(uid string) string {
	return filepath.Join(config.GetDownloadPath(), "uper_avatar", fmt.Sprintf("%v.jpg", uid))
}

func GetNotePosterPath(uid, noteID string) string {
	return filepath.Join(config.GetDownloadPath(), "note_poster", uid, fmt.Sprintf("%v.jpg", noteID))
}
