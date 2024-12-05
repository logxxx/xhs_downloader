package remote_work

import (
	"github.com/logxxx/utils/netutil"
	"github.com/logxxx/xhs_downloader/biz/blog/blogmodel"
	"github.com/logxxx/xhs_downloader/biz/mydp"
	"testing"
)

func TestStartWaitForWork(t *testing.T) {
	//https://www.xiaohongshu.com/explore/674ee40a00000000070243d6?xsec_token=AB_bLHT0yR9Zkw8HwmdqFpKDGXjIsRNroVeVyjUSPM4jU=&xsec_source=pc_feed
	mydp.SendWork("https://www.xiaohongshu.com/explore/674ee40a00000000070243d6?xsec_token=AB_bLHT0yR9Zkw8HwmdqFpKDGXjIsRNroVeVyjUSPM4jU=&xsec_source=pc_feed",
		"674ee40a00000000070243d6",
		"AB_bLHT0yR9Zkw8HwmdqFpKDGXjIsRNroVeVyjUSPM4jU=")
}

func TestStartRecvRemoteWorkResult(t *testing.T) {
	workResult := &blogmodel.ParseBlogResp{}
	code, err := netutil.HttpGet("http://47.119.170.71:8088/recv_work_result", workResult)
	t.Logf("recv_work_result code:%v err:%v resp:%+v", code, err, workResult)
}
