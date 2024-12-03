package blogmodel

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
