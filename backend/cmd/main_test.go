package main

import (
	"fmt"
	"github.com/logxxx/utils"
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

func Extract(content string, begin string, end string) (resp string) {
	beginIdx := strings.Index(content, begin)
	if beginIdx < 0 {
		return ""
	}

	if end == "" {
		return content[beginIdx+len(begin):]
	}

	endIdx := strings.Index(content[beginIdx+len(begin):], end)
	if endIdx < 0 {
		return ""
	}

	resp = content[beginIdx+len(begin) : beginIdx+len(begin)+endIdx]

	return resp
}
