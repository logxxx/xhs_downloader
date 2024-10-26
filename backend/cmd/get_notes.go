package main

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/download"
	"github.com/logxxx/xhs_downloader/biz/mydp"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/config"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

func StartDownloadRecrentlyNotes() {
	cookie.ChangeCookie()

	lastIDKey := fmt.Sprintf("StartDownloadRecrentlyNotes-last_id")
	lastID := 0
	storage.GetStorage().DB().Get("common", lastIDKey, &lastID)
	log.Printf("lastID:%v", lastID)

	for {

		storage.GetStorage().EachUper(func(u model.Uper, currCount, totalCount int) (e error) {

			if currCount < lastID {
				return
			}

			storage.GetStorage().DB().Set("common", lastIDKey, currCount)

			log.Infof("StartDownloadRecrentlyNotes EachUper:%v %v/%v", u.Name, currCount, totalCount)

			if u.IsDelete {
				log.Printf("uper deleted")
				return
			}

			if currCount%10 == 0 {
				cookie.ChangeCookie()
			}

			_, parseNotes, err := mydp.GetNotes(u.UID, cookie.GetCookie(), true)
			if err != nil {
				return
			}

			uperChanged := false
			oldNoteCount := len(u.Notes)
			for i, n := range parseNotes {

				ok := u.AddNote(n.NoteID)
				if ok {
					uperChanged = true
				} else {
					log.Printf("AddNote already has:%v", n.Title)
				}

				dbNote := model.Note{
					NoteID:    n.NoteID,
					UperUID:   u.UID,
					Title:     n.Title,
					URL:       n.URL,
					PosterURL: n.Poster,
					LikeCount: n.LikeCount,
				}
				if ok {
					insertOrUpdate, err := storage.GetStorage().InsertOrUpdateNote(dbNote)
					if err != nil {
						log.Printf("InsertOrUpdateNote err:%v dbNote:%+v", err, dbNote)
						continue
					}
					_ = insertOrUpdate
					log.Printf("InsertOrUpdateNote succ(%v/%v): %+v(%v)", i+1, len(parseNotes), dbNote.Title, insertOrUpdate)
				}

				download.DownloadNote(dbNote, true)
			}

			if uperChanged {
				log.Infof("StartDownloadRecrentlyNotes update uper:%+v => %v", oldNoteCount, len(u.Notes))
				storage.GetStorage().InsertOrUpdateUper(u)
			}

			return
		})
	}
}

func DownloadUperFirstPageNotes(uid string) {

	uper := storage.GetStorage().GetUper(0, uid)

	parseUper, parseNotes, err := mydp.GetNotes(uid, cookie.GetCookie(), true)
	if err != nil {
		log.Printf("get parseNotes err:%v uid:%v", err, uid)
		if strings.Contains(err.Error(), "change account") {
			cookie.ChangeCookie()
		}
		return
	}
	log.Printf("parseUper [%v_%v] get [%v] parseNotes", parseUper.UID, parseUper.Name, len(parseNotes))

	allParseNotes := []string{}
	for _, n := range parseNotes {
		allParseNotes = append(allParseNotes, n.NoteID)
	}
	allParseNotes = utils.RemoveEmpty(utils.RemoveDuplicate(allParseNotes))
	uper = model.Uper{
		ID:               uper.ID,
		UID:              parseUper.UID,
		Name:             parseUper.Name,
		Area:             parseUper.Area,
		AvatarURL:        parseUper.AvatarURL,
		IsGirl:           parseUper.IsGirl,
		Desc:             parseUper.Desc,
		Tags:             parseUper.Tags,
		FansCount:        parseUper.FansCount,
		ReceiveLikeCount: parseUper.ReceiveLikeCount,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
		Notes:            allParseNotes,
	}
	_, err = storage.GetStorage().InsertOrUpdateUper(uper)
	if err != nil {
		log.Printf("InsertOrUpdateUper err:%v uper:%+v", err, uper)
		return
	}
	//log.Printf("InsertOrUpdateUper succ:%+v result:%v", uper, result)

	for i, n := range parseNotes {

		dbNote := model.Note{
			NoteID:    n.NoteID,
			UperUID:   uid,
			Title:     n.Title,
			URL:       n.URL,
			PosterURL: n.Poster,
			LikeCount: n.LikeCount,
		}
		insertOrUpdate, err := storage.GetStorage().InsertOrUpdateNote(dbNote)
		if err != nil {
			log.Printf("InsertOrUpdateNote err:%v dbNote:%+v", err, dbNote)
			continue
		}
		_ = insertOrUpdate
		log.Printf("InsertOrUpdateNote succ(%v/%v): %+v(%v)", i+1, len(parseNotes), dbNote.Title, insertOrUpdate)

		download.DownloadNote(dbNote, true)

		DownloadPoster(dbNote)
	}

	DownloadUperAvatar(uper, config.GetDownloadPath())
}

func StartGetNotes() {

	upers := getAllUpers()
	log.Printf("get %v upers", len(upers))

	continueNoNoteCount := 0
	downloadedCount := 0
	//reachLast := false
	for i, u := range upers {

		//if u == "65fb9db5000000000b00ec0b" {
		//	reachLast = true
		//}
		//
		//if !reachLast {
		//	continue
		//}

		if downloadedCount > 500 && i > 0 && i%50 == 0 {
			log.Printf("sleep for i%%10==0")
			time.Sleep(1 * time.Minute)
		}

		if downloadedCount > 500 && i > 0 && i%100 == 0 {
			log.Printf("change cookie for i%%100==0")
			cookie.ChangeCookie()
		}

		log.Printf("deal parseUper %v/%v %v", i+1, len(upers), u)

		uper := storage.GetStorage().GetUper(0, u)
		log.Printf("IS_NEW:%v NOTES:%v CTIME:%v", uper.ID == 0, len(uper.Notes), uper.CreateTime.Format("01/02 15:04:05"))

		if len(uper.Notes) != 0 && len(uper.Notes) != 14 {
			continue
		}
		//if storage.GetStorage().IsUperScanned(u) {
		//	continue
		//}
		parseUper, parseNotes, err := mydp.GetNotes(u, cookie.GetCookie(), false)
		if err != nil {
			log.Printf("get parseNotes err:%v uid:%v", err, u)
			if strings.Contains(err.Error(), "change account") {
				cookie.ChangeCookie()
			}
			continue
		}
		log.Printf("parseUper [%v_%v] get [%v] parseNotes", parseUper.UID, parseUper.Name, len(parseNotes))

		if len(parseNotes) <= len(utils.RemoveDuplicate(uper.Notes)) {
			continue
		}

		if len(parseNotes) == 0 {
			continueNoNoteCount++
			cookie.ChangeCookie()
		} else {
			continueNoNoteCount = 0
			downloadedCount += len(parseNotes)
		}

		if continueNoNoteCount > 5 {
			for i := 0; i < 60; i++ {
				log.Printf("sleep %v/%v for continueNoNoteCount > 5", i+1, 60)
			}
			continueNoNoteCount = 0
		}

		allParseNotes := []string{}
		for _, n := range parseNotes {
			allParseNotes = append(allParseNotes, n.NoteID)
		}
		modelUper := model.Uper{
			UID:              parseUper.UID,
			Name:             parseUper.Name,
			Area:             parseUper.Area,
			AvatarURL:        parseUper.AvatarURL,
			IsGirl:           parseUper.IsGirl,
			Desc:             parseUper.Desc,
			Tags:             parseUper.Tags,
			FansCount:        parseUper.FansCount,
			ReceiveLikeCount: parseUper.ReceiveLikeCount,
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
			Notes:            allParseNotes,
		}
		result, err := storage.GetStorage().InsertOrUpdateUper(modelUper)
		if err != nil {
			log.Printf("InsertOrUpdateUper err:%v parseUper:%+v", err, modelUper)
			continue
		}
		log.Printf("InsertOrUpdateUper succ:%+v result:%v", modelUper, result)
		storage.GetStorage().SetUperScanned(u)

		//failedReason, err := storage.GetStorage().UperAddNote(parseUper.UID, allParseNotes...)
		//if err != nil {
		//	log.Printf("UperAddNote err:%v uid:%v noteid:%v", err, parseUper.UID, allParseNotes)
		//} else if failedReason != "" {
		//	log.Printf("UperAddNote failed:%v uid:%v noteid:%v", failedReason, parseUper.UID, allParseNotes)
		//} else {
		//	//log.Printf("UperAddNote succ. uid:%v noteid:%v", parseUper.UID, n.NoteID)
		//}

		for i, n := range parseNotes {

			dbNote := model.Note{
				NoteID:    n.NoteID,
				UperUID:   u,
				Title:     n.Title,
				URL:       n.URL,
				PosterURL: n.Poster,
				LikeCount: n.LikeCount,
			}
			insertOrUpdate, err := storage.GetStorage().InsertOrUpdateNote(dbNote)
			if err != nil {
				log.Printf("InsertOrUpdateNote err:%v dbNote:%+v", err, dbNote)
				continue
			}
			_ = insertOrUpdate
			log.Printf("InsertOrUpdateNote succ(%v/%v): %+v(%v)", i+1, len(parseNotes), dbNote.Title, insertOrUpdate)

			DownloadPoster(dbNote)
		}

		DownloadUperAvatar(modelUper, config.GetDownloadPath())

	}
}

func DownloadPoster(n model.Note) {
	//log.Printf("DownloadPoster title:%v", n.Title)

	posterPath := filepath.Join(config.GetDownloadPath(), "note_poster", n.UperUID, fmt.Sprintf("%v.jpg", n.NoteID))

	if utils.HasFile(posterPath) {
		return
	}

	code, resp, err := netutil.HttpGetRaw(n.PosterURL)
	if err != nil {
		log.Printf("DownloadPoster netutil.HttpGetRaw err:%v resp:%v", err, resp)
		return
	}

	if code != 200 || len(resp) <= 1024 {
		log.Printf("DownloadPoster netutil.HttpGetRaw failed. code:%v resp:%v note:%+v", code, resp, n)
		return
	}

	err = fileutil.WriteToFile(resp, posterPath)
	if err != nil {
		log.Printf("WriteToFile err:%v resp:%v", err, resp)
		return
	}
	//log.Printf("DownloadPoster WriteToFile succ:%v len(resp):%v", posterPath, utils.GetShowSize(int64(len(resp))))

	//time.Sleep(1 * time.Second)

	return
}

func getAllUpers() []string {
	allProfiles := []string{}
	allProfilesMap := map[string]bool{}
	fileutil.ReadByLine("chore/upers.txt", func(s string) (e error) {
		if allProfilesMap[s] {
			return
		}
		if len(s) != 24 {
			return
		}
		allProfilesMap[s] = true
		allProfiles = append(allProfiles, s)
		return nil

	})
	return allProfiles
}

func StartScanMyShoucang() {
	for {
		upers, notes, err := mydp.ScanMyShoucang(cookie.GetCookie(), 20)
		if err != nil {
			log.Errorf("StartScanMyShoucang ScanMyShoucang err:%v", err)
		}
		for _, note := range notes {
			download.DownloadNoteByID(note)
		}

		for _, u := range upers {
			time.Sleep(1 * time.Minute)
			DownloadUperFirstPageNotes(u)
		}

		time.Sleep(1 * time.Hour)
	}
}
