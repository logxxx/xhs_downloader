package storage

import (
	"context"
	"errors"
	storm "github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/logxxx/utils/logutil"
	"github.com/logxxx/xhs_downloader/model"
	"log"
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
	db, err := storm.Open("chore/core.db")
	if err != nil {
		panic(err)
	}
	s.db = db
}

func (s *Storage) SetNoteDownloaded(noteID string) {
	s.db.Set("downloaded_note", noteID, time.Now().Unix())
}

func (s *Storage) SetUperScanned(uid string) {
	s.db.Set("scanned_uper", uid, time.Now().Unix())
}

func (s *Storage) IsUperScanned(uid string) bool {
	ok, _ := s.db.KeyExists("scanned_uper", uid)
	return ok
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
	s.db.From("note").One("NoteID", noteID, &w)
	return w
}

func (s *Storage) UpdateNote(w model.Note) error {
	return s.db.From("note").Update(&w)
}

func (s *Storage) InsertOrUpdateNote(w model.Note) (insertOrUpdate string, err error) {

	n := s.GetNote(w.NoteID)

	if n.ID <= 0 {
		w.CreateTime = time.Now()
		w.UpdateTime = time.Now()
		err = s.db.From("note").Save(&w)
		if err != nil {
			log.Printf("InsertOrUpdateNote Save err:%v w:%+v", err, w)
		} else {
			insertOrUpdate = "insert"
		}
	} else {
		w.ID = n.ID
		w.UpdateTime = time.Now()
		err = s.db.From("note").Update(&w)
		if err != nil {
			log.Printf("InsertOrUpdateNote Update err:%v w:%+v", err, w)
		} else {
			insertOrUpdate = "update"
		}
	}

	return

}

func (s *Storage) GetNotesByUper(uperUID string) (resp []model.Note) {
	s.db.From("note").Find("UperUID", uperUID, &resp)
	return
}

func (s *Storage) UperAddNote(uperUID string, noteID string) (failedReason string, err error) {

	if uperUID == "" {
		err = errors.New("empty uperUID")
		return
	}

	if noteID == "" {
		err = errors.New("empty noteID")
		return
	}

	u := s.GetUper(0, uperUID)
	if u.ID <= 0 {
		err = errors.New("uper not found")
		return
	}
	for _, w := range u.Notes {
		if w == noteID {
			failedReason = "noteID already exists"
		}
	}
	u.Notes = append(u.Notes, noteID)
	u.UpdateTime = time.Now()
	err = s.db.From("uper").Update(&u)
	if err != nil {
		return
	}
	return "", nil
}

func (s *Storage) GetUperTotalCount() int {
	resp, _ := s.db.From("uper").Count(&model.Uper{})
	return resp
}

func (s *Storage) GetNoteTotalCount() int {
	resp, _ := s.db.From("note").Count(&model.Uper{})
	return resp
}

func (s *Storage) InsertOrUpdateUper(input model.Uper) (string, error) {
	u := s.GetUper(0, input.UID)
	if u.ID > 0 {
		input.UpdateTime = time.Now()
		return "update", s.db.From("uper").Update(&input)
	} else {
		input.CreateTime = time.Now()
		input.UpdateTime = time.Now()
		return "insert", s.db.From("uper").Save(&input)
	}

}

func (s *Storage) GetUper(id int64, uid string) (resp model.Uper) {
	if id > 0 {
		err := s.db.From("uper").One("ID", id, &resp)
		if err != nil {
			//log.Printf("GetUper By id err:%v id:%v", err, id)
		}
	} else if uid != "" {
		err := s.db.From("uper").One("UID", uid, &resp)
		if err != nil {
			//log.Printf("GetUper By Uid err:%v uid:%v", err, uid)
		}
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

func (s *Storage) GetAllNotes(ctx context.Context) (resp []model.Note) {
	s.db.From("note").All(&resp)
	return
}

func (s *Storage) GetAllUpers(ctx context.Context) (resp []model.Uper) {
	s.db.From("uper").All(&resp)
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
	err := s.db.From("note").Select(qs...).OrderBy("ID").Reverse().Limit(limit).Find(&resp)
	if err != nil {
		logger.Errorf("Find err:%v", err)
	}
	if len(resp) > 0 {
		nextToken = strconv.Itoa(int(resp[len(resp)-1].ID))
	}
	return
}
