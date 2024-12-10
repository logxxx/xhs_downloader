package blogmodel

type ParseNote struct {
	NoteID          string
	URL             string
	UperName        string
	UperUID         string
	Title           string
	Poster          string
	LikeCount       int
	MideaSimpleInfo string
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
