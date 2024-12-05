package main

import (
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/model"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestRun(t *testing.T) {
	work := model.Work{NoteID: "test_noteid_1"}
	resp := map[string]interface{}{}
	code, err := netutil.HttpPost("http://47.119.170.71:8088/send_work", work, &resp)
	log.Printf("send_work code:%v resp:%+v err:%v", code, resp, err)
}

func TestRun2(t *testing.T) {
	work := &model.Work{}
	code, err := netutil.HttpGet("http://47.119.170.71:8088/recv_work", work)
	log.Printf("recv_work code:%v resp:%+v err:%v", code, work, err)
}

func TestRun3(t *testing.T) {
	code := 0
	var err error
	result := &blogmodel.ParseBlogResp{NoteID: "test_send_work_result_1"}
	resp := map[string]interface{}{}
	code, err = netutil.HttpPost("http://47.119.170.71:8088/send_work_result", result, &resp)
	log.Printf("send_work_result code:%v resp:%+v err:%v", code, resp, err)

	result2 := &blogmodel.ParseBlogResp{}
	code, err = netutil.HttpGet("http://47.119.170.71:8088/recv_work_result", result2)
	log.Printf("recv_work_result code:%v resp:%+v err:%v", code, result2, err)
}
