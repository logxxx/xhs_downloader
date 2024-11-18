package thumb

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/ffmpeg"
	"github.com/logxxx/utils/fileutil"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func MakeThumb(downloadPath string) {
	filePath := downloadPath

	fileSize := utils.GetFileSize(filePath)
	if fileSize <= 2*1024*1024 {
		return
	}

	fileType := "video"
	if IsImage(filePath) {
		fileType = "image"
	}

	if fileType == "video" {
		makeThumbCore(filePath)
	} else {
		cmd := `ffmpeg -i %v -y -vf scale=480:-1 %v`
		ffp := ffmpeg.FFProbe("ffprobe")
		fInfo, _ := ffp.NewVideoFile(filePath)
		if fInfo != nil && fInfo.Width < fInfo.Height {
			cmd = `ffmpeg -i %v -y -vf scale=-1:480 %v`
		}
		cmd = fmt.Sprintf(cmd, filePath, filepath.Join(filepath.Dir(filePath), ".thumb", filepath.Base(filePath)))

		if !utils.HasFile(filepath.Join(filepath.Dir(filePath), ".thumb")) {
			os.MkdirAll(filepath.Join(filepath.Dir(filePath), ".thumb"), 0755)
		}
		runCommand(cmd)
	}

}

func makeThumbCore(filePath string) error {

	log.Printf("make thumb %v", filePath)
	_, err := ffmpeg.GenePreviewVideoSlice(filePath, func(vInfo *ffmpeg.VideoFile) ffmpeg.GenePreviewVideoSliceOpt {

		segNum := 3
		segDur := 3
		skipStart := 1
		skipEnd := 3
		if vInfo.Duration < 15 {
			segNum = 2
		}
		if vInfo.Duration > 2*60 {
			skipStart = 10
			skipEnd = 30
		}
		if vInfo.Duration > 5*60 {
			skipStart = 30
			skipEnd = 30
			segNum = 5
		}

		return ffmpeg.GenePreviewVideoSliceOpt{
			//ToPath:      filePath + ".thumb.mp4",
			ToPath:      filepath.Join(filepath.Dir(filePath), ".thumb", filepath.Base(filePath)),
			SegNum:      segNum,
			SegDuration: segDur,
			SkipStart:   skipStart,
			SkipEnd:     skipEnd,
		}
	})
	return err
}

func runCommand(command string) (output []byte, err error) {
	//log.Printf("runCommand:%v", command)
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("runCommand err:%v output:%v", err, string(output))
		return
	}
	//log.Printf("runCommand output:%v", string(output))
	return
}

var AllowImageExts = []string{".bmp", ".jpg", ".jpeg", ".webp", ".png", ".tif", ".gif", ".pcx", ".tga", ".exif", ".fpx", ".svg", ".psd", ".cdr", ".pcd", ".dxf", ".ufo", ".eps", ".ai", ".raw", ".WMF", ".webp", ".avif", ".apng"}

func IsImage(fileName string) bool {
	if fileutil.IsDir(fileName) {
		return false
	}
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)

	return utils.Contains(ext, AllowImageExts)
}

func GenePreviewVideo(filePath string, toPath string) error {

	fpb := ffmpeg.FFProbe("ffprobe")
	vInfo, err := fpb.NewVideoFile(filePath)
	if err != nil {
		log.Errorf("GenePreviewVideo NewVideoFile err:%v", err)
		return err
	}
	//log.Infof("height:%v width:%v", vInfo.Height, vInfo.Width)

	os.MkdirAll(filepath.Dir(toPath), 0755)

	height := vInfo.Height
	width := vInfo.Width

	min := 640
	if vInfo.Height > vInfo.Width { //竖屏

		for {
			if height <= min {
				break
			}
			height /= 2
			width /= 2
		}

	} else {
		for {
			if width <= min {
				break
			}
			height /= 2
			width /= 2
		}
	}

	if width%2 != 0 {
		width -= 1
	}

	if height%2 != 0 {
		height -= 1
	}

	scale := fmt.Sprintf("%v:%v", width, height)

	command := `ffmpeg -y -i %v -ss 00:00:05 -to 10 -vf scale=%v -pix_fmt yuv420p -level 4.2 -crf 30 -threads 8 -strict -2 %v`
	command = fmt.Sprintf(command, filePath, scale, toPath)
	output, err := runCommand(command)
	log.Debugf("GenePreviewVideo command:%v output:%v err:%v", command, string(output), err)
	if err != nil {
		return err
	}
	return nil
}

func GeneVideoShot(filePath string, toPath string) error {

	fpb := ffmpeg.FFProbe("ffprobe")
	vInfo, err := fpb.NewVideoFile(filePath)
	if err != nil {
		log.Errorf("GenePreviewVideo NewVideoFile err:%v", err)
		return err
	}
	//log.Infof("height:%v width:%v", vInfo.Height, vInfo.Width)

	os.MkdirAll(filepath.Dir(toPath), 0755)

	height := vInfo.Height
	width := vInfo.Width

	min := 640
	if vInfo.Height > vInfo.Width { //竖屏

		for {
			if height <= min {
				break
			}
			height /= 2
			width /= 2
		}

	} else {
		for {
			if width <= min {
				break
			}
			height /= 2
			width /= 2
		}
	}

	if width%2 != 0 {
		width -= 1
	}

	if height%2 != 0 {
		height -= 1
	}

	scale := fmt.Sprintf("%v:%v", width, height)

	command := `ffmpeg -y -i %v -ss 00:00:05 -vf scale=%v -pix_fmt yuv420p -level 4.2 -crf 30 -threads 8 -strict -2 -frames:v 120 %v`
	command = fmt.Sprintf(command, filePath, scale, toPath)
	output, err := runCommand(command)
	log.Debugf("GeneVideoShot command:%v output:%v err:%v", command, string(output), err)
	if err != nil {
		return err
	}
	return nil
}
