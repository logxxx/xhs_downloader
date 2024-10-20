package main

import (
	"errors"
	"fmt"
	"github.com/logxxx/xhs_downloader/biz/thumb"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/randutil"
)

var (
	//user-agent
	uaList = []string{
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36Chrome 17.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0Firefox 4.0.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
	}
)

func GetDownloadRealPath(req ParseBlogResp, idx int, mediaType string, downloadPath string) string {

	downloadPath = filepath.Join(downloadPath, req.UserID)

	shortTitle := utils.ShortTitle(req.Title)
	if len(shortTitle) <= 0 {
		shortTitle = fmt.Sprintf("%v", time.Now().Unix())
	}
	fileTitle := fmt.Sprintf("%v_%v", req.UserID, req.NoteID)
	if mediaType == "image" && idx > 0 {
		fileTitle += fmt.Sprintf("_%v", idx)
	}
	//log.Printf("fileTitle:%v", fileTitle)
	suffix := ".jpg"
	if mediaType == "video" || mediaType == "live" {
		suffix = ".mp4"
	}
	if idx > 0 {
		suffix = fmt.Sprintf("_%v%v", idx, suffix)
	}
	if mediaType == "video" {
		downloadPath = filepath.Join(downloadPath, "video")
	}
	if mediaType == "live" {
		downloadPath = filepath.Join(downloadPath, "live")
	}

	fileName := fmt.Sprintf("%v%v", fileTitle, suffix)
	fileRealPath, _ := fileutil.ReplaceInvalidChar(filepath.Join(downloadPath, fileName), "x")

	return fileRealPath
}

func downloadMedia(req ParseBlogResp, idx int, downloadPath string, useBackup bool) (err error, canRetry bool) {

	m := req.Medias[idx]

	reqURL := m.URL
	if m.BackupURL != "" && useBackup {
		reqURL = m.BackupURL
	}

	httpReq, _ := http.NewRequest("GET", reqURL, nil)
	httpReq.Header.Set("user-agent", uaList[rand.Intn(len(uaList))])

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("download GET err:%v", err)
		return
	}

	defer func() {
		httpResp.Body.Close()
	}()

	if httpResp.ContentLength > 300*1024*1024 {
		log.Printf("download GET err:%v", "file size too large")
		err = errors.New("file too large")
		return
	}

	if httpResp.ContentLength <= 50*1024 {
		log.Printf("download GET err:%v", "file size too small")
		//err = errors.New("file too small")
		//return
		if httpResp.ContentLength == 0 {
			err = errors.New("file body empty")
			canRetry = true
			return
		}
	}

	fileRealPath := GetDownloadRealPath(req, idx, m.Type, downloadPath)

	localFileSize := utils.GetFileSize(fileRealPath)
	if localFileSize == httpResp.ContentLength {
		log.Printf("already downloaded:%v size:%v", fileRealPath, utils.GetShowSize(localFileSize))
		return
	}

	log.Printf("Downloading %v: %v contentlen:%v", m.Type, m.URL, utils.GetShowSize(httpResp.ContentLength))

	os.MkdirAll(filepath.Dir(fileRealPath), 0755)
	for {
		if !utils.HasFile(fileRealPath) {
			break
		}
		time.Sleep(100 * time.Millisecond)
		fileRealPath = filepath.Join(filepath.Dir(fileRealPath), fmt.Sprintf("%v_%v", randutil.RandStr(4), filepath.Base(fileRealPath)))
	}
	os.Remove(fileRealPath)
	f, err := os.OpenFile(fileRealPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		log.Printf("download OpenFile err:%v path:%v", err, fileRealPath)
		return
	}
	defer func() {
		f.Close()
		if canRetry {
			os.Remove(fileRealPath)
		}
	}()

	_, err = io.Copy(f, httpResp.Body)
	if err != nil {
		log.Printf("download Copy err:%v path:%v", err, fileRealPath)
		canRetry = true
		return
	}

	m.DownloadPath = fileRealPath

	return
}

var (
	downloadRetryCount = 0
)

func Download(req ParseBlogResp, downloadPath string, splitByDate bool) (resp []*Media) {

	if len(req.Medias) == 0 {
		log.Printf("**** NOTHING TO DOWNLOAD ****")
		return
	}

	if splitByDate {
		downloadPath = filepath.Join(downloadPath, fmt.Sprintf("%v", time.Now().Format("20060102")))
	}

	//log.Printf("Downloading blog %+v", req)

	for i := range req.Medias {
		retryTimes := 0
		err, canRetry := downloadMedia(req, i, downloadPath, false)
		if err == nil {
			thumb.MakeThumb(req.Medias[i].DownloadPath)
			continue
		}

		if !canRetry {
			continue
		}

		log.Printf("downloadMedia err:%v", err)
		retryTimes++
		if retryTimes > 5 {
			continue
		}
		downloadRetryCount++
		log.Printf("downloadMedia retry![%v]", downloadRetryCount)
		err, _ = downloadMedia(req, i, downloadPath, true)

		if err == nil {
			thumb.MakeThumb(req.Medias[i].DownloadPath)
		}

	}

	return req.Medias
}
