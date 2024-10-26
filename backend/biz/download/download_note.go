package download

import (
	"errors"
	"github.com/logxxx/utils"
	"github.com/logxxx/xhs_downloader/biz/black"
	"github.com/logxxx/xhs_downloader/biz/blog"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func DownloadNote(n model.Note, canChangeCookieWhenRetry bool) (result string) {
	logger := log.WithField("func_name", "StartDownloadNote")

	//counter[n.UperUID]++
	//if counter[n.UperUID] > 20 {
	//	//logger.Infof("download enough")
	//	return
	//}

	if !strings.HasPrefix(n.URL, "https:") {
		n.URL = "https://www.xiaohongshu.com" + n.URL
	}

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

	parseResp, err := blog.ParseBlog(n.URL, "")
	if err != nil {
		log.Errorf("ParseBlog err:%v", err)
		return
	}

	if len(parseResp.Medias) == 0 {
		log.Infof("*** find Medias not exists, check AGAIN!!!")
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
		log.Infof("*** NO MEDIA ***")
		n.DownloadTime = time.Now()
		n.DownloadNothing = true
		storage.GetStorage().UpdateNote(n)
		return
	}

	log.Infof("start download: %v", n.URL)
	resp := Download(parseResp, "E:/xhs_downloader_output", true)
	//resp := Download(parseResp, "chore/download/notes_by_uper", true)

	isChanged := false
	for _, m := range resp {
		if m.DownloadPath == "" {
			continue
		}
		isChanged = true
		switch m.Type {
		case "image":
			n.Images = append(n.Images, m.DownloadPath)
		case "video":
			n.Video = m.DownloadPath
		case "live":
			n.Lives = append(n.Lives, m.DownloadPath)
		}
	}

	if isChanged {
		n.DownloadTime = time.Now()

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
			log.Infof("update file size:%v", utils.GetShowSize(totalSize))
			n.FileSize = totalSize
		}

		_, err = storage.GetStorage().InsertOrUpdateNote(n)
		if err != nil {
			log.Errorf("InsertOrUpdateNote err:%v n:%+v", err, n)
		}

		//newNote := storage.GetStorage().GetNote(n.NoteID)
		//log.Infof("after update, note:%+v", newNote)

	}

	return
}

func DownloadNoteByID(note string) (err error) {
	dbNote := storage.GetStorage().GetNote(note)
	if dbNote.IsDownloaded() {
		return errors.New("downloaded")
	}

	blog, err := blog.ParseBlog(note, cookie.GetCookie())
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

	DownloadNote(dbNote, true)

	return

}
