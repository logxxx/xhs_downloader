package blog

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func GetHomePage(cookie, xs string) (resp []blogmodel.ParseNote, err error) {
	html, err := getExploreHtml(cookie, xs)
	//html := GetExploreHTML(cookie)
	if err != nil {
		log.Infof("GetExplores getExploreHtml err:%v", err)
		return
	}

	fileutil.WriteToFile([]byte(html), "xhs.html")

	resp, err = parseHomeNotes(html)
	if err != nil {
		log.Infof("GetExplores parseHomeNotes err:%v", err)
		return
	}

	return
}

func getExploreHtml(cookie string, xs string) (resp string, err error) {
	return getHtml("https://www.xiaohongshu.com/explore?channel_type=web_note_detail_r10&channel_id=homefeed_recommend", cookie, xs)
}

func getHtml(api, cookie string, xs string) (resp string, err error) {

	//req, err := http.NewRequest("GET", "https://www.xiaohongshu.com/explore", nil)
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authority", "www.xiaohongshu.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Cookie", cookie)
	//req.Header.Set("X-S", xs)
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.xiaohongshu.com/explore")
	req.Header.Set("Origin", "https://www.xiaohongshu.com")
	req.Header.Set("Access-Control-Request-Headers", "batch,biz-type,content-type")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Batch", "true")
	req.Header.Set("Biz-Type", "apm_fe")
	//req.Header.Set("Cookie", `abRequestId=3014709a-be7e-50aa-868d-2d7f962337e6; a1=18ed5e61a91cbrmfh8g6zp14zwnuomrum0fdxr09050000296410; webId=51ebec3ba7911f452e97f6f6d94b1978; gid=yYdf2dKW8J6YyYdf2dKy08F8jySDF6ixYkKuij8Uy4uEIf287l6F7C888JjK4y88f4Y408fj; webBuild=4.27.7; web_session=040069b0a5792a12e752bf2c91344bc9ce199e; xsecappid=xhs-pc-web; acw_tc=22225b604b3c170067936a7724d01db140eb396ea4f04f4cb70d33f49ace5409; websectiga=10f9a40ba454a07755a08f27ef8194c53637eba4551cf9751c009d9afb564467; sec_poison_id=4f5e502b-759c-496d-aa3b-32aa12ab2e78; unread={%22ub%22:%2266a3cf3f000000002701019b%22%2C%22ue%22:%2266a52cc3000000002701d94e%22%2C%22uc%22:34}`)

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if httpResp.StatusCode != 200 {
		err = fmt.Errorf("invalid code:%v", httpResp.StatusCode)
		return
	}

	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return
	}

	resp = string(respBytes)

	return
}

func parseHomeNotes(content string) (resp []blogmodel.ParseNote, err error) {

	//log.Infof("html:%v", content)

	query, err := goquery.NewDocumentFromReader(bytes.NewBufferString(content))
	if err != nil {
		return
	}

	resp = make([]blogmodel.ParseNote, 100)

	query.Find(".author-wrapper .name").Each(func(i int, s *goquery.Selection) {

		//log.Printf("name i:%v Text:%v", i, s.Text())
		resp[i].UperName = s.Text()
	})

	query.Find(".like-wrapper .count").Each(func(i int, s *goquery.Selection) {

		//log.Printf("like i:%v Text:%v", i, s.Text())
		resp[i].LikeCount = int(utils.ToI64(s.Text()))
	})

	query.Find(".title span").Each(func(i int, s *goquery.Selection) {

		//log.Printf("title i:%v Text:%v", i, s.Text())
		resp[i].Title = s.Text()
	})

	query.Find(".author-wrapper a").Each(func(i int, s *goquery.Selection) {

		href, _ := s.Attr("href")

		//log.Printf("href i:%v Text:%v", i, href)
		resp[i].UperUID = utils.Extract(href, "/user/profile/", "?")
	})

	query.Find(".cover img").Each(func(i int, s *goquery.Selection) {

		href, _ := s.Attr("src")

		//log.Printf("img i:%v Text:%v", i, href)
		resp[i].Poster = href
	})

	query.Find("a.cover.ld.mask").Each(func(i int, s *goquery.Selection) {

		href, _ := s.Attr("href")

		//log.Printf("noteURL i:%v Text:%v", i, href)
		resp[i].URL = strings.ReplaceAll("https://www.xiaohongshu.com"+href, `\u0026`, "&")
		resp[i].NoteID = utils.Extract(href, "/explore/", "?")
	})

	newResp := []blogmodel.ParseNote{}

	for _, e := range resp {
		if e.NoteID == "" {
			continue
		}
		newResp = append(newResp, e)
	}

	resp = newResp

	for i, e := range resp {
		fmt.Printf("%v/%v %+v\n", i+1, len(resp), e)
	}

	return
}
