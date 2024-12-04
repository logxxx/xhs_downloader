package model

import (
	"github.com/logxxx/utils"
	"time"
)

type Uper struct {
	ID                   int64  `storm:"id,increment"` // 自增ID
	UID                  string `storm:"unique"`
	Name                 string `storm:"index"`
	Area                 string
	AvatarURL            string
	IsGirl               bool
	Desc                 string
	Notes                []string
	NotesLastUpdateTime  time.Time
	HomeTags             []string //up给自己打的主页tag
	Tags                 []string `storm:"index"` //我打的tag
	FansCount            int
	ReceiveLikeCount     int
	IsLike               bool
	IsDelete             bool
	IsBanned             bool
	MyTags               []string
	CreateTime           time.Time `storm:"index"`
	UpdateTime           time.Time `storm:"index"`
	GalleryEmptyLastTime time.Time
}

func (u *Uper) AddNote(n string) bool {
	for _, u := range u.Notes {
		if u == n {
			return false
		}
	}
	u.Notes = append(u.Notes, n)
	return true
}

func (u *Uper) RemoveNote(n string) bool {
	resp := []string{}
	isChanged := false
	for _, u := range u.Notes {
		if u == n {
			isChanged = true
			continue
		}
		resp = append(resp, u)
	}
	u.Notes = resp
	return isChanged
}

func (u *Uper) HasTag(input string) bool {
	for _, t := range u.Tags {
		if t == input {
			return true
		}
	}
	return false
}

type Note struct {
	ID              int64  `storm:"id,increment"` // 自增ID
	NoteID          string `storm:"unique"`
	URL             string
	PosterURL       string
	UperUID         string `storm:"index"`
	Title           string
	Content         string
	DownloadTime    time.Time `storm:"index"`
	DownloadNothing bool
	Video           string
	VideoURL        string
	Images          []string
	ImageURLs       []string
	Lives           []string
	LiveURLs        []string
	LikeCount       int
	IsLike          bool
	IsDelete        bool
	Tags            []string  `storm:"index"`
	WorkCreateTime  time.Time `storm:"index"`
	CreateTime      time.Time `storm:"index"`
	UpdateTime      time.Time `storm:"index"`
	FileSize        int64     `storm:"index"`
	FileSizeReverse int64     `storm:"index"`
}

func (n *Note) HasTag(t string) bool {
	for _, tag := range n.Tags {
		if tag == t {
			return true
		}
	}
	return false
}

func (n *Note) IsDownloaded() bool {
	if n.ID <= 0 {
		return false
	}

	if n.Video != "" && utils.HasFile(n.Video) {
		return true
	}

	if len(n.Images) > 0 && utils.HasFile(n.Images[0]) {
		return true
	}

	if len(n.Lives) > 0 && utils.HasFile(n.Lives[0]) {
		return true
	}

	return false
}

type Work struct {
	BlogURL   string
	NoteID    string
	XSecToken string
}
