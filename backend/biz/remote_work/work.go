package remote_work

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/utils/runutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/mydp"
	"github.com/logxxx/xhs_downloader/biz/queue"
	"github.com/logxxx/xhs_downloader/model"
	"github.com/logxxx/xhs_downloader/utils"
	log "github.com/sirupsen/logrus"
	"moul.io/http2curl"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetFeedApiCookie() string {
	if os.Getenv("FEED_COOKIE") == "cookie1" {
		return cookie.GetCookie1()
	}
	if os.Getenv("FEED_COOKIE") == "cookie2" {
		return cookie.GetCookie2()
	}
	if os.Getenv("FEED_COOKIE") == "cookie3" {
		return cookie.GetCookie3()
	}
	panic(fmt.Sprintf("invalid FEED_COOKIE: [%s]", os.Getenv("FEED_COOKIE")))
}

func init() {
	if utils.IsWorker() {
		fmt.Printf("progress IS WORKER, so StartWaitForWork()\n")
		runutil.GoRunSafe(StartWaitForWork)
	} else {
		fmt.Printf("progress IS WORKER FALSE, so StartRecvRemoteWorkResult()\n")
		runutil.GoRunSafe(StartRecvRemoteWorkResult)
	}

}

func StartWaitForWork() {
	round := 0
	for {
		if round != 0 {
			time.Sleep(10 * time.Second)
		}
		round++
		work := &model.Work{}
		_, err := netutil.HttpGet("http://47.119.170.71:8088/recv_work", work)
		if err != nil {
			log.Errorf("recv work err:%v", err)
			continue
		}
		if work.NoteID == "" {
			continue
		}
		log.Infof("get work:%+v", work)
		xs, xt, err := mydp.GetXsXt(work.NoteID, work.XSecToken)
		if err != nil {
			log.Errorf("GetXsXt err:%v", err)
			mydp.SendWork(work.BlogURL, work.NoteID, work.XSecToken)
			continue
		}
		if xs == "" || xt <= 0 {
			log.Errorf("invalid xs:%v xt:%v", xs, xt)
			mydp.SendWork(work.BlogURL, work.NoteID, work.XSecToken)
			continue
		}

		noteID := work.NoteID
		xsecToken := work.XSecToken
		blogURL := work.BlogURL

		reqHeader := getFeedApiHeaders()

		reqContent := `{"source_note_id":"%v","image_formats":["jpg","webp","avif"],"extra":{"need_body_topic":"1"},"xsec_source":"pc_user","xsec_token":"%v"}`
		reqContent = fmt.Sprintf(reqContent, noteID, xsecToken)
		fileutil.WriteToFile([]byte(reqContent), "req_body.json")
		reqBuf := bytes.NewBufferString(reqContent)
		//log.Printf("START REQUEST FEED url:%v reqBody:%v", ev.Request.URL, reqContent)
		reqURL := "https://edith.xiaohongshu.com/api/sns/web/v1/feed"
		httpReq, _ := http.NewRequest("POST", reqURL, reqBuf)
		for k, v := range reqHeader {
			httpReq.Header.Set(k, fmt.Sprintf("%v", v))
		}
		httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
		httpReq.Header.Set("Origin", "https://www.xiaohongshu.com")
		httpReq.Header.Set("referer", "https://www.xiaohongshu.com/")
		httpReq.Header.Set("content-length", "")
		httpReq.Header.Set("cookie", GetFeedApiCookie())
		httpReq.Header.Set("X-s", xs)
		httpReq.Header.Set("X-t", fmt.Sprintf("%v", xt))

		curl, err := http2curl.GetCurlCommand(httpReq)
		if err == nil {
			fileutil.WriteToFile([]byte(curl.String()), "curl.txt")
		}

		respCode, respBytes, err2 := netutil.HttpDo(httpReq)
		if err2 != nil {
			log.Errorf("call feed api err:%v", err)
		}
		_ = respCode
		log.Printf("HttpDo respCode:%v resp:%v err:%v", respCode, string(respBytes), err)

		feedResp := &blogmodel.FeedResp{}

		if strings.Contains(string(respBytes), "访问频次异常") {
			log.Errorf("访问频次异常")
			mydp.SendWork(work.BlogURL, work.NoteID, work.XSecToken)
			panic(err)
			continue
		}

		json.Unmarshal(respBytes, feedResp)

		parseResult := mydp.ConvFeedResp2ParseResult(blogURL, feedResp)

		fileutil.WriteJsonToFile(parseResult, fmt.Sprintf("work_result_%v.json", time.Now().Format("20060102_150405")))

		_, err = netutil.HttpPost("http://47.119.170.71:8088/send_work_result", parseResult, nil)
		if err != nil {
			log.Errorf("send_work_result err:%v", err)
			mydp.SendWork(work.BlogURL, work.NoteID, work.XSecToken)
		} else {
			log.Infof("send_work_result SUCC!")
		}

	}
}

func StartRecvRemoteWorkResult() {
	for {
		time.Sleep(10 * time.Second)
		workResult := &blogmodel.ParseBlogResp{}
		code, err := netutil.HttpGet("http://47.119.170.71:8088/recv_work_result", workResult)
		log.Infof("recv_work_result code:%v err:%v", code, err)
		if len(workResult.Medias) <= 0 {
			log.Infof("recv_work_result get empty media.")
			continue
		}
		log.Infof("recv_work_result get media:%v", workResult.GetMediaSimpleInfo())
		queue.Push("parse_blog", workResult)
	}
}

func getFeedApiHeaders() http.Header {

	req, err := http.NewRequest("POST", "https://edith.xiaohongshu.com/api/sns/web/v1/feed", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", "abRequestId=18d8450d-628a-5dcd-936d-01b5b40c8276; a1=19179e61ba9rsl23v6l4my9i37wq8py2vjeo1rfl850000293818; webId=e9976c88abe83d72a6350bd21221909a; gid=yjyWjdKJ8DhfyjyWjdKyDVFi0jFM1JqhK146d1Yvj3qWEl28AYUvJy888JjqYyY8i4WJqjWf; web_session=0400697999f01a77b8f86ff0c4344ba1154db9; webBuild=4.46.0; websectiga=16f444b9ff5e3d7e258b5f7674489196303a0b160e16647c6c2b4dcb609f4134; sec_poison_id=bed7d925-7a50-4579-82b2-2576f4e327dc; acw_tc=0a4a870a17334058586182715eb5200875645c08f887dd816930ccbfe60684; xsecappid=xhs-pc-web; unread={%22ub%22:%22672ef96f000000001b02fea4%22%2C%22ue%22:%226743d72f0000000002018605%22%2C%22uc%22:29}")
	req.Header.Set("Origin", "https://www.xiaohongshu.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://www.xiaohongshu.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("X-B3-Traceid", "01484543459be224")
	req.Header.Set("X-S", "XYW_eyJzaWduU3ZuIjoiNTYiLCJzaWduVHlwZSI6IngyIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6ImI4YWIxYzhiYTM1NjUxMDcyYzBkYWVjMWU0YWQyNmQ0NDkzZDNmODg2ZTNiZGE5NTllNzg1YjNlOWNhZjBjNmFhZjFiMDJlZTQxZWQ2Zjg2ZTA4M2VkZjM4NjkzMTc2YzI2YWJjNjIyMjNkODdlN2JlMTZjNTlkMzc3MDJjZmMwM2VjOTY4ZjkwZjMyMGQ1N2U4ZGEzODRmYWY2MTVjMjI5YjkwODM4NDY2OGY3Mzg2NTA5YjMyZjMwMDZkYTZmNGU3NDNhNjJhNTJmNTg0YTIwNjcxNzliNThiYzg0ZWY0NDBiMGMyZjA3MjNiZjJlYjBiYzY4MzQ5ZWY2YTQyMWFjN2ExZmQ4YmFmZmIyZmVlYmRkMzRhMjE2OGMwYWFlYzUxZmEzMWE4NjU0YmUwMDRlOWZiMjgxYjdiNDFlNjQ3MzZmN2Q3N2MxZmJkMjBiOGYxZTgxOTUyOWI2MGYzNDFkYjk0OTdkNzI2MWNjMDhhMjQ2MjY4M2E1OTBjNmMyNWIxODdlZWQ5ZDM2ZmJlMjZlZDZjODA1OTNhZjJjM2NlIn0=")
	req.Header.Set("X-S-Common", "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")
	req.Header.Set("X-T", "1733405894218")
	req.Header.Set("X-Xray-Traceid", "c9cb8572fa0169b35e2f4b4352df8cf4")

	return req.Header
}
