package main

import (
	"errors"
	"fmt"
	"github.com/asdine/storm/v3/q"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	cookie2 "github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/download"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/biz/xhs"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

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
	//storage.GetStorage().DB().Get("common", lastIDKey, &lastID)
	log.Printf("StartDownloadNote-lastID:%v", lastID)

	continueMediasZeroCount := 0
	cookie := cookie2.GetCookie3()
	storage.GetStorage().EachNoteBySelect(50200, func(n model.Note, currCount, totalCount int) (e error) {

		if n.IsDelete {
			log.Printf("StartDownloadNote-Note is Deleted:%v", n.Title)
			return
		}

		//result := download.DownloadNote(n, false, false)
		result := download.ParseNoteAndSaveSourceURL(currCount, n, cookie)
		log.Infof("StartDownloadNote EachNoteBySelect(%v/%v):title:%+v url:%v result:%v", currCount, totalCount, n.Title, n.URL, result)

		if currCount%10 == 0 {
			//if cookie == cookie2.GetCookie1() {
			//	cookie = cookie2.GetCookie2()
			//} else {
			//	cookie = cookie2.GetCookie1()
			//}
			storage.GetStorage().DB().Set("common", lastIDKey, currCount)
		}

		if result != "downloaded" && result != "NoteDisappeared" {
			time.Sleep(10 * time.Second)
			if currCount%100 == 0 {
				time.Sleep(10 * time.Minute)
			}
		}

		if strings.HasPrefix(result, "medias=") {
			if result == "medias=0" {
				continueMediasZeroCount++
			} else {
				continueMediasZeroCount = 0
			}
		}

		if continueMediasZeroCount > 10 {
			log.Printf("*************** continueMediasZeroCount TOO MANY:%v", continueMediasZeroCount)
			return errors.New("continueMediasZeroCount too many")
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
