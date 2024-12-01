package blogutil

import (
	"encoding/json"
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

func ParseNoteHTML(req string) (resp blogmodel.ParseBlogResp, err error) {

	content := req

	if strings.Contains(content, "你访问的页面不见了") {
		resp.IsNoteDisappeared = true
		return
	}

	//ursor":""}}}}</script><div

	content = utils.Extract(req, "window.__INITIAL_STATE__=", "</script></body></html>")
	if content == "" {
		content = utils.Extract(req, "window.__INITIAL_STATE__=", "</script><div")
	}
	if content == "" {
		reason := ""
		if strings.Contains(content, "您当前系统版本过低，请升级后再试") {
			reason = "您当前系统版本过低，请升级后再试"
		} else {
			reason = "DONT KNOW WHY:" + content
		}
		resp.FailedReason = reason
		//err = errors.New("parse failed")
		return
	}

	content = strings.ReplaceAll(content, "undefined", `null`)
	//log.Infof("content:%v", content)

	noteResp := &blogmodel.NoteResp{}
	noteRespForGallery := &blogmodel.NoteRespForWorkGallery{}
	//log.Printf("here Unmarshal start")
	if content != "" {
		err = json.Unmarshal([]byte(content), noteResp)
		//log.Printf("here Unmarshal end")
		if err != nil {
			log.Printf("ParseBlog Unmarshal noteResp err:%v data:%v", err, content)
			return
		}
		err = json.Unmarshal([]byte(content), noteRespForGallery)
		//log.Printf("here Unmarshal end")
		if err != nil {
			log.Printf("ParseBlog Unmarshal noteRespForGallery err:%v data:%v", err, content)
			return
		}
	}

	//log.Printf("noteResp:%+v", noteResp)
	fileutil.WriteToFile([]byte(JsonToStringGrace(noteRespForGallery)), "note_content.json")

	noteDetailCount := 0

	for _, notes := range noteRespForGallery.User.Notes {
		for i, note := range notes {
			fileutil.WriteToFile([]byte(JsonToStringGrace(note)), fmt.Sprintf("note%v.json", i+1))
		}
	}

	for _, noteDetail := range noteResp.Note.NoteDetailMap {
		noteDetailCount++
		//log.Infof(">>>>>>>>>>> note%v:%+v", noteDetailCount, noteDetail.Note)
		if resp.NoteID == "" && noteDetail.Note.NoteID != "" {
			resp.NoteID = noteDetail.Note.NoteID
		}
		if resp.Content == "" && noteDetail.Note.Desc != "" {
			resp.Content = noteDetail.Note.Desc
		}

		if noteDetail.Note.Title != "" && resp.Title == "" {

			resp.Title = noteDetail.Note.Title

		}
		if noteDetail.Note.User.Nickname != "" && resp.Author == "" {
			resp.Author = noteDetail.Note.User.Nickname
			resp.UserID = noteDetail.Note.User.UserID
		}
		for _, imgInfo := range noteDetail.Note.ImageList {
			if imgInfo.URL != "" {
				resp.Medias = append(resp.Medias, blogmodel.Media{
					Type: "image",
					URL:  imgInfo.URL,
				})
			}

			for _, elem := range imgInfo.InfoList {
				//log.Infof("elem%v:%+v", i+1, elem)
				if (elem.ImageScene == "CRD_WM_JPG" || elem.ImageScene == "WB_DFT") && elem.URL != "" {
					//log.Infof("find img:%v", elem.URL)
					resp.Medias = append(resp.Medias, blogmodel.Media{
						Type: "image",
						URL:  elem.URL,
					})
				}
			}
		}

		for i := range resp.Medias {
			m := &resp.Medias[i]
			if m.Type == "image" && (strings.Contains(m.URL, "!nd_dft_wlteh_webp_3") || strings.Contains(m.URL, "!nd_dft_wgth_webp_3")) {
				startIdx := strings.LastIndex(m.URL, "/")
				if startIdx <= 0 {
					log.Printf("get high image failed: startIdx <= 0:%v", m.URL)
					continue
				}
				id := ""
				if strings.Contains(m.URL, "!nd_dft_wlteh_webp_3") {
					id = utils.Extract(m.URL[startIdx:], "/", "!nd_dft_wlteh_webp_3")
				} else {
					id = utils.Extract(m.URL[startIdx:], "/", "!nd_dft_wgth_webp_3")
				}

				m.BackupURL = m.URL
				m.URL = fmt.Sprintf("https://ci.xiaohongshu.com/%v?imageView2/2/w/format/png", id)
				//log.Printf("set high img:%v", m.URL)
			} else {
				log.Printf("get high image failed: not contains nd_dft_wlteh_webp_3:%v", m.URL)
			}
		}

		masterURL := ""
		if len(noteDetail.Note.Video.Media.Stream.H264) > 0 {
			masterURL = noteDetail.Note.Video.Media.Stream.H264[0].MasterURL
		}

		for _, elem := range noteDetail.Note.ImageList {
			live := ""
			for _, h := range elem.Stream.H264 {
				if h.MasterURL != "" {
					live = h.MasterURL
					break
				}
				for _, u := range h.BackupUrls {
					if strings.Contains(u, ".mp4") || strings.Contains(u, ".mov") {
						if live == "" || strings.Contains(u, "sign") {
							live = u
						}
					}
				}
			}
			if live != "" {
				resp.Medias = append(resp.Medias, blogmodel.Media{
					Type: "live",
					URL:  live,
				})
			}
		}

		origKey := noteDetail.Note.Video.Consumer.OriginVideoKey

		if masterURL == "" || origKey == "" {
			continue
		}

		masterURLObj, _ := url.Parse(masterURL)
		videoURL := strings.TrimSuffix(masterURL, masterURLObj.Path) + "/" + origKey

		if strings.Contains(videoURL, "sns-video-qc.xhscdn.com") {
			log.Printf("XXXXXX FIND LOW QUALITY VIDEO:%v", videoURL)
			continue
		}

		resp.Medias = append(resp.Medias, blogmodel.Media{
			Type: "video",
			URL:  videoURL,
		})

	}

	log.Printf("get %v medias", len(resp.Medias))

	if resp.Title == "" && resp.Content != "" {
		resp.Title = resp.Content
	}

	movieMedias := []blogmodel.Media{}
	for _, m := range resp.Medias {
		if m.Type == "video" {
			movieMedias = append(movieMedias, m)
		}
	}
	if len(movieMedias) > 0 { //如果有了视频，则图片是视频封面，不需要
		resp.Medias = movieMedias
	}

	/*
		for _, url := range downloadURLs {
			code, imgData, err := netutil.HttpGetRaw(url.URL)
			if err != nil || code != 200 {
				continue
			}
			fileTitle := fmt.Sprintf("%v_%v", utils.ShortTitle(author), utils.ShortTitle(title))
			log.Infof("fileTitle:%v", fileTitle)
			suffix := ".jpg"
			if url.Type == "video" {
				suffix = ".mp4"
			}
			fileDir := fmt.Sprintf("output/%v", time.Now().Format("20060102"))
			fileName := fmt.Sprintf("%v%v", fileTitle, suffix)
			fileutil.WriteToFileWithRename(imgData, fileDir, fileName)
		}

	*/
	return
}

func JsonToStringGrace(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return ""
	}
	return string(jsonData)
}
