package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Media struct {
	Type         string `json:"type,omitempty"`
	URL          string `json:"url,omitempty"`
	DownloadPath string `json:"download_path,omitempty"`
}

type ParseBlogResp struct {
	Time    string   `json:"time,omitempty"`
	BlogURL string   `json:"blog_url,omitempty"`
	Author  string   `json:"author,omitempty"`
	UserID  string   `json:"user_id,omitempty"`
	Title   string   `json:"title,omitempty"`
	Content string   `json:"content,omitempty"`
	Medias  []*Media `json:"medias,omitempty"`
	NoteID  string   `json:"note_id,omitempty"`
}

func ParseBlog(reqURL string) (resp ParseBlogResp, err error) {

	resp.BlogURL = reqURL
	resp.Time = time.Now().Format("20060102 15:04:05")

	httpReq := getHttpReq(reqURL, "", "")
	code, httpResp, err := netutil.HttpDo(httpReq)
	if err != nil {
		return
	}
	if code != 200 {
		err = fmt.Errorf("invalid code:%v", code)
		return
	}

	//fileutil.WriteToFile(httpResp, fmt.Sprintf("test_live_%v.html", time.Now().Format("20060102_150405")))
	fileutil.WriteToFile(httpResp, fmt.Sprintf("test_live.html"))

	content := utils.Extract(string(httpResp), "window.__INITIAL_STATE__=", "</script></body></html>")
	if content == "" {
		log.Infof("ParseBlog Extract empty! resp:%v", string(httpResp))
		err = errors.New("parse failed")
		return
	}
	//log.Infof("content:%v", content)

	content = strings.ReplaceAll(content, "undefined", `null`)

	noteResp := &NoteResp{}
	err = json.Unmarshal([]byte(content), noteResp)
	if err != nil {
		log.Infof("ParseBlog Unmarshal err:%v data:%v", err, content)
		return
	}
	//log.Infof("NoteResp:%+v", noteResp)

	for _, noteDetail := range noteResp.Note.NoteDetailMap {
		//log.Infof(">>>>>>>>>>> note:%+v", noteDetail.Note)
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
				resp.Medias = append(resp.Medias, &Media{
					Type: "image",
					URL:  imgInfo.URL,
				})
			}

			for _, elem := range imgInfo.InfoList {
				//log.Infof("elem%v:%+v", i+1, elem)
				if (elem.ImageScene == "CRD_WM_JPG" || elem.ImageScene == "WB_DFT") && elem.URL != "" {
					log.Infof("find img:%v", elem.URL)
					resp.Medias = append(resp.Medias, &Media{
						Type: "image",
						URL:  elem.URL,
					})
				}
			}
		}

		for _, m := range resp.Medias {
			if m.Type == "image" && strings.Contains(m.URL, "!nd_dft_wlteh_webp_3") {
				startIdx := strings.LastIndex(m.URL, "/")
				if startIdx <= 0 {
					continue
				}
				id := utils.Extract(m.URL[startIdx:], "/", "!nd_dft_wlteh_webp_3")
				m.URL = fmt.Sprintf("https://ci.xiaohongshu.com/%v?imageView2/2/w/format/png", id)
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
				log.Infof("FIND LIVE!!!!:%v", live)
				resp.Medias = append(resp.Medias, &Media{
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
		resp.Medias = append(resp.Medias, &Media{
			Type: "video",
			URL:  videoURL,
		})

	}

	if resp.Title == "" && resp.Content != "" {
		resp.Title = resp.Content
	}

	movieMedias := []*Media{}
	for _, m := range resp.Medias {
		if m.Type == "video" || m.Type == "live" {
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

func getHttpReq(reqURL string, cookie, xs string) (resp *http.Request) {

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authority", "www.xiaohongshu.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	//req.Header.Set("Cookie", "acw_tc=43dee22040a3a93a7b8f4c694b6f716dc60537b360be07b4ddd6d9f99b630c51; abRequestId=ad5fe3e5-add1-56e7-ac15-27afc1bf6251; webBuild=4.6.0; a1=18e4653d9eegkvo56f6buobnslh00ww0eh749peu650000298220; webId=319268cd5a2e38ff03d9fb61e8327559; web_session=030037a2c6008139c934b9128d224ada1de7d5; gid=yYd4K2qd8fd0yYd4K2qfjxUdddkThl2KiKD7W7KlDIM1x42888EE8j888JjYJJ88WKqDfSd4; websectiga=16f444b9ff5e3d7e258b5f7674489196303a0b160e16647c6c2b4dcb609f4134; sec_poison_id=2fc2009d-0d03-4640-81e8-ff57f44ce7a7; xsecappid=xhs-pc-web")
	//req.Header.Set("Cookie", "acw_tc=148ac47105c4e8d751a7bad32e1b81c4fe837e9935724d59d339cb6e664df2f2; a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; webBuild=4.27.7; web_session=040069b0a5792a12e7525e7690344b620c9270; xsecappid=login; websectiga=8886be45f388a1ee7bf611a69f3e174cae48f1ea02c0f8ec3256031b8be9c7ee; sec_poison_id=677b63c4-6474-4807-b88a-f658344f4542; unread={%22ub%22:%2266a4e3ce000000000d031eca%22%2C%22ue%22:%2266a2641c0000000005004446%22%2C%22uc%22:39}")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://www.xiaohongshu.com/web-login/captcha?redirectPath=https%3A%2F%2Fwww.xiaohongshu.com%2Fexplore%2F65ea72b00000000003036e39&callFrom=web&biz=sns_web&verifyUuid=4167e15f-dc20-47f5-9da7-9699d0137505*XaiGvPwp&verifyType=102&verifyBiz=461")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Google Chrome\";v=\"122\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("X-S", xs)
	return req
}

type NoteResp struct {
	User struct {
		LoggedIn  bool `json:"loggedIn"`
		Activated bool `json:"activated"`
		UserInfo  struct {
		} `json:"userInfo"`
		Follow       []any `json:"follow"`
		UserPageData struct {
		} `json:"userPageData"`
		ActiveTabKey           int     `json:"activeTabKey"`
		Notes                  [][]any `json:"notes"`
		IsFetchingNotes        []bool  `json:"isFetchingNotes"`
		TabScrollTop           []int   `json:"tabScrollTop"`
		UserFetchingStatus     any     `json:"userFetchingStatus"`
		UserNoteFetchingStatus any     `json:"userNoteFetchingStatus"`
		BannedInfo             struct {
			Code      int    `json:"code"`
			ShowAlert bool   `json:"showAlert"`
			Reason    string `json:"reason"`
		} `json:"bannedInfo"`
		FirstFetchNote bool `json:"firstFetchNote"`
		NoteQueries    []struct {
			Num     int    `json:"num"`
			Cursor  string `json:"cursor"`
			UserID  string `json:"userId"`
			HasMore bool   `json:"hasMore"`
		} `json:"noteQueries"`
	} `json:"user"`
	Note struct {
		PrevRouteData struct {
		} `json:"prevRouteData"`
		PrevRoute     string `json:"prevRoute"`
		CommentTarget struct {
		} `json:"commentTarget"`
		IsImgFullscreen bool   `json:"isImgFullscreen"`
		GotoPage        string `json:"gotoPage"`
		FirstNoteID     string `json:"firstNoteId"`
		AutoOpenNote    bool   `json:"autoOpenNote"`
		TopCommentID    string `json:"topCommentId"`
		NoteDetailMap   map[string]struct {
			Comments struct {
				List               []any  `json:"list"`
				Cursor             string `json:"cursor"`
				HasMore            bool   `json:"hasMore"`
				Loading            bool   `json:"loading"`
				FirstRequestFinish bool   `json:"firstRequestFinish"`
			} `json:"comments"`
			CurrentTime int64 `json:"currentTime"`
			Note        struct {
				User struct {
					Avatar   string `json:"avatar"`
					UserID   string `json:"userId"`
					Nickname string `json:"nickname"`
				} `json:"user"`
				InteractInfo struct {
					CommentCount   string `json:"commentCount"`
					ShareCount     string `json:"shareCount"`
					Followed       bool   `json:"followed"`
					Relation       string `json:"relation"`
					Liked          bool   `json:"liked"`
					LikedCount     string `json:"likedCount"`
					Collected      bool   `json:"collected"`
					CollectedCount string `json:"collectedCount"`
				} `json:"interactInfo"`
				ImageList []struct {
					URL      string `json:"url"`
					TraceID  string `json:"traceId"`
					InfoList []struct {
						ImageScene string `json:"imageScene"`
						URL        string `json:"url"`
					} `json:"infoList"`
					FileID string `json:"fileId"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
					Stream struct {
						H264 []struct {
							StreamDesc string `json:"streamDesc"`
							//Ssim          int      `json:"ssim"`
							Width         int      `json:"width"`
							Duration      int      `json:"duration"`
							VideoBitrate  int      `json:"videoBitrate"`
							StreamType    int      `json:"streamType"`
							VideoCodec    string   `json:"videoCodec"`
							DefaultStream int      `json:"defaultStream"`
							AudioDuration int      `json:"audioDuration"`
							Rotate        int      `json:"rotate"`
							BackupUrls    []string `json:"backupUrls"`
							HdrType       int      `json:"hdrType"`
							Psnr          int      `json:"psnr"`
							QualityType   string   `json:"qualityType"`
							Weight        int      `json:"weight"`
							Format        string   `json:"format"`
							Size          int      `json:"size"`
							AvgBitrate    int      `json:"avgBitrate"`
							Vmaf          int      `json:"vmaf"`
							MasterURL     string   `json:"masterUrl"`
							Height        int      `json:"height"`
							Volume        int      `json:"volume"`
							VideoDuration int      `json:"videoDuration"`
							AudioCodec    string   `json:"audioCodec"`
							AudioChannels int      `json:"audioChannels"`
							Fps           int      `json:"fps"`
							AudioBitrate  int      `json:"audioBitrate"`
						} `json:"h264"`
						H265 []any `json:"h265"`
						Av1  []any `json:"av1"`
					} `json:"stream"`
				} `json:"imageList"`
				Video struct {
					Image struct {
						ThumbnailFileid  string `json:"thumbnailFileid"`
						FirstFrameFileid string `json:"firstFrameFileid"`
					} `json:"image"`
					Capa struct {
						Duration int `json:"duration"`
					} `json:"capa"`
					Consumer struct {
						OriginVideoKey string `json:"originVideoKey"`
					} `json:"consumer"`
					Media struct {
						Stream struct {
							H264 []struct {
								StreamDesc string `json:"streamDesc"`
								//Ssim          int      `json:"ssim"`
								Width         int      `json:"width"`
								Duration      int      `json:"duration"`
								VideoBitrate  int      `json:"videoBitrate"`
								StreamType    int      `json:"streamType"`
								VideoCodec    string   `json:"videoCodec"`
								DefaultStream int      `json:"defaultStream"`
								AudioDuration int      `json:"audioDuration"`
								Rotate        int      `json:"rotate"`
								BackupUrls    []string `json:"backupUrls"`
								HdrType       int      `json:"hdrType"`
								Psnr          int      `json:"psnr"`
								QualityType   string   `json:"qualityType"`
								Weight        int      `json:"weight"`
								Format        string   `json:"format"`
								Size          int      `json:"size"`
								AvgBitrate    int      `json:"avgBitrate"`
								Vmaf          int      `json:"vmaf"`
								MasterURL     string   `json:"masterUrl"`
								Height        int      `json:"height"`
								Volume        int      `json:"volume"`
								VideoDuration int      `json:"videoDuration"`
								AudioCodec    string   `json:"audioCodec"`
								AudioChannels int      `json:"audioChannels"`
								Fps           int      `json:"fps"`
								AudioBitrate  int      `json:"audioBitrate"`
							} `json:"h264"`
							H265 []any `json:"h265"`
							Av1  []any `json:"av1"`
						} `json:"stream"`
						VideoID int64 `json:"videoId"`
						Video   struct {
							Duration    int    `json:"duration"`
							Md5         string `json:"md5"`
							HdrType     int    `json:"hdrType"`
							DrmType     int    `json:"drmType"`
							StreamTypes []int  `json:"streamTypes"`
							BizName     int    `json:"bizName"`
							BizID       string `json:"bizId"`
						} `json:"video"`
					} `json:"media"`
				} `json:"video"`
				Time       int64  `json:"time"`
				IPLocation string `json:"ipLocation"`
				NoteID     string `json:"noteId"`
				Type       string `json:"type"`
				Desc       string `json:"desc"`
				AtUserList []any  `json:"atUserList"`
				ShareInfo  struct {
					UnShare bool `json:"unShare"`
				} `json:"shareInfo"`
				Title   string `json:"title"`
				TagList []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
				} `json:"tagList"`
				LastUpdateTime int64 `json:"lastUpdateTime"`
			} `json:"note"`
		} `json:"noteDetailMap"`
		ServerRequestInfo struct {
			State     string `json:"state"`
			ErrorCode int    `json:"errorCode"`
			ErrMsg    string `json:"errMsg"`
		} `json:"serverRequestInfo"`
		Volume            int `json:"volume"`
		RecommendVideoMap struct {
		} `json:"recommendVideoMap"`
		VideoFeedType  string `json:"videoFeedType"`
		Rate           int    `json:"rate"`
		NoteFromSource string `json:"noteFromSource"`
	} `json:"note"`
}
