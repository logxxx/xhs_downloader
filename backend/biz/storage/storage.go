package storage

import (
	"context"
	"errors"
	storm "github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/index"
	"github.com/asdine/storm/v3/q"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/logutil"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"math"
	"path/filepath"
	"reflect"
	"sort"
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

func (s *Storage) DB() *storm.DB {
	return s.db
}

func (s *Storage) initDB() {
	log.Printf("initDB start")
	db, err := storm.Open("chore/core.db")
	if err != nil {
		panic(err)
	}
	log.Printf("initDB succ")
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

	if w.ID > 0 {
		err = s.db.From("note").Update(&w)
		if err != nil {
			log.Printf("InsertOrUpdateNote Update err:%v w:%+v", err, w)
		} else {
			insertOrUpdate = "update"
		}
		return
	}

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

func (s *Storage) UperAddNote(uperUID string, noteIDs ...string) (failedReason string, err error) {

	if uperUID == "" {
		err = errors.New("empty uperUID")
		return
	}

	if len(noteIDs) <= 0 {
		err = errors.New("empty noteID")
		return
	}

	u := s.GetUper(0, uperUID)
	if u.ID <= 0 {
		err = errors.New("uper not found")
		return
	}

	u.Notes = append(u.Notes, noteIDs...)
	u.Notes = utils.RemoveEmpty(utils.RemoveDuplicate(u.Notes))

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
	resp, _ := s.db.From("note").Count(&model.Note{})
	return resp
}

func (s *Storage) InsertOrUpdateUper(input model.Uper) (string, error) {
	u := s.GetUper(0, input.UID)
	if u.ID > 0 {
		input.ID = u.ID
		input.UpdateTime = time.Now()
		return "update", s.db.From("uper").Update(&input)
	} else {
		input.CreateTime = time.Now()
		input.UpdateTime = time.Now()
		return "insert", s.db.From("uper").Save(&input)
	}

}

func (s *Storage) GetUper(id int64, uid string) (resp model.Uper) {

	defer func() {
		newNotes := utils.RemoveEmpty(utils.RemoveDuplicate(resp.Notes))
		if len(newNotes) != len(resp.Notes) {
			log.Printf("GetUper %v fix notes:%v=>%v", resp.Name, len(resp.Notes), len(newNotes))
			resp.Notes = newNotes
			s.db.From("uper").Update(&resp)
		}
	}()

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
	Uid          string
	FilterDelete bool
	OnlyLike     bool
	Tags         []string
	WithNoTag    bool
}

func (s *Storage) GetUpers(ctx context.Context, opt GetUpersOpt, limit int, token string) (resp []model.Uper, nextToken string) {
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
	if opt.WithNoTag {
		qs = append(qs, q.Eq("Tags", nil))
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

func (s *Storage) EachUper(fn func(n model.Uper, currCount, totalCount int) (e error)) (err error) {

	total := s.GetUperTotalCount()

	currCount := 0
	lastID := int64(0)
	for {
		upers := []model.Uper{}
		options := []func(*index.Options){storm.Limit(100)}
		err = s.db.From("uper").Range("ID", lastID+1, int64(math.MaxInt64), &upers, options...)
		if err != nil {
			return
		}
		if len(upers) <= 0 {
			break
		}

		if upers[0].ID == lastID {
			break
		}

		for _, n := range upers {
			currCount++
			err = fn(n, currCount, total)
			if err != nil {
				return err
			}
			lastID = n.ID
		}

	}

	return

}

func (s *Storage) EachNoteBySelect(skip int, fn func(n model.Note, currCount, totalCount int) (e error)) (err error) {

	total := s.GetNoteTotalCount()

	lastID := int64(0)
	currCount := skip
	round := 0
	for {
		notes := []model.Note{}
		ms := []q.Matcher{}
		if lastID > 0 {
			ms = append(ms, q.Lt("ID", lastID))
		}
		round++

		db := s.db.From("note").Select(ms...).OrderBy("ID").Limit(100).Reverse()
		if round == 1 && skip > 0 {
			db = db.Skip(skip)
		}
		err := db.Find(&notes)
		if err != nil {
			return err
		}
		if len(notes) <= 0 {
			break
		}
		lastID = notes[len(notes)-1].ID
		for _, n := range notes {
			currCount++
			fn(n, currCount, total)
		}
	}
	return nil
}

var videoCache = []model.Note{}

func tryGetVideoFromCache() (resp model.Note) {
	if len(videoCache) <= 1 {
		return
	}

	for i, v := range videoCache {
		thumbPath := filepath.Join(filepath.Dir(v.Video), ".thumb", filepath.Base(v.Video))
		if utils.HasFile(thumbPath) {
			videoCache = videoCache[i+1:]
			return v
		}
	}

	return

}

func (s *Storage) GetOneVideoNoteBySize2(token string) (resp model.Note, nextToken string, err error) {

	cacheV := tryGetVideoFromCache()
	if cacheV.ID > 0 {
		log.Printf("GetOneVideoNoteBySize2 tryGetVideoFromCache succ:%+v", cacheV.Video)
		return cacheV, token, nil
	}

	ms := []q.Matcher{
		q.Eq("IsDelete", false),
		q.Not(q.Eq("Video", "")),
		q.Eq("Tags", nil),
		//q.Gt("FileSize", 1),
	}

	if token != "" {
		lastID, _ := strconv.Atoi(token)
		if lastID > 0 {
			ms = append(ms, q.Gt("ID", lastID))
		}
	}

	step := 100

	resps := []model.Note{}
	s.db.From("note").Select(ms...).Limit(step).Find(&resps)

	if len(resps) > 0 {
		nextToken = strconv.Itoa(int(resps[len(resps)-1].ID))
	}

	sort.Slice(resps, func(i, j int) bool {
		return resps[i].FileSize > resps[j].FileSize
	})
	if len(resps) >= 1 {
		resp = resps[0]
	}

	if len(resps) > 1 {
		videoCache = resps[1:]
	}

	return
}

func (s *Storage) GetOneVideoNoteBySize(token string) (resp model.Note, nextToken string, err error) {
	ms := []q.Matcher{
		q.Eq("IsDelete", false),
		q.Not(q.Eq("Video", "")),
		q.Eq("Tags", nil),
	}
	if token != "" {
		lastID, _ := strconv.Atoi(token)
		if lastID > 0 {
			ms = append(ms, q.Gt("FileSizeReverse", int64(lastID)))
		}
	} else {
		ms = append(ms, q.Gt("FileSizeReverse", 1))
	}
	err = s.db.From("note").Select(ms...).OrderBy("FileSizeReverse").First(&resp)
	if resp.FileSizeReverse > 0 {
		nextToken = strconv.FormatInt(resp.FileSizeReverse, 10)
	}

	return
}

func (s *Storage) GetOneNote(token string, t string) (resp model.Note, nextToken string, err error) {
	ms := []q.Matcher{
		q.Not(q.Eq("DownloadTime", nil)),
		q.Eq("Tags", nil),
	}
	switch t {
	case "image":
		ms = append(ms, q.Not(q.Eq("Images", nil)))
	case "live":
		ms = append(ms, q.Not(q.Eq("Lives", nil)))
	case "video":
		ms = append(ms, q.Not(q.Eq("Video", "")))
	}

	if token != "" {
		lastID, _ := strconv.Atoi(token)
		if lastID > 0 {
			ms = append(ms, q.Gt("ID", int64(lastID)))
		}
	}
	err = s.db.From("note").Select(ms...).First(&resp)
	if resp.ID > 0 {
		nextToken = strconv.Itoa(int(resp.ID))
	}

	return
}

func (s *Storage) GetNotesByPage(limit int, token string) (resp []model.Note, nextToken string, err error) {
	ms := []q.Matcher{
		//q.Not(q.Eq("DownloadTime", nil)),
	}

	if token != "" {
		lastID, _ := strconv.Atoi(token)
		if lastID > 0 {
			ms = append(ms, q.Lt("ID", int64(lastID)))
		}
	}

	err = s.db.From("note").Select(ms...).OrderBy("ID").Limit(limit).Find(&resp)

	if len(resp) > 0 {
		nextToken = strconv.Itoa(int(resp[len(resp)-1].ID))
	}

	return
}

func (s *Storage) EachNoteByRange(skip int, fn func(n model.Note, currCount, totalCount int) (e error)) (err error) {

	total := s.GetNoteTotalCount()

	currCount := 0
	lastID := int64(0)
	round := 0
	for {
		round++
		notes := []model.Note{}
		options := []func(*index.Options){
			storm.Limit(100),
		}
		if round == 1 && skip > 0 {
			options = append(options, storm.Skip(skip))
		}
		err = s.db.From("note").Range("ID", lastID+1, int64(math.MaxInt64), &notes, options...)
		if err != nil {
			return
		}
		if len(notes) <= 0 {
			break
		}

		if notes[0].ID == lastID {
			break
		}

		for _, n := range notes {
			currCount++
			err = fn(n, currCount, total)
			if err != nil {
				return err
			}
			lastID = n.ID
		}

	}

	return

}

func (s *Storage) GetNotes(ctx context.Context, opt GetUpersOpt, limit int, token string) (resp []model.Note, nextToken string) {
	logger, _ := logutil.CtxLog(ctx, "GetNotes")
	qs := []q.Matcher{}
	if opt.OnlyLike {
		qs = append(qs, q.Eq("IsLike", true))
	}
	if opt.FilterDelete {
		qs = append(qs, q.Eq("IsDelete", false))
	}
	if opt.Uid != "" {
		qs = append(qs, q.Eq("UperUID", opt.Uid))
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
