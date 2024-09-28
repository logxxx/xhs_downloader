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
	Tags             []string
	FansCount        int
	ReceiveLikeCount int
	IsLike           bool
	IsDelete         bool
	MyTags           []string
	CreateTime       time.Time `storm:"index"`
	UpdateTime       time.Time `storm:"index"`
}

type Note struct {
	ID             int64  `storm:"id,increment"` // 自增ID
	NoteID         string `storm:"unique"`
	URL            string
	PosterURL      string
	UperUID        string `storm:"index"`
	Title          string
	Content        string
	Video          string
	Images         []string
	LikeCount      int
	IsLike         bool
	IsDelete       bool
	WorkCreateTime time.Time `storm:"index"`
	CreateTime     time.Time `storm:"index"`
	UpdateTime     time.Time `storm:"index"`
}
