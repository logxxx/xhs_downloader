package main

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/xhs_downloader/biz/thumb"
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
	reqURL := "https://www.xiaohongshu.com/explore/66c6da5e000000001d03a049?xsec_token=AB7zvzDsrdgf3TEetNvxrhUxVQsn3jX43PdsUQKtHZpT4=&xsec_source=pc_user"
	reqURL = "https://www.xiaohongshu.com/discovery/item/66c6da5e000000001d03a049?xsec_token=AB7zvzDsrdgf3TEetNvxrhUxVQsn3jX43PdsUQKtHZpT4=&xsec_source=pc_user"
	//reqURL = "https://www.xiaohongshu.com/discovery/item/66c6da5e000000001d03a049"
	reqURL = "https://www.xiaohongshu.com/explore/670e46760000000021002695?xsec_token=ABSwOPnSQyQzPoes8C28EX4-qxBEI8wTA5xQW3U24n0fQ=&xsec_source=pc_feed&source=404"
	reqURL = "https://www.xiaohongshu.com/explore/671092d5000000002401a6dc?xsec_token=AB3S2QG8dzwwSE7BDEXwRpnf8P_QE6AVkXyFwRJE9XRic=&xsec_source=pc_user"
	resp, err := ParseBlog(reqURL, rawCookie)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp:%+v", resp)

	Download(resp, "", false)
}

func TestGeneVideoShot(t *testing.T) {
	thumb.GeneVideoShot("N:\\output_bili\\395358743\\雅乐大人_BV1A2421o7PY_面对牢弟偷吃零食牢雅的惩罚是_1.mp4",
		"N:\\output_bili\\395358743\\雅乐大人_BV1A2421o7PY_面对牢弟偷吃零食牢雅的惩罚是_1.mp4.thumb.mp4")
}
