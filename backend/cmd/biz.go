package main

import (
	"fmt"
	"github.com/asdine/storm/v3/q"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/xhs_downloader/biz/black"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/biz/xhs"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"path/filepath"
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

	parseResp, err := ParseBlog(n.URL, "")
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
			retryCookie = cookie
		}
		parseResp, _ = ParseBlog(reqURL, retryCookie)
	}

	if len(parseResp.Medias) == 0 {
		log.Infof("*** NO MEDIA ***")
		n.DownloadTime = time.Now()
		n.DownloadNothing = true
		storage.GetStorage().UpdateNote(n)
		return
	}

	log.Infof("start download: %v", n.URL)
	resp := Download(parseResp, "N:/xhs_downloader_output", true)
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

		newNote := storage.GetStorage().GetNote(n.NoteID)
		log.Infof("after update, note:%+v", newNote)

	}

	return
}

func StartFillFileSize() {

	updated := 0
	updatedTotalSize := int64(0)

	log.Printf("StartFillFileSize start")
	round := 0

	limit := 500
	lastID := int64(0)
	for {

		round++

		log.Printf("scan %v updated:%v totalSize:%v", round, updated, utils.GetShowSize(updatedTotalSize))

		ms := []q.Matcher{
			q.Eq("IsDelete", false),
			q.Not(q.Eq("Video", "")),
			q.Gt("ID", lastID),
		}
		resps := []model.Note{}
		err := storage.GetStorage().DB().From("note").Select(ms...).Limit(limit).Find(&resps)
		if err != nil {
			log.Errorf("Find err:%v", err)
			return
		}
		//log.Printf("round%v get %v resps", round, len(resps))

		if len(resps) > 0 {
			lastID = resps[len(resps)-1].ID
		}

		for i, elem := range resps {
			_ = i

			if elem.FileSize > 0 {
				//log.Printf("skip%v:%+v", i+1, utils.GetShowSize(elem.FileSize))
				continue
			}

			//	log.Printf("deal%v:%+v", i+1, elem)

			totalSize := int64(0)
			for _, elem := range elem.Images {
				totalSize += utils.GetFileSize(elem)
			}
			for _, elem := range elem.Lives {
				totalSize += utils.GetFileSize(elem)
			}

			if elem.Video != "" {
				totalSize += utils.GetFileSize(elem.Video)
			}

			elem.FileSize = totalSize

			if elem.FileSize <= 0 {
				continue
			}

			elem.FileSizeReverse = 1024*1024*1024 - elem.FileSize

			//log.Infof("start update note:%v size:%v", elem.Title, utils.GetShowSize(elem.FileSize))
			err = storage.GetStorage().UpdateNote(elem)
			if err != nil {
				log.Errorf("InsertOrUpdateNote err:%v req:%+v", err, elem)
				return
			}

			//log.Infof("update note:%v size:%v", elem.Title, utils.GetShowSize(elem.FileSize))

			updated++
			updatedTotalSize += elem.FileSize

		}

		if len(resps) < limit {
			return
		}

	}

}

func StartDownloadNote() {

	//counter := map[string]int{}

	lastIDKey := fmt.Sprintf("StartDownloadNote-lastID")
	lastID := 0
	storage.GetStorage().DB().Get("common", lastIDKey, &lastID)
	log.Printf("StartDownloadNote-lastID:%v", lastID)

	storage.GetStorage().EachNoteBySelect(lastID, func(n model.Note, currCount, totalCount int) (e error) {

		if n.IsDelete {
			log.Printf("StartDownloadNote-Note is Deleted:%v", n.Title)
			return
		}

		result := DownloadNote(n, false)
		log.Infof("StartDownloadNote EachNoteBySelect(%v/%v):%+v result:%v", currCount, totalCount, n.Title, result)

		if currCount%10 == 0 {
			storage.GetStorage().DB().Set("common", lastIDKey, currCount)
		}

		return
	})

}

func tryRefreshUperInfo(uid string) {
	if uid == "" {
		return
	}
	uper := storage.GetStorage().GetUper(0, uid)
	if uper.ID > 0 {
		return
	}
	uper = xhs.GetUperInfo(uid)
	if uper.UID == "" {
		return
	}
	storage.GetStorage().InsertOrUpdateUper(uper)
}

func DownloadUperAvatar(u model.Uper, to string) (err error) {
	if u.AvatarURL == "" {
		return
	}

	path := filepath.Join(to, "uper_avatar", fmt.Sprintf("%v.jpg", u.UID))

	if utils.HasFile(path) {
		return
	}

	time.Sleep(1 * time.Second)

	code, resp, err := netutil.HttpGetRaw(u.AvatarURL)
	if err != nil {
		log.Printf("EachNoteByRange netutil.HttpGetRaw err:%v url:%v resp:%v", err, u.AvatarURL, resp)
		return
	}

	if code != 200 || len(resp) <= 1024 {
		log.Printf("EachUper netutil.HttpGetRaw failed. code:%v resp:%v", code, resp)
		return fmt.Errorf("invalid status:%v", code)
	}

	err = fileutil.WriteToFile(resp, path)
	if err != nil {
		log.Printf("WriteToFile err:%v resp:%v", err, resp)
		return
	}

	return
}

func CrontabDownloadUperAvatar() {

	storage.GetStorage().EachUper(func(u model.Uper, currCount, totalCount int) (e error) {

		log.Printf("CrontabDownloadUperAvatar progress %v/%v name:%v", currCount, totalCount, u.Name)

		DownloadUperAvatar(u, config.GetDownloadPath())

		return

	})
}

/*
func DownloadNotePoster() {

	failedCount := 0
	storage.GetStorage().EachNoteByRange(0, func(n model.Note, currCount, totalCount int) (e error) {

		log.Printf("DownloadNotePoster progress %v/%v title:%v", currCount, totalCount, n.Title)

		posterPath := filepath.Join(config.GetDownloadPath(), "note_poster", n.UperUID, fmt.Sprintf("%v.jpg", n.NoteID))

		if utils.HasFile(posterPath) {
			return
		}

		code, resp, err := netutil.HttpGetRaw(n.PosterURL)
		if err != nil {
			log.Printf("EachNoteByRange netutil.HttpGetRaw err:%v resp:%v", err, resp)
			return
		}

		if code != 200 || len(resp) <= 1024 {
			failedCount++
			log.Printf("EachNoteByRange netutil.HttpGetRaw failed. code:%v resp:%v", code, resp)
			if failedCount > 3 {
				panic(failedCount)
			}
			return
		}
		failedCount = 0

		err = fileutil.WriteToFile(resp, posterPath)
		if err != nil {
			log.Printf("WriteToFile err:%v resp:%v", err, resp)
			return
		}
		log.Printf("EachNoteByRange WriteToFile succ:%v len(resp):%v", posterPath, utils.GetShowSize(int64(len(resp))))

		time.Sleep(1 * time.Second)

		return

	})
}

*/
