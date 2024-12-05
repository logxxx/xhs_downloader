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
	"strings"
	"time"
)

func init() {
	if utils.IsWorker() {
		runutil.GoRunSafe(StartWaitForWork)
	} else {
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
			continue
		}
		if xs == "" || xt <= 0 {
			log.Errorf("invalid xs:%v xt:%v", xs, xt)
			continue
		}

		noteID := work.NoteID
		xsecToken := work.XSecToken
		blogURL := work.BlogURL

		reqHeader := map[string]string{}

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
		httpReq.Header.Set("cookie", cookie.GetCookie3())
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
			continue
		}

		json.Unmarshal(respBytes, feedResp)

		parseResult := mydp.ConvFeedResp2ParseResult(blogURL, feedResp)

		fileutil.WriteJsonToFile(parseResult, fmt.Sprintf("work_result_%v.json", time.Now().Format("20060102_150405")))

		_, err = netutil.HttpPost("http://47.119.170.71:8088/send_work_result", parseResult, nil)
		if err != nil {
			log.Errorf("send_work_result err:%v", err)
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
