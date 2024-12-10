package blogmodel

import (
	"fmt"
	"time"
)

type Media struct {
	Type         string `json:"type,omitempty"`
	URL          string `json:"url,omitempty"`
	DownloadPath string `json:"download_path,omitempty"` //只有下载到本地成功了，这个字段才有值！
}

type ParseBlogResp struct {
	Time              string  `json:"time,omitempty"`
	BlogURL           string  `json:"blog_url,omitempty"`
	Author            string  `json:"author,omitempty"`
	UserID            string  `json:"user_id,omitempty"`
	Title             string  `json:"title,omitempty"`
	Content           string  `json:"content,omitempty"`
	Medias            []Media `json:"medias,omitempty"`
	NoteID            string  `json:"note_id,omitempty"`
	LikeCount         int
	NoteCreateTime    time.Time
	Tags              []string
	IsNoteDisappeared bool
	FailedReason      string
	IsFromRemote      bool
	UseCookie         string
	Uper              ParseUper
}

func (resp *ParseBlogResp) GetMediaSimpleInfo() string {
	imgCount := 0
	liveCount := 0
	videoCount := 0
	for _, m := range resp.Medias {
		if m.Type == "image" {
			imgCount++
		}
		if m.Type == "live" {
			liveCount++
		}
		if m.Type == "video" {
			videoCount++
		}
	}
	return fmt.Sprintf("%v[%vI%vV%vL]", imgCount+videoCount+liveCount, imgCount, videoCount, liveCount)

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

// 从用户的作品列表点击作品，弹窗的结果
type NoteRespForWorkGallery2 struct {
	User struct {
		LoggedIn  bool `json:"loggedIn"`
		Activated bool `json:"activated"`
		UserInfo  struct {
			RedID    string `json:"red_id"`
			UserID   string `json:"user_id"`
			Nickname string `json:"nickname"`
			Desc     string `json:"desc"`
			Gender   int    `json:"gender"`
			Images   string `json:"images"`
			Imageb   string `json:"imageb"`
			Guest    bool   `json:"guest"`
			UserID0  string `json:"userId"`
			RedID0   string `json:"redId"`
		} `json:"userInfo"`
		Follow       []any `json:"follow"`
		UserPageData struct {
			TabPublic struct {
				Collection     bool `json:"collection"`
				CollectionNote struct {
					Lock    bool `json:"lock"`
					Count   int  `json:"count"`
					Display bool `json:"display"`
				} `json:"collectionNote"`
				CollectionBoard struct {
					Display bool `json:"display"`
					Lock    bool `json:"lock"`
					Count   int  `json:"count"`
				} `json:"collectionBoard"`
			} `json:"tabPublic"`
			ExtraInfo struct {
				Fstatus   string `json:"fstatus"`
				BlockType string `json:"blockType"`
			} `json:"extraInfo"`
			Result struct {
				Success bool   `json:"success"`
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"result"`
			BasicInfo struct {
				Images     string `json:"images"`
				RedID      string `json:"redId"`
				Gender     int    `json:"gender"`
				IPLocation string `json:"ipLocation"`
				Desc       string `json:"desc"`
				Imageb     string `json:"imageb"`
				Nickname   string `json:"nickname"`
			} `json:"basicInfo"`
			Interactions []struct {
				Type  string `json:"type"`
				Name  string `json:"name"`
				Count string `json:"count"`
			} `json:"interactions"`
			Tags []struct {
				Icon    string `json:"icon"`
				Name    string `json:"name"`
				TagType string `json:"tagType"`
			} `json:"tags"`
		} `json:"userPageData"`
		ActiveTab struct {
			Key      int    `json:"key"`
			Index    int    `json:"index"`
			Query    string `json:"query"`
			Label    string `json:"label"`
			Lock     bool   `json:"lock"`
			SubTabs  any    `json:"subTabs"`
			FeedType int    `json:"feedType"`
		} `json:"activeTab"`
		Notes [][]struct {
			ID       string `json:"id"`
			NoteCard struct {
				Type         string `json:"type"`
				DisplayTitle string `json:"displayTitle"`
				User         struct {
					NickName string `json:"nickName"`
					Avatar   string `json:"avatar"`
					UserID   string `json:"userId"`
					Nickname string `json:"nickname"`
				} `json:"user"`
				InteractInfo struct {
					LikedCount string `json:"likedCount"`
					Sticky     bool   `json:"sticky"`
					Liked      bool   `json:"liked"`
				} `json:"interactInfo"`
				Cover struct {
					TraceID  string `json:"traceId"`
					InfoList []struct {
						URL        string `json:"url"`
						ImageScene string `json:"imageScene"`
					} `json:"infoList"`
					URLPre     string `json:"urlPre"`
					URLDefault string `json:"urlDefault"`
					FileID     string `json:"fileId"`
					Height     int    `json:"height"`
					Width      int    `json:"width"`
					URL        string `json:"url"`
				} `json:"cover"`
				NoteID    string `json:"noteId"`
				XsecToken string `json:"xsecToken"`
			} `json:"noteCard"`
			Index       int    `json:"index"`
			Exposed     bool   `json:"exposed"`
			SsrRendered bool   `json:"ssrRendered"`
			XsecToken   string `json:"xsecToken"`
		} `json:"notes"`
		IsFetchingNotes        []bool   `json:"isFetchingNotes"`
		TabScrollTop           []int    `json:"tabScrollTop"`
		UserFetchingStatus     string   `json:"userFetchingStatus"`
		UserNoteFetchingStatus []string `json:"userNoteFetchingStatus"`
		BannedInfo             struct {
			UserID       string `json:"userId"`
			ServerBanned bool   `json:"serverBanned"`
			Code         int    `json:"code"`
			ShowAlert    bool   `json:"showAlert"`
			Reason       string `json:"reason"`
			API          string `json:"api"`
		} `json:"bannedInfo"`
		FirstFetchNote bool `json:"firstFetchNote"`
		NoteQueries    []struct {
			Num     int    `json:"num"`
			Cursor  string `json:"cursor"`
			UserID  string `json:"userId"`
			Page    int    `json:"page"`
			HasMore bool   `json:"hasMore"`
		} `json:"noteQueries"`
		PageScrolled int  `json:"pageScrolled"`
		ActiveSubTab any  `json:"activeSubTab"`
		IsOwnBoard   bool `json:"isOwnBoard"`
	} `json:"user"`
}
type NoteRespForWorkGallery struct {
	Global struct {
		AppSettings struct {
			NotificationInterval int  `json:"notificationInterval"`
			PrefetchTimeout      int  `json:"prefetchTimeout"`
			PrefetchRedisExpires int  `json:"prefetchRedisExpires"`
			RetryFeeds           bool `json:"retryFeeds"`
			GrayModeConfig       struct {
				Global    bool     `json:"global"`
				DateRange []string `json:"dateRange"`
				GreyRule  struct {
					Layout struct {
						Enable bool     `json:"enable"`
						Pages  []string `json:"pages"`
					} `json:"layout"`
					Pages []string `json:"pages"`
				} `json:"greyRule"`
				DisableLikeNotes  []string `json:"disableLikeNotes"`
				DisableSearchHint bool     `json:"disableSearchHint"`
			} `json:"grayModeConfig"`
			Nio         bool `json:"NIO"`
			ICPInfoList []struct {
				Label string `json:"label"`
				Link  string `json:"link,omitempty"`
				Title string `json:"title,omitempty"`
			} `json:"ICPInfoList"`
			DisableBanAlert string `json:"disableBanAlert"`
		} `json:"appSettings"`
		SupportWebp            bool   `json:"supportWebp"`
		ServerTime             int64  `json:"serverTime"`
		GrayMode               bool   `json:"grayMode"`
		Referer                string `json:"referer"`
		PwaAddDesktopPrompt    any    `json:"pwaAddDesktopPrompt"`
		FirstVisitURL          any    `json:"firstVisitUrl"`
		EasyAccessModalVisible struct {
			AddDesktopGuide bool `json:"addDesktopGuide"`
			CollectGuide    bool `json:"collectGuide"`
			KeyboardList    bool `json:"keyboardList"`
			MiniWindowGuide bool `json:"miniWindowGuide"`
		} `json:"easyAccessModalVisible"`
		CurrentLayout        string `json:"currentLayout"`
		FullscreenLocking    bool   `json:"fullscreenLocking"`
		FeedbackPopupVisible bool   `json:"feedbackPopupVisible"`
		TrackFps             bool   `json:"trackFps"`
		SupportAVIF          bool   `json:"supportAVIF"`
		ImgFormatCollect     struct {
			Ssr []string `json:"ssr"`
			Csr []string `json:"csr"`
		} `json:"imgFormatCollect"`
		IsUndertake bool `json:"isUndertake"`
	} `json:"global"`
	User struct {
		LoggedIn  bool `json:"loggedIn"`
		Activated bool `json:"activated"`
		UserInfo  struct {
			RedID    string `json:"red_id"`
			UserID   string `json:"user_id"`
			Nickname string `json:"nickname"`
			Desc     string `json:"desc"`
			Gender   int    `json:"gender"`
			Images   string `json:"images"`
			Imageb   string `json:"imageb"`
			Guest    bool   `json:"guest"`
			UserID0  string `json:"userId"`
			RedID0   string `json:"redId"`
		} `json:"userInfo"`
		Follow       []any `json:"follow"`
		UserPageData struct {
			TabPublic struct {
				Collection     bool `json:"collection"`
				CollectionNote struct {
					Lock    bool `json:"lock"`
					Count   int  `json:"count"`
					Display bool `json:"display"`
				} `json:"collectionNote"`
				CollectionBoard struct {
					Display bool `json:"display"`
					Lock    bool `json:"lock"`
					Count   int  `json:"count"`
				} `json:"collectionBoard"`
			} `json:"tabPublic"`
			ExtraInfo struct {
				Fstatus   string `json:"fstatus"`
				BlockType string `json:"blockType"`
			} `json:"extraInfo"`
			Result struct {
				Success bool   `json:"success"`
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"result"`
			BasicInfo struct {
				Images     string `json:"images"`
				RedID      string `json:"redId"`
				Gender     int    `json:"gender"`
				IPLocation string `json:"ipLocation"`
				Desc       string `json:"desc"`
				Imageb     string `json:"imageb"`
				Nickname   string `json:"nickname"`
			} `json:"basicInfo"`
			Interactions []struct {
				Type  string `json:"type"`
				Name  string `json:"name"`
				Count string `json:"count"`
			} `json:"interactions"`
			Tags []struct {
				Icon    string `json:"icon"`
				Name    string `json:"name"`
				TagType string `json:"tagType"`
			} `json:"tags"`
		} `json:"userPageData"`
		ActiveTab struct {
			Key      int    `json:"key"`
			Index    int    `json:"index"`
			Query    string `json:"query"`
			Label    string `json:"label"`
			Lock     bool   `json:"lock"`
			SubTabs  any    `json:"subTabs"`
			FeedType int    `json:"feedType"`
		} `json:"activeTab"`
		Notes [][]struct {
			ID       string `json:"id"`
			NoteCard struct {
				Type         string `json:"type"`
				DisplayTitle string `json:"displayTitle"`
				User         struct {
					NickName string `json:"nickName"`
					Avatar   string `json:"avatar"`
					UserID   string `json:"userId"`
					Nickname string `json:"nickname"`
				} `json:"user"`
				InteractInfo struct {
					LikedCount string `json:"likedCount"`
					Sticky     bool   `json:"sticky"`
					Liked      bool   `json:"liked"`
				} `json:"interactInfo"`
				Cover struct {
					TraceID  string `json:"traceId"`
					InfoList []struct {
						URL        string `json:"url"`
						ImageScene string `json:"imageScene"`
					} `json:"infoList"`
					URLPre     string `json:"urlPre"`
					URLDefault string `json:"urlDefault"`
					FileID     string `json:"fileId"`
					Height     int    `json:"height"`
					Width      int    `json:"width"`
					URL        string `json:"url"`
				} `json:"cover"`
				NoteID    string `json:"noteId"`
				XsecToken string `json:"xsecToken"`
			} `json:"noteCard"`
			Index       int    `json:"index"`
			Exposed     bool   `json:"exposed"`
			SsrRendered bool   `json:"ssrRendered"`
			XsecToken   string `json:"xsecToken"`
		} `json:"notes"`
		IsFetchingNotes        []bool   `json:"isFetchingNotes"`
		TabScrollTop           []int    `json:"tabScrollTop"`
		UserFetchingStatus     string   `json:"userFetchingStatus"`
		UserNoteFetchingStatus []string `json:"userNoteFetchingStatus"`
		BannedInfo             struct {
			UserID       string `json:"userId"`
			ServerBanned bool   `json:"serverBanned"`
			Code         int    `json:"code"`
			ShowAlert    bool   `json:"showAlert"`
			Reason       string `json:"reason"`
			API          string `json:"api"`
		} `json:"bannedInfo"`
		FirstFetchNote bool `json:"firstFetchNote"`
		NoteQueries    []struct {
			Num     int    `json:"num"`
			Cursor  string `json:"cursor"`
			UserID  string `json:"userId"`
			Page    int    `json:"page"`
			HasMore bool   `json:"hasMore"`
		} `json:"noteQueries"`
		PageScrolled int  `json:"pageScrolled"`
		ActiveSubTab any  `json:"activeSubTab"`
		IsOwnBoard   bool `json:"isOwnBoard"`
	} `json:"user"`
	Board struct {
		BoardListData struct {
		} `json:"boardListData"`
		IsLoadingBoardList bool `json:"isLoadingBoardList"`
		BoardDetails       struct {
		} `json:"boardDetails"`
		BoardFeedsMap struct {
		} `json:"boardFeedsMap"`
		BoardPageStatus string `json:"boardPageStatus"`
		UserBoardList   []any  `json:"userBoardList"`
	} `json:"board"`
	Login struct {
		LoginMethod any    `json:"loginMethod"`
		From        string `json:"from"`
		ShowLogin   bool   `json:"showLogin"`
		Agreed      bool   `json:"agreed"`
		ShowTooltip bool   `json:"showTooltip"`
		LoginData   struct {
			Phone    string `json:"phone"`
			AuthCode string `json:"authCode"`
		} `json:"loginData"`
		Errors struct {
			Phone    string `json:"phone"`
			AuthCode string `json:"authCode"`
		} `json:"errors"`
		QrData struct {
			Backend struct {
				QrID string `json:"qrId"`
				Code string `json:"code"`
			} `json:"backend"`
			Image  string `json:"image"`
			Status string `json:"status"`
		} `json:"qrData"`
		Counter                   any    `json:"counter"`
		InAntiSpamChecking        bool   `json:"inAntiSpamChecking"`
		RecentFrom                string `json:"recentFrom"`
		IsObPagesVisible          bool   `json:"isObPagesVisible"`
		ObPageFillInProgress      any    `json:"obPageFillInProgress"`
		VerificationCodeStartTime int    `json:"verificationCodeStartTime"`
		AgeSelectValue            string `json:"ageSelectValue"`
		HobbySelectValue          []any  `json:"hobbySelectValue"`
		GenderSelectValue         any    `json:"genderSelectValue"`
		InSpamCheckSendAuthCode   bool   `json:"inSpamCheckSendAuthCode"`
		IsRegFusing               bool   `json:"isRegFusing"`
		LoginStep                 int    `json:"loginStep"`
		IsLogining                bool   `json:"isLogining"`
		LoginPadMountedTime       int    `json:"loginPadMountedTime"`
		LoginTips                 string `json:"loginTips"`
		IsRiskUser                bool   `json:"isRiskUser"`
		CloseLoginModal           bool   `json:"closeLoginModal"`
		TraceID                   string `json:"traceId"`
		InAntiSpamCheckLogin      bool   `json:"inAntiSpamCheckLogin"`
	} `json:"login"`
	Feed struct {
		Query struct {
			CursorScore       string `json:"cursorScore"`
			Num               int    `json:"num"`
			RefreshType       int    `json:"refreshType"`
			NoteIndex         int    `json:"noteIndex"`
			UnreadBeginNoteID string `json:"unreadBeginNoteId"`
			UnreadEndNoteID   string `json:"unreadEndNoteId"`
			UnreadNoteCount   int    `json:"unreadNoteCount"`
			Category          string `json:"category"`
			SearchKey         string `json:"searchKey"`
			NeedNum           int    `json:"needNum"`
			ImageFormats      []any  `json:"imageFormats"`
			NeedFilterImage   bool   `json:"needFilterImage"`
		} `json:"query"`
		IsFetching     bool   `json:"isFetching"`
		IsError        bool   `json:"isError"`
		FeedsWrapper   any    `json:"feedsWrapper"`
		UndertakeNote  any    `json:"undertakeNote"`
		Feeds          []any  `json:"feeds"`
		CurrentChannel string `json:"currentChannel"`
		UnreadInfo     struct {
			CachedFeeds       []any  `json:"cachedFeeds"`
			UnreadBeginNoteID string `json:"unreadBeginNoteId"`
			UnreadEndNoteID   string `json:"unreadEndNoteId"`
			UnreadNoteCount   int    `json:"unreadNoteCount"`
			Timestamp         int    `json:"timestamp"`
		} `json:"unreadInfo"`
		ValidIds struct {
			NoteIds []any `json:"noteIds"`
		} `json:"validIds"`
		MfStatistics struct {
			Timestamp     int `json:"timestamp"`
			VisitTimes    int `json:"visitTimes"`
			ReadFeedCount int `json:"readFeedCount"`
		} `json:"mfStatistics"`
		Channels          any  `json:"channels"`
		IsResourceDisplay bool `json:"isResourceDisplay"`
		IsActivityEnd     bool `json:"isActivityEnd"`
		CancelFeedRequest bool `json:"cancelFeedRequest"`
		PrefetchID        any  `json:"prefetchId"`
		MfRequestMetaData struct {
			Start   any `json:"start"`
			Lasting any `json:"lasting"`
		} `json:"mfRequestMetaData"`
		PlaceholderFeeds  []any `json:"placeholderFeeds"`
		FeedsCacheLogInfo struct {
			Flag                      string `json:"flag"`
			ErrorCode                 int    `json:"errorCode"`
			IsHitMfCache              bool   `json:"isHitMfCache"`
			SSRDocumentChecked        bool   `json:"SSRDocumentChecked"`
			SSRDocumentCheckedSuccess bool   `json:"SSRDocumentCheckedSuccess"`
		} `json:"feedsCacheLogInfo"`
		IsUsingPlaceholderFeeds  bool   `json:"isUsingPlaceholderFeeds"`
		PlaceholderFeedsConsumed bool   `json:"placeholderFeedsConsumed"`
		IsReplace                bool   `json:"isReplace"`
		IsFirstSuccessFetched    bool   `json:"isFirstSuccessFetched"`
		ImgNoteFilterStatus      string `json:"imgNoteFilterStatus"`
		SsrRequestStatus         int    `json:"ssrRequestStatus"`
		SsrRenderExtra           string `json:"ssrRenderExtra"`
	} `json:"feed"`
	Layout struct {
		LayoutInfoReady    bool `json:"layoutInfoReady"`
		Columns            int  `json:"columns"`
		ColumnsWithSidebar int  `json:"columnsWithSidebar"`
		Gap                struct {
			Vertical   int `json:"vertical"`
			Horizontal int `json:"horizontal"`
		} `json:"gap"`
		ColumnWidth      int    `json:"columnWidth"`
		InteractionWidth int    `json:"interactionWidth"`
		WidthType        string `json:"widthType"`
		BufferRow        int    `json:"bufferRow"`
	} `json:"layout"`
	Search struct {
		State         string `json:"state"`
		SearchContext struct {
			Keyword  string `json:"keyword"`
			Page     int    `json:"page"`
			PageSize int    `json:"pageSize"`
			SearchID string `json:"searchId"`
			Sort     string `json:"sort"`
			NoteType int    `json:"noteType"`
			ExtFlags []any  `json:"extFlags"`
		} `json:"searchContext"`
		Feeds               []any  `json:"feeds"`
		SearchValue         string `json:"searchValue"`
		Suggestions         []any  `json:"suggestions"`
		UserInputSugTrigger string `json:"userInputSugTrigger"`
		KeywordFrom         int    `json:"keywordFrom"`
		SearchRecFilter     []any  `json:"searchRecFilter"`
		SearchRecFilterTag  any    `json:"searchRecFilterTag"`
		SearchFeedsWrapper  any    `json:"searchFeedsWrapper"`
		CurrentSearchType   string `json:"currentSearchType"`
		HintWord            struct {
			Title             string `json:"title"`
			SearchWord        string `json:"searchWord"`
			HintWordRequestID string `json:"hintWordRequestId"`
			Type              string `json:"type"`
		} `json:"hintWord"`
		SugType             any `json:"sugType"`
		QueryTrendingInfo   any `json:"queryTrendingInfo"`
		QueryTrendingParams struct {
			Source               string `json:"source"`
			SearchType           string `json:"searchType"`
			LastQuery            string `json:"lastQuery"`
			LastQueryTime        int    `json:"lastQueryTime"`
			WordRequestSituation string `json:"wordRequestSituation"`
			HintWord             string `json:"hintWord"`
			HintWordType         string `json:"hintWordType"`
			HintWordRequestID    string `json:"hintWordRequestId"`
		} `json:"queryTrendingParams"`
		QueryTrendingFetched bool `json:"queryTrendingFetched"`
		OneboxInfo           struct {
		} `json:"oneboxInfo"`
		HasMore                 bool   `json:"hasMore"`
		FirstEnterSearchPage    bool   `json:"firstEnterSearchPage"`
		UserLists               []any  `json:"userLists"`
		FetchUserListsStatus    string `json:"fetchUserListsStatus"`
		IsFetchingUserLists     bool   `json:"isFetchingUserLists"`
		HasMoreUser             bool   `json:"hasMoreUser"`
		SearchCplID             any    `json:"searchCplId"`
		WordRequestID           any    `json:"wordRequestId"`
		HistoryList             []any  `json:"historyList"`
		SearchPageHasPrevRoute  bool   `json:"searchPageHasPrevRoute"`
		SearchHotSpots          []any  `json:"searchHotSpots"`
		HotspotQueryNoteStep    string `json:"hotspotQueryNoteStep"`
		HotspotQueryNoteIndex   int    `json:"hotspotQueryNoteIndex"`
		CanShowHotspotQueryNote bool   `json:"canShowHotspotQueryNote"`
		ForceHotspotSearch      bool   `json:"forceHotspotSearch"`
		SearchCardHotSpots      []any  `json:"searchCardHotSpots"`
		IsHotspotSearch         bool   `json:"isHotspotSearch"`
	} `json:"search"`
	Activity struct {
		IsOpen     bool   `json:"isOpen"`
		CurrentURL string `json:"currentUrl"`
		EntryList  []any  `json:"entryList"`
	} `json:"activity"`
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
		NoteDetailMap   struct {
			Null struct {
				Comments struct {
					List               []any  `json:"list"`
					Cursor             string `json:"cursor"`
					HasMore            bool   `json:"hasMore"`
					Loading            bool   `json:"loading"`
					FirstRequestFinish bool   `json:"firstRequestFinish"`
				} `json:"comments"`
				CurrentTime int `json:"currentTime"`
				Note        struct {
				} `json:"note"`
			} `json:"null"`
		} `json:"noteDetailMap"`
		ServerRequestInfo struct {
			State     string `json:"state"`
			ErrorCode int    `json:"errorCode"`
			ErrMsg    string `json:"errMsg"`
		} `json:"serverRequestInfo"`
		Volume            int `json:"volume"`
		RecommendVideoMap struct {
		} `json:"recommendVideoMap"`
		VideoFeedType string `json:"videoFeedType"`
		Rate          int    `json:"rate"`
		CurrentNoteID any    `json:"currentNoteId"`
		MediaWidth    int    `json:"mediaWidth"`
		NoteHeight    int    `json:"noteHeight"`
	} `json:"note"`
	NioStore struct {
		CollectionListDataSource any `json:"collectionListDataSource"`
		Error                    any `json:"error"`
	} `json:"nioStore"`
	Notification struct {
		IsFetching        bool `json:"isFetching"`
		ActiveTabKey      int  `json:"activeTabKey"`
		NotificationCount struct {
			UnreadCount int `json:"unreadCount"`
			Mentions    int `json:"mentions"`
			Likes       int `json:"likes"`
			Connections int `json:"connections"`
		} `json:"notificationCount"`
		NotificationMap struct {
			Mentions struct {
				MessageList []any  `json:"messageList"`
				HasMore     bool   `json:"hasMore"`
				Cursor      string `json:"cursor"`
			} `json:"mentions"`
			Likes struct {
				MessageList []any  `json:"messageList"`
				HasMore     bool   `json:"hasMore"`
				Cursor      string `json:"cursor"`
			} `json:"likes"`
			Connections struct {
				MessageList []any  `json:"messageList"`
				HasMore     bool   `json:"hasMore"`
				Cursor      string `json:"cursor"`
			} `json:"connections"`
		} `json:"notificationMap"`
	} `json:"notification"`
}

type FeedResp struct {
	Code    int    `json:"code,omitempty"`
	Success bool   `json:"success,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Data    struct {
		CursorScore string `json:"cursor_score,omitempty"`
		Items       []struct {
			NoteCard struct {
				ImageList []struct {
					Height     int    `json:"height,omitempty"`
					Width      int    `json:"width,omitempty"`
					URL        string `json:"url,omitempty"`
					URLPre     string `json:"url_pre,omitempty"`
					URLDefault string `json:"url_default,omitempty"`
					LivePhoto  bool   `json:"live_photo,omitempty"`
					FileID     string `json:"file_id,omitempty"`
					TraceID    string `json:"trace_id,omitempty"`
					InfoList   []struct {
						ImageScene string `json:"image_scene,omitempty"`
						URL        string `json:"url,omitempty"`
					} `json:"info_list,omitempty"`
					Stream struct {
						H264 []struct {
							BackupUrls []string `json:"backup_urls,omitempty"`
							MasterURL  string   `json:"master_url,omitempty"`
						} `json:"h264,omitempty"`
						H265 []struct {
							BackupUrls []string `json:"backup_urls,omitempty"`
							MasterURL  string   `json:"master_url,omitempty"`
						} `json:"h265,omitempty"`
						H266 []struct {
							BackupUrls []string `json:"backup_urls,omitempty"`
							MasterURL  string   `json:"master_url,omitempty"`
						} `json:"h266,omitempty"`
						Av1 []struct {
							BackupUrls []string `json:"backup_urls,omitempty"`
							MasterURL  string   `json:"master_url,omitempty"`
						} `json:"av1,omitempty"`
					} `json:"stream"`
				} `json:"image_list,omitempty"`
				Time       int64  `json:"time,omitempty"`
				Title      string `json:"title,omitempty"`
				IPLocation string `json:"ip_location,omitempty"`
				ShareInfo  struct {
					UnShare bool `json:"un_share,omitempty"`
				} `json:"share_info,omitempty"`
				Video struct {
					Image struct {
						FirstFrameFileid string `json:"first_frame_fileid,omitempty"`
						ThumbnailFileid  string `json:"thumbnail_fileid,omitempty"`
					} `json:"image,omitempty"`
					Capa struct {
						Duration int `json:"duration,omitempty"`
					} `json:"capa,omitempty"`
					Consumer struct {
						OriginVideoKey string `json:"origin_video_key,omitempty"`
					} `json:"consumer,omitempty"`
					Media struct {
						VideoID int64 `json:"video_id,omitempty"`
						Video   struct {
							StreamTypes []int  `json:"stream_types,omitempty"`
							BizName     int    `json:"biz_name,omitempty"`
							BizID       string `json:"biz_id,omitempty"`
							Duration    int    `json:"duration,omitempty"`
							Md5         string `json:"md5,omitempty"`
							HdrType     int    `json:"hdr_type,omitempty"`
							DrmType     int    `json:"drm_type,omitempty"`
						} `json:"video,omitempty"`
						Stream struct {
							Av1  []any `json:"av1,omitempty"`
							H264 []struct {
								Size          int      `json:"size,omitempty"`
								Psnr          int      `json:"psnr,omitempty"`
								Format        string   `json:"format,omitempty"`
								Volume        int      `json:"volume,omitempty"`
								VideoCodec    string   `json:"video_codec,omitempty"`
								AudioBitrate  int      `json:"audio_bitrate,omitempty"`
								Rotate        int      `json:"rotate,omitempty"`
								HdrType       int      `json:"hdr_type,omitempty"`
								DefaultStream int      `json:"default_stream,omitempty"`
								Width         int      `json:"width,omitempty"`
								AudioCodec    string   `json:"audio_codec,omitempty"`
								MasterURL     string   `json:"master_url,omitempty"`
								Ssim          int      `json:"ssim,omitempty"`
								StreamType    int      `json:"stream_type,omitempty"`
								Duration      int      `json:"duration,omitempty"`
								VideoBitrate  int      `json:"video_bitrate,omitempty"`
								QualityType   string   `json:"quality_type,omitempty"`
								AvgBitrate    int      `json:"avg_bitrate,omitempty"`
								AudioDuration int      `json:"audio_duration,omitempty"`
								BackupUrls    []string `json:"backup_urls,omitempty"`
								Vmaf          int      `json:"vmaf,omitempty"`
								StreamDesc    string   `json:"stream_desc,omitempty"`
								Weight        int      `json:"weight,omitempty"`
								Height        int      `json:"height,omitempty"`
								Fps           int      `json:"fps,omitempty"`
								VideoDuration int      `json:"video_duration,omitempty"`
								AudioChannels int      `json:"audio_channels,omitempty"`
							} `json:"h264,omitempty"`
							H265 []struct {
								Height        int      `json:"height,omitempty"`
								AudioCodec    string   `json:"audio_codec,omitempty"`
								AudioChannels int      `json:"audio_channels,omitempty"`
								Rotate        int      `json:"rotate,omitempty"`
								QualityType   string   `json:"quality_type,omitempty"`
								StreamType    int      `json:"stream_type,omitempty"`
								StreamDesc    string   `json:"stream_desc,omitempty"`
								Format        string   `json:"format,omitempty"`
								VideoBitrate  int      `json:"video_bitrate,omitempty"`
								Weight        int      `json:"weight,omitempty"`
								DefaultStream int      `json:"default_stream,omitempty"`
								BackupUrls    []string `json:"backup_urls,omitempty"`
								HdrType       int      `json:"hdr_type,omitempty"`
								MasterURL     string   `json:"master_url,omitempty"`
								Ssim          int      `json:"ssim,omitempty"`
								Duration      int      `json:"duration,omitempty"`
								AvgBitrate    int      `json:"avg_bitrate,omitempty"`
								Size          int      `json:"size,omitempty"`
								VideoCodec    string   `json:"video_codec,omitempty"`
								Vmaf          int      `json:"vmaf,omitempty"`
								VideoDuration int      `json:"video_duration,omitempty"`
								AudioBitrate  int      `json:"audio_bitrate,omitempty"`
								Volume        int      `json:"volume,omitempty"`
								Fps           int      `json:"fps,omitempty"`
								Width         int      `json:"width,omitempty"`
								AudioDuration int      `json:"audio_duration,omitempty"`
								Psnr          float64  `json:"psnr,omitempty"`
							} `json:"h265,omitempty"`
							H266 []any `json:"h266,omitempty"`
						} `json:"stream,omitempty"`
					} `json:"media,omitempty"`
				} `json:"video,omitempty"`
				NoteID       string `json:"note_id,omitempty"`
				Desc         string `json:"desc,omitempty"`
				InteractInfo struct {
					Followed       bool   `json:"followed,omitempty"`
					Relation       string `json:"relation,omitempty"`
					Liked          bool   `json:"liked,omitempty"`
					LikedCount     string `json:"liked_count,omitempty"`
					Collected      bool   `json:"collected,omitempty"`
					CollectedCount string `json:"collected_count,omitempty"`
					CommentCount   string `json:"comment_count,omitempty"`
					ShareCount     string `json:"share_count,omitempty"`
				} `json:"interact_info,omitempty"`
				AtUserList     []any  `json:"at_user_list,omitempty"`
				LastUpdateTime int64  `json:"last_update_time,omitempty"`
				Type           string `json:"type,omitempty"`
				User           struct {
					UserID   string `json:"user_id,omitempty"`
					Nickname string `json:"nickname,omitempty"`
					Avatar   string `json:"avatar,omitempty"`
				} `json:"user,omitempty"`
				TagList []struct {
					ID   string `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
					Type string `json:"type,omitempty"`
				} `json:"tag_list,omitempty"`
			} `json:"note_card,omitempty"`
			ID        string `json:"id,omitempty"`
			ModelType string `json:"model_type,omitempty"`
		} `json:"items,omitempty"`
		CurrentTime int64 `json:"current_time,omitempty"`
	} `json:"data,omitempty"`
}

type AutoGenerated struct {
	Msg  string `json:"msg,omitempty"`
	Data struct {
		CursorScore string `json:"cursor_score,omitempty"`
		Items       []struct {
			ID        string `json:"id,omitempty"`
			ModelType string `json:"model_type,omitempty"`
			NoteCard  struct {
				AtUserList []any  `json:"at_user_list,omitempty"`
				IPLocation string `json:"ip_location,omitempty"`
				Desc       string `json:"desc,omitempty"`
				User       struct {
					UserID   string `json:"user_id,omitempty"`
					Nickname string `json:"nickname,omitempty"`
					Avatar   string `json:"avatar,omitempty"`
				} `json:"user,omitempty"`
				Title        string `json:"title,omitempty"`
				InteractInfo struct {
					Liked          bool   `json:"liked,omitempty"`
					LikedCount     string `json:"liked_count,omitempty"`
					Collected      bool   `json:"collected,omitempty"`
					CollectedCount string `json:"collected_count,omitempty"`
					CommentCount   string `json:"comment_count,omitempty"`
					ShareCount     string `json:"share_count,omitempty"`
					Followed       bool   `json:"followed,omitempty"`
					Relation       string `json:"relation,omitempty"`
				} `json:"interact_info,omitempty"`
				ImageList []struct {
					Height     int    `json:"height,omitempty"`
					Width      int    `json:"width,omitempty"`
					URL        string `json:"url,omitempty"`
					URLPre     string `json:"url_pre,omitempty"`
					URLDefault string `json:"url_default,omitempty"`
					LivePhoto  bool   `json:"live_photo,omitempty"`
					FileID     string `json:"file_id,omitempty"`
					TraceID    string `json:"trace_id,omitempty"`
					InfoList   []struct {
						ImageScene string `json:"image_scene,omitempty"`
						URL        string `json:"url,omitempty"`
					} `json:"info_list,omitempty"`
					Stream struct {
						H264 []struct {
							BackupUrls []string `json:"backup_urls,omitempty"`
							MasterURL  string   `json:"master_url,omitempty"`
						} `json:"h264,omitempty"`
						H265 []any `json:"h265,omitempty"`
						H266 []any `json:"h266,omitempty"`
						Av1  []any `json:"av1,omitempty"`
					} `json:"stream,omitempty"`
					Stream0 struct {
					} `json:"stream,omitempty"`
				} `json:"image_list,omitempty"`
				TagList []struct {
					Name string `json:"name,omitempty"`
					Type string `json:"type,omitempty"`
					ID   string `json:"id,omitempty"`
				} `json:"tag_list,omitempty"`
				Time           int64  `json:"time,omitempty"`
				LastUpdateTime int64  `json:"last_update_time,omitempty"`
				NoteID         string `json:"note_id,omitempty"`
				Type           string `json:"type,omitempty"`
				ShareInfo      struct {
					UnShare bool `json:"un_share,omitempty"`
				} `json:"share_info,omitempty"`
			} `json:"note_card,omitempty"`
		} `json:"items,omitempty"`
		CurrentTime int64 `json:"current_time,omitempty"`
	} `json:"data,omitempty"`
	Code    int  `json:"code,omitempty"`
	Success bool `json:"success,omitempty"`
}
