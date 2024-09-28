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
	"time"
)

func StartDownload() {

	allNotes := []string{}
	fileutil.ReadByLine("chore/all_notes.txt", func(s string) error {
		if s == "" {
			return nil
		}
		allNotes = append(allNotes, s)
		return nil
	})

	if len(allNotes) <= 0 {
		log.Printf("StartDownload return: all notes is empty")
		return
	}

	for _, n := range allNotes {
		if storage.GetStorage().IsNoteDownloaded(n) {
			continue
		}
		noteInfo, err := xhs.GetNote(n)
		if err != nil {
			log.Errorf("xhs.GetNote err:%v input:%v", err, n)
			continue
		}

		tryRefreshUperInfo(noteInfo.UperUID)

		noteInfo, err = DownloadNote(noteInfo, config.GetConfig().DownloadPath)
		if err != nil {
			log.Errorf("xhs.GetNote err:%v noteInfo:%+v", err, noteInfo)
			continue
		}

		_, err = storage.GetStorage().InsertOrUpdateNote(noteInfo)
		if err != nil {
			log.Errorf("InsertOrUpdateNote err:%v noteInfo:%+v", err, noteInfo)
			continue
		}

	}

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
