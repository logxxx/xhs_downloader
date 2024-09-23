package model

import "time"

type Uper struct {
	ID         int64
	UID        string
	Name       string
	AvatarURL  string
	Sex        string
	Desc       string
	Works      []string
	Tags       []string
	IsLike     bool
	IsDelete   bool
	CreateTime time.Time
	UpdateTime time.Time
}

type Work struct {
	ID             int64
	WorkID         string
	UperUID        string
	Title          string
	Content        string
	Video          string
	Images         []string
	IsLike         bool
	IsDelete       bool
	WorkCreateTime time.Time
	CreateTime     time.Time
	UpdateTime     time.Time
}
