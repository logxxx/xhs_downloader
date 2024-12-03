package queue

import (
	"fmt"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/randutil"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func init() {

}

func Push(queueName string, obj interface{}) (err error) {

	msgPath := fmt.Sprintf("chore/queue/%v/%v_%v", queueName, time.Now().Format("20060102_150405"), randutil.RandStr(8))

	return fileutil.WriteJsonToFile(obj, msgPath)
}

var (
	popCache = []string{}
	popLock  sync.Mutex
)

func Pop(queueName string, obj interface{}) (err error) {

	logger := log.WithField("func_name", "Pop")

	popLock.Lock()
	defer popLock.Unlock()

	for i, elem := range popCache {
		if !utils.HasFile(elem) {
			continue
		}
		//logger.Infof("pop from popCache(%v)", len(popCache))
		err = fileutil.ReadJsonFile(elem, obj)
		if err != nil {
			logger.Errorf("Pop ReadJsonFile from popCache err:%v firstFile:%v", err, elem)
			return
		}

		popCacheLen := len(popCache)
		_ = popCacheLen

		if len(popCache) > 1 {
			popCache = popCache[i+1:]
		} else {
			//logger.Infof("clean popCache")
			popCache = []string{}
		}

		//logger.Infof("after pop. popCache len:%v->%v", popCacheLen, len(popCache))

		os.Remove(elem)
		return
	}
	popCache = []string{}

	rootDir := fmt.Sprintf("chore/queue/%v", queueName)

	if !utils.HasFile(rootDir) {
		err = os.MkdirAll(rootDir, 0755)
		if err != nil {
			logger.Errorf("Pop MkdirAll err:%v rootDir:%v", err, rootDir)
			return
		}
	}

	for {
		files, e := os.ReadDir(rootDir)
		if e != nil {
			logger.Errorf("Pop ReadDir err:%v rootDir:%v", err, rootDir)
			return e
		}

		if len(files) <= 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		//log.Infof("pop from ReadDir(%v)", len(files))

		filePath := filepath.Join(rootDir, files[0].Name())

		err = fileutil.ReadJsonFile(filePath, obj)
		if err != nil {
			logger.Errorf("Pop ReadJsonFile err:%v filePath:%v", err, filePath)
			return err
		}

		os.Remove(filePath)

		if len(files) == 1 {
			break
		}

		for _, elem := range files[1:] {
			popCache = append(popCache, filepath.Join(rootDir, elem.Name()))
		}

		//log.Infof("after append to popCache, len:%v", len(popCache))

		break

	}

	return
}
