package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/black"
	"github.com/logxxx/xhs_downloader/biz/blog"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func ParseNoteAndSaveSourceURL(idx int, n model.Note, cookie string) (result string) {

	if n.IsDownloaded() {
		return "downloaded"
	}

	type SaveInfo struct {
		Idx       int
		UperUID   string
		NoteID    string
		URL       string
		Medias    []blogmodel.Media
		ParseTime string
	}

	parseResp, err := blog.ParseBlog(n.URL, cookie)
	if err != nil {
		log.Errorf("ParseBlog err:%v", err)
		return
	}

	if parseResp.IsNoteDisappeared {
		result = "NoteDisappeared"
		return
	}

	if len(parseResp.Medias) > 0 {
		saveInfo := SaveInfo{
			Idx:       idx,
			UperUID:   n.UperUID,
			NoteID:    n.NoteID,
			URL:       n.URL,
			Medias:    parseResp.Medias,
			ParseTime: time.Now().Format("0102_150405"),
		}

		saveInfoData, _ := json.Marshal(saveInfo)
		fileutil.AppendToFile("parse_result.json", fmt.Sprintf("%v,\n", string(saveInfoData)))
	}

	return fmt.Sprintf("medias=%v", len(parseResp.Medias))
}

func DownloadNote(n model.Note, directlyUseCookie bool, canChangeCookieWhenRetry bool) (result string) {
	logger := log.WithField("func_name", "StartDownloadNote")

	//counter[n.UperUID]++
	//if counter[n.UperUID] > 20 {
	//	//logger.Infof("download enough")
	//	return
	//}

	if !n.DownloadTime.IsZero() {
		//logger.Infof("Downloaded")

		if n.FileSize <= 0 {
			totalSize := int64(0)
			for _, elem := range n.Images {
				totalSize += utils.GetFileSize(elem)
			}
			for _, elem := range n.Lives {
				totalSize += utils.GetFileSize(elem)
			}

			if n.Video != "" {
				totalSize += utils.GetFileSize(n.Video)
			}
			n.FileSize = totalSize

			if n.FileSize > 0 {
				log.Infof("update file size:%v", utils.GetShowSize(n.FileSize))
				storage.GetStorage().InsertOrUpdateNote(n)
			}

		}

		if n.FileSizeReverse <= 0 && n.FileSize > 0 {
			n.FileSizeReverse = 1024*1024*1024 - n.FileSize
			storage.GetStorage().InsertOrUpdateNote(n)
		}

		result = "!n.DownloadTime.IsZero()"
		return
	}

	if n.DownloadNothing {
		result = "n.DownloadNothing"
		return
	}

	if n.URL == "" {
		result = "n.URL empty"
		return
	}

	reason := black.HitBlack(n.Title, n.URL)
	if reason != "" {
		logger.Infof("title HIT BLACK:%v", reason)
		result = "title HIT BLACK"
		return
	}

	reason = black.HitBlack(n.Content, n.URL)
	if reason != "" {
		logger.Infof("content HIT BLACK:%v", reason)
		result = "content HIT BLACK"
		return
	}

	logger = log.WithFields(log.Fields{
		"title":    n.Title,
		"uper_uid": n.UperUID,
	})

	input := ""
	if directlyUseCookie {
		input = cookie.GetCookie3()
	}
	parseResp, err := blog.ParseBlog(n.URL, input)
	if err != nil {
		log.Errorf("ParseBlog err:%v", err)
		return
	}

	if len(parseResp.Medias) == 0 {
		log.Infof("*** DownloadNote find Medias not exists, check AGAIN")
		elems := strings.Split(n.URL, "/")
		reqURL := "https://www.xiaohongshu.com/explore/" + elems[len(elems)-1]
		log.Infof("reqURL:%v", reqURL)
		retryCookie := ""
		if canChangeCookieWhenRetry {
			retryCookie = cookie.GetCookie()
		}
		parseResp, _ = blog.ParseBlog(reqURL, retryCookie)
	}

	if len(parseResp.Medias) == 0 {
		log.Infof("*** DownloadNote NO MEDIA ***")
		n.DownloadTime = time.Now()
		n.DownloadNothing = true
		storage.GetStorage().UpdateNote(n)
		return
	}

	log.Infof("start download: %v", n.URL)
	resp := Download(parseResp, "E:/xhs_downloader_output", true, false)
	//resp := Download(parseResp, "chore/download/notes_by_uper", true)
	log.Infof("download resp:%+v", resp)

	UpdateDownloadRespToDB(model.Uper{}, n, resp)

	return
}

func UpdateDownloadRespToDB(u model.Uper, n model.Note, parseResults []blogmodel.Media) {

	//log.Infof("UpdateDownloadRespToDB uper:%+v note:%+v parseResults(%v):%+v", u, n, len(parseResults), parseResults)
	log.Infof("UpdateDownloadRespToDB start. note:%v", n.Title)
	defer func() {
		log.Infof("UpdateDownloadRespToDB finish")
	}()

	isChanged := false
	for _, m := range parseResults {
		if m.DownloadPath == "" {
			continue
		}
		isChanged = true
		switch m.Type {
		case "image":
			n.ImageURLs = append(n.ImageURLs, m.URL)
			n.Images = append(n.Images, m.DownloadPath)
		case "video":
			n.VideoURL = m.URL
			n.Video = m.DownloadPath
		case "live":
			n.LiveURLs = append(n.LiveURLs, m.URL)
			n.Lives = append(n.Lives, m.DownloadPath)
		}
	}

	if isChanged {
		n.DownloadTime = time.Now()

		//if n.FileSize <= 0 {
		//	totalSize := int64(0)
		//	for _, elem := range n.Images {
		//		totalSize += utils.GetFileSize(elem)
		//	}
		//	for _, elem := range n.Lives {
		//		totalSize += utils.GetFileSize(elem)
		//	}
		//
		//	if n.Video != "" {
		//		totalSize += utils.GetFileSize(n.Video)
		//	}
		//	log.Infof("update file size:%v", utils.GetShowSize(totalSize))
		//	n.FileSize = totalSize
		//}

		_, err := storage.GetStorage().InsertOrUpdateNote(n)
		if err != nil {
			log.Errorf("InsertOrUpdateNote err:%v n:%+v", err, n)
		}

		//newNote := storage.GetStorage().GetNote(n.NoteID)
		//log.Infof("after update, note:%+v", newNote)

	} else {
		log.Infof("no change, no need to update")
	}

	dbU := storage.GetStorage().GetUper(0, u.UID)
	if dbU.ID > 0 {
		u.GalleryEmptyLastTime = dbU.GalleryEmptyLastTime
		u.Notes = dbU.Notes
	}

	u.AddNote(n.NoteID)

	if u.UID != "" {
		result, err := storage.GetStorage().InsertOrUpdateUper(u)
		log.Infof("InsertOrUpdateUper input:%+v result:%v err:%v", u, result, err)
	}

}

func DownloadNoteByID(note string) (err error) {
	dbNote := storage.GetStorage().GetNote(note)
	if dbNote.IsDownloaded() {
		return errors.New("downloaded")
	}

	blog, err := blog.ParseBlog(note, cookie.GetCookie3())
	if err != nil {
		log.Errorf("StartScanMyShoucang ParseBlog err:%v note:%v", err, note)
		return
	}

	dbNote = model.Note{
		NoteID:    blog.NoteID,
		URL:       note,
		PosterURL: "",
		UperUID:   blog.UserID,
		Title:     blog.Title,
		Content:   blog.Content,
	}
	for _, m := range blog.Medias {
		if m.Type == "video" {
			dbNote.Video = m.URL
		}
		if m.Type == "live" {
			dbNote.Lives = append(dbNote.Lives, m.URL)
		}
		if m.Type == "image" {
			dbNote.Images = append(dbNote.Images, m.URL)
		}
	}

	DownloadNote(dbNote, false, true)

	return

}
