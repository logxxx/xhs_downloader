package mydp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/logxxx/utils"
	"github.com/logxxx/xhs_downloader/biz/black"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

func GetCtxWithCancel() (context.Context, func()) {
	var options []chromedp.ExecAllocatorOption
	options = append(options, chromedp.DisableGPU)
	options = append(options, chromedp.Flag("ignore-certificate-errors", true))
	options = append(options, chromedp.Flag("disable-web-security", true))
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

func GetNotes(uid, cookie string, onlyOnePage bool) (uper ParseUper, notes []ParseNote, err error) {

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

				roundUper, roundNotes, err := ParseHtml(content)
				if err != nil {
					log.Printf("ParseHtml err:%v", err)
					return err
				}
				notesStr := utils.JsonToString(roundNotes)
				if lastRoundNotes == notesStr {
					sameCount++
					log.Printf("SAME %v", sameCount)
					if len(lastRoundNotes) != 14 {
						break
					}
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
					if black.HitBlack(n.Title, n.URL) != "" {
						continue
					}
					notes = append(notes, n)
				}

				log.Printf("round %v get %v notes", round, len(roundNotes))

				if onlyOnePage {
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

func ScanMyShoucang(cookie string, pageCount int) (upers, works []string, err error) {

	shoucangURL := "https://www.xiaohongshu.com/user/profile/61d13a62000000001000b704?tab=fav&subTab=note"

	ctx, cancel := GetCtxWithCancel()
	defer cancel()

	lastUpers := []string{}

	ctx = context.WithValue(ctx, "XHS_COOKIE", cookie)

	chromedp.Run(ctx,
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

				chromedp.Sleep(1 * time.Second).Do(ctx)
				content := ""
				chromedp.InnerHTML(`document.querySelectorAll('.tab-content-item')[1]`, &content, chromedp.ByJSPath).Do(ctx)

				err := chromedp.ScrollIntoView("document.querySelectorAll('.tab-content-item')[1].lastElementChild", chromedp.ByJSPath).Do(ctx)
				if err != nil {
					panic(err)
				}

				pageUpers := utils.ExtractAll(content, `href="/user/profile/`, `?`, false)
				if reflect.DeepEqual(pageUpers, lastUpers) {
					log.Printf("deep equal:%v", pageUpers)
					return nil
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

				log.Printf("round %v get %v newUper(%v)", round, newCount, len(upers))

				pageWorks := utils.ExtractAll(content, `target="_self" href="/user/profile/`, `xsec_source=pc_user"`, false)

				newWorkCount := 0
				for _, w := range pageWorks {
					if utils.Contains(w, works) {
						continue
					}
					newWorkCount++
					works = append(works, strings.ReplaceAll("https://www.xiaohongshu.com/user/profile/"+w+"xsec_source=pc_user", "&amp;", "&"))
				}

				log.Printf("round %v get %v newWork(%v)", round, newWorkCount, len(works))

			}

			return nil
		}),
	)

	return
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
