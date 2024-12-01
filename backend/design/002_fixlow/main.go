package main

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/blog"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"log"
	"strings"
	"time"
)

func main() {
	type SaveInfo struct {
		Idx       int
		UperUID   string
		NoteID    string
		URL       string
		Medias    []blogmodel.Media
		ParseTime string
	}
	infos := []SaveInfo{}

	fileutil.ReadJsonFile("D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\002_fixlow\\only_fixed.json", &infos)

	urlArr := []string{}
	urls := map[string]bool{}

	for _, info := range infos {
		for _, m := range info.Medias {
			urls[m.URL] = true
		}
	}

	log.Printf("get %v urls", len(urls))

	count := 0
	round := 0
	for u := range urls {
		urlArr = append(urlArr, u)
		count++
	}

	round++
	fileutil.WriteToFile([]byte(strings.Join(urlArr, "\n")), fmt.Sprintf("D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\001extract_parse_result\\urls_%v.txt", round))

}

func main1() {
	type SaveInfo struct {
		Idx       int
		UperUID   string
		NoteID    string
		URL       string
		Medias    []blogmodel.Media
		ParseTime string
	}
	infos := []SaveInfo{}

	fileutil.ReadJsonFile("D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\002_fixlow\\parse_result.json", &infos)

	log.Printf("get %v infos", len(infos))

	fixedSaveInfo := []SaveInfo{}

	needFixCount := 0
	for _, elem := range infos {
		needReparse := false
		for _, m := range elem.Medias {
			if m.Type != "video" {
				continue
			}
			if strings.Contains(m.URL, "sns-video-qc.xhscdn.com") {
				needReparse = true
				break
			}
		}
		if !needReparse {
			continue
		}
		needFixCount++
	}

	for i, elem := range infos {
		needReparse := false
		for _, m := range elem.Medias {
			if m.Type != "video" {
				continue
			}
			if strings.Contains(m.URL, "sns-video-qc.xhscdn.com") {
				needReparse = true
				break
			}
		}
		if !needReparse {
			continue
		}

		resp, err := blog.ParseBlog(elem.URL, cookie.GetCookie())
		if err != nil {
			panic(err)
		}
		log.Printf("Fix(%v/%v) ParseBlog resp:%+v", len(fixedSaveInfo), needFixCount, resp)

		saveInfo := SaveInfo{
			Idx:       elem.Idx,
			UperUID:   elem.UperUID,
			NoteID:    elem.NoteID,
			URL:       elem.URL,
			Medias:    resp.Medias,
			ParseTime: time.Now().Format("0102_150405"),
		}

		fixedSaveInfo = append(fixedSaveInfo, saveInfo)

		infos[i] = saveInfo

		fileutil.WriteJsonToFile(fixedSaveInfo, "D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\002_fixlow\\only_fixed.json")
		fileutil.WriteJsonToFile(infos, "D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\002_fixlow\\parse_result_fixed.json")

	}

}
