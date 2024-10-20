package proto

type ApiGetUpersResp struct {
	Data  []ApiUperInfo `json:"data,omitempty"`
	Token string        `json:"token,omitempty"`
}

type ApiGetUperResp struct {
	Data ApiUperInfo `json:"data,omitempty"`
}

type ApiUperInfo struct {
	UID      string   `json:"uid,omitempty"`
	Name     string   `json:"name,omitempty"`
	Desc     string   `json:"desc,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Notes    []string `json:"notes,omitempty"`
	Avatar   string   `json:"avatar,omitempty"`
	IsDelete bool     `json:"is_delete,omitempty"`
}

type ApiGetUperNotesResp struct {
	Data  []ApiUperNote `json:"data,omitempty"`
	Token string        `json:"token,omitempty"`
}

type ApiUperNote struct {
	UperUID   string   `json:"uper_uid,omitempty"`
	NoteID    string   `json:"note_id,omitempty"`
	Poster    string   `json:"poster,omitempty"`
	Title     string   `json:"title,omitempty"`
	Content   string   `json:"content,omitempty"`
	Video     string   `json:"video,omitempty"`
	Images    []string `json:"images,omitempty"`
	Lives     []string `json:"lives,omitempty"`
	Tags      []string `json:"tags"`
	ShowSize  string   `json:"show_size,omitempty"`
	IsDeleted bool     `json:"is_deleted,omitempty"`
}

type ApiGetOneNoteResp struct {
	Data  ApiUperNote `json:"data"`
	Token string      `json:"token,omitempty"`
}
