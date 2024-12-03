package crontab

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/download"
	"github.com/logxxx/xhs_downloader/biz/mydp"
	"github.com/logxxx/xhs_downloader/biz/queue"
	"github.com/logxxx/xhs_downloader/biz/storage"
	utils2 "github.com/logxxx/xhs_downloader/utils"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	IsPaused = false
)

func StartScanMyShoucang() {

	for {

		upers := []string{}

		upersFile := "chore/myshoucang_50_upers.txt"

		if utils.HasFile(upersFile) {
			fileutil.ReadByLine(upersFile, func(s string) error {
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
			upers, notes, err = mydp.ScanMyShoucang(cookie.GetCookie1(), -1)
			if err != nil {
				log.Errorf("StartScanMyShoucang ScanMyShoucang err:%v", err)
			}

			fileutil.WriteToFile([]byte(strings.Join(upers, "\n")), upersFile)
		}

		for _, note := range notes {
			err = download.DownloadNoteByID(note)
			if err != nil {
				log.Errorf("DownloadNoteByID err:%v note:%v", err, note)
			}
		}

		hit := false
		_ = hit

		continueEmptyWorkTimes := 0
		for i, u := range upers {
			log.Infof("scan upers %v/%v:%v", i+1, len(upers), u)
			for {
				if IsPaused {
					log.Printf("pausing...")
					time.Sleep(10 * time.Second)
				} else {
					break
				}
			}

			if len(u) != 24 {
				continue
			}

			if u == "5c500b7a000000001803676c" {
				hit = true
				continue
			}
			if !hit {
				continue
			}

			if utils2.IsBlackUid(u) {
				log.Infof("is black uid:%v", u)
				continue
			}

			log.Printf("Start scan uper %v/%v: %v", i+1, len(upers), u)

			dbUper := storage.GetStorage().GetUper(0, u)
			if dbUper.ID > 0 && time.Since(dbUper.GalleryEmptyLastTime).Hours() < 3*24 {
				log.Printf("GalleryEmptyLastTime too recent:%v", dbUper.GalleryEmptyLastTime.Format("2006-01-02 15:04:05"))
				continue
			}

			if dbUper.ID > 0 && time.Since(dbUper.NotesLastUpdateTime).Hours() < 3*24 {
				log.Printf("NotesLastUpdateTime too recent:%v", dbUper.NotesLastUpdateTime.Format("2006-01-02 15:04:05"))
				continue
			} else {
				log.Printf("NotesLastUpdateTime:%v", dbUper.NotesLastUpdateTime)
			}

			if dbUper.IsBanned {
				log.Printf("Uper is BANNED")
				continue
			}

			cookie.ChangeCookie()

			uperURL := fmt.Sprintf("https://www.xiaohongshu.com/user/profile/%v?channel_type=web_note_detail_r10&parent_page_channel_type=web_profile_board", u)

			fileutil.AppendToFile("download_report.txt", fmt.Sprintf("\n%v/%v %v\n", i+1, len(upers), uperURL))

			noteResp, _ := mydp.GetNotes2(u, cookie.GetCookie(), func(parseResult blogmodel.ParseBlogResp) {

				queue.Push("parse_blog", parseResult)

			})

			record := []string{"\n----------------------------------------", fmt.Sprintf("%v [%v/%v] %v %v_NOTES", time.Now().Format("2006/01/02 15:04:05"), i+1, len(upers), u, noteResp.NoteCount)}
			record = append(record, noteResp.Records...)
			fileutil.AppendToFile("getnotes2_record.txt", strings.Join(record, "\n"))

			if noteResp.IsUperBanned {
				dbU := storage.GetStorage().GetUper(0, u)
				dbU.IsBanned = true
				result, err := storage.GetStorage().InsertOrUpdateUper(dbU)
				log.Infof("Refresh IsUperBanned. uid:%v result:%v err:%v", u, result, err)
				continue
			}

			if noteResp.IsGalleryEmpty {
				dbU := storage.GetStorage().GetUper(0, u)
				dbU.UID = u
				dbU.GalleryEmptyLastTime = time.Now()
				result, err := storage.GetStorage().InsertOrUpdateUper(dbU)
				log.Infof("Refresh GalleryEmptyLastTime. uid:%v result:%v err:%v", u, result, err)
			}

			log.Infof("UPER [%v/%v]%v GET %v NOTES", i+1, len(upers), u, noteResp.NoteCount)

			if !noteResp.IsGalleryEmpty && noteResp.NoteCount <= 0 {
				log.Infof("GET EMPTY NOTE")
				continueEmptyWorkTimes++
			} else {
				continueEmptyWorkTimes = 0
				dbU := storage.GetStorage().GetUper(0, u)
				if dbU.ID > 0 {
					dbU.NotesLastUpdateTime = time.Now()
					storage.GetStorage().InsertOrUpdateUper(dbU)
				}
			}
			if continueEmptyWorkTimes > 10 {
				log.Infof("GET TOO MANY EMPTY NOTE")
				return
			}

			time.Sleep(10 * time.Second)

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
