package web

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/reqresp"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/biz/thumb"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	"github.com/logxxx/xhs_downloader/proto"
	log "github.com/sirupsen/logrus"
	"os"
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

func GetDistDir() string {
	return "D:\\mytest\\mywork\\xhs_downloader\\frontend\\dist"
}

func InitWeb() {

	g := gin.Default()

	g.Use(reqresp.Cors())
	g.Use(gzip.Gzip(gzip.DefaultCompression))

	g.StaticFile("/", GetDistDir())
	g.StaticFS("/dist", gin.Dir(GetDistDir(), true))

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

	g.GET("/one_video", func(c *gin.Context) {
		token := c.Query("token")
		dbNote, nextToken, err := storage.GetStorage().GetOneVideoNoteBySize2(token)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		thumbPath := filepath.Join(filepath.Dir(dbNote.Video), ".thumb", filepath.Base(dbNote.Video))
		if !utils.HasFile(thumbPath) {
			log.Printf("make thumb for :%v ", thumbPath)
			thumb.MakeThumb(dbNote.Video)
			log.Printf("make thumb finish :%v ", thumbPath)
		}

		apiResp := proto.ApiGetOneNoteResp{
			Data: proto.ApiUperNote{
				UperUID: dbNote.UperUID,
				NoteID:  dbNote.NoteID,
				//Poster:   utils.B64(GetNotePosterPath(dbNote.UperUID, dbNote.NoteID)),
				Title:    dbNote.Title,
				Content:  dbNote.Content,
				Video:    utils.B64(dbNote.Video),
				Images:   ArrayB64(dbNote.Images),
				Lives:    ArrayB64(dbNote.Lives),
				ShowSize: utils.GetShowSize(dbNote.FileSize),
			},
			Token: nextToken,
		}

		c.JSON(200, apiResp)
	})

	g.GET("/one_note", func(c *gin.Context) {
		token := c.Query("token")
		t := c.Query("type")
		dbNote, nextToken, err := storage.GetStorage().GetOneNote(token, t)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		apiResp := proto.ApiGetOneNoteResp{
			Data: proto.ApiUperNote{
				UperUID:  dbNote.UperUID,
				NoteID:   dbNote.NoteID,
				Poster:   utils.B64(GetNotePosterPath(dbNote.UperUID, dbNote.NoteID)),
				Title:    dbNote.Title,
				Content:  dbNote.Content,
				Video:    utils.B64(dbNote.Video),
				Images:   ArrayB64(dbNote.Images),
				Lives:    ArrayB64(dbNote.Lives),
				ShowSize: utils.GetShowSize(dbNote.FileSize),
			},
			Token: nextToken,
		}

		c.JSON(200, apiResp)
	})

	g.GET("/each_note", func(c *gin.Context) {
		token := c.Query("token")
		resp, nextToken, err := storage.GetStorage().GetNotesByPage(1, token)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}
		apiResp := &GetNotesResp{
			Data:  resp,
			Token: nextToken,
		}
		c.JSON(200, apiResp)
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

		resp := proto.ApiUperInfo{
			UID:    u.UID,
			Name:   u.Name,
			Desc:   u.Desc,
			Tags:   u.Tags,
			Notes:  u.Notes,
			Avatar: utils.B64(GetUperAvatarPath(u.UID)),
		}

		with := c.Query("with")
		if with == "withoutNotes" {
			resp.Notes = nil
		}

		reqresp.MakeResp(c, resp)
	})

	g.GET("/debug/uper", func(c *gin.Context) {
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

		thumbPath := filepath.Join(filepath.Dir(filePath), ".thumb", filepath.Base(filePath))
		if utils.HasFile(thumbPath) {
			log.Printf("get file use thumb:%v", thumbPath)
			filePath = thumbPath
		} else {
			if utils.GetFileSize(filePath) > 1*1024*1024 {
				log.Printf("make thumb online:%v to:%v", filePath, thumbPath)
				thumb.GeneVideoShot(filePath, thumbPath)
				filePath = thumbPath
			}
		}

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
				NoteID:    dbNote.NoteID,
				Poster:    utils.B64(GetNotePosterPath(dbNote.UperUID, dbNote.NoteID)),
				Title:     dbNote.Title,
				Content:   dbNote.Content,
				Video:     "",  //todo
				Images:    nil, //todo
				Lives:     nil, //todo
				Tags:      dbNote.Tags,
				IsDeleted: dbNote.IsDelete,
			}
			resp.Data = append(resp.Data, note)
		}

		reqresp.MakeResp(c, resp)

	})

	g.GET("/uper/delete", func(c *gin.Context) {
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

		dbUper.IsDelete = true
		storage.GetStorage().InsertOrUpdateUper(dbUper)

		for _, note := range dbUper.Notes {
			dbNote := storage.GetStorage().GetNote(note)
			DeleteNote(dbNote)
		}

		reqresp.MakeRespOk(c)

	})

	g.GET("/note/delete", func(c *gin.Context) {
		noteID := c.Query("note_id")

		dbNote := storage.GetStorage().GetNote(noteID)
		if dbNote.ID <= 0 {
			reqresp.MakeErrMsg(c, errors.New("note not found"))
			return
		}

		DeleteNote(dbNote)

		reqresp.MakeRespOk(c)

	})

	g.GET("/note/add_tag", func(c *gin.Context) {
		noteID := c.Query("note_id")
		if noteID == "" {
			reqresp.MakeErrMsg(c, errors.New("empty noteID"))
			return
		}
		tag := c.Query("tag")
		if tag == "" {
			reqresp.MakeErrMsg(c, errors.New("empty tag"))
			return
		}

		dbNote := storage.GetStorage().GetNote(noteID)
		if dbNote.ID <= 0 {
			reqresp.MakeErrMsg(c, errors.New("note not found"))
			return
		}

		if dbNote.HasTag(tag) {
			reqresp.MakeRespOk(c)
			return
		}

		dbNote.Tags = append(dbNote.Tags, tag)

		err := storage.GetStorage().UpdateNote(dbNote)
		if err != nil {
			reqresp.MakeErrMsg(c, err)
			return
		}

		reqresp.MakeRespOk(c)

	})

	g.GET("/uper/add_tag", func(c *gin.Context) {
		uid := c.Query("uid")
		if uid == "" {
			reqresp.MakeErrMsg(c, errors.New("empty uid"))
			return
		}
		tag := c.Query("tag")
		if tag == "" {
			reqresp.MakeErrMsg(c, errors.New("empty tag"))
			return
		}

		dbUper := storage.GetStorage().GetUper(0, uid)
		if dbUper.ID <= 0 {
			reqresp.MakeErrMsg(c, errors.New("uper not found"))
			return
		}

		if !dbUper.HasTag(tag) {
			dbUper.Tags = append(dbUper.Tags, tag)

			_, err := storage.GetStorage().InsertOrUpdateUper(dbUper)
			if err != nil {
				reqresp.MakeErrMsg(c, err)
				return
			}
		}

		for _, note := range dbUper.Notes {
			dbNote := storage.GetStorage().GetNote(note)
			if dbNote.ID <= 0 {
				continue
			}
			if dbNote.HasTag(tag) {
				continue
			}
			dbNote.Tags = append(dbNote.Tags, tag)
			storage.GetStorage().UpdateNote(dbNote)
		}

		reqresp.MakeRespOk(c)

	})

	g.GET("/upers", func(c *gin.Context) {
		token := c.Query("token")
		limitStr := c.Query("limit")
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 10
		}

		with := c.Query("with")
		opt := storage.GetUpersOpt{}
		if with == "withNoTag" {
			opt.WithNoTag = true
		}
		dbUpers, nextToken := storage.GetStorage().GetUpers(c, opt, limit, token)
		resp := &proto.ApiGetUpersResp{
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

	g.GET("/debug/note", func(c *gin.Context) {
		noteID := c.Query("note_id")

		dbNote := storage.GetStorage().GetNote(noteID)
		if dbNote.ID <= 0 {
			reqresp.MakeErrMsg(c, errors.New("note not found"))
			return
		}

		reqresp.MakeResp(c, dbNote)
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
		port = 6080
	}

	g.Run(fmt.Sprintf(":%v", port))
}

func GetUperAvatarPath(uid string) string {
	return filepath.Join(config.GetDownloadPath(), "uper_avatar", fmt.Sprintf("%v.jpg", uid))
}

func GetNotePosterPath(uid, noteID string) string {
	return filepath.Join(config.GetDownloadPath(), "note_poster", uid, fmt.Sprintf("%v.jpg", noteID))
}

func ArrayB64(input []string) (output []string) {
	for _, elem := range input {
		output = append(output, utils.B64(elem))
	}
	return
}

func DeleteNote(dbNote model.Note) {

	go func() {
		time.Sleep(10 * time.Second)
		for _, elem := range dbNote.Images {
			if elem != "" && utils.HasFile(elem) {
				log.Printf("remove Image:%v", elem)
				os.Remove(elem)
				thumbFilePath := filepath.Join(filepath.Dir(elem), ".thumb", filepath.Base(elem))
				if utils.HasFile(thumbFilePath) {
					os.Remove(thumbFilePath)
				}
			}
		}

		for _, elem := range dbNote.Lives {
			if elem != "" && utils.HasFile(elem) {
				log.Printf("remove Live:%v", elem)
				os.Remove(elem)
				thumbFilePath := filepath.Join(filepath.Dir(elem), ".thumb", filepath.Base(elem))
				if utils.HasFile(thumbFilePath) {
					os.Remove(thumbFilePath)
				}
			}
		}

		if dbNote.Video != "" && utils.HasFile(dbNote.Video) {
			log.Printf("remove Video:%v", dbNote.Video)
			os.Remove(dbNote.Video)
			thumbFilePath := filepath.Join(filepath.Dir(dbNote.Video), ".thumb", filepath.Base(dbNote.Video))
			if utils.HasFile(thumbFilePath) {
				os.Remove(thumbFilePath)
			}
		}

		dbNote.IsDelete = true

		storage.GetStorage().UpdateNote(dbNote)
	}()
}
