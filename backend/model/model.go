package model

import "time"

type Uper struct {
	ID               int64     `storm:"id,increment" json:"id,omitempty"` // 自增ID
	UID              string    `storm:"unique" json:"uid,omitempty"`
	Name             string    `storm:"index" json:"name,omitempty"`
	Area             string    `json:"area,omitempty"`
	AvatarURL        string    `json:"avatar_url,omitempty"`
	IsGirl           bool      `json:"is_girl,omitempty"`
	Desc             string    `json:"desc,omitempty"`
	Notes            []string  `json:"notes,omitempty"`
	Tags             []string  `json:"tags,omitempty"`
	FansCount        int       `json:"fans_count,omitempty"`
	ReceiveLikeCount int       `json:"receive_like_count,omitempty"`
	IsLike           bool      `json:"is_like,omitempty"`
	IsDelete         bool      `json:"is_delete,omitempty"`
	MyTags           []string  `json:"my_tags,omitempty"`
	StarCount        int       `json:"star_count,omitempty"`
	CreateTime       time.Time `storm:"index" json:"create_time"`
	UpdateTime       time.Time `storm:"index" json:"update_time"`
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
