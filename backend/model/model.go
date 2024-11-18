package model

import "time"

type Uper struct {
	ID               int64  `storm:"id,increment"` // 自增ID
	UID              string `storm:"unique"`
	Name             string `storm:"index"`
	Area             string
	AvatarURL        string
	IsGirl           bool
	Desc             string
	Notes            []string
	Tags             []string `storm:"index"`
	FansCount        int
	ReceiveLikeCount int
	IsLike           bool
	IsDelete         bool
	MyTags           []string
	CreateTime       time.Time `storm:"index"`
	UpdateTime       time.Time `storm:"index"`
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
	Images          []string
	Lives           []string
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
	if n.DownloadTime.IsZero() {
		return false
	}
	if len(n.Images) <= 0 && len(n.Lives) <= 0 && n.Video == "" {
		return false
	}
	return true
}
