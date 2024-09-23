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

func (s *Storage) UpdateWorkByWorkID(workID string, h func(w model.Work) model.Work) (err error) {
	w := s.GetWork(workID)
	if w.ID <= 0 {
		return
	}
	newW := h(w)
	if reflect.DeepEqual(w, newW) {
		return
	}
	return s.UpdateWork(newW)
}

func (s *Storage) GetWork(workID string) model.Work {
	w := model.Work{}
	s.db.From("work").One("WorkID", workID, &w)
	return w
}

func (s *Storage) UpdateWork(w model.Work) error {
	return s.db.From("work").Save(w)
}

func (s *Storage) InsertWork(w model.Work) error {
	return s.db.From("work").Update(w)
}

func (s *Storage) GetWorks(uperUID string) (resp []model.Work) {
	s.db.From("work").Find("UperUID", uperUID, &resp)
	return
}

func (s *Storage) UperAddWork(uperUID string, workID string) error {
	u := s.GetUper(0, uperUID)
	if u.ID <= 0 {
		return errors.New("uper not found")
	}
	for _, w := range u.Works {
		if w == workID {
			return nil
		}
	}
	u.Works = append(u.Works, workID)
	u.UpdateTime = time.Now()
	s.db.Update(u)
	return nil
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
	FilterBlack bool
	OnlyLike    bool
}

func (s *Storage) GetUpers(ctx context.Context, opt GetUpersOpt, limit int, token string) (nextToken string, resp []model.Uper) {
	logger, _ := logutil.CtxLog(ctx, "GetUpers")
	qs := []q.Matcher{}
	if opt.OnlyLike {
		qs = append(qs, q.Eq("IsLike", true))
	}
	if opt.FilterBlack {
		qs = append(qs, q.Eq("IsBlack", false))
	}
	if token != "" {
		id, _ := strconv.Atoi(token)
		if id > 0 {
			qs = append(qs, q.Lt("ID", id))
		}
	}
	err := s.db.From("uper").Select(qs...).OrderBy("ID").Reverse().Limit(limit).Find(&resp)
	if err != nil {
		logger.Errorf("Range err:%v", err)
	}
	if len(resp) > 0 {
		nextToken = strconv.Itoa(int(resp[len(resp)-1].ID))
	}
	return
}
