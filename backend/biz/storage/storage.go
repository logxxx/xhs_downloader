package storage

import (
	"context"
	"errors"
	storm "github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/logxxx/utils/logutil"
	"github.com/logxxx/xhs_downloader/model"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var (
	_storage     *Storage
	_storageLock sync.Mutex
)

type Storage struct {
	db *storm.DB
}

func GetStorage() *Storage {
	if _storage != nil {
		return _storage
	}
	_storageLock.Lock()
	defer _storageLock.Unlock()
	if _storage != nil {
		return _storage
	}

	s := &Storage{}
	s.initDB()

	_storage = s

	return _storage

}

func (s *Storage) initDB() {
	db, err := storm.Open("chore/xhs.db")
	if err != nil {
		panic(err)
	}
	s.db = db
}

func (s *Storage) SetNoteDownloaded(noteID string) {
	s.db.Set("downloaded_note", noteID, time.Now().Unix())
}

func (s *Storage) IsNoteDownloaded(noteID string) bool {
	ok, _ := s.db.KeyExists("downloaded_note", noteID)
	return ok
}

func (s *Storage) UpdateNoteBynoteID(noteID string, h func(w model.Note) model.Note) (err error) {
	w := s.GetNote(noteID)
	if w.ID <= 0 {
		return
	}
	newW := h(w)
	if reflect.DeepEqual(w, newW) {
		return
	}
	return s.UpdateNote(newW)
}

func (s *Storage) GetNote(noteID string) model.Note {
	w := model.Note{}
	s.db.From("Note").One("noteID", noteID, &w)
	return w
}

func (s *Storage) UpdateNote(w model.Note) error {
	return s.db.From("Note").Save(w)
}

func (s *Storage) InsertNote(w model.Note) error {
	return s.db.From("Note").Update(w)
}

func (s *Storage) GetNotesByUper(uperUID string) (resp []model.Note) {
	s.db.From("Note").Find("UperUID", uperUID, &resp)
	return
}

func (s *Storage) UperAddNote(uperUID string, noteID string) error {
	u := s.GetUper(0, uperUID)
	if u.ID <= 0 {
		return errors.New("uper not found")
	}
	for _, w := range u.Notes {
		if w == noteID {
			return nil
		}
	}
	u.Notes = append(u.Notes, noteID)
	u.UpdateTime = time.Now()
	return s.db.Update(u)
}

func (s *Storage) InsertUper(input model.Uper) error {
	input.CreateTime = time.Now()
	input.UpdateTime = time.Now()
	return s.db.Save(input)
}

func (s *Storage) GetUper(id int64, uid string) (resp model.Uper) {
	if id > 0 {
		s.db.One("ID", id, &resp)
	} else if uid != "" {
		s.db.One("UID", uid, &resp)
	}
	return
}

type GetUpersOpt struct {
	FilterDelete bool
	OnlyLike     bool
	Tags         []string
}

func (s *Storage) GetUpers(ctx context.Context, opt GetUpersOpt, limit int, token string) (nextToken string, resp []model.Uper) {
	logger, _ := logutil.CtxLog(ctx, "GetUpers")
	qs := []q.Matcher{}
	if opt.OnlyLike {
		qs = append(qs, q.Eq("IsLike", true))
	}
	if opt.FilterDelete {
		qs = append(qs, q.Eq("IsDelete", false))
	}
	if len(opt.Tags) > 0 {
		qs = append(qs, q.In("Tags", opt.Tags))
	}
	if token != "" {
		id, _ := strconv.Atoi(token)
		if id > 0 {
			qs = append(qs, q.Lt("ID", id))
		}
	}
	err := s.db.From("uper").Select(qs...).OrderBy("ID").Reverse().Limit(limit).Find(&resp)
	if err != nil {
		logger.Errorf("Find err:%v", err)
	}
	if len(resp) > 0 {
		nextToken = strconv.Itoa(int(resp[len(resp)-1].ID))
	}
	return
}

func (s *Storage) GetNotes(ctx context.Context, opt GetUpersOpt, limit int, token string) (nextToken string, resp []model.Note) {
	logger, _ := logutil.CtxLog(ctx, "GetNotes")
	qs := []q.Matcher{}
	if opt.OnlyLike {
		qs = append(qs, q.Eq("IsLike", true))
	}
	if opt.FilterDelete {
		qs = append(qs, q.Eq("IsDelete", false))
	}
	if token != "" {
		id, _ := strconv.Atoi(token)
		if id > 0 {
			qs = append(qs, q.Lt("ID", id))
		}
	}
	err := s.db.From("Note").Select(qs...).OrderBy("ID").Reverse().Limit(limit).Find(&resp)
	if err != nil {
		logger.Errorf("Find err:%v", err)
	}
	if len(resp) > 0 {
		nextToken = strconv.Itoa(int(resp[len(resp)-1].ID))
	}
	return
}
