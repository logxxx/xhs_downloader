package main

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/biz/xhs"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

func StartDownloadNote() {

	counter := map[string]int{}

	storage.GetStorage().EachNote(func(n model.Note, currCount, totalCount int) (e error) {

		logger := log.WithField("func_name", "StartDownloadNote")

		counter[n.UperUID]++
		if counter[n.UperUID] > 20 {
			//logger.Infof("download enough")
			return
		}

		if !n.DownloadTime.IsZero() {
			//logger.Infof("Downloaded")
			return
		}
		if n.URL == "" {
			logger.Infof("URL empty")
			return
		}

		time.Sleep(1 * time.Second)

		logger.Infof("download start")

		uper := storage.GetStorage().GetUper(0, n.UperUID)

		logger = log.WithFields(log.Fields{
			"currCount":  currCount,
			"totalCount": totalCount,
			"title":      n.Title,
			"uper":       uper.Name,
		})

		if !strings.HasPrefix(n.URL, "https:") {
			n.URL = "https://www.xiaohongshu.com" + n.URL
		}

		parseResp, err := ParseBlog(n.URL)
		if err != nil {
			log.Errorf("ParseBlog err:%v", err)
			return
		}
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
			err = storage.GetStorage().UpdateNote(n)
			if err != nil {
				log.Errorf("UpdateNote err:%v n:%+v", err, n)
			}

			newNote := storage.GetStorage().GetNote(n.NoteID)
			log.Infof("after update, note:%+v", newNote)

		}

		return
	})

}

func DownloadNote(n model.Note, downloadPath string) (resp model.Note, err error) {
	return
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
		log.Printf("EachNote netutil.HttpGetRaw err:%v url:%v resp:%v", err, u.AvatarURL, resp)
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

func DownloadNotePoster() {

	failedCount := 0
	storage.GetStorage().EachNote(func(n model.Note, currCount, totalCount int) (e error) {

		log.Printf("DownloadNotePoster progress %v/%v title:%v", currCount, totalCount, n.Title)

		posterPath := filepath.Join(config.GetDownloadPath(), "note_poster", n.UperUID, fmt.Sprintf("%v.jpg", n.NoteID))

		if utils.HasFile(posterPath) {
			return
		}

		code, resp, err := netutil.HttpGetRaw(n.PosterURL)
		if err != nil {
			log.Printf("EachNote netutil.HttpGetRaw err:%v resp:%v", err, resp)
			return
		}

		if code != 200 || len(resp) <= 1024 {
			failedCount++
			log.Printf("EachNote netutil.HttpGetRaw failed. code:%v resp:%v", code, resp)
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
		log.Printf("EachNote WriteToFile succ:%v len(resp):%v", posterPath, utils.GetShowSize(int64(len(resp))))

		time.Sleep(1 * time.Second)

		return

	})
}
