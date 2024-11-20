package main

import (
	"github.com/logxxx/utils/runutil"
	"github.com/logxxx/xhs_downloader/biz/black"
	"github.com/logxxx/xhs_downloader/biz/web"
	"time"
)

func main() {

	black.Init("chore/black.txt", "chore/white.txt")

	//runutil.GoRunSafe(StartGetNotes)

	//runutil.GoRunSafe(StartDownloadRecrentlyNotes)
	//runutil.GoRunSafe(StartDownloadNote)

	runutil.GoRunSafe(StartScanMyShoucang)

	//runutil.GoRunSafe(FixFailedVideo)

	//runutil.GoRunSafe(StartFillFileSize)

	//runutil.GoRunSafe(CrontabDownloadUperAvatar)

	//runutil.GoRunSafe(DownloadNotePoster)

	runutil.GoRunSafe(web.InitWeb)

	for {
		time.Sleep(10 * time.Second)
	}
}
