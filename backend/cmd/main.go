package cmd

import "time"

func main() {
	StartDownload()

	InitWeb()

	for {
		time.Sleep(10 * time.Second)
	}
}
