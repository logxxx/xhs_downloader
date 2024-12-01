package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/logxxx/utils/netutil"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type UploadFileResp struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	Bytes         int    `json:"bytes"`
	CreatedAt     int    `json:"created_at"`
	Filename      string `json:"filename"`
	Purpose       string `json:"purpose"`
	Status        string `json:"status"`
	StatusDetails string `json:"status_details"`
}

func main() {
	run("D:\\mytest\\mywork\\xhs_downloader\\backend\\design\\005_gpt\\test1.png")
}

func run(imgFilePath string) {
	reqURL := "https://api.moonshot.cn/v1/chat/completions"

	imgFile, err := os.Open(imgFilePath)
	if err != nil {
		return
	}
	defer imgFile.Close()

	// 创建一个multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	contentType := w.FormDataContentType()
	// 添加文件字段
	fw, err := w.CreateFormFile("file", filepath.Base(imgFilePath))
	if err != nil {
		panic(err)
	}
	if _, err = io.Copy(fw, imgFile); err != nil {
		panic(err)
	}

	// 结束multipart form
	w.Close()

	uploadFileURL := "https://api.moonshot.cn/v1/files"
	httpReq1, _ := http.NewRequest("POST", uploadFileURL, &b)
	httpReq1.Header.Set("Authorization", "Bearer sk-EWjjAI2HOZyBQy9PLjynNi3IKwZXYEmAtrLtJadENmqYlKST")
	httpReq1.Header.Set("Content-Type", contentType)

	respCode, respBytes, err := netutil.HttpDo(httpReq1)
	if err != nil {
		return
	}
	log.Printf("upload file respCode:%v resp:%v", respCode, string(respBytes))

	uploadFileResp := &UploadFileResp{}

	err = json.Unmarshal(respBytes, uploadFileResp)
	if err != nil {
		return
	}

	//获取文件内容
	getFileInfoURL := fmt.Sprintf("https://api.moonshot.cn/v1/files/%v/content", uploadFileResp.ID)
	log.Printf("getFileInfoURL:%v", getFileInfoURL)
	httpReq2, _ := http.NewRequest("GET", getFileInfoURL, nil)
	httpReq2.Header.Set("Authorization", "Bearer sk-EWjjAI2HOZyBQy9PLjynNi3IKwZXYEmAtrLtJadENmqYlKST")

	respCode, respBytes, err = netutil.HttpDo(httpReq2)
	if err != nil {
		return
	}
	log.Printf("get file content respCode:%v resp:%v", respCode, string(respBytes))

	return

	content := `{
    "model": "moonshot-v1-8k",
	"refs": ["%v"],
    "messages": [
        {
            "role": "system",
            "content": "MM智能助手"
        },
        { 
			"role": "user",
			"content": "这个图片里有什么"
		}
    ],
    "temperature": 0.3
}`

	content = fmt.Sprintf(content, uploadFileResp.ID)

	reqBuf := bytes.NewBufferString(content)

	httpReq, _ := http.NewRequest("POST", reqURL, reqBuf)
	httpReq.Header.Set("Authorization", "Bearer sk-EWjjAI2HOZyBQy9PLjynNi3IKwZXYEmAtrLtJadENmqYlKST")
	httpReq.Header.Set("Content-Type", "application/json")

	respCode, respBytes, err = netutil.HttpDo(httpReq)
	if err != nil {
		return
	}
	log.Printf("respCode:%v resp:%v", respCode, string(respBytes))
	return
}
