package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/logxxx/utils"
	"github.com/logxxx/utils/ffmpeg"
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
	log.Printf("fileTitle:%v", fileTitle)
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

func downloadMedia(req ParseBlogResp, idx int, downloadPath string) (err error, canRetry bool) {

	m := req.Medias[idx]

	log.Printf("Downloading %v: %v", m.Type, m.URL)

	httpReq, _ := http.NewRequest("GET", m.URL, nil)
	httpReq.Header.Set("user-agent", uaList[rand.Intn(len(uaList))])

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("download GET err:%v", err)
		return
	}

	defer func() {
		httpResp.Body.Close()
	}()

	log.Printf("contentlen:%v", utils.GetShowSize(httpResp.ContentLength))

	if httpResp.ContentLength > 200*1024*1024 {
		log.Printf("download GET err:%v", "file size too large")
		err = errors.New("file too large")
		//return
	}

	if httpResp.ContentLength <= 50*1024 {
		log.Printf("download GET err:%v", "file size too small")
		//err = errors.New("file too small")
		//return
	}

	fileRealPath := GetDownloadRealPath(req, idx, m.Type, downloadPath)

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

func Download(req ParseBlogResp, downloadPath string, splitByDate bool) (resp []*Media) {

	log.Printf("Downloading blog %+v", req)

	if splitByDate {
		downloadPath = filepath.Join(downloadPath, fmt.Sprintf("%v", time.Now().Format("20060102")))
	}

	for i := range req.Medias {
		retryTimes := 0
		err, canRetry := downloadMedia(req, i, downloadPath)
		time.Sleep(2 * time.Second)
		if err != nil {
			log.Printf("downloadMedia err:%v", err)
			if canRetry {
				retryTimes++
				if retryTimes > 5 {
					continue
				}
				log.Printf("downloadMedia retry!")
				err, _ = downloadMedia(req, i, downloadPath)
			}
		}

		if err != nil {
			continue
		}

		filePath := req.Medias[i].DownloadPath
		fileSize := utils.GetFileSize(filePath)
		if fileSize > 3*1024*1024 {
			if req.Medias[i].Type == "video" || req.Medias[i].Type == "live" {
				makeThumb(filePath)
			} else {
				cmd := `ffmpeg -i %v -y -vf scale=480:-1 %v`
				ffp := ffmpeg.FFProbe("ffprobe")
				fInfo, _ := ffp.NewVideoFile(filePath)
				if fInfo != nil && fInfo.Width < fInfo.Height {
					cmd = `ffmpeg -i %v -y -vf scale=-1:480 %v`
				}
				cmd = fmt.Sprintf(cmd, filePath, filepath.Join(filepath.Dir(filePath), ".thumb", filepath.Base(filePath)))

				if !utils.HasFile(filepath.Join(filepath.Dir(filePath), ".thumb")) {
					os.MkdirAll(filepath.Join(filepath.Dir(filePath), ".thumb"), 0755)
				}
				runCommand(cmd)
			}
		}

	}

	return req.Medias
}

func makeThumb(filePath string) error {

	log.Printf("make thumb %v", filePath)
	_, err := ffmpeg.GenePreviewVideoSlice(filePath, func(vInfo *ffmpeg.VideoFile) ffmpeg.GenePreviewVideoSliceOpt {

		segNum := 3
		segDur := 3
		skipStart := 1
		skipEnd := 3
		if vInfo.Duration < 15 {
			segNum = 2
		}
		if vInfo.Duration > 2*60 {
			skipStart = 10
			skipEnd = 30
		}
		if vInfo.Duration > 5*60 {
			skipStart = 30
			skipEnd = 30
			segNum = 5
		}

		return ffmpeg.GenePreviewVideoSliceOpt{
			//ToPath:      filePath + ".thumb.mp4",
			ToPath:      filepath.Join(filepath.Dir(filePath), ".thumb", filepath.Base(filePath)),
			SegNum:      segNum,
			SegDuration: segDur,
			SkipStart:   skipStart,
			SkipEnd:     skipEnd,
		}
	})
	return err
}

func runCommand(command string) (output []byte, err error) {
	log.Printf("runCommand:%v", command)
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("runCommand err:%v output:%v", err, string(output))
		return
	}
	//log.Printf("runCommand output:%v", string(output))
	return
}
