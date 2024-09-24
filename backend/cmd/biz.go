package main

import (
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/biz/xhs"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
)

type DownloadNoteResp struct {
}

func StartDownload() {

	allNotes := []string{}
	fileutil.ReadByLine("chore/all_notes.txt", func(s string) error {
		if s == "" {
			return nil
		}
		allNotes = append(allNotes, s)
		return nil
	})

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

		err = storage.GetStorage().InsertNote(noteInfo)
		if err != nil {
			log.Errorf("InsertNote err:%v noteInfo:%+v", err, noteInfo)
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
	storage.GetStorage().InsertUper(uper)
}
