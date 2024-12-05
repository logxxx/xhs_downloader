package blog

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/utils/randutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/blog/blogutil"
	cookie2 "github.com/logxxx/xhs_downloader/biz/cookie"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func GetHtmlByApi(reqURL, cookie string) (resp []byte) {
	httpReq := getHttpReq(reqURL, cookie, "")
	code, httpResp, err := netutil.HttpDo(httpReq)
	if err != nil {
		log.Errorf("HttpDo err:%v", err)
		return
	}
	if code != 200 {
		log.Errorf("HttpDo invalid code:%v", code)
		err = fmt.Errorf("invalid code:%v", code)
		return
	}

	return httpResp
}

func ParseBlog(reqURL, cookie string) (resp blogmodel.ParseBlogResp, err error) {

	//log.Printf("start ParseBlog:%v", reqURL)

	defer func() {
		//log.Printf("finish ParseBlog:%v", reqURL)
	}()

	if !strings.HasPrefix(reqURL, "https:") {
		reqURL = "https://www.xiaohongshu.com" + reqURL
	}

	defer func() {
		imgCount := 0
		videoCount := 0
		liveCount := 0
		for _, m := range resp.Medias {
			if m.Type == "image" {
				imgCount++
			}
			if m.Type == "live" {
				liveCount++
			}
			if m.Type == "video" {
				videoCount++
			}
		}
		//log.Infof("ParseBlog url:%v get %vI%vV%vL total:%v *** useCookie:%v ***", reqURL, imgCount, videoCount, liveCount, imgCount+videoCount+liveCount, cookie2.GetCookieName(resp.UseCookie))
	}()

	remoteURL := fmt.Sprintf("http://47.119.170.71:8088/parse_blog?blog_url=%v&trace_id=%v", reqURL, randutil.RandStr(8))
	//remoteURL := fmt.Sprintf("http://127.0.0.1:8088/parse_blog?blog_url=%v&trace_id=%v", reqURL, randutil.RandStr(8))
	//log.Printf("remoteURL:%v", remoteURL)
	remoteReq, _ := http.NewRequest("GET", remoteURL, nil)
	remoteReq.Header.Set("mycookie", cookie)
	code, err := netutil.HttpReqGet(remoteReq, &resp)
	if code == 200 {
		//log.Infof("****** GET BLOG INFO FROM REMOTE ******")
		resp.IsFromRemote = true
		return
	}
	log.Infof("GET BLOG INFO FROM REMOTE failed. code:%v err:%v", code, err)

	return ParseBlogCore(reqURL, cookie)

}

func ParseBlogCore(reqURL, cookie string) (resp blogmodel.ParseBlogResp, err error) {
	//log.Printf("Start PraseBolg:%v", reqURL)

	//httpResp := GetHtmlByChromedp(reqURL, cookie)
	httpResp := GetHtmlByApi(reqURL, cookie)
	//log.Printf("GetHtmlByApi finish")

	//fileutil.WriteToFile(httpResp, fmt.Sprintf("test_live_%v.html", time.Now().Format("20060102_150405")))
	fileutil.WriteToFile(httpResp, fmt.Sprintf("test_live.html"))

	//else if strings.Contains(string(httpResp), "你访问的页面不见了") {
	//			reason = "note disappear"
	//		}

	//log.Printf("here Extract finish")

	resp, err = blogutil.ParseNoteHTML(string(httpResp))
	if err != nil {
		return
	}
	if resp.FailedReason != "" {
		log.Infof("ParseBlog Failed! url:%v resp:%v Cookie:%v", reqURL, resp.FailedReason, cookie2.GetCookieName(cookie))
	}

	resp.BlogURL = reqURL
	resp.Time = time.Now().Format("20060102 15:04:05")

	return resp, nil
}

func getHttpReq(reqURL string, cookie, xs string) (resp *http.Request) {

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authority", "www.xiaohongshu.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	//req.Header.Set("Cookie", "acw_tc=43dee22040a3a93a7b8f4c694b6f716dc60537b360be07b4ddd6d9f99b630c51; abRequestId=ad5fe3e5-add1-56e7-ac15-27afc1bf6251; webBuild=4.6.0; a1=18e4653d9eegkvo56f6buobnslh00ww0eh749peu650000298220; webId=319268cd5a2e38ff03d9fb61e8327559; web_session=030037a2c6008139c934b9128d224ada1de7d5; gid=yYd4K2qd8fd0yYd4K2qfjxUdddkThl2KiKD7W7KlDIM1x42888EE8j888JjYJJ88WKqDfSd4; websectiga=16f444b9ff5e3d7e258b5f7674489196303a0b160e16647c6c2b4dcb609f4134; sec_poison_id=2fc2009d-0d03-4640-81e8-ff57f44ce7a7; xsecappid=xhs-pc-web")
	//req.Header.Set("Cookie", "acw_tc=148ac47105c4e8d751a7bad32e1b81c4fe837e9935724d59d339cb6e664df2f2; a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; webBuild=4.27.7; web_session=040069b0a5792a12e7525e7690344b620c9270; xsecappid=login; websectiga=8886be45f388a1ee7bf611a69f3e174cae48f1ea02c0f8ec3256031b8be9c7ee; sec_poison_id=677b63c4-6474-4807-b88a-f658344f4542; unread={%22ub%22:%2266a4e3ce000000000d031eca%22%2C%22ue%22:%2266a2641c0000000005004446%22%2C%22uc%22:39}")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://www.xiaohongshu.com/web-login/captcha?redirectPath=https%3A%2F%2Fwww.xiaohongshu.com%2Fexplore%2F65ea72b00000000003036e39&callFrom=web&biz=sns_web&verifyUuid=4167e15f-dc20-47f5-9da7-9699d0137505*XaiGvPwp&verifyType=102&verifyBiz=461")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Google Chrome\";v=\"122\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("X-S", xs)
	return req
}
