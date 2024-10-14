package proto

type ApiGetUperInfoResp struct {
	Data  []ApiUperInfo `json:"data,omitempty"`
	Token string        `json:"token,omitempty"`
}

type ApiUperInfo struct {
	UID    string   `json:"uid,omitempty"`
	Name   string   `json:"name,omitempty"`
	Desc   string   `json:"desc,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Avatar string   `json:"avatar,omitempty"`
}

type ApiGetUperNotesResp struct {
	Data  []ApiUperNote `json:"data,omitempty"`
	Token string        `json:"token,omitempty"`
}

type ApiUperNote struct {
	NoteID  string   `json:"note_id,omitempty"`
	Poster  string   `json:"poster,omitempty"`
	Title   string   `json:"title,omitempty"`
	Content string   `json:"content,omitempty"`
	Video   string   `json:"video,omitempty"`
	Images  []string `json:"images,omitempty"`
	Lives   []string `json:"lives,omitempty"`
}
