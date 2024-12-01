package mydp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/go-vgo/robotgo"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/utils/randutil"
	"github.com/logxxx/utils/runutil"
	"github.com/logxxx/xhs_downloader/biz/black"
	"github.com/logxxx/xhs_downloader/biz/blog"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/blog/blogutil"
	cookie2 "github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/storage"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"moul.io/http2curl"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetCtxWithCancel() (context.Context, func()) {
	var options []chromedp.ExecAllocatorOption
	//options = append(options, chromedp.DisableGPU)
	options = append(options, chromedp.Flag("ignore-certificate-errors", true))
	options = append(options, chromedp.Flag("disable-web-security", true))
	options = append(options, chromedp.Flag("enable-automation", false))                       //防止监测webdriver
	options = append(options, chromedp.Flag("disable-blink-features", "AutomationControlled")) //禁用blink特征

	//Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
	//options = append(options, chromedp.Flag("blink-settings", "imagesEnabled=false"))
	//options = append(options, chromedp.Headless)
	actX, _ := chromedp.NewExecAllocator(context.Background(), options...)

	ctx, cancel := chromedp.NewContext(actX, chromedp.WithErrorf(func(s string, i ...interface{}) {
		return
	}))
	return ctx, cancel
}

func SetCookie(ctx context.Context) error {

	rawCookie := ctx.Value("XHS_COOKIE")
	cookie, ok := rawCookie.(string)
	if !ok || cookie == "" {
		panic("empty cookie")
	}

	elems := strings.Split(cookie, "; ")
	for _, e := range elems {
		kv := strings.Split(e, "=")
		if len(kv) != 2 {
			continue
		}
		k := kv[0]
		v := kv[1]

		err := network.SetCookie(k, v).WithDomain(".xiaohongshu.com").
			//WithHTTPOnly(true).
			Do(ctx)
		if err != nil {
			log.Printf("network.SetCookie err:%v cookie:%v=%v", err, k, v)
		} else {
			//log.Printf("set cookie:%v=%v", k, v)
		}
	}
	return nil
}

func ConvImageUrlToHighQuality(imgURL string) (resp string) {

	resp = imgURL

	suffix := ""

	if strings.Contains(imgURL, "!nd_dft_wlteh_webp_3") {
		suffix = "!nd_dft_wlteh_webp_3"
	}

	if strings.Contains(imgURL, "!nd_dft_wgth_webp_3") {
		suffix = "!nd_dft_wgth_webp_3"
	}

	if strings.Contains(imgURL, "!nd_prv_wlteh_webp_3") {
		suffix = "!nd_prv_wlteh_webp_3"
	}

	if suffix == "" {
		return
	}

	startIdx := strings.LastIndex(imgURL, "/")
	if startIdx <= 0 {
		return
	}

	id := utils.Extract(imgURL[startIdx:], "/", suffix)

	return fmt.Sprintf("https://ci.xiaohongshu.com/%v?imageView2/2/w/format/png", id)

}

type GetNotes2Resp struct {
	NoteCount      int
	IsGalleryEmpty bool
	IsUperBanned   bool
	IsHitRisk      bool
	Records        []string
}

func convFeedResp2ParseResult(blogURL string, feedResp *blogmodel.FeedResp) (resp blogmodel.ParseBlogResp) {

	if len(feedResp.Data.Items) <= 0 {
		log.Printf("resp.Data.Items IS EMPTY")
		return
	}

	if reason := black.HitBlack(feedResp.Data.Items[0].NoteCard.Title, feedResp.Data.Items[0].NoteCard.Desc); reason != "" {
		log.Printf("HIT BLACK:%v", reason)
		return
	}

	parseResult := blogmodel.ParseBlogResp{
		Time:           time.Now().Format("20060102 15:04:05"),
		BlogURL:        blogURL,
		Author:         feedResp.Data.Items[0].NoteCard.User.Nickname,
		UserID:         feedResp.Data.Items[0].NoteCard.User.UserID,
		Title:          feedResp.Data.Items[0].NoteCard.Title,
		Content:        feedResp.Data.Items[0].NoteCard.Desc,
		NoteID:         feedResp.Data.Items[0].NoteCard.NoteID,
		LikeCount:      int(utils.ToI64(feedResp.Data.Items[0].NoteCard.InteractInfo.LikedCount)),
		NoteCreateTime: time.Unix(0, feedResp.Data.Items[0].NoteCard.Time),
	}

	for _, tag := range feedResp.Data.Items[0].NoteCard.TagList {
		parseResult.Tags = append(parseResult.Tags, tag.Name)
	}

	respBytesGrace := blogutil.JsonToStringGrace(feedResp)
	fileutil.WriteToFile([]byte(respBytesGrace), "feed_resp.json")
	for _, item := range feedResp.Data.Items {
		videoURL := ""
		if len(item.NoteCard.Video.Media.Stream.H264) > 0 {
			videoURL = item.NoteCard.Video.Media.Stream.H264[0].MasterURL
			if videoURL == "" && len(item.NoteCard.Video.Media.Stream.H264[0].BackupUrls) > 0 {
				videoURL = item.NoteCard.Video.Media.Stream.H264[0].BackupUrls[0]
			}
		}
		if videoURL == "" && len(item.NoteCard.Video.Media.Stream.H265) > 0 {
			videoURL = item.NoteCard.Video.Media.Stream.H265[0].MasterURL
			if videoURL == "" && len(item.NoteCard.Video.Media.Stream.H265[0].BackupUrls) > 0 {
				videoURL = item.NoteCard.Video.Media.Stream.H265[0].BackupUrls[0]
			}
		}
		origKey := item.NoteCard.Video.Consumer.OriginVideoKey
		if videoURL != "" && origKey != "" {

			videoURLObj, _ := url.Parse(videoURL)
			videoURL = strings.TrimSuffix(videoURL, videoURLObj.Path) + "/" + origKey

			media := blogmodel.Media{
				Type: "video",
				URL:  videoURL,
			}
			parseResult.Medias = append(parseResult.Medias, media)
		} else {
			for _, elem := range item.NoteCard.ImageList {

				media := blogmodel.Media{
					Type: "image",
					URL:  ConvImageUrlToHighQuality(elem.URLDefault),
				}
				parseResult.Medias = append(parseResult.Medias, media)

				if elem.LivePhoto {
					liveURL := ""
					if len(elem.Stream.H264) > 0 {
						if elem.Stream.H264[0].MasterURL != "" {
							liveURL = elem.Stream.H264[0].MasterURL
						} else if len(elem.Stream.H264[0].BackupUrls) > 0 {
							liveURL = elem.Stream.H264[0].BackupUrls[0]
						}
					}
					if liveURL == "" {
						if len(elem.Stream.H265) > 0 {
							liveURL = elem.Stream.H265[0].MasterURL
						}
					}

					if liveURL != "" {
						media := blogmodel.Media{
							Type: "live",
							URL:  liveURL,
						}
						parseResult.Medias = append(parseResult.Medias, media)
					}
				}
			}
		}
	}

	respBytesGrace = blogutil.JsonToStringGrace(parseResult)
	fileutil.WriteToFile([]byte(respBytesGrace), "parse_media.json")

	return parseResult
}

func GetNotes2(uid, cookie string, parseResultHandler func(parseUper ParseUper, parseResult blogmodel.ParseBlogResp)) (resp GetNotes2Resp, err error) {

	uperURL := fmt.Sprintf("https://www.xiaohongshu.com/user/profile/%v?channel_type=web_note_detail_r10&parent_page_channel_type=web_profile_board", uid)

	log.Infof("GetNotes2 uperURL:%v useCookie:%v", uperURL, cookie2.GetCookieName(cookie))

	ctx, cancel := GetCtxWithCancel()
	go func() {
		time.Sleep(3000 * time.Second)
		cancel()
	}()
	defer cancel()

	ctx = context.WithValue(ctx, "XHS_COOKIE", cookie)

	page := 1

	xsecToken := ""
	noteID := ""
	downloadFinishMsg := ""

	downloaded := map[string]bool{} //key: note_id
	uniqNoteIDMap := map[string]bool{}

	defer func() {
		resp.NoteCount = len(uniqNoteIDMap)
	}()

	parseUperInfo := ParseUper{}

	continueParseBlogFailedCount := 0

	// 创建一个chromedp的实例并设置监听网络请求的选项
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:

			if !strings.Contains(ev.Request.URL, "web/v1/feed") {
				return
			}

			defer func() {
				downloadFinishMsg = noteID
				log.Infof("ListenTarget for feed finish. set downloadFinishMsg = noteID = %v", downloadFinishMsg)
			}()

			reqContent := `{"source_note_id":"%v","image_formats":["jpg","webp","avif"],"extra":{"need_body_topic":"1"},"xsec_source":"pc_user","xsec_token":"%v"}`
			reqContent = fmt.Sprintf(reqContent, noteID, xsecToken)
			fileutil.WriteToFile([]byte(reqContent), "req_body.json")
			reqBuf := bytes.NewBufferString(reqContent)
			//log.Printf("START REQUEST FEED url:%v reqBody:%v", ev.Request.URL, reqContent)
			httpReq, _ := http.NewRequest("POST", ev.Request.URL, reqBuf)
			for k, v := range ev.Request.Headers {
				httpReq.Header.Set(k, fmt.Sprintf("%v", v))
			}
			httpReq.Header.Set("Cookie", cookie)

			curl, err := http2curl.GetCurlCommand(httpReq)
			if err == nil {
				fileutil.WriteToFile([]byte(curl.String()), "curl")
			}

			respCode, respBytes, err := netutil.HttpDo(httpReq)
			_ = respCode
			//log.Printf("HttpDo respCode:%v resp:%v err:%v", respCode, string(respBytes), err)

			feedResp := &blogmodel.FeedResp{}

			if strings.Contains(string(respBytes), "访问频次异常") {
				resp.IsHitRisk = true
				resp.Records = append(resp.Records, "访问频次异常")
				cancel()
			}

			json.Unmarshal(respBytes, feedResp)

			parseResult := convFeedResp2ParseResult(ev.Request.URL, feedResp)

			if len(parseResult.Medias) > 0 {
				continueParseBlogFailedCount = 0
			}

			reportContent := fmt.Sprintf(" t:%v like:%v title:%v blogURL:%v",
				time.Now().Format("01/02 15:04"), parseResult.LikeCount, parseResult.Title, parseResult.BlogURL)

			reportContent = fmt.Sprintf("feedApi进行下载(%v)", parseResult.GetMediaSimpleInfo()) + reportContent
			fileutil.AppendToFile("download_report.txt", reportContent)

			resp.Records = append(resp.Records, fmt.Sprintf("\t-%v noteID:%v media:%v scene:FeedApi", len(resp.Records)+1, noteID, parseResult.GetMediaSimpleInfo()))

			parseResultHandler(parseUperInfo, parseResult)

		case *network.EventResponseReceived:
			//if strings.Contains(ev.Response.URL, "web/v1/feed") {
			//	log.Printf("接口响应捕获: %v, RespHeaders:%+v ev.Response.RequestHeaders:%+v", ev.Response.URL, ev.Response.Headers, ev.Response.RequestHeaders)
			//}

		}
	})

	err = chromedp.Run(ctx,
		chromedp.ActionFunc(SetCookie),
		chromedp.Navigate(uperURL),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {

			if parseUperInfo.UID == "" {
				content := ""
				chromedp.InnerHTML(`document.querySelector('html')`, &content, chromedp.ByJSPath).Do(ctx)

				//接相关投诉该账户违反
				if strings.Contains(content, "接相关投诉该账户违反") {
					log.Infof("******************** 违规账号 ********************")
					resp.IsGalleryEmpty = true
					resp.Records = append(resp.Records, "违规账号")
					return nil
				}

				if strings.Contains(content, "没有发布任何内容") {
					log.Infof("******************** 没有发布任何内容 ********************")
					resp.IsGalleryEmpty = true
					resp.Records = append(resp.Records, "没有发布任何内容")
					return nil
				}

				parseUperInfo, _, _ = ParseHtml(uid, content)
			}

			robotgo.ScrollDir(1, "down")

			lastDownloadedNoteID := ""
			continueDownloadedTimes := 0
			lastPage := page
			lastFirstAndLastNoteID := ""
			for {

				if continueDownloadedTimes > 10 {
					log.Printf("RETURN: UPER RECENTLY HAS NO NEW NOTE")
					//resp.Records = append(resp.Records, "RETURN: UPER RECENTLY HAS NO NEW NOTE")
					//break
				}

				currNotes := []*cdp.Node{}
				err = chromedp.Run(ctx, chromedp.Nodes(".cover.ld.mask", &currNotes))
				if err != nil {
					log.Errorf("chromedp.Nodes .cover.ld.mask err:%v", err)
					resp.Records = append(resp.Records, fmt.Sprintf("chromedp.Nodes .cover.ld.mask err:%v", err))
					break
				}
				log.Printf("page:%v currNotes(%v):", page, len(currNotes))

				if len(currNotes) <= 0 {
					break
				}

				if lastPage != page {
					lastPage = page
					firstHref, _ := currNotes[0].Attribute("href")
					firstNoteID := utils.Extract(firstHref, fmt.Sprintf("/user/profile/%v/", uid), "?")
					log.Printf("--first: %v", firstNoteID)
					lastHref, _ := currNotes[len(currNotes)-1].Attribute("href")
					lastNoteID := utils.Extract(lastHref, fmt.Sprintf("/user/profile/%v/", uid), "?")
					log.Printf("--last: %v", lastNoteID)
					currFirstAndLastNoteID := firstNoteID + lastNoteID
					if currFirstAndLastNoteID == lastFirstAndLastNoteID {
						log.Infof("翻页后没有新元素，所以退出")
						break
					}
					lastFirstAndLastNoteID = currFirstAndLastNoteID
				}

				var currRoundNode *cdp.Node
				currRoundNodeIdx := -1

				if lastDownloadedNoteID == "" {
					currRoundNode = currNotes[0]
					currRoundNodeIdx = 0
				} else {
					for i, note := range currNotes {
						href, _ := note.Attribute("href")

						xsecToken = utils.Extract(href, "xsec_token=", "&")
						noteID = utils.Extract(href, fmt.Sprintf("/user/profile/%v/", uid), "?")

						if noteID == lastDownloadedNoteID {
							if i < len(currNotes)-1 {
								currRoundNodeIdx = i + 1
								currRoundNode = currNotes[currRoundNodeIdx]
							}
							break
						}

					}
				}

				if currRoundNode == nil {
					log.Infof("往下翻页")
					err = chromedp.ScrollIntoView("document.querySelector('#userPostedFeeds').lastElementChild", chromedp.ByJSPath).Do(ctx)
					if err != nil {
						log.Errorf("ScrollIntoView err:%v", err)
						resp.Records = append(resp.Records, fmt.Sprintf("ScrollIntoView err:%v", err))
					}
					page++
					time.Sleep(3 * time.Second)
					continue
				}

				href, _ := currRoundNode.Attribute("href")
				xsecToken = utils.Extract(href, "xsec_token=", "&")
				noteID = utils.Extract(href, fmt.Sprintf("/user/profile/%v/", uid), "?")
				lastDownloadedNoteID = noteID
				log.Infof("get curr note: idx=%v href=%v", currRoundNodeIdx, noteID)

				uniqNoteIDMap[noteID] = true

				// --------------- some mime start -------------------

				title := ""
				poster := ""
				likeCountStr := ""

				runutil.GoRunSafe(func() {
					selector := fmt.Sprintf(`.note-item:nth-child(%v) .title`, currRoundNodeIdx+1)
					err = chromedp.InnerHTML(selector, &title, chromedp.ByQueryAll).Do(ctx)
					if err != nil {
						log.Errorf("get title err:%v sel:%v", err, selector)
					}

					if title != "" {
						title = utils.Extract(title, ">", "<")
					}
					log.Infof("Get title:%v", title)
				})

				runutil.GoRunSafe(func() {
					imgNodes := []*cdp.Node{}
					selector := fmt.Sprintf(".note-item:nth-child(%v) img", currRoundNodeIdx+1)
					err = chromedp.Run(ctx, chromedp.Nodes(selector, &imgNodes))
					if err != nil {
						log.Errorf("get imgNodes err:%v sel:%v", err, selector)
					}

					for _, img := range imgNodes {
						src, _ := img.Attribute("src")
						if strings.Contains(src, "sns-webpic-qc") {
							poster = src
							break
						}
					}
					log.Infof("Get poster:%v", poster)
				})

				runutil.GoRunSafe(func() {
					selector := fmt.Sprintf(`.note-item:nth-child(%v) .count`, currRoundNodeIdx+1)
					err = chromedp.InnerHTML(selector, &likeCountStr, chromedp.ByQueryAll).Do(ctx)
					if err != nil {
						log.Errorf("get likeCountStr err:%v sel:%v", err, selector)
					}
					log.Infof("Get likeCountStr:%v", likeCountStr)
				})

				time.Sleep(2 * time.Second)

				reportContent := fmt.Sprintf(" t:%v like:%v title:%v blogURL:%v poster:%v\n",
					time.Now().Format("01/02 15:04"), likeCountStr, title, "https://www.xiaohongshu.com/"+href, poster)

				likeCount, _ := strconv.Atoi(likeCountStr)

				if reason := black.HitBlack(title, href); reason != "" {
					reportContent = fmt.Sprintf("跳过下载(命中黑字:%v)", reason) + reportContent
					log.Infof(reportContent)
					fileutil.AppendToFile("download_report.txt", reportContent)
					continue
				}

				//len(uniqNoteIDMap) > 2: 最近的可能还没收到点赞。姑且先下了
				if len(uniqNoteIDMap) > 2 && !strings.Contains(likeCountStr, "万") && likeCount < 10 && !black.IsWhite(title) {
					reportContent = "跳过下载(点赞太少)" + reportContent
					log.Infof(reportContent)
					fileutil.AppendToFile("download_report.txt", reportContent)
					continue
				}

				// --------------- some mime end -------------------

				dbNote := storage.GetStorage().GetNote(noteID)
				if dbNote.IsDownloaded() {
					log.Printf("NOTE(%v/%v) DB DOWNLOADED:%v %v %v", currRoundNodeIdx+1, noteID, len(currNotes), dbNote.ID, dbNote.Title)
					continueDownloadedTimes++
					reportContent = "跳过下载(已下载过)" + reportContent
					fileutil.AppendToFile("download_report.txt", reportContent)
					continue
				}

				if downloaded[noteID] {
					log.Printf("NOTE(%v/%v) MEM DOWNLOADED:%v", currRoundNodeIdx+1, len(currNotes), noteID)
					continue
				}

				continueDownloadedTimes = 0

				//如果能通过html拿到，就不要触发feed接口了
				blogURL := fmt.Sprintf("https://www.xiaohongshu.com/explore/%v?xsec_token=%v&xsec_source=pc_feed", noteID, xsecToken)
				parseResp, err := blog.ParseBlog(blogURL, cookie)
				if err == nil && len(parseResp.Medias) > 0 {
					log.Infof("ParseBlog SUCC. len(media):%v", len(parseResp.Medias))
					parseResultHandler(ParseUper{
						Name:             parseResp.Author,
						Desc:             "",
						UID:              parseResp.UserID,
						Area:             "",
						IsGirl:           false,
						FansCount:        0,
						ReceiveLikeCount: 0,
						AvatarURL:        "",
						Tags:             nil,
					}, parseResp)
					time.Sleep(1 * time.Second)
					isRemote := ""
					if parseResp.IsFromRemote {
						isRemote = "_REMOTE"
					}
					resp.Records = append(resp.Records, fmt.Sprintf("\t-%v noteID:%v media:%v scene:blog.ParseBlog%v", len(resp.Records)+1, noteID, parseResp.GetMediaSimpleInfo(), isRemote))

					reportContent = fmt.Sprintf("html进行下载(%v)[cookie=%v]", parseResp.GetMediaSimpleInfo(), parseResp.UseCookie) + reportContent
					fileutil.AppendToFile("download_report.txt", reportContent)
					log.Infof("html解析成功: %v %v", noteID, parseResp.GetMediaSimpleInfo())
					continueParseBlogFailedCount = 0
					continue
				} else {
					log.Infof("html解析失败:%v", noteID)
					continueParseBlogFailedCount++
					if continueParseBlogFailedCount > 5 {
						reportContent = fmt.Sprintf("解析失败次数过多 %v\n", uid)
						fileutil.AppendToFile("download_report.txt", reportContent)
						log.Infof("解析失败次数过多")
						break
					}
					//continue //风控，不能往下走了
				}

				log.Printf("extract noteID:%v xsec_token:%v", noteID, xsecToken)

				//************************** hook start *****************************
				evalResp := map[string]interface{}{}
				errMsg := chromedp.Evaluate(`window._webmsxyw("/api/sns/web/v1/feed",{"source_note_id":"6746f2560000000007032141","image_formats":["jpg","webp","avif"],"extra":{"need_body_topic":"1"},"xsec_source":"pc_user","xsec_token":"ABlDmvLp27Z6py2807XtRA-O6QMnIu6ZSyTXwSukgS1uo="})`, &evalResp).Do(ctx).Error()
				log.Printf("evalResp:%+v errMsg:%v", evalResp, errMsg)
				time.Sleep(2 * time.Second)
				os.Exit(1)
				//************************** hook end *****************************

				downloadFinishMsg = ""

				time.Sleep(1 * time.Second)
				target := fmt.Sprintf("document.querySelectorAll('.note-item')[%v]", currRoundNodeIdx)
				log.Printf("start click target:%v", target)
				err = chromedp.Click(target, chromedp.ByJSPath).Do(ctx)
				if err != nil {
					log.Printf("click target err:%v target:%v", err, target)
					resp.Records = append(resp.Records, fmt.Sprintf("click target err:%v target:%v", err, target))
					return err
				}
				log.Printf("finish click target:%v", target)

				//执行click后，会触发上面的Listener,监听到feed接口调用，并执行下载;那边下载完后，会设置isDownloading=false

				//err = chromedp.Click(`document.querySelector('.close-circle')`, chromedp.ByJSPath).Do(ctx)
				//if err != nil {
				//	log.Printf("click close-circle err:%v", err)
				//}

				isDownloadingRound := 0

				for {
					time.Sleep(1 * time.Second)
					isDownloadingRound++

					if isDownloadingRound%3 == 0 {
						log.Printf("press right")
						robotgo.KeyDown("right")
					} else {
						log.Printf("press left")
						robotgo.KeyDown("left")
					}

					log.Printf("is downloading [%v] %v...", noteID, isDownloadingRound)
					if isDownloadingRound > 30 {
						log.Printf("wait for download finish TIMEOUT")
						resp.Records = append(resp.Records, fmt.Sprintf("wait for download finish TIMEOUT:%v", noteID))
						return nil
					}
					if downloadFinishMsg != "" {
						log.Printf("wait for download finish succ!round=%v noteID=%v", isDownloadingRound, downloadFinishMsg)
						break
					}
				}

				err = chromedp.KeyEvent(kb.Escape).Do(ctx)
				if err != nil {
					log.Printf("KeyEvent(kb.Escape) err:%v", err)
					resp.Records = append(resp.Records, fmt.Sprintf("KeyEvent(kb.Escape) err:%v", err))
					return err
				}

				time.Sleep(3 * time.Second)

			}

			return nil
		}),
	)

	return
}

func GetNotes(uid, cookie string, pageCount int) (uper ParseUper, notes []ParseNote, err error) {

	uperURL := fmt.Sprintf("https://www.xiaohongshu.com/user/profile/%v?channel_type=web_note_detail_r10&parent_page_channel_type=web_profile_board", uid)

	ctx, cancel := GetCtxWithCancel()
	go func() {
		time.Sleep(300 * time.Second)
		cancel()
	}()
	defer cancel()

	ctx = context.WithValue(ctx, "XHS_COOKIE", cookie)

	err = chromedp.Run(ctx,
		chromedp.ActionFunc(SetCookie),
		chromedp.Sleep(1*time.Second),
		chromedp.Navigate(uperURL),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {

			round := 0
			lastRoundNotes := ""

			sameCount := 0
			for {
				round++
				chromedp.Sleep(time.Second * time.Duration((2 + rand.Intn(3)))).Do(ctx)
				content := ""
				//chromedp.InnerHTML(`document.querySelector('div.feeds-tab-container')`, &content, chromedp.ByJSPath).Do(ctx)
				chromedp.InnerHTML(`document.querySelector('html')`, &content, chromedp.ByJSPath).Do(ctx)

				if strings.Contains(content, "接相关投诉该账户违反") || strings.Contains(content, "还没有发布任何内容") {
					return nil
				}

				if strings.Contains(content, "访问频次异常") {
					return errors.New("need change account")
				}

				chromedp.ScrollIntoView("document.querySelector('#userPostedFeeds').lastElementChild", chromedp.ByJSPath).Do(ctx)

				roundUper, roundNotes, err := ParseHtml(uid, content)
				if err != nil {
					log.Printf("ParseHtml err:%v", err)
					return err
				}
				notesStr := utils.JsonToString(roundNotes)
				if lastRoundNotes == notesStr {
					sameCount++
					log.Printf("SAME %v works:%v", sameCount, len(roundNotes))
					if len(lastRoundNotes) != 14 {
						break
					}
					if sameCount < 5 {
						continue
					}
					log.Printf("loop %v break because notesStr is same.", round)
					//log.Printf("1:%v", lastRoundNotes)
					//log.Printf("2:%v", notesStr)
					break
				}
				lastRoundNotes = notesStr

				if uper.Name == "" {
					uper = roundUper
				}

				for _, n := range roundNotes {
					has := false
					for _, old := range notes {
						if old.NoteID == n.NoteID {
							has = true
							break
						}
					}
					if has {
						continue
					}
					if black.HitBlack(n.Title, n.URL) != "" {
						continue
					}
					notes = append(notes, n)
				}

				log.Printf("round %v get %v notes", round, len(roundNotes))

				if pageCount > 0 && round > pageCount {
					break
				}

			}

			return nil
		}),
	)

	//fmt.Printf("get %v explores\n", len(allExplores))
	//for i := range allExplores {
	//	fmt.Printf("%v: %v\n", i+1, allExplores[i])
	//}

	uper.UID = uid

	return
}

type ParseNote struct {
	NoteID    string
	Title     string
	URL       string
	Poster    string
	LikeCount int
}

type ParseUper struct {
	Name             string
	Desc             string
	UID              string
	Area             string
	IsGirl           bool
	FansCount        int
	ReceiveLikeCount int
	AvatarURL        string
	Tags             []string
}

func ScanMyFav(cookie string, pageCount int) (upers []string, err error) {

	logger := log.WithField("trace_id", randutil.RandStr(8))

	favURL := "https://www.xiaohongshu.com/user/profile/61d13a62000000001000b704?tab=liked"

	ctx, cancel := GetCtxWithCancel()
	defer cancel()

	lastUpers := []string{}

	fileName := fmt.Sprintf("my_fav_%v.txt", time.Now().Format("20060102_150405"))

	ctx = context.WithValue(ctx, "XHS_COOKIE", cookie)
	worksDeepEqualCount := 0
	chromedp.Run(ctx,
		chromedp.ActionFunc(func(c context.Context) error {

			headers := network.Headers{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
				"accept-language": "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7",
				"referer":         "https://www.xiaohongshu.com/",
				"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			}
			return network.SetExtraHTTPHeaders(headers).Do(c)
		}),
		chromedp.ActionFunc(SetCookie),
		chromedp.Navigate(favURL),
		chromedp.Sleep(3*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {

			round := 0

			for {
				round++
				log.Printf("ScanMyFav round %v", round)

				if pageCount >= 0 && round > pageCount {
					return nil
				}

				time.Sleep(time.Duration((1 + rand.Intn(2))) * time.Second)
				content := ""
				chromedp.InnerHTML(`document.querySelectorAll('.tab-content-item')[2]`, &content, chromedp.ByJSPath).Do(ctx)

				log.Printf("scrolling page %v...", round+1)
				//err := chromedp.ScrollIntoView("document.querySelectorAll('.tab-content-item')[1].lastElementChild", chromedp.ByJSPath).Do(ctx)
				err := chromedp.ScrollIntoView("document.querySelectorAll('.note-item')[document.querySelectorAll('.note-item').length-1]", chromedp.ByJSPath).Do(ctx)
				if err != nil {
					panic(err)
				}

				pageUpers := utils.ExtractAll(content, `href="/user/profile/`, `?`, false)
				if reflect.DeepEqual(pageUpers, lastUpers) {
					worksDeepEqualCount++
					log.Printf("works deep equal[%v] %v", len(pageUpers), worksDeepEqualCount)
					if worksDeepEqualCount > 10 {
						return nil
					}

				}
				lastUpers = pageUpers

				newCount := 0
				for _, p := range pageUpers {
					if len(p) != 24 {
						continue
					}
					if utils.Contains(p, upers) {
						continue
					}
					newCount++
					upers = append(upers, p)
				}

				logger.Printf("round %v get newUper(%v/%v)", round, len(upers), newCount)

				fileutil.WriteToFile([]byte(strings.Join(upers, "\n")), fileName)
			}

			return nil
		}),
	)

	return
}

func ScanMyShoucang(cookie string, pageCount int) (upers, works []string, err error) {

	logger := log.WithField("trace_id", randutil.RandStr(8))

	shoucangURL := "https://www.xiaohongshu.com/user/profile/61d13a62000000001000b704?tab=fav&subTab=note"

	ctx, cancel := GetCtxWithCancel()
	defer cancel()

	lastUpers := []string{}

	ctx = context.WithValue(ctx, "XHS_COOKIE", cookie)
	worksDeepEqualCount := 0
	chromedp.Run(ctx,
		chromedp.ActionFunc(func(c context.Context) error {

			headers := network.Headers{
				"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
				"accept-language": "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7",
				"referer":         "https://www.xiaohongshu.com/explore/67212cf1000000001d03ab89?xsec_token=ABd59NCaXTdJOzaocEwHW1rzCoKqDQUZrelRtfQKyoLn0=&xsec_source=pc_feed",
				"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			}
			return network.SetExtraHTTPHeaders(headers).Do(c)
		}),
		chromedp.ActionFunc(SetCookie),
		chromedp.Navigate(shoucangURL),

		chromedp.ActionFunc(func(ctx context.Context) error {

			round := 0

			for {
				round++
				log.Printf("ScanMyShoucang round %v", round)

				if pageCount >= 0 && round > pageCount {
					return nil
				}

				time.Sleep(time.Duration((1 + rand.Intn(3))) * time.Second)
				content := ""
				chromedp.InnerHTML(`document.querySelectorAll('.tab-content-item')[1]`, &content, chromedp.ByJSPath).Do(ctx)

				log.Printf("scrolling page %v...", round+1)
				//err := chromedp.ScrollIntoView("document.querySelectorAll('.tab-content-item')[1].lastElementChild", chromedp.ByJSPath).Do(ctx)
				err := chromedp.ScrollIntoView("document.querySelectorAll('.note-item')[document.querySelectorAll('.note-item').length-1]", chromedp.ByJSPath).Do(ctx)
				if err != nil {
					panic(err)
				}

				pageUpers := utils.ExtractAll(content, `href="/user/profile/`, `?`, false)
				if reflect.DeepEqual(pageUpers, lastUpers) {
					worksDeepEqualCount++
					log.Printf("works deep equal[%v] %v", len(pageUpers), worksDeepEqualCount)
					if worksDeepEqualCount > 10 {
						return nil
					}

				}
				lastUpers = pageUpers

				newCount := 0
				for _, p := range pageUpers {
					if len(p) != 24 {
						continue
					}
					if utils.Contains(p, upers) {
						continue
					}
					newCount++
					upers = append(upers, p)
				}

				pageWorks := utils.ExtractAll(content, `target="_self" href="/user/profile/`, `xsec_source=pc_user"`, false)

				newWorkCount := 0
				for _, w := range pageWorks {
					if utils.Contains(w, works) {
						continue
					}
					newWorkCount++
					works = append(works, strings.ReplaceAll("https://www.xiaohongshu.com/user/profile/"+w+"xsec_source=pc_user", "&amp;", "&"))
				}

				logger.Printf("round %v get newUper(%v/%v) newWork(%v/%v)", round, len(upers), newCount, newWorkCount, len(works))

			}

			return nil
		}),
	)

	return
}

func ParseHtml(uid string, content string) (uper ParseUper, notes []ParseNote, err error) {

	d, err := goquery.NewDocumentFromReader(bytes.NewBufferString(content))
	if err != nil {
		return
	}

	d.Find("div").Each(func(i int, div *goquery.Selection) {

		if div.HasClass("tag-item") {
			uper.Tags = append(uper.Tags, div.Text())
		}

		if div.HasClass("user-name") {
			uper.Name = div.Text()
		}

		if div.HasClass("user-desc") {
			uper.Desc = div.Text()
		}

	})

	d.Find("use").Each(func(i int, use *goquery.Selection) {
		href, _ := use.Attr("href")
		//log.Printf("use.href:%v", href)
		if href == "#female" {
			uper.IsGirl = true
		}
	})

	tmpCount := 0
	d.Find("span").Each(func(i int, span *goquery.Selection) {

		if span.HasClass("user-IP") {
			uper.Area = strings.ReplaceAll(span.Text(), " IP属地：", "")
		}

		if span.HasClass("count") {
			tmpCount = int(utils.ToI64(span.Text()))
		}
		if span.HasClass("shows") {
			switch span.Text() {
			case "粉丝":
				uper.FansCount = tmpCount
			case "获赞与收藏":
				uper.ReceiveLikeCount = tmpCount
			}
		}
	})

	d.Find("img").Each(func(i int, img *goquery.Selection) {
		if img.HasClass("user-image") {
			src, _ := img.Attr("src")
			if src != "" {
				src = utils.Extract(src, "", "?")
				if src != "" {
					uper.AvatarURL = src
					//log.Printf("set user avatar url:%v", uper.AvatarURL)
				}

			}

		}
	})

	d.Find("section").Each(func(i int, s *goquery.Selection) {

		note := ParseNote{}

		//feeds-tab-container
		if !s.HasClass("note-item") {
			return
		}

		imgSrc, _ := s.Find("img").Attr("src")
		//log.Printf("封面:%v", imgSrc)
		note.Poster = imgSrc

		//log.Printf("i:%v s:%+v", i, s.Text())

		s.Find("a").Each(func(i int, a *goquery.Selection) {

			if note.Title == "" {
				if a.HasClass("title") {
					title := a.Find("span").Text()
					//log.Printf("find title:%v", title)
					note.Title = title
				}
			}

			if note.URL == "" {
				href, _ := a.Attr("href")
				if strings.Contains(href, "xsec_token") {
					note.URL = href
					elems := strings.Split(note.URL, "/")
					if len(elems) > 4 {
						note.NoteID = utils.Extract(elems[4], "", "?")
					}
				}
			}

		})

		s.Find("span").Each(func(i int, span *goquery.Selection) {
			if span.HasClass("count") {
				note.LikeCount = int(utils.ToI64(span.Text()))
			}
		})

		notes = append(notes, note)
	})

	//log.Printf("uper:%+v notes:%v", uper, len(notes))
	//for i, n := range notes {
	//	log.Printf("%v: %+v", i+1, n)
	//}

	uper.UID = uid

	return
}

func GetHtmlByChromedp(reqURL, cookie string) (resp []byte) {
	ctx, cancel := GetCtxWithCancel()
	go func() {
		time.Sleep(300 * time.Second)
		cancel()
	}()
	defer cancel()

	ctx = context.WithValue(ctx, "XHS_COOKIE", cookie)

	content := ""
	chromedp.Run(ctx,
		chromedp.ActionFunc(SetCookie),
		chromedp.Navigate(reqURL),
		chromedp.Sleep(10*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {

			//chromedp.InnerHTML(`document.querySelector('div.feeds-tab-container')`, &content, chromedp.ByJSPath).Do(ctx)
			chromedp.InnerHTML(`document.querySelector('html')`, &content, chromedp.ByJSPath).Do(ctx)

			return nil
		}),
	)
	return []byte(content)
}
