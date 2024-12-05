package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/netutil"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"time"
)

type Input struct {
	Module   string      `json:"module"`
	Function string      `json:"function"`
	Args     interface{} `json:"args"`
}

type CreateTaskResp struct {
	TaskId string `json:"task_id"`
}

type TaskInfoResp struct {
	Status        int    `json:"status"`
	ErrnoCode     int    `json:"errno_code"`
	Subpath       string `json:"subpath"`
	InfoHash      string `json:"info_hash"`
	DownloadSize  int64  `json:"download_size"`
	DownloadSpeed int64  `json:"download_speed"`
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	Gcid          string `json:"gcid"`
	Origin        struct {
		RecvBytes int `json:"recv_bytes"`
	} `json:"origin"`
	P2S struct {
		RecvBytes int `json:"recv_bytes"`
	} `json:"p2s"`
	P2P struct {
		RecvBytes int `json:"recv_bytes"`
	} `json:"p2p"`
	Dcdn struct {
		RecvBytes int `json:"recv_bytes"`
	} `json:"dcdn"`
	Index        int `json:"index"`
	SubFilecount int `json:"sub_filecount"`
}

func Download(ctx context.Context, url string, path string, isWait bool) (resp interface{}, err error) {

	//log.Infof("webhook.Download start. url:%v path:%v", url, path)

	createTaskReq := Input{
		Module:   "download",
		Function: "Create",
		Args: map[string]interface{}{
			"url":       url,
			"file_name": filepath.Base(path),
			"base_path": filepath.Dir(path),
		},
	}
	createResp := &CreateTaskResp{}

	err = callWebhook(createTaskReq, createResp)
	if err != nil {
		log.Errorf("callWebhook err:%v createTaskReq:%+v", err, createTaskReq)
		return
	}

	//log.Printf("create task resp:%+v", createResp)

	if createResp.TaskId == "" {
		err = errors.New("empty createResp.TaskId")
		return
	}

	startTaskReq := Input{
		Module:   "download",
		Function: "Start",
		Args: map[string]interface{}{
			"task_id": createResp.TaskId,
		},
	}

	startTaskResp := map[string]interface{}{}
	err = callWebhook(startTaskReq, &startTaskResp)
	if err != nil {
		log.Errorf("callWebhook err:%v startTaskReq:%+v", err, createTaskReq)
		return
	}

	if !isWait {
		return
	}

	taskInfoReq := Input{
		Module:   "download",
		Function: "Info",
		Args: map[string]interface{}{
			"task_id": createResp.TaskId,
		},
	}

	round := 0
	for {

		time.Sleep(10 * time.Second)

		round++

		taskInfoResp := &TaskInfoResp{}
		err = callWebhook(taskInfoReq, taskInfoResp)
		if err != nil {
			log.Errorf("callWebhook err:%v taskInfoReq:%+v", err, taskInfoReq)
			time.Sleep(1 * time.Second)
			if round > 5 {
				log.Infof("WAIT FOR ROUND TOO LONG 1")
				break
			}
			continue
		}
		//log.Infof("round%v taskInfo:%+v", round, taskInfoResp)
		if taskInfoResp.FileSize <= 0 {
			taskInfoResp.FileSize = -1
		}
		log.Infof("round%v progress (%v/%v)%.2f%% file:%v",
			round, utils.GetShowSize(taskInfoResp.DownloadSize), utils.GetShowSize(taskInfoResp.FileSize),
			float64(taskInfoResp.DownloadSize)/float64(taskInfoResp.FileSize)*100, taskInfoResp.FileName)

		if taskInfoResp.DownloadSize == taskInfoResp.FileSize {
			log.Infof("DOWNLOAD COMPLETE: %v", path)
			break
		}
		if round > 10 {
			log.Infof("WAIT FOR ROUND TOO LONG")
			break
		}

	}

	return
}

func callWebhook(input Input, resp interface{}) (err error) {
	reqURL := "http://127.0.0.1:10600/webhook"

	reqBytes, _ := json.Marshal(input)
	//log.Infof("webhook req:%+v", string(reqBytes))

	reqBuf := bytes.NewBuffer(reqBytes)

	httpReq, _ := http.NewRequest("POST", reqURL, reqBuf)

	respCode, respBytes, err := netutil.HttpDo(httpReq)
	if err != nil {
		return
	}
	_ = respCode
	//log.Printf("respCode:%v resp:%v", respCode, string(respBytes))

	if resp != nil {
		err = json.Unmarshal(respBytes, resp)
		if err != nil {
			return
		}
	}

	return
}
