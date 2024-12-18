package download

import (
	"context"
	"fmt"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/webhook"
	log "github.com/sirupsen/logrus"
	"io"

	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func GetDownloadRealPath(req blogmodel.ParseBlogResp, idx int, mediaType string, downloadPath string) string {

	downloadPath = filepath.Join(downloadPath, req.UserID)

	shortTitle := utils.ShortTitle(req.Title)
	if len(shortTitle) <= 0 {
		shortTitle = fmt.Sprintf("%v", time.Now().Unix())
	}
	fileTitle := fmt.Sprintf("%v_%v", req.UserID, req.NoteID)
	//if mediaType == "image" && idx > 0 {
	//	fileTitle += fmt.Sprintf("_%v", idx)
	//}
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

func downloadMediaByThunder(req blogmodel.ParseBlogResp, idx int, downloadPath string) (err error) {

	m := &req.Medias[idx]

	if utils.HasFile(downloadPath) {
		log.Infof("ALREADY HAS FILE:%v", downloadPath)
		m.DownloadPath = downloadPath
		return
	}

	_, err = webhook.Download(context.Background(), m.URL, downloadPath, false)
	if err == nil {
		m.DownloadPath = downloadPath
	}
	return
}

func downloadMediaByHttp(scene string, req blogmodel.ParseBlogResp, idx int, downloadPath string) (err error) {

	m := &req.Medias[idx]

	fileRealPath := GetDownloadRealPath(req, idx, m.Type, downloadPath)

	if utils.HasFile(fileRealPath) {
		log.Printf("ALREADY DOWNLOADED:%v", fileRealPath)
		return
	}

	httpReq, _ := http.NewRequest("GET", m.URL, nil)
	httpReq.Header.Set("user-agent", uaList[rand.Intn(len(uaList))])

	httpResp2, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("download GET err:%v", err)
		return
	}

	defer func() {
		httpResp2.Body.Close()
	}()

	os.MkdirAll(filepath.Dir(fileRealPath), 0755)
	f, err := os.OpenFile(fileRealPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		log.Printf("download OpenFile err:%v path:%v reqURL:%v", err, fileRealPath, m.URL)
		return
	}
	defer func() {
		f.Close()
	}()

	_, err = io.Copy(f, httpResp2.Body)
	if err != nil {
		log.Printf("download io.Copy err:%v path:%v reqURL:%v", err, fileRealPath, m.URL)
		return
	}

	log.Printf("[%v]Download Finish %v[%v] %v %v", scene, m.Type, idx, utils.GetShowFileSize(fileRealPath), filepath.Base(fileRealPath))

	m.DownloadPath = fileRealPath

	return
}

func downloadMediaByWebHook(scene string, req blogmodel.ParseBlogResp, idx int, downloadPath string, mustUseLocal bool) (err error) {
	m := &req.Medias[idx]

	fileRealPath := GetDownloadRealPath(req, idx, m.Type, downloadPath)

	if utils.HasFile(fileRealPath) {
		log.Printf("DOWNLOADED:%v", fileRealPath)
		return
	}

	log.Printf("[%v]Downloading %v[%v]", scene, m.Type, idx)

	_, err = webhook.Download(context.Background(), m.URL, fileRealPath, false)
	if err != nil {
		log.Printf("webhook.Download err:%v url:%v path:%v", err, m.URL, fileRealPath)
	}

	m.DownloadPath = fileRealPath

	return
}

func downloadMedia(scene string, req blogmodel.ParseBlogResp, idx int, downloadPath string, mustUseLocal bool) (err error) {

	m := &req.Medias[idx]

	reqURL := m.URL

	httpReq, _ := http.NewRequest("HEAD", reqURL, nil)
	httpReq.Header.Set("user-agent", uaList[rand.Intn(len(uaList))])

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("download HEAD err:%v", err)
		return
	}

	if httpResp.ContentLength > 700*1024*1024 {
		log.Printf("download GET warning:%v size:%v", "file size too large", utils.GetShowSize(httpResp.ContentLength))
		return
	}

	fileRealPath := GetDownloadRealPath(req, idx, m.Type, downloadPath)

	localFileSize := utils.GetFileSize(fileRealPath)
	if localFileSize == httpResp.ContentLength {
		log.Printf("already downloaded:%v size:%v", fileRealPath, utils.GetShowSize(localFileSize))
		return
	}

	log.Printf("[%v]Downloading %v[%v] len:%v", scene, m.Type, idx, utils.GetShowSize(httpResp.ContentLength))

	if !mustUseLocal {
		_, err = webhook.Download(context.Background(), m.URL, fileRealPath, false)
		if err == nil {
			m.DownloadPath = fileRealPath
			log.Printf("Download by webhook SUCC")
			return
		}
		log.Printf("Download by webhook err:%v", err)
		panic(err)
	}

	httpReq, _ = http.NewRequest("GET", reqURL, nil)
	httpReq.Header.Set("user-agent", uaList[rand.Intn(len(uaList))])

	httpResp2, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("download GET err:%v", err)
		return
	}

	defer func() {
		httpResp2.Body.Close()
	}()

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
		log.Printf("download OpenFile err:%v path:%v reqURL:%v", err, fileRealPath, reqURL)
		return
	}
	defer func() {
		f.Close()
	}()

	_, err = io.Copy(f, httpResp2.Body)
	if err != nil {
		log.Printf("download io.Copy err:%v path:%v reqURL:%v", err, fileRealPath, reqURL)
		elems := strings.Split(req.Title, "\n")
		if len(elems) > 0 {
			req.Title = elems[0]
		}
		fileutil.AppendToFile("download_failed.txt", fmt.Sprintf("%v\n%v\n%v\n%v\n", req.Title, req.BlogURL, fileRealPath, reqURL))
		webhook.Download(context.Background(), m.URL, fileRealPath, false)
		return
	}

	m.DownloadPath = fileRealPath

	log.Infof("DOWNLOAD BY ORIG SUCC! path:%v reqURL:%v", fileRealPath, reqURL)

	return
}

var (
	downloadRetryCount = 0
)

func DownloadToHome(scene string, req blogmodel.ParseBlogResp, downloadPath string, splitByDate bool, forceUseLocal bool) (resp []blogmodel.Media) {
	//log.Printf("Download start:%v %v", req.Title, req.BlogURL)

	if len(req.Medias) == 0 {
		log.Printf("**** [%v]Download: NOTHING TO DOWNLOAD ****", scene)
		return
	}

	if splitByDate {
		downloadPath = filepath.Join(downloadPath, fmt.Sprintf("%v", time.Now().Format("20060102")), "home")
	}

	//log.Printf("Downloading to:%v", downloadPath)

	for idx := range req.Medias {
		i := idx
		m := &req.Medias[idx]
		if m.Type == "image" || m.Type == "live" {
			err := downloadMediaByHttp(scene, req, i, downloadPath)
			if err == nil {
				continue
			}
		}
		downloadMediaByWebHook(scene, req, i, downloadPath, forceUseLocal)
	}

	return req.Medias
}

func Download(scene string, req blogmodel.ParseBlogResp, downloadPath string, splitByDate bool, forceUseLocal bool) (resp []blogmodel.Media) {
	//log.Printf("Download start:%v %v", req.Title, req.BlogURL)

	if len(req.Medias) == 0 {
		log.Printf("**** [%v]Download: NOTHING TO DOWNLOAD ****", scene)
		return
	}

	if splitByDate {
		downloadPath = filepath.Join(downloadPath, fmt.Sprintf("%v", time.Now().Format("20060102")))
	}

	//log.Printf("Downloading to:%v", downloadPath)

	for idx := range req.Medias {
		i := idx
		m := &req.Medias[idx]
		if m.Type == "image" || m.Type == "live" {
			err := downloadMediaByHttp(scene, req, i, downloadPath)
			if err == nil {
				continue
			}
		}
		downloadMediaByWebHook(scene, req, i, downloadPath, forceUseLocal)
	}

	return req.Medias
}
