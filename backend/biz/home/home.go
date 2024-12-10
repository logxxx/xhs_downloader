package home

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/xhs_downloader/biz/black"
	"github.com/logxxx/xhs_downloader/biz/blog"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	cookie2 "github.com/logxxx/xhs_downloader/biz/cookie"
	"github.com/logxxx/xhs_downloader/biz/download"
	"github.com/logxxx/xhs_downloader/biz/queue"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	currCookie              = ""
	useCookieFailedLastTime = time.Time{}
)

func GetParseBlogCookie() string {
	return currCookie
}

func ChangeParseBlogCookie(reqCookie string, parseResult blogmodel.ParseBlogResp) {
	if len(parseResult.Medias) > 0 {
		return
	}

	//解析结果为空，是否要使用cookie3?
	if reqCookie == cookie2.GetCookie3() {
		//使用cookie3解析结果为空。
		currCookie = ""
		useCookieFailedLastTime = time.Now()
		return
	}

	//使用空cookie解析为空。

	if time.Since(useCookieFailedLastTime).Minutes() <= 30 {
		return
	}
	log.Infof("*** USE COOKIE3 TO PARSE MEDIA ***")
	useCookieFailedLastTime = time.Time{}

	currCookie = cookie2.GetCookie3()
}

func StartDownloadHome() {
	xs := `XYW_eyJzaWduU3ZuIjoiNTYiLCJzaWduVHlwZSI6IngyIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6IjdkZDRkMjY2YjFjZGFmNTJlOTZjZGM2ZTY3ZDVjZjE4NWI3OGJkMzdkZjk2ODUzMWFlMzIwMTg4MDNjZDcwOTM1MGY4NTA5MGU0ZDAxMjdjNjgwMzU5MzI1MmQ0MDZmZmU2MjAxOGZhZmFkNDhjYTU0ZWQxY2VhZWQ0YzQzNTA2YmQ0MzViYmIzNzdkMDU4ZWNhYTkwODNkMmQ4YTZlMWJiMmY2OTRiNGE2MDQ3ZWVmYzZjYjFhYmRlZGE1NDg3MjExMWFkOWE1NDA2NGEyYjI5ZTViMDdmMWFjZWE0MDJlMmQyNGQyN2M2ZmI0MmM2ZGEzY2Q2N2MyMzczNTY2ZDFjNTIyYjE5MWJjMmM3MTUwYmNlNDE0Y2NmZDYxYWFiZTAxMTg1MmIxZWY0MWY4MjNlN2MwZGQ0ZGFhYjBiOTk1ODc5M2ZlYWVmZmE5MWU3Y2ZkNGI4Nzg4OWJiZDFmYWNmYjAyYTNkODk3YzhmODFiM2Q2YTYxZGU2Y2NhNjQ1YTIzZDkwYjMwNDAyNjc0OWM1MWI3NDA4OGRkY2QxZjQ3In0=`

	logger := log.WithField("func_name", "StartDownloadHome")

	waitMinute := 10
	round := 0

	deadCookie := `a1=190f57a60ce1pzrfezgs740ln6bhaw5sew2wopupy50000121723; webId=8946bc0ba9fb796d38d7e710072b6e12; gid=yj8i2W0Wy8dYyj8i2W0K8EU7SdyUuFidukMWJUv481IKDE28x0E2Ml888yJyWJq8jfyWSKWW; abRequestId=8946bc0ba9fb796d38d7e710072b6e12; x-user-id-creator.xiaohongshu.com=61d13a62000000001000b704; customerClientId=585309193620957; customer-sso-sid=68c517443427088577997624430688006f2b05a8; access-token-creator.xiaohongshu.com=customer.creator.AT-68c517443427088576637418a24cacjkef53n74u; galaxy_creator_session_id=SBxEIH8svbyAoCqQ3n0NhVz1NWMTEBnVS593; galaxy.creator.beaker.session.id=1733057920148045297125; webBuild=4.46.0; acw_tc=0a4a2dad17335396757193034efa7d65f9fe4647fdcc0139bf6e7ec2bcc63e; xsecappid=xhs-pc-web; unread={%22ub%22:%226447b1bb000000001300f0c2%22%2C%22ue%22:%2264531fc10000000027028711%22%2C%22uc%22:25}; websectiga=82e85efc5500b609ac1166aaf086ff8aa4261153a448ef0be5b17417e4512f28; sec_poison_id=2091ef1e-0904-40cf-9cdf-248b61a92d15; web_session=040069b0a5792a12e7529ef266354b791309b4`

	for {
		round++
		logger = logger.WithField("round", round)
		if round != 1 {
			for i := 0; i < waitMinute; i++ {
				logger.Infof("waiting %v/%v...", i+1, waitMinute)
				time.Sleep(time.Minute)
			}
		}
		resp, err := blog.GetHomePage(deadCookie, xs)
		if err != nil {
			logger.Errorf("GetHomePage err:%v", err)
		}
		logger.Infof("get %v notes", len(resp))

		for i, e := range resp {

			if black := black.HitBlack(fmt.Sprintf("%v", e.Title), ""); black != "" {
				continue
			}

			logger = logger.WithField("note_idx", i)
			cookie := GetParseBlogCookie()
			parseResult, err := blog.ParseBlog(e.URL, cookie)
			ChangeParseBlogCookie(cookie, parseResult)
			if err != nil {
				logger.Errorf("ParseBlog err:%v note:%+v", err, e)
				queue.Push("home_parse_failed", e)
				continue
			}
			resp[i].MideaSimpleInfo = parseResult.GetMediaSimpleInfo()
			log.Infof("ParseBlog get %v medias. cookie:%v title:%v", parseResult.GetMediaSimpleInfo(), cookie2.GetCookieName(cookie), e.Title)
			if len(parseResult.Medias) <= 0 {
				queue.Push("home_parse_failed", e)
				continue
			}
			parseResult.LikeCount = e.LikeCount
			parseResult.Author = e.UperName
			parseResult.UserID = e.UperUID
			parseResult.Title = e.Title
			parseResult.BlogURL = e.URL

			downloadResp := download.DownloadToHome("StartDownloadHome", parseResult, "E:/xhs_downloader_output", true, false)

			download.UpdateDownloadRespToDB(model.Uper{
				UID:              parseResult.Uper.UID,
				Name:             parseResult.Uper.Name,
				Area:             parseResult.Uper.Area,
				AvatarURL:        parseResult.Uper.AvatarURL,
				IsGirl:           parseResult.Uper.IsGirl,
				Desc:             parseResult.Uper.Desc,
				HomeTags:         parseResult.Uper.Tags,
				FansCount:        parseResult.Uper.FansCount,
				ReceiveLikeCount: parseResult.Uper.ReceiveLikeCount,
			}, model.Note{
				NoteID:         parseResult.NoteID,
				URL:            parseResult.BlogURL,
				UperUID:        parseResult.UserID,
				Title:          parseResult.Title,
				Content:        parseResult.Content,
				DownloadTime:   time.Now(),
				LikeCount:      parseResult.LikeCount,
				Tags:           parseResult.Tags,
				WorkCreateTime: parseResult.NoteCreateTime,
			}, downloadResp)

		}

		if len(resp) > 0 {
			fileutil.WriteJsonToFile(resp, fmt.Sprintf("chore/home_page/%v.json", time.Now().Format("0102_150405")))
		}

	}

}
