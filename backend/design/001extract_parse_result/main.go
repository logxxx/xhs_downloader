package main

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/blog"
	log "github.com/sirupsen/logrus"
	"strings"
)

func main() {
	type SaveInfo struct {
		Idx       int
		UperUID   string
		NoteID    string
		URL       string
		Medias    []blog.Media
		ParseTime string
	}
	infos := []SaveInfo{}

	fileutil.ReadJsonFile("D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\001extract_parse_result\\parse_result.json", &infos)

	log.Printf("get %v infos", len(infos))

	urls := map[string]bool{}

	for _, info := range infos {
		for _, m := range info.Medias {
			urls[m.URL] = true
		}
	}

	log.Printf("get %v urls", len(urls))

	urlArr := []string{}
	count := 0
	round := 0
	for u := range urls {
		urlArr = append(urlArr, u)
		count++
		if count == 1000 {
			round++
			fileutil.WriteToFile([]byte(strings.Join(urlArr, "\n")), fmt.Sprintf("D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\001extract_parse_result\\urls_%v.txt", round))
			urlArr = []string{}
			count = 0
		}
	}

	round++
	fileutil.WriteToFile([]byte(strings.Join(urlArr, "\n")), fmt.Sprintf("D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\001extract_parse_result\\urls_%v.txt", round))

}
