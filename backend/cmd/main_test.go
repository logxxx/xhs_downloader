package main

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/xhs_downloader/biz/blog"
	"github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/download"
	"github.com/logxxx/xhs_downloader/biz/mydp"
	"github.com/logxxx/xhs_downloader/biz/thumb"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"
)

func TestParseBlog(t *testing.T) {

	type M struct {
		Type string
		URL  string
	}

	m := M{
		Type: "image",
		URL:  "http://sns-webpic-qc.xhscdn.com/202409221656/e777717b87c6367eecc1c022237418ad/1000g0082dikqm1egs0305p09ih2k4lj22u59lug!nd_dft_wlteh_webp_3",
	}

	if m.Type == "image" && strings.Contains(m.URL, "!nd_dft_wlteh_webp_3") {
		startIdx := strings.LastIndex(m.URL, "/")
		if startIdx <= 0 {
			t.Fatal()
		}
		id := utils.Extract(m.URL[startIdx:], "/", "!nd_dft_wlteh_webp_3")
		m.URL = fmt.Sprintf("https://ci.xiaohongshu.com/%v?imageView2/2/w/format/png", id)
	}

	if m.URL != "https://ci.xiaohongshu.com/1000g0082dikqm1egs0305p09ih2k4lj22u59lug?imageView2/2/w/format/png" {
		t.Fatal(m.URL)
	}
}

func TestExtract(t *testing.T) {
	data, err := os.ReadFile("test_live.html")
	if err != nil {
		t.Fatal(err)
	}

	content := utils.Extract(string(data), "window.__INITIAL_STATE__=", "</script></body></html>")
	t.Logf("content:%v", content)
}

func TestParseBlog2(t *testing.T) {
	reqURL := `https://www.xiaohongshu.com/explore/67190e600000000024017615?xsec_token=ABoFitQaIom1egZSk8FNBYh8loEv-WWS29fw8fj2cUcyU=&xsec_source=pc_user`

	elems := strings.Split(reqURL, "\n")

	for _, e := range elems {
		resp, err := blog.ParseBlog(e, cookie.GetCookie())
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("resp:%+v", resp)

		download.Download(resp, "", false)
	}

}

func TestGeneVideoShot(t *testing.T) {
	thumb.GeneVideoShot("N:\\output_bili\\395358743\\雅乐大人_BV1A2421o7PY_面对牢弟偷吃零食牢雅的惩罚是_1.mp4",
		"N:\\output_bili\\395358743\\雅乐大人_BV1A2421o7PY_面对牢弟偷吃零食牢雅的惩罚是_1.mp4.thumb.mp4")
}

func TestScanMyShoucang(t *testing.T) {
	upers, works, _ := mydp.ScanMyShoucang(cookie.GetCookie(), 1)
	log.Printf("upers(%v):%v \n works(%v)", len(upers), upers, len(works))
	for i, w := range works {
		log.Printf("%v: %v", i, w)
	}
}
