package main

import (
	"github.com/logxxx/utils/runutil"
	"github.com/logxxx/xhs_downloader/biz/web"
	"time"
)

func main() {
	//runutil.GoRunSafe(StartGetNotes)
	runutil.GoRunSafe(StartDownloadNote)

	//runutil.GoRunSafe(CrontabDownloadUperAvatar)

	//runutil.GoRunSafe(DownloadNotePoster)

	//runutil.GoRunSafe(StartDownloadNote)

	runutil.GoRunSafe(web.InitWeb)

	for {
		time.Sleep(10 * time.Second)
	}
}
