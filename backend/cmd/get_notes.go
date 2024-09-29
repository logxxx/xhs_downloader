package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/storage"
	"github.com/logxxx/xhs_downloader/model"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	cookie     = ""
	rawCookie2 = `a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; customer-sso-sid=68c51739881124315146420967fc9e9fcaf5e8c0; x-user-id-creator.xiaohongshu.com=61d13a62000000001000b704; customerClientId=585309193620957; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517398811247446431509vcoizkfl5iohtgtp; xsecappid=xhs-pc-web; webBuild=4.35.0; web_session=040069b0a5792a12e752d7b1c5344b7498bd20; unread={%22ub%22:%2266f04d37000000001201221c%22%2C%22ue%22:%2266f0c934000000000c01a94f%22%2C%22uc%22:28}; acw_tc=0037368b76eb084b57220da6544683b6c66ab107526d14e9f88dc22415626e2a; websectiga=2a3d3ea002e7d92b5c9743590ebd24010cf3710ff3af8029153751e41a6af4a3; sec_poison_id=70aa5d78-36b6-471c-805c-5185f5f12273`
	rawCookie  = `abRequestId=18d8450d-628a-5dcd-936d-01b5b40c8276; xsecappid=xhs-pc-web; a1=19179e61ba9rsl23v6l4my9i37wq8py2vjeo1rfl850000293818; webId=e9976c88abe83d72a6350bd21221909a; gid=yjyWjdKJ8DhfyjyWjdKyDVFi0jFM1JqhK146d1Yvj3qWEl28AYUvJy888JjqYyY8i4WJqjWf; webBuild=4.35.0; websectiga=29098a4cf41f76ee3f8db19051aaa60c0fc7c5e305572fec762da32d457d76ae; sec_poison_id=57e4b6a5-dba8-4c3a-a77f-7f521ba20348; acw_tc=2ed945131369672a44b2137468ace1efdf97f8c723fac726363060a0f692a992; unread={%22ub%22:%22646f35ff0000000013008d18%22%2C%22ue%22:%2264a276b50000000034015247%22%2C%22uc%22:26}; web_session=0400697999f01a77b8f86ff0c4344ba1154db9`
)

func changeCookie() {
	if cookie == rawCookie {
		cookie = rawCookie2
	} else {
		cookie = rawCookie
	}
}

func StartGetNotes() {

	cookie = rawCookie

	upers := getAllUpers()
	log.Printf("get %v upers", len(upers))

	continueNoNoteCount := 0
	downloadedCount := 0
	for i, u := range upers {

		if downloadedCount > 500 && i > 0 && i%10 == 0 {
			log.Printf("sleep for i%%10==0")
			time.Sleep(1 * time.Minute)
		}

		if downloadedCount > 500 && i > 0 && i%100 == 0 {
			log.Printf("change cookie for i%%100==0")
			changeCookie()
		}

		log.Printf("deal parseUper %v/%v %v", i+1, len(upers), u)

		uper := storage.GetStorage().GetUper(0, u)
		if len(uper.Notes) != 0 && len(uper.Notes) != 14 {
			continue
		}
		//if storage.GetStorage().IsUperScanned(u) {
		//	continue
		//}
		parseUper, parseNotes, err := getNotes(u, cookie)
		if err != nil {
			log.Printf("get parseNotes err:%v uid:%v", err, u)
			continue
		}
		log.Printf("parseUper [%v_%v] get [%v] parseNotes", parseUper.UID, parseUper.Name, len(parseNotes))

		if len(parseNotes) <= len(utils.RemoveDuplicate(uper.Notes)) {
			continue
		}

		if len(parseNotes) == 0 {
			continueNoNoteCount++
			changeCookie()
		} else {
			continueNoNoteCount = 0
			downloadedCount += len(parseNotes)
		}

		if continueNoNoteCount > 5 {
			for i := 0; i < 60; i++ {
				log.Printf("sleep %v/%v for continueNoNoteCount > 5", i+1, 60)
			}
			continueNoNoteCount = 0
		}

		modelUper := model.Uper{
			UID:              parseUper.UID,
			Name:             parseUper.Name,
			Area:             parseUper.Area,
			AvatarURL:        parseUper.AvatarURL,
			IsGirl:           parseUper.IsGirl,
			Desc:             parseUper.Desc,
			Tags:             parseUper.Tags,
			FansCount:        parseUper.FansCount,
			ReceiveLikeCount: parseUper.ReceiveLikeCount,
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
		}
		result, err := storage.GetStorage().InsertOrUpdateUper(modelUper)
		if err != nil {
			log.Printf("InsertOrUpdateUper err:%v parseUper:%+v", err, modelUper)
			continue
		}
		log.Printf("InsertOrUpdateUper succ:%+v result:%v", modelUper, result)
		storage.GetStorage().SetUperScanned(u)

		allParseNotes := []string{}
		for _, n := range parseNotes {
			allParseNotes = append(allParseNotes, n.NoteID)
		}
		failedReason, err := storage.GetStorage().UperAddNote(parseUper.UID, allParseNotes...)
		if err != nil {
			log.Printf("UperAddNote err:%v uid:%v noteid:%v", err, parseUper.UID, allParseNotes)
		} else if failedReason != "" {
			log.Printf("UperAddNote failed:%v uid:%v noteid:%v", failedReason, parseUper.UID, allParseNotes)
		} else {
			//log.Printf("UperAddNote succ. uid:%v noteid:%v", parseUper.UID, n.NoteID)
		}

		for _, n := range parseNotes {

			dbNote := model.Note{
				NoteID:    n.NoteID,
				UperUID:   u,
				Title:     n.Title,
				URL:       n.URL,
				PosterURL: n.Poster,
				LikeCount: n.LikeCount,
			}
			insertOrUpdate, err := storage.GetStorage().InsertOrUpdateNote(dbNote)
			if err != nil {
				log.Printf("InsertOrUpdateNote err:%v dbNote:%+v", err, dbNote)
				continue
			}
			_ = insertOrUpdate
			//log.Printf("InsertOrUpdateNote succ: %+v insertOrUpdate:%v", dbNote, insertOrUpdate)
		}
	}
}

func getAllUpers() []string {
	allProfiles := []string{}
	allProfilesMap := map[string]bool{}
	fileutil.ReadByLine("chore/upers.txt", func(s string) (e error) {
		if allProfilesMap[s] {
			return
		}
		if len(s) != 24 {
			return
		}
		allProfilesMap[s] = true
		allProfiles = append(allProfiles, s)
		return nil

	})
	return allProfiles
}

func getNotes(uid, cookie string) (uper ParseUper, notes []ParseNote, err error) {

	uperURL := fmt.Sprintf("https://www.xiaohongshu.com/user/profile/%v?channel_type=web_note_detail_r10&parent_page_channel_type=web_profile_board", uid)

	ctx, cancel := getCtxWithCancel()
	go func() {
		time.Sleep(300 * time.Second)
		cancel()
	}()
	defer cancel()

	ctx = context.WithValue(ctx, "XHS_COOKIE", cookie)

	chromedp.Run(ctx,
		chromedp.ActionFunc(setCookie),
		chromedp.Sleep(5*time.Second),
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

				chromedp.ScrollIntoView("document.querySelector('#userPostedFeeds').lastElementChild", chromedp.ByJSPath).Do(ctx)

				roundUper, roundNotes, err := ParseHtml(content)
				if err != nil {
					log.Printf("ParseHtml err:%v", err)
					return err
				}
				notesStr := utils.JsonToString(roundNotes)
				if lastRoundNotes == notesStr {
					sameCount++
					log.Printf("SAME %v", sameCount)
					if sameCount < 3 {
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
					notes = append(notes, n)
				}

				log.Printf("round %v get %v notes", round, len(roundNotes))

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

func getCtxWithCancel() (context.Context, func()) {
	var options []chromedp.ExecAllocatorOption
	options = append(options, chromedp.DisableGPU)
	options = append(options, chromedp.Flag("ignore-certificate-errors", true))
	options = append(options, chromedp.Flag("disable-web-security", true))
	//Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
	options = append(options, chromedp.Flag("blink-settings", "imagesEnabled=false"))
	//options = append(options, chromedp.Headless)
	actX, _ := chromedp.NewExecAllocator(context.Background(), options...)

	ctx, cancel := chromedp.NewContext(actX, chromedp.WithErrorf(func(s string, i ...interface{}) {
		return
	}))
	return ctx, cancel
}

func setCookie(ctx context.Context) error {

	rawCookie := ctx.Value("XHS_COOKIE")
	cookie, ok := rawCookie.(string)
	if !ok || cookie == "" {
		panic("no cookie")
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

func ParseHtml(content string) (uper ParseUper, notes []ParseNote, err error) {

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

	return
}
