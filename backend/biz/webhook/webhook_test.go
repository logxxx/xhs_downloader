package webhook

import (
	"context"
	"testing"
)

var (
	ctx = context.Background()
)

func TestDownload(t *testing.T) {
	resp, err := Download(ctx,
		"http://sns-video-bd.xhscdn.com/spectrum/1040g35831ahnihicmu705ok9ndo8cb6cf7ug3sg",
		"D:\\mytest\\mywork\\xhs_downloader\\backend\\biz\\webhook\\1\\2\\3\\test.mp4",
		false)
	if err != nil {
		panic(err)
	}
	t.Logf("resp:%+v", resp)
}

func TestListTasks(t *testing.T) {
	req := Input{
		Module:   "download",
		Function: "Infos",
		Args: map[string]interface{}{
			"task_id": "",
		},
	}
	resp := map[string]interface{}{}

	err := callWebhook(req, &resp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("resp:%+v", resp)
}
