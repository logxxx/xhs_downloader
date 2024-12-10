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

func TestFixFailedVideo(t *testing.T) {
	FixFailedVideo()
}

func TestExtract(t *testing.T) {
	data, err := os.ReadFile("test_live.html")
	if err != nil {
		t.Fatal(err)
	}

	content := utils.Extract(string(data), "window.__INITIAL_STATE__=", "</script></body></html>")
	t.Logf("content:%v", content)
}

func TestParseBlog4(t *testing.T) {
	reqURL := `
https://www.xiaohongshu.com/explore/674c67c8000000000703bfa4?xsec_token=ABQHtq6kGx6nsesfTwiV81A-I_DZph8-0JAaLpth98H_Q=&xsec_source=pc_feed
`

	elems := strings.Split(reqURL, "\n")

	log.Printf("get %v elems", len(elems))

	for _, e := range elems {
		if e == "" {
			continue
		}
		resp, err := blog.ParseBlog(e, cookie.GetCookie2())
		if err != nil {
			t.Logf("err:%+v", err)
			//t.Fatal(err)
			continue
		}
		t.Logf("ParseBlog resp:%+v", resp)

		download.Download("TestParseBlog4", resp, "", true, true)
	}

}

func TestGeneVideoShot(t *testing.T) {
	thumb.GeneVideoShot("E:\\test\\1.mp4",
		"E:\\test\\1.mp4.thumb.mp4")
}

func TestScanMyShoucang(t *testing.T) {
	upers, works, _ := mydp.ScanMyShoucang(cookie.GetCookie(), 1)
	log.Printf("upers(%v):%v \n works(%v)", len(upers), upers, len(works))
	for i, w := range works {
		log.Printf("%v: %v", i, w)
	}
}

func TestExtractUIDByURL(t *testing.T) {
	resp := ExtractUIDByURL("https://www.xiaohongshu.com/user/profile/5c26e25b0000000006012115/66878434000000001e013b34?xsec_token=ABxNqPSOjYrmndfjK5aHZSpCbnjomEoNZY_0KSEG1F9SM=&xsec_source=pc_user")
	t.Logf("uid:%v", resp)
}

func TestConvImageUrlToHighQuality(t *testing.T) {
	resp := mydp.ConvImageUrlToHighQuality("http://sns-webpic-qc.xhscdn.com/202411231845/0b2dd039a9b292a197ba1af2f8d5b653/1040g2sg31ae01hvrna0g5n2uu3lk0lr4toa9bv0!nd_prv_wlteh_webp_3")
	t.Logf("resp: %v", resp)
}
