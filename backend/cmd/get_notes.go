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
	"os"
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

			_, parseNotes, err := mydp.GetNotes(u.UID, cookie.GetCookie(), 1)
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

				download.DownloadNote(dbNote, false, true)
			}

			if uperChanged {
				log.Infof("StartDownloadRecrentlyNotes update uper:%+v => %v", oldNoteCount, len(u.Notes))
				storage.GetStorage().InsertOrUpdateUper(u)
			}

			return
		})
	}
}

func DownloadUperNPageNotes(uid string, n int) {

	uper := storage.GetStorage().GetUper(0, uid)

	parseUper, parseNotes, err := mydp.GetNotes(uid, cookie.GetCookie(), n)
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

		dbNote := storage.GetStorage().GetNote(n.NoteID)
		if dbNote.IsDownloaded() {
			log.Printf("note downloaded: %v %v", n.NoteID, dbNote.Title)
			continue
		}

		dbNote = model.Note{
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

		result := download.DownloadNote(dbNote, true, true)
		log.Printf("Downlaod Result:%v", result)

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
		parseUper, parseNotes, err := mydp.GetNotes(u, cookie.GetCookie(), -1)
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

		upers := []string{}

		if utils.HasFile("myshoucang_50_upers.txt") {
			fileutil.ReadByLine("myshoucang_50_upers.txt", func(s string) error {
				if s == "" {
					return nil
				}
				upers = append(upers, s)
				return nil
			})
			log.Printf("get %v upers from 50_upers.txt", len(upers))
		}

		notes := []string{}
		var err error
		if len(upers) <= 0 {
			upers, notes, err = mydp.ScanMyShoucang(cookie.GetCookie1(), 50)
			if err != nil {
				log.Errorf("StartScanMyShoucang ScanMyShoucang err:%v", err)
			}

			fileutil.WriteToFile([]byte(strings.Join(upers, "\n")), "myshoucang_50_upers.txt")
		}

		for _, note := range notes {
			err = download.DownloadNoteByID(note)
			if err != nil {
				log.Errorf("DownloadNoteByID err:%v note:%v", err, note)
			}
		}

		hit := false
		_ = hit
		for i, u := range upers {

			if u == "5dc0eb66000000000100502b" {
				hit = true
			}
			if !hit {
				//continue
			}

			if u == "5c26e25b0000000006012115" {
				continue
			}

			log.Printf("Start scan uper %v/%v: %v", i+1, len(upers), u)

			uper := storage.GetStorage().GetUper(0, u)
			if uper.ID > 0 {
				log.Printf("uper already has")
				//continue
			}

			DownloadUperNPageNotes(u, -1)
			cookie.ChangeCookie()
			time.Sleep(1 * time.Minute)
		}

		time.Sleep(1 * time.Hour)

		//for {
		//	if time.Now().Hour() > 8 {
		//		break
		//	}
		//	log.Printf("ScanMyShoucang skip: not time")
		//	time.Sleep(10 * time.Minute)
		//}
	}
}

func FixFailedVideo() {
	type Info struct {
		Title        string
		NoteURL      string
		DownlaodPath string
		VideoURL     string
	}
	notes := []Info{}
	note := Info{}
	line := 0
	fileutil.ReadByLine("D:\\mytest\\mywork\\xhs_downloader\\backend\\cmd\\download_failed.txt", func(s string) (e error) {
		line++
		//富家千金风即视感美甲太有氛围感了美甲💅
		//https://www.xiaohongshu.com/explore/663c8ae6000000001e0311db?xsec_token=ABvbFhKcuagTZtKpgvCXZSQGNmmP0NZGToqLLf9eRET6Q=&xsec_source=pc_user
		//E:\xhs_downloader_output\20241030\5b012f40e8ac2b46e32d32a2\video\5b012f40e8ac2b46e32d32a2_663c8ae6000000001e0311db.mp4
		//http://sns-video-bd.xhscdn.com/pre_post/1040g0cg312ips0n5080g4a5g28nk0cl2ssfuikg
		if line%4 == 0 {
			note.Title = s
			return
		}
		if line%4 == 2 {
			note.NoteURL = s
			return
		}
		if line%4 == 3 {
			note.DownlaodPath = s
			return
		}
		if line%4 == 1 {
			note.VideoURL = s
			notes = append(notes, note)
			note = Info{}
			return
		}
		return
	})

	for i, n := range notes {
		log.Printf("FixFailedVideo note %v/%v: %+v", i+1, len(notes), n)
		if utils.GetFileSize(n.DownlaodPath) > 1024 {
			log.Printf("ALREADY HAS FILE!")
			continue
		}
		os.Remove(n.DownlaodPath)
		_, respData, err := netutil.HttpGetRaw(n.VideoURL)
		if err != nil {
			log.Errorf("HttpGetRaw err:%v url:%v", err, n.VideoURL)
			continue
		}
		fileutil.WriteToFile(respData, n.DownlaodPath)

		dbNote := storage.GetStorage().GetNote(ExtractNoteIDByURL(n.NoteURL))
		dbNote.Video = n.DownlaodPath
		dbNote.DownloadNothing = false
		dbNote.DownloadTime = time.Now()
		err = storage.GetStorage().UpdateNote(dbNote)
		if err != nil {
			log.Errorf("UpdateNote err:%v dbNote:%+v", err, dbNote)
		}
	}
}

func ExtractNoteIDByURL(noteURL string) string {
	return utils.Extract(noteURL, "/", "?")
}
